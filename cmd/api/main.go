package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladfc/ghira/internal/graphql"
	"github.com/vladfc/ghira/internal/graphql/resolver"
	"github.com/vladfc/ghira/internal/user"

	"github.com/vladfc/ghira/internal/config"
	"github.com/vladfc/ghira/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Log.Level,
	}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// db, err := database.Open(ctx, cfg.Database)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	now := time.Now().UTC()

	seedUsers := []user.User{
		{
			ID:           "8b7d8a0a-6c7e-4bb2-bb6a-68c73a4b6b01",
			Email:        "olena.koval@example.com",
			Username:     "olena",
			PasswordHash: "seed_password_hash_1",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "1c4a3df7-4a5d-41c5-8708-3a8a94f52c21",
			Email:        "maksym.shevchenko@example.com",
			Username:     "maksym",
			PasswordHash: "seed_password_hash_2",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		{
			ID:           "bfa5ef65-7401-47dd-b96e-c0d16f91e5a4",
			Email:        "daria.melnyk@example.com",
			Username:     "daria",
			PasswordHash: "seed_password_hash_3",
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}

	userRepo := user.NewMemoryRepository(seedUsers)
	userService := user.NewService(userRepo)

	graphqlHandler := graphql.NewHandler(&resolver.Resolver{
		UserService: userService,
	})
	srv := server.New(cfg.HTTP, logger, graphqlHandler)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("http server listening", "addr", cfg.HTTP.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	logger.Info("http server stopped")
	return nil
}
