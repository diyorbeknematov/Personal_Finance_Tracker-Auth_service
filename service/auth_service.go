package service

import (
	"auth-service/models"
	"auth-service/storage"
	"log/slog"
	"time"
)

type AuthService interface {
	RegisterUser(user models.RegisterUser) (*models.Response, error)
	EmailExists(email string) (bool, error)
	LoginUser(login models.LoginUserReq) (*models.User, error)
	DeleteUser(id string) (*models.Response, error)
	ResetPassword(reset models.ResetPassword) (*models.Response, error)
	SaveRefreshToken(refreshToken models.RefreshToken) (*models.Response, error)
	InvalidateRefreshToken(email string) (*models.Response, error)
	IsRefreshTokenValid(email string) (bool, error)
	UpdateUserRoles(manage models.ManageUserRoles) (*models.Response, error)

	AddTokenBlacklist(token string, expirationTime time.Duration) (*models.Response, error)
	IsTokenBlacklisted(token string) (bool, error)
	StoreCode(email, code string, expirationTime time.Duration) (*models.Response, error)
	IsCodeValid(email, code string) (bool, error)
}

type authServiceImpl struct {
	storage storage.IStorage
	logger  *slog.Logger
}

func NewAuthService(storage storage.IStorage, logger *slog.Logger) AuthService {
	return &authServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *authServiceImpl) RegisterUser(user models.RegisterUser) (*models.Response, error) {
	resp, err := s.storage.AuthRepository().RegisterUser(user)
	if err != nil {
		s.logger.Error("RegisterUser error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) EmailExists(email string) (bool, error) {
	resp, err := s.storage.AuthRepository().EmailExists(email)
	if err != nil {
		s.logger.Error("EmailExists error", "error", err)
		return false, err
	}
	return resp, nil
}

func (s *authServiceImpl) LoginUser(login models.LoginUserReq) (*models.User, error) {
	resp, err := s.storage.AuthRepository().LoginUser(login)
	if err != nil {
		s.logger.Error("LoginUser error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) DeleteUser(id string) (*models.Response, error) {
	resp, err := s.storage.AuthRepository().LogOutUser(id)
	if err != nil {
		s.logger.Error("LogOutUser error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) ResetPassword(reset models.ResetPassword) (*models.Response, error) {
	resp, err := s.storage.AuthRepository().ResetPassword(reset.Email, reset.Password)
	if err != nil {
		s.logger.Error("ResetPassword error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) SaveRefreshToken(refreshToken models.RefreshToken) (*models.Response, error) {
	resp, err := s.storage.AuthRepository().SaveRefreshToken(refreshToken)
	if err != nil {
		s.logger.Error("SaveRefreshToken error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) InvalidateRefreshToken(email string) (*models.Response, error) {
	resp, err := s.storage.AuthRepository().InvalidateRefreshToken(email)
	if err != nil {
		s.logger.Error("InvalidateRefreshToken error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) IsRefreshTokenValid(email string) (bool, error) {
	resp, err := s.storage.AuthRepository().IsRefreshTokenValid(email)
	if err != nil {
		s.logger.Error("IsRefreshTokenValid error", "error", err)
		return false, err
	}
	return resp, nil
}

func (s *authServiceImpl) UpdateUserRoles(manage models.ManageUserRoles) (*models.Response, error) {
	resp, err := s.storage.AuthRepository().ManageUserRoles(manage.Email, manage.Role)
	if err != nil {
		s.logger.Error("ManageUserRoles error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) AddTokenBlacklist(token string, expirationTime time.Duration) (*models.Response, error) {
	resp, err := s.storage.RedisStore().AddTokenBlacklist(token, expirationTime)
	if err != nil {
		s.logger.Error("AddTokenBlacklist error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) IsTokenBlacklisted(token string) (bool, error) {
	resp, err := s.storage.RedisStore().IsTokenBlacklisted(token)
	if err != nil {
		s.logger.Error("IsTokenBlacklisted error", "error", err)
		return false, err
	}
	return resp, nil
}

func (s *authServiceImpl) StoreCode(email, code string, expirationTime time.Duration) (*models.Response, error) {
	resp, err := s.storage.RedisStore().StoreCode(email, code, expirationTime)
	if err != nil {
		s.logger.Error("StoreCode error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *authServiceImpl) IsCodeValid(email, code string) (bool, error) {
	resp, err := s.storage.RedisStore().IsCodeValid(email, code)
	if err != nil {
		s.logger.Error("IsCodeValid error", "error", err)
		return false, err
	}
	return resp, nil
}
