package service

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/redis"
	"eazy-quizy-auth/pkg/utils"
	"fmt"
	"time"

	"go.uber.org/zap"
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
		return "", nil, err
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
		return err
	}

	remainingTTL := time.Until(claims.ExpiresAt.Time)
	return s.redis.Set(ctx, "jwt_blacklist:"+token, "1", remainingTTL)
}

func (s *JWTService) IsTokenValid(ctx context.Context, token string) (bool, error) {
	_, err := s.redis.Get(ctx, "jwt_blacklist:"+token)
	if err == nil {
		return false, nil
	}

	_, err = utils.ParseJWT(token, s.Jwt.Secret)
	return err == nil, nil
}
