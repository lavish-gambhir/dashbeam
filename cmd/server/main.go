package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lavish-gambhir/dashbeam/cmd/server/handlers"
	"github.com/lavish-gambhir/dashbeam/internal/config"
	"github.com/lavish-gambhir/dashbeam/pkg/logger"
	"github.com/lavish-gambhir/dashbeam/services/ingestion"
	"github.com/lavish-gambhir/dashbeam/shared/database"
	"github.com/lavish-gambhir/dashbeam/shared/middleware"
)

type App struct {
	config *config.AppConfig
	pool   *pgxpool.Pool
	server *http.Server

	ingestionSvc ingestion.Service
}

func index(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "===dashbeam===")
}

func setupApp(cfg *config.AppConfig, pool *pgxpool.Pool, logger *slog.Logger) *App {
	mux := http.NewServeMux()
	ingestionService := ingestion.New()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:           middleware.Logging(logger)(middleware.Cors(mux)),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	app := &App{
		config:       cfg,
		pool:         pool,
		server:       server,
		ingestionSvc: ingestionService,
	}

	app.registerRoutes()

	return app
}

func (a *App) registerRoutes() {
	http.HandleFunc("/", index)
	http.HandleFunc("/healthz", handlers.HealthCheckHandler)
	http.HandleFunc("/readyz", handlers.ReadyzHandler)
}

func (a *App) Start(ctx context.Context, logger *slog.Logger) <-chan error {
	errC := make(chan error)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("=== dashbeam shutting down ===")
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer func() {
			a.pool.Close()
			stop()
			cancel()
			close(errC)
		}()
		a.server.SetKeepAlivesEnabled(false)
		if err := a.server.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("=== dashbeam shut down complete ===")
	}()

	go func() {
		logger.Info("Listening and serving", "addr", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			errC <- err
		}
	}()

	return errC
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	// load deps
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to laod config: %v", err)
	}
	pool, err := database.NewPostgres(ctx, cfg)
	logger := logger.NewSlogger(cfg)

	app := setupApp(cfg, pool, logger)
	if err := <-app.Start(ctx, logger); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
