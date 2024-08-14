package middleware

import (
	"auth-service/api/token"
	"auth-service/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAuthenticated(service service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Cookie'dan tokenni olish
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Authorization cookie not found",
			})
			ctx.Abort()
			return
		}

		ok, err := service.IsTokenBlacklisted(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Unauthorized",
			})
			ctx.Abort()
			return
		}
		if ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Token is blacklisted",
			})
			ctx.Abort()
			return
		}

		// JWT tokenni tekshirish va tasdiqlash
		claims, err := token.ExtractAndValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"Error": "Invalid token",
			})
			ctx.Abort()
			return
		}

		// Foydalanuvchi ma'lumotlarini context ga qo'shish
		ctx.Set("claims", claims)

		ctx.Next()
	}
}

func LogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Request received",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		c.Next()

		logger.Info("Response sent",
			slog.Int("status", c.Writer.Status()),
		)
	}
}
