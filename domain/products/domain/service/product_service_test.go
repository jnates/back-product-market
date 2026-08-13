package service

import (
	"context"
	"testing"

	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product := &model.Product{
		ProductName:        "Test Product",
		ProductAmount:      10,
		ProductUserCreated: 1,
	}

	created := &model.Product{
		ProductID:          1,
		ProductName:        "Test Product",
		ProductAmount:      10,
		ProductUserCreated: 1,
	}

	mockRepo.On("CreateProduct", context.Background(), product).Return(created, nil)

	res, err := productService.CreateProduct(context.Background(), product)
	assertions.NoError(err)
	assertions.Equal(created, res)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "CreateProduct", 1)
}

func TestGetProductNotFound(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	mockRepo.On("GetProduct", context.Background(), int64(2)).Return(nil, apperrors.ErrNotFound)

	res, err := productService.GetProduct(context.Background(), 2)
	assertions.ErrorIs(err, apperrors.ErrNotFound)
	assertions.Nil(res)
}

func TestGetProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product := &model.Product{
		ProductID:          1,
		ProductName:        "Test Product",
		ProductAmount:      10,
		ProductUserCreated: 1,
	}

	mockRepo.On("GetProduct", context.Background(), int64(1)).Return(product, nil)

	res, err := productService.GetProduct(context.Background(), 1)
	assertions.NoError(err)
	assertions.Equal(product, res)
}

func TestGetProducts(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	products := []*model.Product{
		{ProductID: 1, ProductName: "Test Product 1", ProductAmount: 10, ProductUserCreated: 1},
		{ProductID: 2, ProductName: "Test Product 2", ProductAmount: 20, ProductUserCreated: 1},
	}

	mockRepo.On("GetProducts", context.Background()).Return(products, nil)

	res, err := productService.GetProducts(context.Background())
	assertions.NoError(err)
	assertions.Equal(products, res)
}

func TestUpdateProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product := &model.Product{
		ProductName:       "Test Product",
		ProductAmount:     20,
		ProductUserModify: 2,
	}

	mockRepo.On("UpdateProduct", context.Background(), int64(1), product).Return(nil)

	err := productService.UpdateProduct(context.Background(), 1, product)
	assertions.NoError(err)
}

func TestDeleteProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	mockRepo.On("DeleteProduct", context.Background(), int64(1)).Return(nil)

	err := productService.DeleteProduct(context.Background(), 1)
	assertions.NoError(err)
}
