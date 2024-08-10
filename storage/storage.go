package storage

import (
	"auth-service/storage/postgres"
	rdb "auth-service/storage/redis"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type IStorage interface {
	AuthRepository() postgres.AuthenticationRepository
	UserRepository() postgres.UserRepository
	RedisStore() rdb.RedisStore
}

type storageImpl struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewUserStorage(db *sql.DB, rdb *redis.Client) IStorage {
	return &storageImpl{db: db, rdb: rdb}
}

func (s *storageImpl) AuthRepository() postgres.AuthenticationRepository {
	return postgres.NewAuthenticationRepository(s.db)
}

func (s *storageImpl) UserRepository() postgres.UserRepository {
	return postgres.NewUserRepository(s.db)
}

func (s *storageImpl) RedisStore() rdb.RedisStore {
	return rdb.NewRedisStore(s.rdb)
}
