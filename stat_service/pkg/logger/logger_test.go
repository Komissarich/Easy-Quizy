package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestNew(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		ctx := context.Background()
		newCtx, err := New(ctx)

		require.NoError(t, err)
		assert.NotNil(t, newCtx)

		logger := GetLoggerFromCtx(newCtx)
		assert.NotNil(t, logger)
	})

	t.Run("logger is set in context", func(t *testing.T) {
		ctx := context.Background()
		newCtx, _ := New(ctx)

		logger := newCtx.Value(LoggerKey)
		assert.NotNil(t, logger)
		_, ok := logger.(*Logger)
		assert.True(t, ok)
	})
}

func TestGetLoggerFromCtx(t *testing.T) {
	t.Run("logger exists in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), LoggerKey, &Logger{logger: zap.NewNop()})
		logger := GetLoggerFromCtx(ctx)
		assert.NotNil(t, logger)
	})

	t.Run("logger missing in context", func(t *testing.T) {
		ctx := context.Background()
		assert.Panics(t, func() {
			GetLoggerFromCtx(ctx)
		})
	})
}

func TestLoggerMethods(t *testing.T) {
	ctx := context.Background()
	mockLogger := zaptest.NewLogger(t)
	testLogger := &Logger{logger: mockLogger}

	t.Run("Info with request ID", func(t *testing.T) {
		ctxWithID := context.WithValue(ctx, RequestID, "test-123")
		testLogger.Info(ctxWithID, "test message", zap.String("key", "value"))
	})

	t.Run("Info without request ID", func(t *testing.T) {
		testLogger.Info(ctx, "test message", zap.String("key", "value"))
	})

	t.Run("Error with request ID", func(t *testing.T) {
		ctxWithID := context.WithValue(ctx, RequestID, "test-123")
		testLogger.Error(ctxWithID, "test error", zap.String("key", "value"))
	})

	t.Run("Error without request ID", func(t *testing.T) {
		testLogger.Error(ctx, "test error", zap.String("key", "value"))
	})

	t.Run("Fatal with request ID", func(t *testing.T) {
		// Since Fatal exits, we can't test it directly in the same process
		// In real code you'd use a mock or test this separately
		t.Skip("Fatal cannot be tested directly")
	})
}

func TestNew_Error(t *testing.T) {
	// Чтобы протестировать ошибку, можно использовать кастомный config
	// Но в текущей реализации zap.NewProduction() не возвращает ошибок в нормальных условиях
	// Этот тест просто показывает как можно протестировать ошибку
	t.Run("returns error when logger creation fails", func(t *testing.T) {
		// В реальном коде нужно мокировать zap.NewProduction()
		// Здесь просто проверяем что функция возвращает ошибку
		originalNewProduction := zapNewProduction
		zapNewProduction = func(options ...zap.Option) (*zap.Logger, error) {
			return nil, assert.AnError
		}
		defer func() { zapNewProduction = originalNewProduction }()

		_, err := New(context.Background())
		assert.Error(t, err)
	})
}

// Для тестирования ошибки
var zapNewProduction = zap.NewProduction
