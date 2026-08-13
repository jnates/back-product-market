// Package service contains the product domain's business logic.
package service

import (
	"context"

	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/domain/products/domain/repository"
)

type productService struct {
	ProductRepository repository.ProductRepository
}

// ProductService exposes the product use cases consumed by the HTTP handlers.
type ProductService interface {
	// CreateProduct creates a new product.
	//
	// Parameters:
	//   - product: the product to create
	//
	// Returns:
	//   - the created product with its generated ID
	//   - an error if persistence fails
	CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error)

	// GetProduct retrieves a single product by ID.
	//
	// Parameters:
	//   - id: the product identifier
	//
	// Returns:
	//   - the matching product
	//   - apperrors.ErrNotFound if no product matches id
	GetProduct(ctx context.Context, id int64) (*model.Product, error)

	// GetProducts retrieves every product.
	//
	// Returns:
	//   - the list of products, empty if none exist
	GetProducts(ctx context.Context) ([]*model.Product, error)

	// UpdateProduct updates an existing product identified by id.
	//
	// Parameters:
	//   - id: the product identifier
	//   - product: the new field values
	//
	// Returns:
	//   - apperrors.ErrNotFound if no product matches id
	UpdateProduct(ctx context.Context, id int64, product *model.Product) error

	// DeleteProduct removes a product identified by id.
	//
	// Parameters:
	//   - id: the product identifier
	//
	// Returns:
	//   - apperrors.ErrNotFound if no product matches id
	DeleteProduct(ctx context.Context, id int64) error
}

// NewProductService builds a ProductService backed by the given repository.
func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{
		ProductRepository: productRepository,
	}
}

// CreateProduct creates a new product by delegating to the repository.
//
// Parameters:
//   - product: the product to create
//
// Returns:
//   - the created product with its generated ID
//   - an error if persistence fails
func (ps *productService) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	return ps.ProductRepository.CreateProduct(ctx, product)
}

// GetProduct retrieves a single product by ID by delegating to the repository.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - the matching product
//   - apperrors.ErrNotFound if no product matches id
func (ps *productService) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	return ps.ProductRepository.GetProduct(ctx, id)
}

// GetProducts retrieves every product by delegating to the repository.
//
// Returns:
//   - the list of products, empty if none exist
func (ps *productService) GetProducts(ctx context.Context) ([]*model.Product, error) {
	return ps.ProductRepository.GetProducts(ctx)
}

// UpdateProduct updates an existing product identified by id by delegating to the repository.
//
// Parameters:
//   - id: the product identifier
//   - product: the new field values
//
// Returns:
//   - apperrors.ErrNotFound if no product matches id
func (ps *productService) UpdateProduct(ctx context.Context, id int64, product *model.Product) error {
	return ps.ProductRepository.UpdateProduct(ctx, id, product)
}

// DeleteProduct removes a product identified by id by delegating to the repository.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - apperrors.ErrNotFound if no product matches id
func (ps *productService) DeleteProduct(ctx context.Context, id int64) error {
	return ps.ProductRepository.DeleteProduct(ctx, id)
}
