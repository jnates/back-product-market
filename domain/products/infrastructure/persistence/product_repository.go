// Package persistence implements the product repository against PostgreSQL.
package persistence

import (
	"context"
	"errors"
	"time"

	"backend_crudgo/domain/products/constants"
	"backend_crudgo/domain/products/domain/model"
	repoDomain "backend_crudgo/domain/products/domain/repository"
	"backend_crudgo/infrastructure/kit/apperrors"

	"github.com/jackc/pgx/v5"
	"github.com/jnates/go-toolkit/tools/querybuilder"
	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"

	"github.com/rs/zerolog/log"
)

type sqlProductRepo struct {
	pool pgxtool.DBPool
}

// NewProductRepository builds a ProductRepository backed by the given connection pool.
func NewProductRepository(pool pgxtool.DBPool) repoDomain.ProductRepository {
	return &sqlProductRepo{pool: pool}
}

// CreateProduct inserts a new product and returns it with its generated ID.
//
// Parameters:
//   - product: the product to persist; ProductUserCreated must be set by the caller
//
// Returns:
//   - the created product including its generated ID and timestamps
//   - an error if the insert fails
func (sr *sqlProductRepo) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
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
	if err := sr.pool.QueryRow(ctx, query, args...).Scan(&productID); err != nil {
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
//   - apperrors.ErrNotFound if no product matches id
func (sr *sqlProductRepo) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
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

	var product ProductDB
	row := sr.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&product.ProductID, &product.ProductName, &product.ProductAmount, &product.ProductPrice,
		&product.ProductUserCreated, &product.ProductDateCreated, &product.ProductUserModify,
		&product.ProductDateModify); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}

	return product.ToModel(), nil
}

// GetProducts retrieves every product.
//
// Returns:
//   - the list of products, empty if none exist
//   - an error if the query fails
func (sr *sqlProductRepo) GetProducts(ctx context.Context) ([]*model.Product, error) {
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

	rows, err := sr.pool.Query(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing select query")
		return nil, err
	}
	defer rows.Close()

	items, err := pgxtool.ScanRows(rows, subLogger, scanProductDB)
	if err != nil {
		return nil, err
	}

	products := make([]*model.Product, 0, len(items))
	for i := range items {
		products = append(products, items[i].ToModel())
	}

	return products, nil
}

// scanProductDB scans a single row into a ProductDB, following column order in constants.Columns.
func scanProductDB(s pgxtool.Scanner) (ProductDB, error) {
	var p ProductDB
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
//   - apperrors.ErrNotFound if no product matches id
func (sr *sqlProductRepo) UpdateProduct(ctx context.Context, id int64, product *model.Product) error {
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

	tag, err := sr.pool.Exec(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing update query")
		return err
	}

	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

// DeleteProduct removes a product identified by id.
//
// Parameters:
//   - id: the product identifier
//
// Returns:
//   - apperrors.ErrNotFound if no product matches id
func (sr *sqlProductRepo) DeleteProduct(ctx context.Context, id int64) error {
	subLogger := log.With().Str("repository", "ProductRepository").Str("method", "DeleteProduct").Logger()

	query, args, errBuilder := querybuilder.NewDeleteBuilder().
		Table(constants.ProductsTable).
		Where(constants.ColumnProductID, querybuilder.OpEqual, id).
		Build()
	if errBuilder != nil {
		subLogger.Error().Err(errBuilder).Msg("error building delete query")
		return errBuilder
	}

	tag, err := sr.pool.Exec(ctx, query, args...)
	if err != nil {
		subLogger.Error().Err(err).Msg("error executing delete query")
		return err
	}

	if tag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
