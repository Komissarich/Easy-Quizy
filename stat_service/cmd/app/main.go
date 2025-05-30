package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"quiz_app/internal/config"
	"quiz_app/internal/middleware"
	"quiz_app/internal/statistics/repository"
	"quiz_app/internal/statistics/service"
	api "quiz_app/pkg/api/v1"
	"quiz_app/pkg/faults"
	"quiz_app/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
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

	// Logger
	ctx, err := logger.New(ctx)
	if err != nil {
		panic(fmt.Errorf("logger error: %v", err))
	}
	log = logger.GetLoggerFromCtx(ctx)
	log.Info(ctx, "logger started")

	// Config
	cfg, err = config.New()
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to load config: %v", zap.Error(err)))

	}
	log.Info(ctx, "config loaded")

	// TCP Connection
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to listen: %v", zap.Error(err)))
	}
	log.Info(ctx, fmt.Sprintf("listening gRPC on port %d", cfg.GRPCPort))

	// Repository
	var repo *repository.Repository

	err = faults.Retry(
		func() error {
			repo, err = repository.NewRepository(ctx, cfg)
			return err
		}, cfg.Postgres.MaxRetries, cfg.Postgres.BaseDelay)

	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to connect to database: %v", zap.Error(err)))
	}
	log.Info(ctx, "new repository created")

	// Server
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggerUnaryInterceptor,
		),
	)
	log.Info(ctx, "server started")

	// Service
	service := service.New(ctx, repo)

	api.RegisterStatisticsServer(server, service)
	log.Info(ctx, "gRPC service started")
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(ctx, fmt.Sprintf("failed to serve: %v", zap.Error(err)))
		}
	}()

	select {
	case <-ctx.Done():
		server.GracefulStop()
		repo.CloseConn()
		log.Info(ctx, "server stopped")
	}
}
