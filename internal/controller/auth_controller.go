package controller

import (
	"context"
	"eazy-quizy-auth/internal/entity"
	"eazy-quizy-auth/internal/service"
	v1 "eazy-quizy-auth/pkg/api/v1"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/utils"
	"errors"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	v1.UnimplementedAuthServiceServer
	authService   service.AuthService
	friendService service.FriendService
	l             *logger.Logger
}

func NewAuthController(authService service.AuthService, friendService service.FriendService, l *logger.Logger) *AuthController {
	return &AuthController{
		authService:   authService,
		friendService: friendService,
		l:             l,
	}
}

func (c *AuthController) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email and password are required")
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	userID, err := c.authService.Register(ctx, user)
	if err != nil {
		c.l.Error("Failed to register user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &v1.RegisterResponse{UserId: userID}, nil
}

func (c *AuthController) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	token, user, err := c.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "login failed")
	}

	return &v1.LoginResponse{
		Token: token,
		User:  convertUserToProto(user),
	}, nil
}

func (c *AuthController) GetMe(ctx context.Context, req *v1.GetMeRequest) (*v1.UserResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	return convertUserToProto(user), nil
}

func (c *AuthController) UpdateMe(ctx context.Context, req *v1.UpdateMeRequest) (*v1.UserResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		user.Password = *req.Password
	}

	if err := c.authService.UpdateUser(ctx, user); err != nil {
		c.l.Error("Failed to update user", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return convertUserToProto(user), nil
}

func (c *AuthController) AddFriend(ctx context.Context, req *v1.AddFriendRequest) (*v1.FriendResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	if err := c.friendService.AddFriend(ctx, user.ID, req.FriendId); err != nil {
		if errors.Is(err, service.ErrCannotAddYourself) {
			return nil, status.Error(codes.InvalidArgument, "cannot add yourself as a friend")
		}
		if errors.Is(err, service.ErrAlreadyFriends) {
			return nil, status.Error(codes.AlreadyExists, "users are already friends")
		}
		c.l.Error("Failed to add friend", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add friend")
	}

	return &v1.FriendResponse{
		Success: true,
		Message: "friend added successfully",
	}, nil
}

func (c *AuthController) RemoveFriend(ctx context.Context, req *v1.RemoveFriendRequest) (*v1.FriendResponse, error) {
	user, err := c.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if err := c.friendService.RemoveFriend(ctx, user.ID, req.FriendId); err != nil {
		if errors.Is(err, service.ErrNotFriends) {
			return nil, status.Error(codes.NotFound, "users are not friends")
		}
		c.l.Error("Failed to remove friend", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to remove friend")
	}

	return &v1.FriendResponse{
		Success: true,
		Message: "friend removed successfully",
	}, nil
}

func (c *AuthController) GetFriends(ctx context.Context, req *v1.GetFriendsRequest) (*v1.FriendsListResponse, error) {
	user, err := c.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	friends, err := c.friendService.GetFriends(ctx, user.ID)
	if err != nil {
		c.l.Error("Failed to get friends", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get friends")
	}

	var protoFriends []*v1.UserResponse
	for _, friend := range friends {
		protoFriends = append(protoFriends, convertUserToProto(friend))
	}

	return &v1.FriendsListResponse{Friends: protoFriends}, nil
}

func (c *AuthController) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	if err := c.authService.Logout(ctx, req.Token); err != nil {
		return nil, status.Error(codes.Internal, "logout failed")
	}
	return &v1.LogoutResponse{Success: true}, nil
}

// Вспомогательная функция для конвертации модели пользователя в protobuf сообщение
func convertUserToProto(user *entity.User) *v1.UserResponse {
	userID := strconv.FormatUint(user.ID, 10)
	return &v1.UserResponse{
		Id:       userID,
		Username: user.Username,
		Email:    user.Email,
		// QuizScore:  user.QuizScore,
		// AuthorRank: user.AuthorRank,
		// PlayerRank: user.PlayerRank,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
