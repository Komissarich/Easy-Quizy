package service

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/utils"
	"errors"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, user *entity.User) (string, error)
	Login(ctx context.Context, email, password string) (string, *entity.User, error)
	ValidateToken(ctx context.Context, token string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	Logout(ctx context.Context, token string) error
}

type authService struct {
	userRepo   repository.UserRepository
	Jwt        *config.JWTConfig
	jwtService *JWTService
	l          *logger.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwt *config.JWTConfig, l *logger.Logger) *authService {
	return &authService{
		userRepo: userRepo,
		Jwt:      jwt,
		l:        l,
	}
}

func (s *authService) Register(ctx context.Context, user *entity.User) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("can't hash password: %w", err)
	}

	userID, err := s.userRepo.SaveUser(ctx, user.Email, passHash)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return "", fmt.Errorf("can't save user: %w", err)
		}

		return "", fmt.Errorf("can't save user: %w", err)
	}

	userIDUint := strconv.Itoa(int(userID))

	return userIDUint, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *entity.User, error) {
	user, err := s.userRepo.User(ctx, email)
	if err != nil {
		return "", nil, fmt.Errorf("can't find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, entity.ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(s.Jwt.Secret, user.ID, user.Email, 1*time.Hour)
	if err != nil {
		return "", nil, fmt.Errorf("can't generate token: %w", err)
	}

	return token, user, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*entity.User, error) {
	return s.jwtService.ValidateToken(ctx, token)
}

func (s *authService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't parse user id: %w", err)
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't find user by id: %w", err)
	}

	return user, nil
}

func (s *authService) UpdateUser(ctx context.Context, user *entity.User) error {
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("can't update user: %w", err)
	}

	return nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.jwtService.InvalidateToken(ctx, token)
}
