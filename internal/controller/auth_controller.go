package controller

import (
	"context"
	"eazy-quizy-auth/internal/service"
	auth "eazy-quizy-auth/pkg/api/v1"
)

type AuthController struct {
	auth.UnimplementedAuthServiceServer
	authService   service.AuthService
	friendService service.FriendService
}

func NewAuthController(authService service.AuthService, friendService service.FriendService) *AuthController {
	return &AuthController{
		authService:   authService,
		friendService: friendService,
	}
}

func (c *AuthController) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	// Implement registration logic
	panic("implement me")
}

func (c *AuthController) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// Implement login logic
	panic("implement me")
}

func (c *AuthController) ValidateToken(ctx context.Context, req *auth.ValidateTokenRequest) (*auth.ValidateTokenResponse, error) {
	// Implement token validation logic
	panic("implement me")
}
