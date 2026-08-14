package mapper_test

import (
	"testing"
	"time"

	"backend_crudgo/products/mapper"
	"backend_crudgo/products/models"

	"github.com/stretchr/testify/assert"
)

func TestToProductModel(t *testing.T) {
	now := time.Now()
	dbProduct := &models.ProductDB{
		ProductID:          1,
		ProductName:        "Test",
		ProductAmount:      10,
		ProductPrice:       5.5,
		ProductUserCreated: 1,
		ProductDateCreated: now,
		ProductUserModify:  1,
		ProductDateModify:  now,
	}

	product := mapper.ToProductModel(dbProduct)

	assert.Equal(t, dbProduct.ProductID, product.ProductID)
	assert.Equal(t, dbProduct.ProductName, product.ProductName)
	assert.Equal(t, dbProduct.ProductAmount, product.ProductAmount)
	assert.Equal(t, dbProduct.ProductPrice, product.ProductPrice)
	assert.Equal(t, dbProduct.ProductUserCreated, product.ProductUserCreated)
	assert.Equal(t, dbProduct.ProductDateCreated, product.ProductDateCreated)
	assert.Equal(t, dbProduct.ProductUserModify, product.ProductUserModify)
	assert.Equal(t, dbProduct.ProductDateModify, product.ProductDateModify)
}
