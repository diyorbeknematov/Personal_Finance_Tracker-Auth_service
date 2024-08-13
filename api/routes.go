package api

import (
	"auth-service/api/handler"
	"auth-service/config"
	"auth-service/service"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth-service/api/docs"
)

type Controller interface {
	SetupRoutes(authService service.AuthService, logger *slog.Logger)
	StartServer(cfg *config.Config) error
}

type controllerImpl struct {
	port   string
	router *gin.Engine
}

func NewController() Controller {
	router := gin.Default()
	return &controllerImpl{router: router}
}

func (c *controllerImpl) StartServer(cfg *config.Config) error {
	if c.port == "" {
		c.port = fmt.Sprintf(":%d", cfg.HTTP_PORT)
	}

	return c.router.Run(c.port)
}

// @title Auth Service
// @version 1.0
// @description Auth service
// @BasePath /api/v1
// @schemes http
// @in header
// @name Authorization
func (c *controllerImpl) SetupRoutes(authService service.AuthService, logger *slog.Logger) {
	h := handler.NewMainHandler(authService, logger)

	c.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router := c.router.Group("/api/v1")

	auth := router.Group("/auth")
	{
		auth.POST("/register", h.AuthHandler().RegisterUser)
		auth.POST("/login", h.AuthHandler().LoginUser)
		auth.POST("/logout", h.AuthHandler().LogOutUser)
		auth.POST("/roles", h.AuthHandler().ManageUserRoles)
		auth.POST("/forgot-password", h.AuthHandler().ForgotPassword)
		auth.POST("/reset-password", h.AuthHandler().ResetPassword)
		auth.POST("/refresh-token", h.AuthHandler().RefreshToken)
	}
}
