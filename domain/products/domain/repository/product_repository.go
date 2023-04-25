package repository

import (
	"context"

	"backend_crudgo/domain/products/domain/model"
	response "backend_crudgo/types"
)

// ProductRepository interfaces handlers products.
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) (*response.CreateResponse, error)
	GetProduct(ctx context.Context, id string) (*response.GenericResponse, error)
	GetProducts(ctx context.Context) (*response.GenericResponse, error)
	UpdateProduct(ctx context.Context, id string, product *model.Product) (*response.GenericResponse, error)
	DeleteProduct(ctx context.Context, id string) (*response.GenericResponse, error)
}
