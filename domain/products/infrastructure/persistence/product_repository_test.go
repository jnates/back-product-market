package persistence_test

//
//import (
//	"backend_crudgo/domain/products/domain/model"
//	"backend_crudgo/domain/products/infrastructure/persistence"
//	"backend_crudgo/infrastructure/database"
//	"context"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestNewProductRepository(t *testing.T) {
//	conn := &database.DataDB{}
//	repo := persistence.NewProductRepository(conn)
//	assert.NotNil(t, repo, "Repository should not be nil")
//	assert.IsType(t, &persistence.SqlProductRepo{}, repo, "Repository should be of type SqlProductRepo")
//}
//
//func TestSqlProductRepo_CreateProductHandler(t *testing.T) {
//	conn := &database.DataDB{}
//	repo := persistence.NewProductRepository(conn)
//
//	product := &model.Product{
//		ProductID:          "12",
//		ProductName:        "Test Product",
//		ProductAmount:      42,
//		ProductUserCreated: 1,
//		ProductUserModify:  1,
//	}
//
//	ctx := context.Background()
//
//	resp, err := repo.CreateProductHandler(ctx, product)
//	assert.NotNil(t, resp, "Response should not be nil")
//	assert.NoError(t, err, "Error should be nil")
//	assert.Equal(t, "Product created", resp.Message, "Response message should be 'Product created'")
//}
//
//func TestSqlProductRepo_GetProductHandler(t *testing.T) {
//	conn := &database.DataDB{}
//	repo := persistence.NewProductRepository(conn)
//
//	productID := "1"
//	ctx := context.Background()
//
//	resp, err := repo.GetProductHandler(ctx, productID)
//	assert.NotNil(t, resp, "Response should not be nil")
//	assert.NoError(t, err, "Error should be nil")
//	assert.Equal(t, "Get product success", resp.Message, "Response message should be 'Get product success'")
//	assert.NotNil(t, resp.Product, "Product should not be nil")
//}
