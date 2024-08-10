package handler

import (
	"auth-service/api/token"
	"auth-service/models"
	"auth-service/pkg/helper"
	"auth-service/service"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ManageUserRoles(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	LogOutUser(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type userHandlerImpl struct {
	authService service.AuthService
	logger      *slog.Logger
}

func NewUserHandler(authService service.AuthService, logger *slog.Logger) UserHandler {
	return &userHandlerImpl{authService: authService, logger: logger}
}

// @Summary Register user
// @Description Registers a new user
// @Accept json
// @Produce json
// @Param user body models.RegisterUser true "User details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /auth/register [POST]
func (h *userHandlerImpl) RegisterUser(ctx *gin.Context) {
	var userReq models.RegisterUser

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		h.logger.Error("BindJSON error", "error", err)
		ctx.JSON(400, models.Error{Message: "Invalid request body"})
		return
	}

	exists, err := h.authService.EmailExists(userReq.Email)
	if err != nil {
		h.logger.Error("EmailExists error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error checking email"})
		return
	}
	if exists {
		ctx.JSON(400, models.Error{Message: "Email already exists"})
		return
	}

	resp, err := h.authService.RegisterUser(userReq)
	if err != nil {
		h.logger.Error("RegisterUser error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error registering user"})
		return
	}

	ctx.JSON(200, resp)
}

// @Summary Login user
// @Description Login a user
// @Accept json
// @Produce json
// @Param user body models.LoginUserReq true "User credentials"
// @Success 200 {object} models.LoginUserResp
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Failure 404 {object} models.Error
// @router /auth/login [post]
func (h *userHandlerImpl) LoginUser(ctx *gin.Context) {
	var userReq models.LoginUserReq

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		h.logger.Error("BindJSON error", "error", err)
		ctx.JSON(400, models.Error{Message: "Invalid request body"})
		return
	}

	exists, err := h.authService.EmailExists(userReq.Email)
	if err != nil {
		h.logger.Error("EmailExists error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error checking email"})
		return
	}
	if !exists {
		ctx.JSON(404, models.Error{Message: "Email not found"})
		return
	}

	user, err := h.authService.LoginUser(userReq)
	if err != nil {
		h.logger.Error("LoginUser error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error logging in"})
		return
	}

	accessToken, err := token.GeneratedJWTTokenAccess(models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
	if err != nil {
		h.logger.Error("GeneratedJwtTokenAccess error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error generating access token"})
		return
	}

	refreshToken, err := token.GeneratedJwtTokenRefresh(models.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	})
	if err != nil {
		h.logger.Error("GeneratedJwtTokenRefresh error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error generating refresh token"})
		return
	}

	_, err = h.authService.SaveRefreshToken(models.RefreshToken{
		Email: user.Email,
		RefreshToken: refreshToken,
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	
	if err!= nil {
        h.logger.Error("SaveRefreshToken error", "error", err)
        ctx.JSON(500, models.Error{Message: "Error saving refresh token"})
        return
    }

	ctx.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	ctx.JSON(200, models.LoginUserResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// @Summary Logout user
// @Description Logout a user
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @Failure 404 {object} models.Error
// @router /auth/logout [post]
func (h *userHandlerImpl) LogOutUser(ctx *gin.Context) {
	val, ok := ctx.Get("claims")
	if !ok {
		h.logger.Error("Token not found in context")
		ctx.JSON(401, models.Error{Message: "Unauthorized"})
		return
	}
	claims, ok := val.(*token.Claims)
	if !ok {
		h.logger.Error("Token claims not found in context")
		ctx.JSON(401, models.Error{Message: "Unauthorized"})
		return
	}

	_, err := h.authService.DeleteUser(claims.ID)
	if err != nil {
		h.logger.Error("DeleteUser error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error logging out"})
		return
	}

	tokenstr := ctx.Request.Header.Get("Authorization")
	_, err = h.authService.AddTokenBlacklist(tokenstr, time.Duration(claims.ExpiresAt))
	if err != nil {
		h.logger.Error("AddTokenBlacklist error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error blacklisting token"})
		return
	}

	_, err = h.authService.InvalidateRefreshToken(claims.Email)
	if err != nil {
		h.logger.Error("InvalidateRefreshToken error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error invalidating refresh token"})
		return
	}

	ctx.SetCookie("access_token", "", -1, "/", "", false, true)
	ctx.JSON(200, models.Response{
		Status:  "success",
		Message: "User logged out",
	})
}

