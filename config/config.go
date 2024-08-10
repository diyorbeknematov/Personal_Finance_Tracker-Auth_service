package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTP_PORT          int    `yaml:"http_port"`
	GRPC_PORT          int    `yaml:"grpc_port"`
	DB_HOST            string `yaml:"db_host"`
	DB_PORT            int    `yaml:"db_port"`
	DB_USER            string `yaml:"db_user"`
	DB_PASSWORD        string `yaml:"db_password"`
	DB_NAME            string `yaml:"db_name"`
	Redis_HOST         string `yaml:"redis_host"`
	Redis_PORT         int    `yaml:"redis_port"`
	Redis_PASSWORD     string `yaml:"redis_password"`
	Redis_DB           int    `yaml:"redis_db"`
	Jwt_SECRET_ACCESS  string `yaml:"jwt_secret"`
	Jwt_SECRET_REFRESH string `yaml:"jwt_secret_refresh"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	config := &Config{}

	config.HTTP_PORT = cast.ToInt(coalesce("HTTP_PORT", 8080))
	config.GRPC_PORT = cast.ToInt(coalesce("GRPC_PORT", 50051))

	config.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	config.DB_PORT = cast.ToInt(coalesce("DB_PORT", 5432))
	config.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	config.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "your_password"))
	config.DB_NAME = cast.ToString(coalesce("DB_NAME", "your_db"))

	config.Redis_HOST = cast.ToString(coalesce("REDIS_HOST", "localhost"))
	config.Redis_PORT = cast.ToInt(coalesce("REDIS_PORT", 6379))
	config.Redis_PASSWORD = cast.ToString(coalesce("REDIS_PASSWORD", ""))
	config.Redis_DB = cast.ToInt(coalesce("REDIS_DB", 0))

	config.Jwt_SECRET_ACCESS = cast.ToString(coalesce("JWT_SECRET_ACCESS", "your_secret_access"))
	config.Jwt_SECRET_REFRESH = cast.ToString(coalesce("JWT_SECRET_REFRESH", "your_secret_refresh"))

	return config
}

func coalesce(key string, defaults interface{}) interface{} {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaults
}
