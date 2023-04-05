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
	CreateProduct(ctx context.Context, product *model.Product) (*response.CreateResponse, error)
	GetProduct(ctx context.Context, id string) (*response.GenericResponse, error)
	GetProducts(ctx context.Context) (*response.GenericResponse, error)
	UpdateProduct(ctx context.Context, id string, product *model.Product) (*response.GenericResponse, error)
	DeleteProduct(ctx context.Context, id string) (*response.GenericResponse, error)
}

func NewProductService(ProductRepository repository.ProductRepository) ProductService {
	return &productService{
		ProductRepository: ProductRepository,
	}
}

func (ps *productService) CreateProduct(ctx context.Context, product *model.Product) (*response.CreateResponse, error) {
	return ps.ProductRepository.CreateProduct(ctx, product)
}

func (ps *productService) GetProduct(ctx context.Context, id string) (*response.GenericResponse, error) {
	return ps.ProductRepository.GetProduct(ctx, id)
}

func (ps *productService) GetProducts(ctx context.Context) (*response.GenericResponse, error) {
	return ps.ProductRepository.GetProducts(ctx)
}

func (ps *productService) UpdateProduct(ctx context.Context, id string, product *model.Product) (*response.GenericResponse, error){
	 return ps.ProductRepository.UpdateProduct(ctx, id, product)
}

func (ps *productService) DeleteProduct(ctx context.Context, id string) (*response.GenericResponse, error)  {
	return ps.ProductRepository.DeleteProduct(ctx, id)
}