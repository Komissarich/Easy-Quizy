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
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthController struct {
	v1.UnimplementedAuthServiceServer
	authService    service.AuthService
	friendService  service.FriendService
	quizzesService service.QuizzesService

	l *logger.Logger
}

func NewAuthController(authService service.AuthService, friendService service.FriendService, quizzesService service.QuizzesService, l *logger.Logger) *AuthController {
	return &AuthController{
		authService:    authService,
		friendService:  friendService,
		quizzesService: quizzesService,
		l:              l,
	}
}

func (c *AuthController) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	if req.Email == "" || !isValidEmail(req.Email) {
		c.l.Warn("Invalid or empty email", zap.String("email", req.Email))
		return nil, status.Error(codes.InvalidArgument, "invalid or empty email")
	}
	if req.Password == "" || len(req.Password) < 8 {
		c.l.Warn("Invalid password", zap.Int("password_length", len(req.Password)))
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	if req.Username == "" {
		c.l.Warn("username is required")
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	user := &entity.User{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
	}

	c.l.Info("Trying to register user", zap.String("email", req.Email), zap.Any("username", req.Username))

	userID, err := c.authService.Register(ctx, user)
	if err != nil {
		c.l.Error("Failed to register user", zap.Error(err))
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "unique constraint") ||
			strings.Contains(errMsg, "duplicate key") ||
			strings.Contains(errMsg, "duplicate entry") ||
			strings.Contains(errMsg, "violates unique") {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	c.l.Info("User registered successfully", zap.String("user_id", userID))

	return &v1.RegisterResponse{UserId: userID}, nil
}

func (c *AuthController) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		c.l.Warn("email and password are required")
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	c.l.Info("Trying to login user", zap.String("email", req.Email))

	token, user, err := c.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		c.l.Error("Failed to login user", zap.Error(err))
		if strings.Contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Unauthenticated, "invalid email or password")
	}

	c.l.Info("User logged in successfully", zap.Uint64("user_id", user.ID))

	return &v1.LoginResponse{
		Token: token,
		User:  convertUserToProto(user),
	}, nil
}

func (c *AuthController) ValidateToken(ctx context.Context, req *v1.ValidateTokenRequest) (*v1.ValidateTokenResponse, error) {
	if req.Token == "" {
		c.l.Warn("token is required")
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	c.l.Info("Trying to validate token", zap.String("token", req.Token))

	user, err := c.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		c.l.Error("Failed to validate token", zap.Error(err))
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	c.l.Info("Token validated successfully", zap.Uint64("user_id", user.ID))

	return &v1.ValidateTokenResponse{
		User: convertUserToProto(user),
	}, nil
}

func (c *AuthController) GetMe(ctx context.Context, req *v1.GetMeRequest) (*v1.UserResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok || user == nil {
		c.l.Error("User not found in context")

		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	c.l.Info("GetMe success", zap.Uint64("user_id", user.ID))

	return convertUserToProto(user), nil
}

func (c *AuthController) UpdateMe(ctx context.Context, req *v1.UpdateMeRequest) (*v1.UserResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok || user == nil {
		c.l.Warn("User not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	c.l.Info("Updating user", zap.Uint64("user_id", user.ID))

	if req.Email != nil {
		if !isValidEmail(*req.Email) {
			c.l.Warn("Invalid email", zap.String("email", *req.Email))
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}
		user.Email = *req.Email
		c.l.Info("Updating email", zap.String("email", *req.Email))
	}
	if req.Password != nil {
		if len(*req.Password) < 8 {
			c.l.Warn("Invalid password", zap.Int("password_length", len(*req.Password)))
			return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
		}
		user.Password = *req.Password
		c.l.Info("Updating password")
	}

	if err := c.authService.UpdateUser(ctx, user); err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return nil, status.Error(codes.AlreadyExists, "email already exists")
		}

		return nil, status.Error(codes.Internal, "failed to update user")
	}

	c.l.Info("User updated successfully", zap.Uint64("user_id", user.ID))

	return convertUserToProto(user), nil
}

func (c *AuthController) AddFriend(ctx context.Context, req *v1.AddFriendRequest) (*v1.FriendResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		c.l.Warn("user not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}
	c.l.Info("AddFriend method called", zap.Uint64("user_id", user.ID), zap.String("friend_id", req.FriendId))

	c.l.Info("User found", zap.Uint64("user_id", user.ID))

	if err := c.friendService.AddFriend(ctx, user.ID, req.FriendId); err != nil {
		if errors.Is(err, service.ErrCannotAddYourself) {
			c.l.Warn("cannot add yourself as a friend")
			return nil, status.Error(codes.InvalidArgument, "cannot add yourself as a friend")
		}
		if errors.Is(err, service.ErrAlreadyFriends) {
			c.l.Warn("users are already friends")
			return nil, status.Error(codes.AlreadyExists, "users are already friends")
		}
		c.l.Error("Failed to add friend", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to add friend")
	}

	c.l.Info("Friend added successfully", zap.Uint64("user_id", user.ID), zap.String("friend_id", req.FriendId))

	return &v1.FriendResponse{
		Success: true,
		Message: "friend added successfully",
	}, nil
}

func (c *AuthController) RemoveFriend(ctx context.Context, req *v1.RemoveFriendRequest) (*v1.FriendResponse, error) {
	user, ok := utils.GetUser(ctx)
	if !ok {
		c.l.Warn("user not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
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
	user, ok := utils.GetUser(ctx)
	if !ok {
		c.l.Warn("user not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
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

func (c *AuthController) AddFavoriteQuiz(ctx context.Context, req *v1.AddFavoriteQuizRequest) (*v1.FavoriteQuizResponse, error) {
	userID, ok := utils.GetUser(ctx)
	if !ok {
		c.l.Warn("user not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	if err := c.quizzesService.AddFavoriteQuiz(userID.ID, req.QuizId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add favorite quiz: %v", err)
	}

	return &v1.FavoriteQuizResponse{Success: true, Message: "Quiz added to favorites"}, nil
}

func (c *AuthController) GetFavoriteQuizzes(ctx context.Context, req *v1.GetFavoriteQuizzesRequest) (*v1.FavoriteQuizzesResponse, error) {
	userID, ok := utils.GetUser(ctx)
	if !ok {
		c.l.Warn("user not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	favoriteQuizzes, err := c.quizzesService.GetFavoriteQuizzes(userID.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get favorite quizzes: %v", err)
	}

	return &v1.FavoriteQuizzesResponse{QuizIds: favoriteQuizzes}, nil
}

func (c *AuthController) RemoveFavoriteQuiz(ctx context.Context, req *v1.RemoveFavoriteQuizRequest) (*v1.FavoriteQuizResponse, error) {
	userID, ok := utils.GetUser(ctx)
	if !ok {
		c.l.Warn("user not found in context")
		return nil, status.Error(codes.Unauthenticated, "user not found in context")
	}

	if err := c.quizzesService.RemoveFavoriteQuiz(userID.ID, req.QuizId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove favorite quiz: %v", err)
	}

	return &v1.FavoriteQuizResponse{Success: true, Message: "Quiz removed from favorites"}, nil
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
		Id:        userID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
