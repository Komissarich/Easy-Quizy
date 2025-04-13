package application

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/internal/controller"
	"eazy-quizy-auth/internal/postgresql/database"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/internal/service"
	auth "eazy-quizy-auth/pkg/api/v1"
	"eazy-quizy-auth/pkg/interceptors"
	"eazy-quizy-auth/pkg/logger"
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Application struct {
	config string
}

func New(config string) *Application {
	return &Application{
		config: config,
	}
}

func (a *Application) Run(ctx context.Context) {
	l, err := logger.Setup()
	if err != nil {
		l.Fatal("failed to setup logger", zap.Error(err))
	}

	cfg, err := config.New()
	if err != nil {
		l.Fatal("failed to load config file", zap.Error(err))
	}

	l.Info("Starting service")

	db, err := database.New(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close(ctx)

	l.Info("database connection established")

	userRepo := repository.NewUserRepository(db.DB)
	friendRepo := repository.NewFriendRepository(db.DB)

	authService := service.NewAuthService(*userRepo, cfg.JWT.Secret)
	friendService := service.NewFriendService(*friendRepo, *userRepo)

	authController := controller.NewAuthController(authService, friendService)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.RoleInterceptor(ctx, []string{"admin", "user"})))
	auth.RegisterAuthServiceServer(grpcServer, authController)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
