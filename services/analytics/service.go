package analytics

import (
	"context"
	"log/slog"
	"time"

	"github.com/lavish-gambhir/dashbeam/shared/config"
	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

type Service interface {
	Start(ctx context.Context) error
	Stop() error
}

type service struct {
	messageQueue streaming.MessageQueue
	processor    *EventProcessor
	logger       *slog.Logger
	config       config.AnalyticsConfig
	stopCh       chan struct{}
}

func New(
	messageQueue streaming.MessageQueue,
	processor *EventProcessor,
	config config.AnalyticsConfig,
	logger *slog.Logger,
) Service {
	return &service{
		messageQueue: messageQueue,
		processor:    processor,
		logger:       logger.With("service", "analytics"),
		config:       config,
		stopCh:       make(chan struct{}),
	}
}

func (s *service) Start(ctx context.Context) error {
	s.logger.Info("starting analytics service")

	// Start consuming events from Redis streams
	go s.consumeEvents(ctx)
	return nil
}

func (s *service) Stop() error {
	s.logger.Info("stopping analytics service")
	close(s.stopCh)
	return nil
}

func (s *service) processAvailableEvents(ctx context.Context) error {
	// Subscribe to all event topics using pattern matching
	s.messageQueue.Subscribe(ctx, "events.*", func(event streaming.Event) error {
		events := []streaming.Event{event}
		return s.processor.ProcessEvents(ctx, events)
	}, nil)

	return nil
}

func (s *service) consumeEvents(ctx context.Context) {
	ticker := time.NewTicker(s.config.ProcessingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("context cancelled, stopping event consumption")
			return
		case <-s.stopCh:
			s.logger.Info("stop signal received, stopping event consumption")
			return
		case <-ticker.C:
			if err := s.processAvailableEvents(ctx); err != nil {
				s.logger.Error("failed to process events", slog.Any("error", err))
			}
		}
	}
}
