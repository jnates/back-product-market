package service

import (
	"backend_crudgo/mocks"
	"context"
	"errors"
	"testing"

	"backend_crudgo/domain/products/domain/model"
	response "backend_crudgo/types"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product := &model.Product{
		ProductID:          "1",
		ProductName:        "Test Product",
		ProductAmount:      10,
		ProductUserCreated: 1,
	}

	expectedResponse := &response.CreateResponse{
		Message: "Product created",
	}

	mockRepo.On("CreateProduct", context.Background(), product).Return(expectedResponse, nil)

	res, err := productService.CreateProduct(context.Background(), product)
	assertions.NoError(err)
	assertions.Equal(expectedResponse, res)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "CreateProduct", 1)
	mockRepo.AssertCalled(t, "CreateProduct", context.Background(), product)
}

func TestGetProductNotFound(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	expectedResponse := &response.GenericResponse{
		Error: "Product not found",
	}

	mockRepo.On("GetProduct", context.Background(), "2").Return(expectedResponse, errors.New("Product not found"))

	res, err := productService.GetProduct(context.Background(), "2")
	assertions.Error(err)
	assertions.Equal(expectedResponse, res)
}

func TestGetProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product := &model.Product{
		ProductID:          "1",
		ProductName:        "Test Product",
		ProductAmount:      10,
		ProductUserCreated: 1,
	}

	expectedResponse := &response.GenericResponse{
		Message: "Get product success",
		Product: product,
	}

	mockRepo.On("GetProduct", context.Background(), "1").Return(expectedResponse, nil)

	res, err := productService.GetProduct(context.Background(), "1")
	assertions.NoError(err)
	assertions.Equal(expectedResponse, res)
}

func TestGetProducts(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product1 := &model.Product{
		ProductID:          "1",
		ProductName:        "Test Product 1",
		ProductAmount:      10,
		ProductUserCreated: 1,
	}

	product2 := &model.Product{
		ProductID:          "2",
		ProductName:        "Test Product 2",
		ProductAmount:      20,
		ProductUserCreated: 1,
	}

	expectedResponse := &response.GenericResponse{
		Message: "Get products success",
		Product: []*model.Product{product1, product2},
	}

	mockRepo.On("GetProducts", context.Background()).Return(expectedResponse, nil)

	res, err := productService.GetProducts(context.Background())
	assertions.NoError(err)
	assertions.Equal(expectedResponse, res)
}

func TestUpdateProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	product := &model.Product{
		ProductID:         "1",
		ProductName:       "Test Product",
		ProductAmount:     20,
		ProductUserModify: 2,
	}

	expectedResponse := &response.GenericResponse{
		Message: "Product updated",
	}

	mockRepo.On("UpdateProduct", context.Background(), "1", product).Return(expectedResponse, nil)

	res, err := productService.UpdateProduct(context.Background(), "1", product)
	assertions.NoError(err)
	assertions.Equal(expectedResponse, res)
}

func TestDeleteProduct(t *testing.T) {
	assertions := assert.New(t)
	mockRepo := mocks.NewProductRepository(t)
	productService := NewProductService(mockRepo)

	expectedResponse := &response.GenericResponse{
		Message: "Product deleted",
	}

	mockRepo.On("DeleteProduct", context.Background(), "1").Return(expectedResponse, nil)

	res, err := productService.DeleteProduct(context.Background(), "1")
	assertions.NoError(err)
	assertions.Equal(expectedResponse, res)
}
