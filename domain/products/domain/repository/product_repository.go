package repository

import (
	"context"

	"backend_crudgo/domain/products/domain/model"
	response "backend_crudgo/types"
)

// ProductRepository interfaces handlers products
type ProductRepository interface {
	CreateProductHandler(ctx context.Context, product *model.Product) (*response.CreateResponse, error)
	GetProductHandler(ctx context.Context, id string) (*response.GenericResponse, error)
	GetProductsHandler(ctx context.Context) (*response.GenericResponse, error)
	UpdateProductHandler(ctx context.Context, id string, product *model.Product) (*response.GenericResponse, error)
	DeleteProductHandler(ctx context.Context, id string) (*response.GenericResponse, error)
}
