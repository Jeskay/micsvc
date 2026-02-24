package server

import (
	"context"
	"net/http"

	"github.com/Jeskay/micsvc/internal/auth"
	"github.com/Jeskay/micsvc/internal/dto"
	"github.com/Jeskay/micsvc/internal/user"
	proto "github.com/Jeskay/micsvc/protos"
)

type UserServer struct {
	userSvc *user.Service
	authSvc *auth.Service
	proto.UnimplementedUsersServer
}

func NewUserGRPC(authSvc *auth.Service, userSvc *user.Service) proto.UsersServer {
	return &UserServer{
		userSvc: userSvc,
		authSvc: authSvc,
	}
}

func (s *UserServer) Login(ctx context.Context, r *proto.LoginRequest) (*proto.LoginResponse, error) {
	ok, token := s.authSvc.Login(r.Id, r.Password)
	if !ok {
		return &proto.LoginResponse{Status: http.StatusUnauthorized}, nil
	}
	return &proto.LoginResponse{Status: http.StatusAccepted, Token: token}, nil
}

func (s *UserServer) Register(ctx context.Context, r *proto.RegisterUser) (*proto.StatusResponse, error) {
	if err := s.authSvc.Register(&dto.User{ID: r.Id, Password: r.Password}); err != nil {
		return &proto.StatusResponse{Code: http.StatusBadGateway}, err
	}
	return &proto.StatusResponse{Code: http.StatusCreated}, nil
}

func (s *UserServer) GetAll(ctx context.Context, r *proto.Empty) (*proto.UsersResponse, error) {
	users, err := s.userSvc.GetAll()
	if err != nil {
		return &proto.UsersResponse{}, err
	}
	pUsers := make([]*proto.UserData, len(users))
	for i, u := range users {
		pUsers[i] = &proto.UserData{Id: u.ID, Name: u.Name, Surname: u.Surname, Email: u.Email}
	}
	return &proto.UsersResponse{Users: pUsers}, err
}

func (s *UserServer) Post(ctx context.Context, r *proto.User) (*proto.StatusResponse, error) {
	err := s.userSvc.Add(&dto.User{
		ID:       r.Id,
		Name:     r.Name,
		Surname:  r.Surname,
		Email:    r.Email,
		Password: r.Password,
	})
	if err != nil {
		return &proto.StatusResponse{Code: http.StatusBadRequest}, err
	}
	return &proto.StatusResponse{Code: http.StatusCreated}, nil
}

func (s *UserServer) Update(ctx context.Context, r *proto.UpdateRequest) (*proto.StatusResponse, error) {
	user := &dto.User{ID: r.Data.Id, Name: r.Data.Name, Surname: r.Data.Surname, Email: r.Data.Email}
	if err := s.userSvc.Update(r.Key.Id, user); err != nil {
		return &proto.StatusResponse{Code: http.StatusBadRequest}, err
	}
	return &proto.StatusResponse{Code: http.StatusOK}, nil
}

func (s *UserServer) Delete(ctx context.Context, r *proto.UserKey) (*proto.StatusResponse, error) {
	if err := s.userSvc.Delete(r.Id); err != nil {
		return &proto.StatusResponse{Code: http.StatusBadRequest}, err
	}
	return &proto.StatusResponse{Code: http.StatusOK}, nil
}
