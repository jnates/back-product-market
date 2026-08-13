// Package repository defines the persistence contract for products.
package repository

import (
	"context"

	"backend_crudgo/domain/products/domain/model"
)

// ProductRepository defines data access operations for products.
// Implementations must map apperrors.ErrNotFound when a product does not exist.
type ProductRepository interface {
	// CreateProduct persists a new product and returns it with its generated ID.
	CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	// GetProduct retrieves a single product by ID.
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	// GetProducts retrieves every product.
	GetProducts(ctx context.Context) ([]*model.Product, error)
	// UpdateProduct updates an existing product identified by id.
	UpdateProduct(ctx context.Context, id int64, product *model.Product) error
	// DeleteProduct removes a product identified by id.
	DeleteProduct(ctx context.Context, id int64) error
}
