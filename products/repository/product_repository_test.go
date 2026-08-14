package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/products/constants"
	"backend_crudgo/products/models"

	"github.com/jackc/pgx/v5"
	"github.com/jnates/go-toolkit/tools/querybuilder"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// anyArgs builds n wildcard argument matchers, used when the exact values
// (e.g. time.Now()) aren't predictable from the test.
func anyArgs(n int) []interface{} {
	args := make([]interface{}, n)
	for i := range args {
		args[i] = pgxmock.AnyArg()
	}
	return args
}

func TestCreateProduct(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	product := &models.Product{ProductName: "Test", ProductAmount: 10, ProductPrice: 5.5, ProductUserCreated: 1}

	query, _, errBuilder := querybuilder.NewInsertBuilder().
		Into(constants.ProductsTable).
		Columns(constants.ColumnProductName, constants.ColumnProductAmount, constants.ColumnProductPrice,
			constants.ColumnProductUserCreated, constants.ColumnProductDateCreated,
			constants.ColumnProductUserModify, constants.ColumnProductDateModify).
		Values(product.ProductName, product.ProductAmount, product.ProductPrice, product.ProductUserCreated,
			time.Now(), product.ProductUserCreated, time.Now()).
		Returning(constants.ColumnProductID).
		Build()
	require.NoError(t, errBuilder)

	pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(anyArgs(7)...).
		WillReturnRows(pgxmock.NewRows([]string{constants.ColumnProductID}).AddRow(int64(1)))

	created, err := repo.CreateProduct(context.Background(), product)

	require.NoError(t, err)
	assert.Equal(t, int64(1), created.ProductID)
	assert.Equal(t, product.ProductName, created.ProductName)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestCreateProductError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	product := &models.Product{ProductName: "Test", ProductUserCreated: 1}

	pool.ExpectQuery(".*").WithArgs(anyArgs(7)...).WillReturnError(errors.New("db down"))

	created, err := repo.CreateProduct(context.Background(), product)

	assert.Error(t, err)
	assert.Nil(t, created)
}

func TestGetProduct(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	now := time.Now()

	query, _, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.ProductsTable).
		Where(constants.ColumnProductID, querybuilder.OpEqual, int64(1)).
		Limit(1).
		Build()
	require.NoError(t, errBuilder)

	pool.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(int64(1)).
		WillReturnRows(pgxmock.NewRows(constants.Columns).
			AddRow(int64(1), "Test", 10, 5.5, 1, now, 1, now))

	product, err := repo.GetProduct(context.Background(), 1)

	require.NoError(t, err)
	assert.Equal(t, int64(1), product.ProductID)
	assert.Equal(t, "Test", product.ProductName)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestGetProductNotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)

	pool.ExpectQuery(".*").WithArgs(anyArgs(1)...).WillReturnError(pgx.ErrNoRows)

	product, err := repo.GetProduct(context.Background(), 99)

	assert.ErrorIs(t, err, customErrors.ErrNotFound)
	assert.Nil(t, product)
}

func TestGetProductQueryError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)

	pool.ExpectQuery(".*").WithArgs(anyArgs(1)...).WillReturnError(errors.New("db down"))

	product, err := repo.GetProduct(context.Background(), 1)

	assert.Error(t, err)
	assert.False(t, errors.Is(err, customErrors.ErrNotFound))
	assert.Nil(t, product)
}

func TestGetProducts(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	now := time.Now()

	query, _, errBuilder := querybuilder.NewSelectBuilder().
		Select(constants.Columns...).
		From(constants.ProductsTable).
		OrderBy(constants.ColumnProductID, querybuilder.Asc).
		Build()
	require.NoError(t, errBuilder)

	pool.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(pgxmock.NewRows(constants.Columns).
			AddRow(int64(1), "Test 1", 10, 5.5, 1, now, 1, now).
			AddRow(int64(2), "Test 2", 20, 6.5, 1, now, 1, now))

	products, err := repo.GetProducts(context.Background())

	require.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, int64(1), products[0].ProductID)
	assert.Equal(t, int64(2), products[1].ProductID)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestGetProductsQueryError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)

	pool.ExpectQuery(".*").WillReturnError(errors.New("db down"))

	products, err := repo.GetProducts(context.Background())

	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestUpdateProduct(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	product := &models.Product{ProductName: "Updated", ProductAmount: 20, ProductUserModify: 2}

	pool.ExpectExec(".*").WithArgs(anyArgs(5)...).WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateProduct(context.Background(), 1, product)

	require.NoError(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestUpdateProductNotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	product := &models.Product{ProductName: "Updated"}

	pool.ExpectExec(".*").WithArgs(anyArgs(5)...).WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err = repo.UpdateProduct(context.Background(), 99, product)

	assert.ErrorIs(t, err, customErrors.ErrNotFound)
}

func TestUpdateProductExecError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)
	product := &models.Product{ProductName: "Updated"}

	pool.ExpectExec(".*").WithArgs(anyArgs(5)...).WillReturnError(errors.New("db down"))

	err = repo.UpdateProduct(context.Background(), 1, product)

	assert.Error(t, err)
	assert.False(t, errors.Is(err, customErrors.ErrNotFound))
}

func TestDeleteProduct(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)

	pool.ExpectExec(".*").WithArgs(anyArgs(1)...).WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.DeleteProduct(context.Background(), 1)

	require.NoError(t, err)
	assert.NoError(t, pool.ExpectationsWereMet())
}

func TestDeleteProductExecError(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)

	pool.ExpectExec(".*").WithArgs(anyArgs(1)...).WillReturnError(errors.New("db down"))

	err = repo.DeleteProduct(context.Background(), 1)

	assert.Error(t, err)
	assert.False(t, errors.Is(err, customErrors.ErrNotFound))
}

func TestDeleteProductNotFound(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	repo := NewProductRepository(pool)

	pool.ExpectExec(".*").WithArgs(anyArgs(1)...).WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err = repo.DeleteProduct(context.Background(), 99)

	assert.ErrorIs(t, err, customErrors.ErrNotFound)
}
