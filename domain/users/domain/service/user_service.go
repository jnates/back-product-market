package service

import (
	"context"

	"backend_crudgo/domain/users/domain/model"
	"backend_crudgo/domain/users/domain/repository"
	response "backend_crudgo/types"
)

type userService struct {
	UserRepository repository.UserRepository
}

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) (*response.CreateResponse, error)
	LoginUser(ctx context.Context, user *model.User) (*response.GenericUserResponse, error)
	GetUser(ctx context.Context, id string) (*response.GenericUserResponse, error)
	GetUsers(ctx context.Context) (*response.GenericUserResponse, error)
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		UserRepository: userRepository,
	}
}

func (ps *userService) CreateUser(ctx context.Context, user *model.User) (*response.CreateResponse, error) {
	return ps.UserRepository.CreateUser(ctx, user)
}

func (ps *userService) GetUser(ctx context.Context, id string) (*response.GenericUserResponse, error) {
	return ps.UserRepository.GetUser(ctx, id)
}

func (ps *userService) LoginUser(ctx context.Context, user *model.User) (*response.GenericUserResponse, error) {
	return ps.UserRepository.LoginUser(ctx, user)
}

func (ps *userService) GetUsers(ctx context.Context) (*response.GenericUserResponse, error) {
	return ps.UserRepository.GetUsers(ctx)
}
