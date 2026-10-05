package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drobyshevv/movie-reservation/internal/app"
	"github.com/drobyshevv/movie-reservation/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	cfg := config.MustLoadConfig()
	log := setupLogger(cfg.Env)

	app, err := app.NewApp(log, *cfg)
	if err != nil {
		log.Error("failed to init app", "error", err)
		return
	}

	chErr := make(chan error, 1)

	go func() {
		chErr <- app.Run()
	}()

	sig, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-chErr:
		log.Error("failed to start app", "error", err)
		return
	case <-sig.Done():
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		err := app.Stop(ctx)
		if err != nil {
			log.Error("failed to stop app", "error", err)
			return
		}
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return log
}
