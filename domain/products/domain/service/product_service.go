package service

import (
	"context"

	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/domain/products/domain/repository"
	response "backend_crudgo/types"
)

type productService struct {
	ProductRepository repository.ProductRepository
}

type ProductService interface {
	CreateProductHandler(ctx context.Context, product *model.Product) (*response.ProductCreateResponse, error)
	GetProductHandler(ctx context.Context, id string) (*response.GenericResponse, error)
	GetProductsHandler(ctx context.Context) (*response.GenericResponse, error)
}

func NewProductService(ProductRepository repository.ProductRepository) ProductService {
	return &productService{
		ProductRepository: ProductRepository,
	}
}

func (ps *productService) CreateProductHandler(ctx context.Context, product *model.Product) (*response.ProductCreateResponse, error) {
	return ps.ProductRepository.CreateProductHandler(ctx, product)
}

func (ps *productService) GetProductHandler(ctx context.Context, id string) (*response.GenericResponse, error) {
	return ps.ProductRepository.GetProductHandler(ctx, id)
}

func (ps *productService) GetProductsHandler(ctx context.Context) (*response.GenericResponse, error) {
	return ps.ProductRepository.GetProductsHandler(ctx)
}
