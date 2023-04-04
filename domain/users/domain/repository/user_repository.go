package repository

import (
	"context"

	"backend_crudgo/domain/users/domain/model"
	response "backend_crudgo/types"
)

//UserRepository interfaces handlers users
type UserRepository interface {
	CreateUserHandler(ctx context.Context, user *model.User) (*response.CreateResponse, error)
	GetUserHandler(ctx context.Context, id string) (*response.GenericUserResponse, error)
	LoginUserHandler(ctx context.Context, user *model.User) (*response.GenericUserResponse, error)
	GetUsersHandler(ctx context.Context) (*response.GenericUserResponse, error)
}
