package server

import (
	"auth-service/config"
	"auth-service/generated/user"
	"auth-service/service"
	"auth-service/storage"
	"fmt"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

func StartServer(logger *slog.Logger, storage storage.IStorage) {
	cfg := config.Load()

	log.Println("Server started")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC_PORT))
	if err != nil {
		logger.Error("Listen error", "error", err)
		log.Fatal(err)
	}

	log.Printf("Listening on %s\n", listener.Addr())

	s := grpc.NewServer()
	userService := service.NewUserService(storage, logger)
	user.RegisterAuthServiceServer(s, userService)

	
	if err := s.Serve(listener); err != nil {
		logger.Error("Serve error", "error", err)
		log.Fatal(err)
	}
}
