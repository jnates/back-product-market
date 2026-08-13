package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend_crudgo/domain/products/constants"
	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/mocks"

	"github.com/labstack/echo/v4"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewProductHandler(t *testing.T) {
	pool, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer pool.Close()

	h := NewProductHandler(pool)

	assert.NotNil(t, h)
	assert.NotNil(t, h.Service)
}

func newProductRequest(t *testing.T, method, target string, body any) (echo.Context, *httptest.ResponseRecorder) {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		assert.NoError(t, err)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, target, reader)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	return echo.New().NewContext(req, rec), rec
}

func TestCreateProductHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		product := model.Product{ProductName: "Test", ProductAmount: 1, ProductUserCreated: 1}
		created := &model.Product{ProductID: 1, ProductName: "Test", ProductAmount: 1, ProductUserCreated: 1}

		mockSvc.On("CreateProduct", mock.Anything, &product).Return(created, nil)

		c, rec := newProductRequest(t, http.MethodPost, "/products", product)

		assert.NoError(t, h.CreateProduct(c))
		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("{invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.NoError(t, h.CreateProduct(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		product := model.Product{ProductName: "Test"}
		mockSvc.On("CreateProduct", mock.Anything, &product).Return(nil, apperrors.ErrConflict)

		c, rec := newProductRequest(t, http.MethodPost, "/products", product)

		assert.NoError(t, h.CreateProduct(c))
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestGetProductHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		product := &model.Product{ProductID: 1, ProductName: "Test"}
		mockSvc.On("GetProduct", mock.Anything, int64(1)).Return(product, nil)

		c, rec := newProductRequest(t, http.MethodGet, "/products/1", nil)
		c.SetParamNames(constants.ID)
		c.SetParamValues("1")

		assert.NoError(t, h.GetProduct(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		c, rec := newProductRequest(t, http.MethodGet, "/products/abc", nil)
		c.SetParamNames(constants.ID)
		c.SetParamValues("abc")

		assert.NoError(t, h.GetProduct(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		mockSvc.On("GetProduct", mock.Anything, int64(99)).Return(nil, apperrors.ErrNotFound)

		c, rec := newProductRequest(t, http.MethodGet, "/products/99", nil)
		c.SetParamNames(constants.ID)
		c.SetParamValues("99")

		assert.NoError(t, h.GetProduct(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestGetProductsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		products := []*model.Product{{ProductID: 1}, {ProductID: 2}}
		mockSvc.On("GetProducts", mock.Anything).Return(products, nil)

		c, rec := newProductRequest(t, http.MethodGet, "/products", nil)

		assert.NoError(t, h.GetProducts(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		mockSvc.On("GetProducts", mock.Anything).Return(nil, assert.AnError)

		c, rec := newProductRequest(t, http.MethodGet, "/products", nil)

		assert.NoError(t, h.GetProducts(c))
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateProductHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		product := model.Product{ProductName: "Updated"}
		mockSvc.On("UpdateProduct", mock.Anything, int64(1), &product).Return(nil)

		c, rec := newProductRequest(t, http.MethodPut, "/products/1", product)
		c.SetParamNames(constants.ID)
		c.SetParamValues("1")

		assert.NoError(t, h.UpdateProduct(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		c, rec := newProductRequest(t, http.MethodPut, "/products/abc", model.Product{})
		c.SetParamNames(constants.ID)
		c.SetParamValues("abc")

		assert.NoError(t, h.UpdateProduct(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader([]byte("{invalid")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames(constants.ID)
		c.SetParamValues("1")

		assert.NoError(t, h.UpdateProduct(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		product := model.Product{ProductName: "Updated"}
		mockSvc.On("UpdateProduct", mock.Anything, int64(1), &product).Return(apperrors.ErrNotFound)

		c, rec := newProductRequest(t, http.MethodPut, "/products/1", product)
		c.SetParamNames(constants.ID)
		c.SetParamValues("1")

		assert.NoError(t, h.UpdateProduct(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestDeleteProductHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		mockSvc.On("DeleteProduct", mock.Anything, int64(1)).Return(nil)

		c, rec := newProductRequest(t, http.MethodDelete, "/products/1", nil)
		c.SetParamNames(constants.ID)
		c.SetParamValues("1")

		assert.NoError(t, h.DeleteProduct(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		c, rec := newProductRequest(t, http.MethodDelete, "/products/abc", nil)
		c.SetParamNames(constants.ID)
		c.SetParamValues("abc")

		assert.NoError(t, h.DeleteProduct(c))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc := mocks.NewProductService(t)
		h := &ProductHandler{Service: mockSvc}

		mockSvc.On("DeleteProduct", mock.Anything, int64(99)).Return(apperrors.ErrNotFound)

		c, rec := newProductRequest(t, http.MethodDelete, "/products/99", nil)
		c.SetParamNames(constants.ID)
		c.SetParamValues("99")

		assert.NoError(t, h.DeleteProduct(c))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
