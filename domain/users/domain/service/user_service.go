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
	CreateUserHandler(ctx context.Context, user *model.User) (*response.CreateResponse, error)
	LoginUserHandler(ctx context.Context, user *model.User) (*response.GenericUserResponse, error)
	GetUserHandler(ctx context.Context, id string) (*response.GenericUserResponse, error)
	GetUsersHandler(ctx context.Context) (*response.GenericUserResponse, error)
}

func NewUserService(UserRepository repository.UserRepository) UserService {
	return &userService{
		UserRepository: UserRepository,
	}
}

func (ps *userService) CreateUserHandler(ctx context.Context, user *model.User) (*response.CreateResponse, error) {
	return ps.UserRepository.CreateUserHandler(ctx, user)
}

func (ps *userService) GetUserHandler(ctx context.Context, id string) (*response.GenericUserResponse, error) {
	return ps.UserRepository.GetUserHandler(ctx, id)
}

func (ps *userService) LoginUserHandler(ctx context.Context, user *model.User) (*response.GenericUserResponse, error) {
	return ps.UserRepository.LoginUserHandler(ctx, user)
}

func (ps *userService) GetUsersHandler(ctx context.Context) (*response.GenericUserResponse, error) {
	return ps.UserRepository.GetUsersHandler(ctx)
}
