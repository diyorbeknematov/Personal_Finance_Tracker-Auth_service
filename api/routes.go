package api

import (
	"auth-service/api/handler"
	"auth-service/api/middleware"
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
		c.port = fmt.Sprintf("auth_app:%d", cfg.HTTP_PORT)
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

	auth1 := router.Group("/auth")
	{
		auth1.POST("/forgot-password", h.AuthHandler().ForgotPassword)
		auth1.POST("/reset-password", h.AuthHandler().ResetPassword)
		auth1.POST("/register", h.AuthHandler().RegisterUser)
		auth1.POST("/login", h.AuthHandler().LoginUser)
	}

	auth := router.Group("/auth", middleware.IsAuthenticated(authService), middleware.LogMiddleware(logger))
	{
		auth.POST("/logout", h.AuthHandler().LogOutUser)
		auth.POST("/roles", h.AuthHandler().ManageUserRoles)
		auth.POST("/refresh-token", h.AuthHandler().RefreshToken)
	}
}
