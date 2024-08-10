package redis

import (
	"auth-service/models"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore interface {
	AddTokenBlacklist(token string, expirationTime time.Duration) (*models.Response, error)
	IsTokenBlacklisted(token string) (bool, error)
	StoreCode(email, code string, expirationTime time.Duration) (*models.Response, error)
	IsCodeValid(email, code string) (bool, error)
}

type redisStoreImpl struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisStore(client *redis.Client) RedisStore {
	return &redisStoreImpl{client: client}
}

func (rdb *redisStoreImpl) AddTokenBlacklist(token string, expirationTime time.Duration) (*models.Response, error) {
	err := rdb.client.Set(ctx, token, "blacklisted", expirationTime).Err()
	if err != nil {
		return &models.Response{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	return &models.Response{
		Status:  "success",
		Message: "Token added to blacklist successfully",
	}, nil
}

func (rdb *redisStoreImpl) IsTokenBlacklisted(token string) (bool, error) {
	val, err := rdb.client.Get(ctx, token).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == "blacklisted", nil
}

func (rdb *redisStoreImpl) StoreCode(email, code string, expirationTime time.Duration) (*models.Response, error) {
	err := rdb.client.Set(ctx, email+":code", code, expirationTime).Err()
	if err != nil {
		return &models.Response{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}

	return &models.Response{
		Status:  "success",
		Message: "Code stored successfully",
	}, nil
}

func (rdb *redisStoreImpl) IsCodeValid(email, code string) (bool, error) {
	val, err := rdb.client.Get(ctx, email+":code").Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return val == code, nil
}
