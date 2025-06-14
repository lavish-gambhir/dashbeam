package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/config"
	"github.com/redis/go-redis/v9"
)

type MessageQueue interface {
	Publish(ctx context.Context, topic string, event Event) error
	Subscribe(ctx context.Context, topic string, handler func(Event) error, opts *SubscribeOptions)
}

type Subscriber struct {
	pubsub  *redis.PubSub
	handler func(Event) error
}

type SubscribeOptions struct {
	Pattern      bool
	MaxRetries   int
	RetryBackoff []time.Duration
}

func DefaultSubscriberOptions() *SubscribeOptions {
	return &SubscribeOptions{
		Pattern:    true,
		MaxRetries: 3,
		RetryBackoff: []time.Duration{
			time.Second,
			time.Second * 5,
			time.Second * 10,
		},
	}
}

type RedisQueue struct {
	client *redis.Client
	logger *slog.Logger

	mu          sync.Mutex
	done        chan struct{}
	subscribers map[string]*Subscriber
}

func NewRedisQueue(ctx context.Context, cfg *config.AppConfig, logger *slog.Logger) (*RedisQueue, error) {
	client := redis.NewClient(&redis.Options{
		Addr:            cfg.Redis.Addr,
		Password:        cfg.Redis.Password,
		DB:              cfg.Redis.DB,
		PoolSize:        cfg.Redis.PoolSize,
		ReadTimeout:     cfg.Redis.Timeout,
		WriteTimeout:    cfg.Redis.Timeout,
		MaxRetryBackoff: time.Second * 2,
		MinRetryBackoff: time.Millisecond * 100,
	})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &RedisQueue{
		client:      client,
		logger:      logger,
		done:        make(chan struct{}),
		subscribers: make(map[string]*Subscriber),
	}, nil
}

func (r *RedisQueue) Publish(ctx context.Context, topic string, event Event) error {
	if event.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		event.ID = id
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return apperr.Wrapf(err, apperr.JSONEncodingFailed, "%s, %v", "failed to marshal event", event.ID)
	}

	pipe := r.client.Pipeline()
	pipe.Publish(ctx, topic, payload)
	eventKey := fmt.Sprintf("event:%s", event.ID)
	pipe.Set(ctx, eventKey, payload, time.Hour*24) // storing the event for replay
	if _, err := pipe.Exec(ctx); err != nil {
		return apperr.Wrapf(err, apperr.RedisPipeExecFailed, "%s, %v", "failed to exec redis pipe, for event", event.ID)
	}
	r.logger.Info("published message", slog.String("topic", topic), slog.String("event_id", event.ID.String()), slog.String("event_type", event.Type.String()))
	return nil
}

func (r *RedisQueue) Subscribe(ctx context.Context, topic string, handler func(Event) error, opts *SubscribeOptions) {
	if opts == nil {
		opts = DefaultSubscriberOptions()
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, exists := r.subscribers[topic]; exists {
		existing.pubsub.Close()
	}

	var pubsub *redis.PubSub
	if opts.Pattern {
		pubsub = r.client.PSubscribe(ctx, topic)
	} else {
		pubsub = r.client.Subscribe(ctx, topic)
	}
	r.subscribers[topic] = &Subscriber{
		pubsub:  pubsub,
		handler: handler,
	}
	go r.handleSubscription(ctx, topic, opts)
}

func (r *RedisQueue) handleSubscription(ctx context.Context, topic string, opts *SubscribeOptions) {
	var consecutiveFailures int
	var backoffIdx int

	maxBackoffAttempts := 5
	for {
		if err := r.processSubscription(ctx, topic); err != nil {
			consecutiveFailures++
			r.logger.Error("subscription processing failed", slog.String("topic", topic), slog.Int("consecutive_failures", consecutiveFailures), slog.Any("err", err))
			var backoff time.Duration
			if backoffIdx < len(opts.RetryBackoff) {
				backoff = opts.RetryBackoff[backoffIdx]
				backoffIdx++
			} else {
				if consecutiveFailures > maxBackoffAttempts {
					r.logger.Error("handleSubscription: exceeded max retry attempts, stopping retries", slog.String("topic", topic))
					r.mu.Lock()
					if sub, exists := r.subscribers[topic]; exists {
						sub.pubsub.Close()
						delete(r.subscribers, topic)
					}
					r.mu.Unlock()
					return
				}
				backoff = opts.RetryBackoff[len(opts.RetryBackoff)-1]
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
				continue
			}
		}
		consecutiveFailures = 0
		backoffIdx = 0
	}
}

func (r *RedisQueue) processSubscription(ctx context.Context, topic string) error {
	r.mu.Lock()
	subscriber, exists := r.subscribers[topic]
	if !exists {
		r.mu.Unlock()
		return fmt.Errorf("subscriber not found for topic:%s", topic)
	}
	r.mu.Unlock()
	ch := subscriber.pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return fmt.Errorf("subscription channel closed")
			}
			var ev Event
			if err := json.Unmarshal([]byte(msg.Payload), &ev); err != nil {
				r.logger.Error("failed to unmarshal event", slog.String("topic", topic), slog.Any("err", err))
				continue
			}
			if err := r.executeWithRetry(ctx, ev, subscriber.handler, topic); err != nil {
				r.logger.Error("handler failed with retries", slog.String("topic", topic), slog.String("event_id", ev.ID.String()), slog.Any("err", err))
				if err := r.storeFailedEvent(ctx, topic, ev, err); err != nil {
					r.logger.Error("failed to store failed event", slog.String("topic", topic), slog.String("event_id", ev.ID.String()), slog.Any("err", err))
				}
			}

		}
	}
}