// @summary Update user role
// @Description Update user role
// @Accept json
// @Produce json
// @Param user body models.ManageUserRoles true "User details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Failure 404 {object} models.Error
// @Router /auth/roles [post]
func (h *userHandlerImpl) ManageUserRoles(ctx *gin.Context) {
	var userReq models.ManageUserRoles

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		h.logger.Error("BindJSON error", "error", err)
		ctx.JSON(400, models.Error{Message: "Invalid request body"})
		return
	}

	resp, err := h.authService.UpdateUserRoles(userReq)
	if err != nil {
		h.logger.Error("UpdateUserRoles error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error updating user roles"})
		return
	}

	ctx.JSON(200, resp)
}

// @Summary Forgot password
// @Description Forgot user password
// @accept json
// @Produce json
// @param user body models.ForgotPassword true "User details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @router /auth/forgot-password [post]
func (h *userHandlerImpl) ForgotPassword(ctx *gin.Context) {
	var userReq models.ForgotPassword

	if err := ctx.ShouldBindJSON(&userReq); err != nil {
		h.logger.Error("BindJSON error", "error", err)
		ctx.JSON(400, models.Error{Message: "Invalid request body"})
		return
	}

	code := strconv.Itoa(rand.Intn(999999) + 100000)

	_, err := h.authService.StoreCode(userReq.Email, code, time.Duration(time.Minute*5))
	if err != nil {
		h.logger.Error("StoreCode error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error storing code"})
		return
	}

	err = helper.SendPasswordResetEmail(userReq.Email, code)
	if err != nil {
		h.logger.Error("SendResetPasswordEmail error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error sending email"})
		return
	}

	ctx.JSON(200, models.Response{
		Status:  "success",
		Message: "Password reset email sent",
	})
}

// @summary Reset password
// @Description Reset user password
// @accept json
// @produce json
// @param resetPassword body models.ResetPassword true "Reset password details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @router /auth/reset-password [post]
func (h *userHandlerImpl) ResetPassword(ctx *gin.Context) {
	var resetPasswordReq models.ResetPassword

	if err := ctx.ShouldBindJSON(&resetPasswordReq); err != nil {
		h.logger.Error("BindJSON error", "error", err)
		ctx.JSON(400, models.Error{Message: "Invalid request body"})
		return
	}

	isvalid, err := h.authService.IsCodeValid(resetPasswordReq.Email, resetPasswordReq.Code)
	if err != nil {
		h.logger.Error("IsCodeValid error", "error", err)
		ctx.JSON(400, models.Error{Message: "Invalid code"})
		return
	}
	if !isvalid {
		ctx.JSON(400, models.Error{Message: "Invalid email or code"})
		return
	}

	resp, err := h.authService.ResetPassword(resetPasswordReq)
	if err != nil {
		h.logger.Error("ResetPassword error", "error", err)
		ctx.JSON(500, models.Error{Message: "Error resetting password"})
		return
	}

	ctx.JSON(200, resp)
}

// @summary Refresh token
// @description Refresh user token
// @accept json
// @produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} models.Error
// @Failure 500 {object} models.Error
// @router /auth/refresh-token [post]
func (h *userHandlerImpl) RefreshToken(ctx *gin.Context) {
	val, ok := ctx.Get("claims")
	if !ok {
		h.logger.Error("Token not found in context")
		ctx.JSON(401, models.Error{Message: "Unauthorized"})
		return
	}
	claims, ok := val.(*token.Claims)
	if !ok {
		h.logger.Error("Token claims not found in context")
		ctx.JSON(401, models.Error{Message: "Unauthorized"})
		return
	}

	is, err := h.authService.IsRefreshTokenValid(claims.Email)
	if err != nil {
		h.logger.Error("IsRefreshTokenValid error", "error", err)
		ctx.JSON(401, models.Error{Message: "Unauthorized"})
		return
	}
	if !is {
		ctx.JSON(401, models.Error{Message: "Invalid refresh token"})
		return
	}

	accessToken, err := token.GeneratedJWTTokenAccess(models.User{
		ID:    claims.ID,
		Email: claims.Email,
		Role:  claims.Role,
	})
	if err != nil {
		h.logger.Error("GeneratedJwtTokenRefresh error", "error", err)
		ctx.JSON(500, models.Error{
			Message: "Error generating refresh token",
		})
		return
	}

	ctx.SetCookie("access_token", accessToken, 3600, "/", "", false, true)
	ctx.JSON(200, gin.H{
		"access_token": accessToken,
	})
}
