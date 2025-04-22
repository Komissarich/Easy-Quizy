package interceptors

import (
	"context"
	"eazy-quizy-auth/internal/service"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/utils"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authService service.AuthService
	l           *logger.Logger
}

func NewAuthInterceptor(authService service.AuthService, l *logger.Logger) *AuthInterceptor {
	return &AuthInterceptor{authService: authService, l: l}
}

func (i *AuthInterceptor) Unary(ctx context.Context) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if info.FullMethod == "/auth.AuthService/Login" ||
			info.FullMethod == "/auth.AuthService/Register" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			i.l.Error("Missing metadata in context")
			return nil, status.Error(codes.Unauthenticated, "metadata is required")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			i.l.Error("Missing authorization header")
			return nil, status.Error(codes.Unauthenticated, "authorization header is required")
		}

		user, err := i.authService.ValidateToken(ctx, authHeaders[0])
		if err != nil {
			i.l.Error("Token validation failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = utils.WithUser(ctx, user)

		i.l.Info("Request authenticated",
			zap.String("method", info.FullMethod),
			zap.Uint64("user_id", user.ID),
		)

		return handler(ctx, req)
	}
}