func (r *RedisQueue) executeWithRetry(ctx context.Context, event Event, handler func(Event) error, topic string) error {
	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		_, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := handler(event); err != nil {
			lastErr = err
			r.logger.Error("handler failed, retrying", slog.String("topic", topic), slog.String("event_id", event.ID.String()), slog.Int("attempt", attempt+1), slog.Any("err", err))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second * time.Duration(attempt+1)):
				continue
			}
		} else {
			return nil
		}
	}
	return lastErr
}

func (r *RedisQueue) storeFailedEvent(ctx context.Context, topic string, event Event, err error) error {
	failedEvent := struct {
		Event      Event     `json:"event"`
		Error      string    `json:"error"`
		FailedAt   time.Time `json:"failed_at"`
		RetryCount int       `json:"retry_count"`
	}{
		Event:      event,
		Error:      err.Error(),
		FailedAt:   time.Now(),
		RetryCount: 0,
	}
	data, err := json.Marshal(failedEvent)
	if err != nil {
		return apperr.Wrapf(err, apperr.JSONEncodingFailed, "%s,topic=%s,event=%s", "failed to marshal event", topic, event.ID.String())
	}
	key := fmt.Sprintf("failed_events=%s:%s", topic, event.ID.String())
	return r.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (r *RedisQueue) Close() error {
	close(r.done)
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, sub := range r.subscribers {
		sub.pubsub.Close()
	}
	return r.client.Close()
}

func (r *RedisQueue) PSubscribe(ctx context.Context, pattern string, handler func(Event) error) {
	opts := DefaultSubscriberOptions()
	opts.Pattern = true
	r.Subscribe(ctx, pattern, handler, opts)
}

func (r *RedisQueue) GetEventHistory(ctx context.Context, eventID uuid.UUID) (*Event, error) {
	eventKey := fmt.Sprintf("event:%s", eventID)
	data, err := r.client.Get(ctx, eventKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, apperr.Wrapf(err, apperr.RedisNoEvent, "id:%v", eventID)
		}
		return nil, apperr.Wrapf(err, apperr.RedisUnknown, "%s,id:%v", "failed to get event", eventID)
	}
	var ev Event
	if err := json.Unmarshal(data, &ev); err != nil {
		return nil, apperr.Wrapf(err, apperr.JSONDecodingFailed, "%s,id:%v", "failed to unmarshal event", eventID)
	}
	return &ev, nil
}

func (r *RedisQueue) ReplayFailedEvents(ctx context.Context, topic string) error {
	pattern := fmt.Sprintf("failed_events:%s:*", topic)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return apperr.Wrapf(err, apperr.RedisUnknown, "%s,topic=%s", "failed to get failed events", topic)
	}
	pipe := r.client.Pipeline()
	for _, key := range keys {
		pipe.Get(ctx, key)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return apperr.Wrapf(err, apperr.RedisUnknown, "%s,topic=%s", "failed to execute redis cmds", topic)
	}
	for _, cmd := range cmds {
		getCmd, ok := cmd.(*redis.StringCmd)
		if !ok {
			return fmt.Errorf("failed to cast cmd to *redisCmd, cmd:%s", cmd)
		}
		data, err := getCmd.Bytes()
		if err != nil {
			r.logger.Error("failed to convert cmd to bytes", slog.String("topic", topic))
			continue
		}
		var failedEvent struct {
			Event Event `json:"event"`
		}
		if err := json.Unmarshal(data, &failedEvent); err != nil {
			r.logger.Error("failed to unmarshal cmd bytes", slog.String("topic", topic))
			continue
		}
		if err := r.Publish(ctx, topic, failedEvent.Event); err != nil {
			r.logger.Error("failed to replay event", slog.String("topic", topic), slog.String("event_id", failedEvent.Event.ID.String()), slog.Any("err", err))
			continue
		}
		if len(cmd.Args()) > 0 {
			cmdString := cmd.Args()[1].(string)
			r.client.Del(ctx, cmdString)
		}
	}
	return nil
}
