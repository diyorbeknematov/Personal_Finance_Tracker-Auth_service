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
	defer listener.Close()

	log.Printf("Listening on %s\n", listener.Addr())

	s := grpc.NewServer()
	user.RegisterAuthServiceServer(s, service.NewUserService(storage, logger))
	err = s.Serve(listener)
	if err != nil {
		logger.Error("Serve error", "error", err)
		log.Fatal(err)
	}
}
