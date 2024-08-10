package handler

import (
	"auth-service/service"
	"log/slog"
)

type MainHandler interface {
	AuthHandler() UserHandler
}

type mainHandlerImpl struct {
	authService service.AuthService
	logger      *slog.Logger
}

func NewMainHandler(authService service.AuthService, logger *slog.Logger) MainHandler {
	return &mainHandlerImpl{authService: authService, logger: logger}
}

func (h *mainHandlerImpl) AuthHandler() UserHandler {
	return NewUserHandler(h.authService, h.logger)
}
