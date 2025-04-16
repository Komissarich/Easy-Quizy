package service

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/redis"
	"eazy-quizy-auth/pkg/utils"
	"errors"
	"fmt"
	"time"
)

type JWTService struct {
	userRepo  repository.UserRepository
	redis     *redis.Client
	jwtSecret string
}

func NewJWTService(userRepo repository.UserRepository, jwtSecret string, redis *redis.Client) *JWTService {
	return &JWTService{
		userRepo:  userRepo,
		redis:     redis,
		jwtSecret: jwtSecret,
	}
}

func (s *JWTService) GenerateToken(ctx context.Context, email, password string) (string, *entity.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, entity.ErrUserNotFound
	}

	token, err := utils.GenerateJWT(s.jwtSecret, user.ID, user.Email, 24*time.Hour)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *JWTService) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	claims, err := utils.ParseJWT(token, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, entity.ErrUserNotFound
	}

	if user.Email != claims.Email {
		return nil, errors.New("token email mismatch")
	}

	return user, nil
}

func (s *JWTService) InvalidateToken(ctx context.Context, token string) error {
	claims, err := utils.ParseJWT(token, s.jwtSecret)
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

	_, err = utils.ParseJWT(token, s.jwtSecret)
	return err == nil, nil
}
