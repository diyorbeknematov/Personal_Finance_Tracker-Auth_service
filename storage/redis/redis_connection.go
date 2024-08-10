package redis

import (
	"auth-service/config"

	"github.com/redis/go-redis/v9"
)

func RedisConnect(cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis_HOST,
		Password: cfg.Redis_PASSWORD,
		DB:       cfg.Redis_DB,
	})

	return client, nil
}
