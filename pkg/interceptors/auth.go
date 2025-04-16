package interceptors

import (
	"context"
	"eazy-quizy-auth/internal/service"
	"eazy-quizy-auth/pkg/utils"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authService service.AuthService
}

func NewAuthInterceptor(authService service.AuthService) *AuthInterceptor {
	return &AuthInterceptor{authService: authService}
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
			return nil, status.Error(codes.Unauthenticated, "metadata not provided")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")

		user, err := i.authService.ValidateToken(ctx, token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
		}

		ctx = utils.WithUser(ctx, user)

		return handler(ctx, req)
	}
}
