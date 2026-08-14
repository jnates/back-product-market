// Package interfaces defines the hexagonal-architecture ports for the products
// domain: ProductRepository, ProductService, and ProductHandler.
//
// revive:disable:var-naming The package name "interfaces" is meaningful in this context.
package interfaces

import (
	"context"

	"backend_crudgo/products/models"

	"github.com/labstack/echo/v4"
)

// ProductRepository defines data access operations for products.
// Implementations must map customErrors.ErrNotFound when a product does not exist.
type ProductRepository interface {
	// CreateProduct persists a new product and returns it with its generated ID.
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	// GetProduct retrieves a single product by ID.
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	// GetProducts retrieves every product.
	GetProducts(ctx context.Context) ([]*models.Product, error)
	// UpdateProduct updates an existing product identified by id.
	UpdateProduct(ctx context.Context, id int64, product *models.Product) error
	// DeleteProduct removes a product identified by id.
	DeleteProduct(ctx context.Context, id int64) error
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
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)

	// GetProduct retrieves a single product by ID.
	//
	// Parameters:
	//   - id: the product identifier
	//
	// Returns:
	//   - the matching product
	//   - customErrors.ErrNotFound if no product matches id
	GetProduct(ctx context.Context, id int64) (*models.Product, error)

	// GetProducts retrieves every product.
	//
	// Returns:
	//   - the list of products, empty if none exist
	GetProducts(ctx context.Context) ([]*models.Product, error)

	// UpdateProduct updates an existing product identified by id.
	//
	// Parameters:
	//   - id: the product identifier
	//   - product: the new field values
	//
	// Returns:
	//   - customErrors.ErrNotFound if no product matches id
	UpdateProduct(ctx context.Context, id int64, product *models.Product) error

	// DeleteProduct removes a product identified by id.
	//
	// Parameters:
	//   - id: the product identifier
	//
	// Returns:
	//   - customErrors.ErrNotFound if no product matches id
	DeleteProduct(ctx context.Context, id int64) error
}

// ProductHandler exposes the HTTP endpoints for products.
type ProductHandler interface {
	// CreateProduct handles POST /products.
	CreateProduct(c echo.Context) error
	// GetProduct handles GET /products/:id.
	GetProduct(c echo.Context) error
	// GetProducts handles GET /products.
	GetProducts(c echo.Context) error
	// UpdateProduct handles PUT /products/:id.
	UpdateProduct(c echo.Context) error
	// DeleteProduct handles DELETE /products/:id.
	DeleteProduct(c echo.Context) error
}
