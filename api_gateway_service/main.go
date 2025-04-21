package main

import (
	"api_gateway/gen/auth_service"
	"api_gateway/gen/quiz_service"
	"api_gateway/gen/stat_service"
	"context" // для QuizService

	// для StatsService
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Подключаемся к gRPC-сервисам
	quizConn, err := grpc.NewClient(
		"quiz_service:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Quiz gRPC service: %v", err)
	}
	defer quizConn.Close()

	statsConn, err := grpc.NewClient(
		"stats_service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Stats gRPC service: %v", err)
	}
	defer statsConn.Close()

	authConn, err := grpc.NewClient(
		"stats_service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to Auth gRPC service: %v", err)
	}
	defer authConn.Close()

	// 2. Создаём Gateway-маршрутизатор
	mux := runtime.NewServeMux()

	// 3. Регистрируем обработчики для обоих сервисов

	if err := quiz_service.RegisterQuizServiceHandler(ctx, mux, quizConn); err != nil {
		log.Fatalf("failed to register Quiz gateway: %v", err)
	}
	if err := stat_service.RegisterStatisticsHandler(ctx, mux, statsConn); err != nil {
		log.Fatalf("failed to register Stats gateway: %v", err)
	}

	if err := auth_service.RegisterAuthServiceHandler(ctx, mux, authConn); err != nil {
		log.Fatalf("failed to register Stats gateway: %v", err)
	}

	// 4. Настраиваем HTTP-сервер
	srv := &http.Server{
		Addr:    ":8085",
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// 5. Запускаем REST-сервер
	log.Println("Starting HTTP gateway on :8085 (serving both Quiz and Stats services)")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
