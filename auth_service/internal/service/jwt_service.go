package service

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/redis"
	"eazy-quizy-auth/pkg/utils"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type JWTService struct {
	userRepo repository.UserRepository
	redis    *redis.Client
	Jwt      *config.JWTConfig
	l        *logger.Logger
}

func NewJWTService(userRepo repository.UserRepository, redis *redis.Client, Jwt *config.JWTConfig, l *logger.Logger) *JWTService {
	return &JWTService{
		userRepo: userRepo,
		redis:    redis,
		Jwt:      Jwt,
		l:        l,
	}
}

func (s *JWTService) GenerateToken(ctx context.Context, email, password string) (string, *entity.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, entity.ErrUserNotFound
	}

	token, err := utils.GenerateJWT(s.Jwt.Secret, user.ID, user.Email)
	if err != nil {
		return "", nil, fmt.Errorf("can't generate token: %w", err)
	}

	return token, user, nil
}

func (s *JWTService) ValidateToken(ctx context.Context, authHeader string) (*entity.User, error) {
	tokenString, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		s.l.Error("Failed to extract token", zap.Error(err))
		return nil, fmt.Errorf("invalid auth header: %w", err)
	}

	claims, err := utils.ParseJWT(tokenString, s.Jwt.Secret)
	if err != nil {
		s.l.Error("Failed to parse token",
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, entity.ErrUserNotFound
	}

	if user.ID == 0 || user.Email == "" {
		return nil, fmt.Errorf("invalid user data")
	}

	return user, nil
}

func (s *JWTService) InvalidateToken(ctx context.Context, token string) error {
	claims, err := utils.ParseJWT(token, s.Jwt.Secret)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	remainingTTL := time.Until(claims.ExpiresAt.Time)

	return s.redis.Set(ctx, "jwt_blacklist:"+token, "1", remainingTTL)
}
