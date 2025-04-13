package interceptors

import (
	"context"

	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	role = "user_role"
)

func RoleInterceptor(ctx context.Context, allowedRoles []string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		role, ok := ctx.Value(role).(string)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "role missing")
		}

		hasAccess := slices.Contains(allowedRoles, role)

		if !hasAccess {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		return handler(ctx, req)
	}
}
