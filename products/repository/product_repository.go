// Package repository implements the product repository against PostgreSQL.
package repository

import (
	"context"
	"errors"
	"time"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/products/constants"
	"backend_crudgo/products/interfaces"
	"backend_crudgo/products/mapper"
	"backend_crudgo/products/models"

	"github.com/jackc/pgx/v5"
	"github.com/jnates/go-toolkit/tools/querybuilder"
	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"

	"github.com/rs/zerolog/log"
)

// productRepository implements interfaces.ProductRepository against PostgreSQL via tools/querybuilder.
type productRepository struct {
	pool pgxtool.DBPool
}

// NewProductRepository builds a ProductRepository backed by the given connection pool.
func NewProductRepository(pool pgxtool.DBPool) interfaces.ProductRepository {
	return &productRepository{pool: pool}
}

// CreateProduct inserts a new product and returns it with its generated ID.
//
// Parameters:
//   - product: the product to persist; ProductUserCreated must be set by the caller
//
// Returns:
//   - the created product including its generated ID and timestamps
//   - an error if the insert fails
func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	subLogger := log.With().Str("repository", "ProductRepository").Str("method", "CreateProduct").Logger()

	now := time.Now()
	query, args, errBuilder := querybuilder.NewInsertBuilder().
		Into(constants.ProductsTable).
		Columns(constants.ColumnProductName, constants.ColumnProductAmount, constants.ColumnProductPrice,
			constants.ColumnProductUserCreated, constants.ColumnProductDateCreated,
			constants.ColumnProductUserModify, constants.ColumnProductDateModify).
		Values(product.ProductName, product.ProductAmount, product.ProductPrice, product.ProductUserCreated,
			now, product.ProductUserCreated, now).
		Returning(constants.ColumnProductID).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building insert query")
		return nil, errBuilder
	}

	var productID int64
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&productID); err != nil {
		subLogger.Error().Err(err).Msg("error executing insert query")
		return nil, err
	}

	created := *product
	created.ProductID = productID
	created.ProductDateCreated = now
	created.ProductUserModify = product.ProductUserCreated
	created.ProductDateModify = now

	return &created, nil
}

// GetProduct retrieves a single product by ID.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - the matching product
//   - customErrors.ErrNotFound if no product matches id
func (r *productRepository) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	subLogger := log.With().Str("repository", "ProductRepository").Str("method", "GetProduct").Logger()

	query, args, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.ProductsTable).
		Where(constants.ColumnProductID, querybuilder.OpEqual, id).
		Limit(1).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building select query")
		return nil, errBuilder
	}

	var product models.ProductDB
	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductPrice,
		&product.ProductUserCreated, &product.ProductDateCreated, &product.ProductUserModify,
		&product.ProductDateModify); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, customErrors.ErrNotFound
		}
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}

	return mapper.ToProductModel(&product), nil
}

// GetProducts retrieves every product.
//
// Returns:
//   - the list of products, empty if none exist
//   - an error if the query fails
func (r *productRepository) GetProducts(ctx context.Context) ([]*models.Product, error) {
	subLogger := log.With().Str("repository", "ProductRepository").Str("method", "GetProducts").Logger()

	query, args, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.ProductsTable).
		OrderBy(constants.ColumnProductID, querybuilder.Asc).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building select query")
		return nil, errBuilder
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}
	defer rows.Close()

	items, err := pgxtool.ScanRows(rows, subLogger, scanProductDB)
	if err != nil {
		return nil, err
	}

	products := make([]*models.Product, 0, len(items))
	for i := range items {
		products = append(products, mapper.ToProductModel(&items[i]))
	}

	return products, nil
}

// scanProductDB scans a single row into a ProductDB, following column order in constants.Columns.
func scanProductDB(s pgxtool.Scanner) (models.ProductDB, error) {
	var p models.ProductDB
	err := s.Scan(&p.ProductID, &p.ProductName, &p.ProductAmount, &p.ProductPrice,
		&p.ProductUserCreated, &p.ProductDateCreated, &p.ProductUserModify, &p.ProductDateModify)
	return p, err
}

// UpdateProduct updates the name, amount and modifier of an existing product.
//
// Parameters:
//   - id: the product identifier
//   - product: carries the new name, amount and the user performing the modification
//
// Returns:
//   - customErrors.ErrNotFound if no product matches id
func (r *productRepository) UpdateProduct(ctx context.Context, id int64, product *models.Product) error {
	subLogger := log.With().Str("repository", "ProductRepository").Str("method", "UpdateProduct").Logger()

	query, args, errBuilder := querybuilder.NewUpdateBuilder().
		Table(constants.ProductsTable).
		Set(constants.ColumnProductName, product.ProductName).
		Set(constants.ColumnProductAmount, product.ProductAmount).
		Set(constants.ColumnProductUserModify, product.ProductUserModify).
		Set(constants.ColumnProductDateModify, time.Now()).
		Where(constants.ColumnProductID, querybuilder.OpEqual, id).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building update query")
		return errBuilder
	}

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing update query")
		return err
	}

	if tag.RowsAffected() == 0 {
		return customErrors.ErrNotFound
	}

	return nil
}

// DeleteProduct removes a product identified by id.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - customErrors.ErrNotFound if no product matches id
func (r *productRepository) DeleteProduct(ctx context.Context, id int64) error {
	subLogger := log.With().Str("repository", "ProductRepository").Str("method", "DeleteProduct").Logger()

	query, args, errBuilder := querybuilder.NewDeleteBuilder().
		Table(constants.ProductsTable).
		Where(constants.ColumnProductID, querybuilder.OpEqual, id).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building delete query")
		return errBuilder
	}

	tag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing delete query")
		return err
	}

	if tag.RowsAffected() == 0 {
		return customErrors.ErrNotFound
	}

	return nil
}
