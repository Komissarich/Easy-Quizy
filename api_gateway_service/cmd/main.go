package main

import (
	"api_gateway/gen/auth_service"
	"api_gateway/gen/quiz_service"
	"api_gateway/gen/stat_service"
	logger "api_gateway/pkg"
	"bytes"
	"context" // для QuizService
	"fmt"
	"io"
	"log"

	// для StatsService

	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Разрешаем запросы с любых источников (для production укажите конкретные домены)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Разрешаем необходимые методы
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Разрешаем необходимые заголовки
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Для предварительных OPTIONS-запросов
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func HealthHeandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Service is healthy!")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method != "GET" {

			// Читаем тело ОДИН раз
			body, err := io.ReadAll(r.Body)

			if err != nil {
				http.Error(w, "Error reading body", http.StatusBadRequest)
				return
			}

			// Восстанавливаем тело для следующих обработчиков
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Логируем
			log.Printf("Raw request body: %s", string(body))

			// Проверяем Content-Type
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
				return
			}
		}
		fmt.Println("nice redirect")
		next.ServeHTTP(w, r)
	})
}
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l, err := logger.Setup()

	if err != nil {
		l.Fatal("failed to setup logger", zap.Error(err))
	}

	// 1. Подключаемся к gRPC-сервисам
	quizConn, err := grpc.NewClient(
		"quiz_service:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Fatal("failed to connect to Quiz gRPC service", zap.Error(err))
	}
	defer quizConn.Close()

	statsConn, err := grpc.NewClient(
		"stat_service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Fatal("failed to connect to Stats gRPC service", zap.Error(err))
	}
	defer statsConn.Close()

	authConn, err := grpc.NewClient(
		"auth_service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		l.Fatal("failed to connect to Auth gRPC service", zap.Error(err))
	}
	defer authConn.Close()

	// 2. Создаём Gateway-маршрутизатор

	rootMux := http.NewServeMux()

	rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	grpcGatewayMux := runtime.NewServeMux(
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Gateway error: %v", err)
			runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, r, err)
		}),
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			// Переносим Authorization header в gRPC метаданные
			token := req.Header.Get("Authorization")
			if token != "" {
				return metadata.New(map[string]string{
					"authorization": token,
				})
			}
			return nil
		}),
	)

	// 3. Регистрируем обработчики для обоих сервисов

	if err := quiz_service.RegisterQuizServiceHandler(ctx, grpcGatewayMux, quizConn); err != nil {
		l.Fatal("failed to register Quiz gateway", zap.Error(err))
	}
	if err := stat_service.RegisterStatisticsHandler(ctx, grpcGatewayMux, statsConn); err != nil {
		l.Fatal("failed to register Stats gateway", zap.Error(err))
	}

	if err := auth_service.RegisterAuthServiceHandler(ctx, grpcGatewayMux, authConn); err != nil {
		l.Fatal("failed to register Auth gateway", zap.Error(err))
	}
	fmt.Println(authConn.GetState().String())
	rootMux.Handle("/", grpcGatewayMux)
	// corsHandler := allowCORS(rootMux)

	// 2. Тестовый endpoint без gRPC
	// rootMux.HandleFunc("/v1/users/register", func(w http.ResponseWriter, r *http.Request) {
	// 	body, _ := io.ReadAll(r.Body)
	// 	repaired, _ := repairJSON(body)
	// 	log.Printf("Raw body: %s", string(body))
	// 	log.Printf("Repaired body: %s", string(repaired))
	// 	client := auth_service.NewAuthServiceClient(authConn)
	// 	resp, err := client.Register(r.Context(), &auth_service.RegisterRequest{
	// 		Email:    "egorkart1@gmail.com", // Заполните из body
	// 		Password: "aaaabbbbb",
	// 	})

	// 	log.Println(resp, err)
	// })
	corsHandler := allowCORS(loggingMiddleware(rootMux))

	srv := &http.Server{
		Addr:    ":8085",
		Handler: corsHandler,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		if err := srv.Shutdown(ctx); err != nil {
			l.Info("HTTP server shutdown error", zap.Error(err))
		}
	}()

	// 5. Запускаем REST-сервер

	l.Info("Starting HTTP gateway on :8085 (serving both Quiz and Stats services)")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		l.Fatal("HTTP server failed", zap.Error(err))
	}
}
