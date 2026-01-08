package redis

import (
	"bizbundl/internal/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var TenantKey = "tenant_id"

func NewRedisClient(cfg *config.Config) *redis.Client {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}

	log.Info().Str("addr", addr).Msg("Connected to Redis")
	return client
}

// Key prefixes the key with the tenant ID from the context
func Key(ctx context.Context, key string) string {
	if ctx == nil {
		return key
	}
	tenantID, ok := ctx.Value(TenantKey).(string)
	if !ok || tenantID == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", tenantID, key)
}
