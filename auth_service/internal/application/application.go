package application

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/internal/controller"
	"eazy-quizy-auth/internal/repository"
	"eazy-quizy-auth/internal/service"
	auth "eazy-quizy-auth/pkg/api/v1"
	"eazy-quizy-auth/pkg/interceptors"
	"eazy-quizy-auth/pkg/logger"
	"eazy-quizy-auth/pkg/postgresql/database"
	"eazy-quizy-auth/pkg/redis"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

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

	redisClient := redis.NewClient(cfg)
	if err := redisClient.Ping(context.Background()); err != nil {
		l.Fatal("Failed to connect to Redis: %v", zap.Error(err))
	}
	userRepo := repository.NewUserRepository(db.DB)
	friendRepo := repository.NewFriendRepository(db.DB)
	quizzesRepo := repository.NewQuizzesRepo(db.DB)

	jwtService := service.NewJWTService(*userRepo, redisClient, &cfg.JWT, l)

	authService := service.NewAuthService(*userRepo, jwtService, l)
	friendService := service.NewFriendService(*friendRepo, *userRepo, l)
	quizzesService := service.NewQuizzesService(*quizzesRepo, l)

	authController := controller.NewAuthController(authService, friendService, quizzesService, *jwtService, l)

	authInterceptor := interceptors.NewAuthInterceptor(authService, l)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		l.Fatal("failed to listen:", zap.Error(err))
	}

	l.Info("listening on port " + cfg.GRPCPort)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary(ctx)),
	)

	auth.RegisterAuthServiceServer(grpcServer, authController)

	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			l.Fatal("failed to serve:", zap.Error(err))
		}
	}()

	l.Info("server started")

	select {
	case <-ctx.Done():
		l.Info("Server graceful stopped")
		grpcServer.GracefulStop()
	}
}
