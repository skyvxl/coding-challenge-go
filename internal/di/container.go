package di

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"learn/internal/config"
	"learn/internal/db"
	"learn/internal/handler"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunApp() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := newConfig()
	if err != nil {
		return err
	}
	pool, err := newDatabase(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()
	log.Println("Database connected")
	q := db.New(pool)
	subscriptionHandler := handler.NewHandler(q)
	err = runServer(ctx, subscriptionHandler.Routes(), cfg)
	if err != nil {
		return err
	}
	return nil
}

func newDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed create connection: %w", err)
	}
	if err = conn.Ping(ctx); err != nil {
		conn.Close()
		if errors.Is(err, context.Canceled) {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed connect to the database: %w", err)
	}
	return conn, nil
}

func newConfig() (*config.Config, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

func runServer(ctx context.Context, h http.Handler, cfg *config.Config) error {
	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serverErrors := make(chan error, 1)
	log.Printf("Server started on %s\n", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("failed to listen and serve: %w", err)
		}
	}()
	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
		log.Println("Server stopped gracefully")
	}
	return nil
}
