// Package service contains the product domain's business logic.
package service

import (
	"context"

	"backend_crudgo/products/interfaces"
	"backend_crudgo/products/models"
)

// productService implements interfaces.ProductService, delegating to the product repository.
type productService struct {
	productRepository interfaces.ProductRepository
}

// NewProductService builds a ProductService backed by the given repository.
func NewProductService(productRepository interfaces.ProductRepository) interfaces.ProductService {
	return &productService{
		productRepository: productRepository,
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
func (s *productService) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	return s.productRepository.CreateProduct(ctx, product)
}

// GetProduct retrieves a single product by ID by delegating to the repository.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - the matching product
//   - customErrors.ErrNotFound if no product matches id
func (s *productService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	return s.productRepository.GetProduct(ctx, id)
}

// GetProducts retrieves every product by delegating to the repository.
//
// Returns:
//   - the list of products, empty if none exist
func (s *productService) GetProducts(ctx context.Context) ([]*models.Product, error) {
	return s.productRepository.GetProducts(ctx)
}

// UpdateProduct updates an existing product identified by id by delegating to the repository.
//
// Parameters:
//   - id: the product identifier
//   - product: the new field values
//
// Returns:
//   - customErrors.ErrNotFound if no product matches id
func (s *productService) UpdateProduct(ctx context.Context, id int64, product *models.Product) error {
	return s.productRepository.UpdateProduct(ctx, id, product)
}

// DeleteProduct removes a product identified by id by delegating to the repository.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - customErrors.ErrNotFound if no product matches id
func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepository.DeleteProduct(ctx, id)
}
