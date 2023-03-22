package repository

import (
	"context"

	"backend_crudgo/domain/products/domain/model"
	response "backend_crudgo/types"
)

//ProductRepository interfaces handlers products
type ProductRepository interface {
	CreateProductHandler(ctx context.Context, product *model.Product) (*response.ProductCreateResponse, error)
	GetProductHandler(ctx context.Context, id string) (*response.ProductResponse, error)
}
