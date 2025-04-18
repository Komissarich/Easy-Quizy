package middleware

import (
	"context"
	"fmt"
	"quiz_app/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequestIDUnaryInterseptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if ctx.Value(logger.RequestID) == nil {
		ctx = context.WithValue(ctx, logger.RequestID, uuid.New().String())
	}
	return handler(ctx, req)
}

func LoggerUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if ctx.Value(logger.LoggerKey) == nil {
		var err error
		ctx, err = logger.New(ctx)
		if err != nil {
			return nil, fmt.Errorf("logger creation error: %v", err)
		}
	}
	ctx = context.WithValue(ctx, logger.RequestID, uuid.New().String())
	logger.GetLoggerFromCtx(ctx).Info(ctx, "request",
		zap.String("method", info.FullMethod),
		zap.Time("request time", time.Now()))
	return handler(ctx, req)
}

func ErrorsUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	smth, err := handler(ctx, req)
	switch status.Code(err) {
	case codes.NotFound:
		logger.GetLoggerFromCtx(ctx).Error(ctx, "NOT_FOUND")
		return smth, nil
	case codes.Canceled:
		logger.GetLoggerFromCtx(ctx).Error(ctx, "CANCELED")
		return smth, nil
	case codes.Unknown:
		logger.GetLoggerFromCtx(ctx).Error(ctx, "UNKNOWN")
		return handler(ctx, req)
	default:
		return handler(ctx, req)
	}
}
