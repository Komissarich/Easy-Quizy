package rest

import (
	"context"
	"fmt"
	"net/http"
	"quiz_app/internal/config"
	"quiz_app/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	api "quiz_app/pkg/api/v1"
)

func Run(ctx context.Context, cfg *config.Config) {
	log := logger.GetLoggerFromCtx(ctx)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := api.RegisterStatisticsHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort),
		opts,
	)
	if err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to start gRPC-gateway: %v", zap.Error(err)))
	}

	log.Info(ctx, "REST proxy started")

	log.Info(ctx, fmt.Sprintf("listening REST on %s:%d", cfg.Host, cfg.HTTPPort))

	if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HTTPPort), mux); err != nil {
		log.Fatal(ctx, fmt.Sprintf("failed to server gRPC-gateway: %v", zap.Error(err)))
	}

}
