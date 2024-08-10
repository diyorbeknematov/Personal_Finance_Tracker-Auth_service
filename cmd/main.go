package main

import (
	"auth-service/api"
	"auth-service/cmd/server"
	"auth-service/config"
	"auth-service/pkg/logs"
	"auth-service/service"
	"auth-service/storage"
	"auth-service/storage/postgres"
	"auth-service/storage/redis"
	"log"
)

func main() {
	log.Println("Server started")
	logger := logs.InitLogger()
	logger.Info("Application started")
	cfg := config.Load()

	db, err := postgres.ConnectDB(cfg)
	if err != nil {
		logger.Error("Database connection error", "error", err)
		log.Fatal(err)
	}

	rdb, err := redis.RedisConnect(cfg)
	if err != nil {
		logger.Error("Redis connection error", "error", err)
		log.Fatal(err)
	}

	storage := storage.NewUserStorage(db, rdb)

	authService := service.NewAuthService(storage, logger)

	go func() {
		log.Println("Stargin GRPC server")
		logger.Info("Starting GRPC server")
		server.StartServer(logger, storage)
	}()

	controller := api.NewController()
	controller.SetupRoutes(authService, logger)
	err = controller.StartServer(cfg)
	if err != nil {
		logger.Error("Start server error", "error", err)
		log.Fatal(err)
	}
}
