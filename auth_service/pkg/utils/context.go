package utils

import (
	"context"
	"eazy-quizy-auth/internal/entity"
)

type contextKey string

const (
	UserKey contextKey = "user"
)

func WithUser(ctx context.Context, user *entity.User) context.Context {
	return context.WithValue(ctx, UserKey, user)
}

func GetUser(ctx context.Context) (*entity.User, bool) {
	user, ok := ctx.Value(UserKey).(*entity.User)
	return user, ok
}
