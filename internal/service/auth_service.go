package service

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/repository"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (*entity.User, error)
	ValidateToken(ctx context.Context, token string) (bool, *entity.User, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *authService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) Register(ctx context.Context, email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("can't hash password: %w", err)
	}

	err = s.userRepo.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return err
		}

		return fmt.Errorf("can't save user: %w", err)
	}

	return nil
}

func (s *authService) Login(ctx context.Context, email string, password string) (*entity.User, error) {
	user, err := s.userRepo.User(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can't find user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("can't compare password: %w", err)
	}

	return user, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (bool, *entity.User, error) {
	// Implement token validation logic
	panic("implement me")
}
