package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Connect(redisURL string) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("failed to parse redis URL: %v", err)
	}
	Client = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("redis connected")
}

func Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return Client.Set(ctx, key, value, ttl).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

func Del(ctx context.Context, keys ...string) error {
	return Client.Del(ctx, keys...).Err()
}

func DelPattern(ctx context.Context, pattern string) error {
	keys, err := Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return Client.Del(ctx, keys...).Err()
}
