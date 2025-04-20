package redis_test

import (
	"context"
	"eazy-quizy-auth/internal/config"
	"eazy-quizy-auth/pkg/redis"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisClient(t *testing.T) {
	cfg := &config.Config{
		Redis: config.RedisConfig{
			Host:     "localhost",
			Port:     "6379",
			Password: "",
			DB:       0,
		},
	}

	client := redis.NewClient(cfg)
	ctx := context.Background()

	t.Run("Ping", func(t *testing.T) {
		err := client.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Set and Get", func(t *testing.T) {
		err := client.Set(ctx, "test_key", "test_value", time.Minute)
		assert.NoError(t, err)

		val, err := client.Get(ctx, "test_key")
		assert.NoError(t, err)
		assert.Equal(t, "test_value", val)
	})

	t.Run("Delete", func(t *testing.T) {
		err := client.Set(ctx, "to_delete", "value", time.Minute)
		assert.NoError(t, err)

		err = client.Delete(ctx, "to_delete")
		assert.NoError(t, err)

		_, err = client.Get(ctx, "to_delete")
		assert.Error(t, err)
	})

	t.Run("Set and Get with expiration", func(t *testing.T) {
		err := client.Set(ctx, "test_key", "test_value", time.Second)
		assert.NoError(t, err)

		time.Sleep(2 * time.Second)

		_, err = client.Get(ctx, "test_key")
		assert.Error(t, err)
	})
}
