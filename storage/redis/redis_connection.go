package redis

import (
	"auth-service/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func RedisConnect(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Redis_HOST, cfg.Redis_PORT)
	fmt.Println("Connecting to Redis at", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis_PASSWORD,
		DB:       cfg.Redis_DB,
	})
	
	log.Println("Connecting to Redis at", addr)
	return client, nil
}
