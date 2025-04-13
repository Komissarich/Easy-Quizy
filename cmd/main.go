package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"quiz_app/internal/config"
	"quiz_app/pkg/logger"
	"quiz_app/pkg/postgres"

	"go.uber.org/zap"
)

var (
	ctx context.Context
	log *logger.Logger
	cfg *config.Config
)

func main() {
	ctx = context.Background()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	ctx, err := logger.New(ctx)
	if err != nil {
		panic(fmt.Errorf("logger error: %v", err))
	}
	log = logger.GetLoggerFromCtx(ctx)
	log.Info(ctx, "logger started")

	cfg, err = config.New()
	if err != nil {
		log.Fatal(ctx, fmt.Sprint("failed to load config", zap.Error(err)))

	}
	log.Info(ctx, "config loaded")

	pool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Error(ctx, fmt.Sprint("failed to connect to postgres", zap.Error(err)))
	} else {
		log.Info(ctx, "connected to postgres")
		log.Info(ctx, fmt.Sprint("pinging postgres: ", pool.Ping(ctx)))
	}
	_, err = net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort))
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to listen: %v", zap.Error(err)))
	}
	select {
	case <-ctx.Done():
		// server.GracefulStop()
		pool.Close()
		log.Info(ctx, "server stopped")
	}
}
