package service

import (
	pb "auth-service/generated/user"
	"auth-service/storage"
	"context"
	"log/slog"
)

type UserService interface {
	GetUserProfile(context.Context, *pb.GetUserProfileReq) (*pb.UserProfile, error)
	UpdateUserProfile(context.Context, *pb.UpdateUserProfileReq) (*pb.UpdateUserProfileResp, error)
	GetUsersList(context.Context, *pb.GetUsersListReq) (*pb.GetUsersListResp, error)
	ChangePassword(context.Context, *pb.ChangePasswordReq) (*pb.ChangePasswordResp, error)
	ValidateToken(context.Context, *pb.ValidateTokenReq) (*pb.ValidateTokenResp, error)
}

type userServiceImpl struct {
	pb.UnimplementedAuthServiceServer
	storage storage.IStorage
	logger  *slog.Logger
}

func NewUserService(storage storage.IStorage, logger *slog.Logger) *userServiceImpl {
	return &userServiceImpl{
		storage: storage,
		logger:  logger,
	}
}

func (s *userServiceImpl) GetUserProfile(ctx context.Context, req *pb.GetUserProfileReq) (*pb.UserProfile, error) {
	resp, err := s.storage.UserRepository().GetUserProfile(req.GetId())
	if err != nil {
		s.logger.Error("GetUserProfile error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *userServiceImpl) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileReq) (*pb.UpdateUserProfileResp, error) {
	resp, err := s.storage.UserRepository().UpdateUserProfile(req)
	if err != nil {
		s.logger.Error("UpdateUserProfile error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *userServiceImpl) GetUsersList(ctx context.Context, req *pb.GetUsersListReq) (*pb.GetUsersListResp, error) {
	resp, err := s.storage.UserRepository().GetUsersList(req)
	if err != nil {
		s.logger.Error("GetUsersList error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *userServiceImpl) ChangePassword(ctx context.Context, req *pb.ChangePasswordReq) (*pb.ChangePasswordResp, error) {
	resp, err := s.storage.UserRepository().ChangePassword(req.GetId(), req)
	if err != nil {
		s.logger.Error("ChangePassword error", "error", err)
		return nil, err
	}
	return resp, nil
}

func (s *userServiceImpl) ValidateToken(ctx context.Context, req *pb.ValidateTokenReq) (*pb.ValidateTokenResp, error) {
	return nil, nil
}
