package app

import (
	"awesomeProject2/internal/config"
	"awesomeProject2/internal/quizzes/quizzes/service"
	"awesomeProject2/internal/quizzes/repository"
	v1 "awesomeProject2/pkg/api/v1"
	"awesomeProject2/pkg/logger"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Application struct {
	config string
}

func New(config string) *Application {
	return &Application{
		config: config,
	}
}

func (app *Application) Run(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	ctx, err := logger.NewLog(ctx)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	l := logger.GetLoggerFromCtx(ctx)
	cfg, err := config.New()
	if err != nil {
		l.Fatal(ctx, err.Error())
		return
	}
	l.Info(ctx, "Start quiz service")
	repo := repository.NewRepository(ctx, cfg)
	defer repo.CloseConn()
	service := service.New(ctx, repo)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.GRPCPort))
	if err != nil {
		l.Fatal(ctx, err.Error())
	}
	l.Info(ctx, "Listen on port "+strconv.Itoa(cfg.GRPCPort))
	grpcServer := grpc.NewServer()
	v1.RegisterQuizServiceServer(grpcServer, service)
	reflection.Register(grpcServer)
	go func() {
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			l.Info(ctx, err.Error())
		}
	}()
	l.Info(ctx, "Start quiz server")
	select {
	case <-ctx.Done():
		l.Info(ctx, "Quiz service gracefully stopped")
		grpcServer.GracefulStop()
	}
}
