//go:build integration

package auth_test

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/internal/controller"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/internal/service"
	v1 "eazy-quizy-auth/pkg/api/v1"
	"eazy-quizy-auth/pkg/interceptors"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/postgresql/database"
	"eazy-quizy-auth/pkg/redis"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func stringPtr(s string) *string {
	return &s
}

type AuthIntegrationTestSuite struct {
	suite.Suite
	client   v1.AuthServiceClient
	conn     *grpc.ClientConn
	ctx      context.Context
	cancel   context.CancelFunc
	server   *grpc.Server
	listener net.Listener
	db       *database.DB
	logger   *logger.Logger
}

func (s *AuthIntegrationTestSuite) SetupSuite() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel

	var err error
	s.logger, err = logger.Setup()
	if err != nil {
		s.T().Fatal("Failed to initialize logger:", err)
	}

	s.db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to initialize in-memory database:", err)
	}

	if err != nil {
		s.T().Fatal("Failed to run migrations:", err)
	}

	redisClient := redis.NewClient(&config.Config{
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	})

	userRepo := repository.NewUserRepository(s.db.DB)
	friendRepo := repository.NewFriendRepository(s.db.DB)

	jwtService := service.NewJWTService(
		*userRepo,
		redisClient,
		&config.JWTConfig{
			Secret: "secret",
			TTL:    time.Hour,
		},
		s.logger,
	)

	authService := service.NewAuthService(*userRepo, jwtService, s.logger)
	friendService := service.NewFriendService(*friendRepo, *userRepo, s.logger)

	authController := controller.NewAuthController(authService, friendService, s.logger)

	authInterceptor := interceptors.NewAuthInterceptor(authService, s.logger)

	s.listener, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		s.T().Fatal("Failed to create listener:", err)
	}

	s.server = grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary(ctx)))
	v1.RegisterAuthServiceServer(s.server, authController)

	go func() {
		if err := s.server.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
			s.T().Fatal("Failed to start gRPC server:", err)
		}
	}()

	addr := s.listener.Addr().String()
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Second)

	s.conn, err = grpc.DialContext(s.ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		s.T().Fatal("Failed to connect to gRPC server:", err)
	}
	s.client = v1.NewAuthServiceClient(s.conn)
}

func (s *AuthIntegrationTestSuite) TearDownSuite() {
	if s.conn != nil {
		s.conn.Close()
	}
	if s.server != nil {
		s.server.Stop()
	}
	if s.listener != nil {
		s.listener.Close()
	}
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
	if s.cancel != nil {
		s.cancel()
	}
}

func TestAuthServiceIntegration(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}

func (s *AuthIntegrationTestSuite) TestFullAuthFlow() {
	t := s.T()
	testUsername := "testuser_" + time.Now().Format("20060102150405")
	testEmail := testUsername + "@example.com"
	testPassword := "securePassword123!"

	registerResp, err := s.client.Register(s.ctx, &v1.RegisterRequest{
		Email:    testEmail,
		Password: testPassword,
		Username: &testUsername,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, registerResp.GetUserId())

	_, err = s.client.Register(s.ctx, &v1.RegisterRequest{
		Username: &testUsername,
		Email:    testEmail,
		Password: testPassword,
	})
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))

	loginResp, err := s.client.Login(s.ctx, &v1.LoginRequest{
		Email:    testEmail,
		Password: testPassword,
	})
	require.NoError(t, err)
	token := loginResp.GetToken()
	t.Logf("Login token: %s", token)
	assert.NotEmpty(t, token)

	ctx := metadata.AppendToOutgoingContext(s.ctx, "Authorization", "Bearer "+token)
	meResp, err := s.client.GetMe(ctx, &v1.GetMeRequest{})
	require.NoError(t, err, "failed to get user data")
	assert.Equal(t, testUsername, meResp.GetUsername())
	assert.Equal(t, testEmail, meResp.GetEmail())

	_, err = s.client.Login(s.ctx, &v1.LoginRequest{
		Email:    testEmail,
		Password: "wrong_password",
	})
	assert.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	_, err = s.client.Login(s.ctx, &v1.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "anypassword",
	})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	ctx = metadata.AppendToOutgoingContext(s.ctx, "Authorization", "Bearer invalid.token.here")
	_, err = s.client.GetMe(ctx, &v1.GetMeRequest{})
	assert.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func (s *AuthIntegrationTestSuite) TestValidation() {
	t := s.T()
	timestamp := time.Now().Format("20060102150405.000")
	testCases := []struct {
		name        string
		req         *v1.RegisterRequest
		expectedErr codes.Code
	}{
		{
			name: "Successful registration with empty username",
			req: &v1.RegisterRequest{
				Email:    "valid1_" + timestamp + "@example.com",
				Password: "password123",
				Username: stringPtr(""),
			},
			expectedErr: codes.OK,
		},
		{
			name: "Successful registration without username",
			req: &v1.RegisterRequest{
				Email:    "valid2_" + timestamp + "@example.com",
				Password: "password123",
			},
			expectedErr: codes.OK,
		},
		{
			name: "Invalid email",
			req: &v1.RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedErr: codes.InvalidArgument,
		},
		{
			name: "Short password",
			req: &v1.RegisterRequest{
				Email:    "valid3_" + timestamp + "@example.com",
				Password: "short",
			},
			expectedErr: codes.InvalidArgument,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := s.client.Register(s.ctx, tc.req)
			if tc.expectedErr == codes.OK {
				assert.NoError(t, err, "expected successful registration")
				assert.NotEmpty(t, resp.GetUserId(), "expected valid user ID")
			} else {
				assert.Error(t, err, "expected error")
				assert.Equal(t, tc.expectedErr, status.Code(err), "expected error code")
			}
		})
	}
}
