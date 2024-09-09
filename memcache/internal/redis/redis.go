package redis_wrapper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

// RedisCache structure to handle Redis operations
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes a new Redis cache instance
func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	// Check if Redis is up and running
	for {
		if _, err := client.Ping(context.Background()).Result(); err != nil {
			log.Error().Err(err).Msg("Failed to connect to Redis")
			time.Sleep(time.Millisecond * 300)
			continue
		}
		break
	}

	return &RedisCache{client: client}
}

// Save saves a short URL and its corresponding long URL in Redis with an expiration time
func (cache *RedisCache) Save(shortURL, longURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := cache.client.Set(ctx, shortURL, longURL, 24*time.Hour).Err()
	if err != nil {
		log.Error().Err(err).Str("short_url", shortURL).Str("long_url", longURL).Msg("Failed to save URL in Redis")
		return err
	}

	log.Info().Str("short_url", shortURL).Str("long_url", longURL).Msg("Saved URL in Redis")
	return nil
}

// Get retrieves the long URL corresponding to the short URL from Redis
func (cache *RedisCache) Get(shortURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	longURL, err := cache.client.Get(ctx, shortURL).Result()
	if err == redis.Nil {
		log.Warn().Str("short_url", shortURL).Msg("Short URL not found in Redis")
		return "", nil // Key does not exist
	} else if err != nil {
		log.Error().Err(err).Str("short_url", shortURL).Msg("Failed to get URL from Redis")
		return "", err
	}

	log.Info().Str("short_url", shortURL).Str("long_url", longURL).Msg("Retrieved long URL from Redis")
	return longURL, nil
}
