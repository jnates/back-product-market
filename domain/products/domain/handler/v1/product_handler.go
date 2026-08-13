// Package v1 exposes the HTTP handlers for the products domain.
package v1

import (
	"net/http"
	"strconv"

	"backend_crudgo/domain/products/constants"
	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/domain/products/domain/service"
	"backend_crudgo/domain/products/infrastructure/persistence"
	"backend_crudgo/infrastructure/kit/apperrors"
	"backend_crudgo/infrastructure/middlewares"

	pgxtool "github.com/jnates/go-toolkit/tools/sqlconnection/pgx"

	"github.com/jnates/go-toolkit/tools/customserver"
	"github.com/labstack/echo/v4"
)

// ProductHandler exposes the HTTP endpoints for products.
type ProductHandler struct {
	Service service.ProductService
}

// NewProductHandler wires the product repository, service and handler.
func NewProductHandler(pool pgxtool.DBPool) *ProductHandler {
	return &ProductHandler{
		Service: service.NewProductService(persistence.NewProductRepository(pool)),
	}
}

// CreateProduct handles POST /products.
//
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @ID CreateProduct
// @Security BearerAuth
// @Param ProductRequest body model.Product true "Product data"
// @Success 201 {object} customserver.GenericResponse{data=model.Product}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products [POST]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	var product model.Product
	if err := c.Bind(&product); err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	created, err := h.Service.CreateProduct(ctx, &product)
	if err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, customserver.GenerateSuccessGenericResponse(created))
}

// GetProduct handles GET /products/:id.
//
// @Description Get a single product by ID
// @Tags Products
// @Produce json
// @ID GetProduct
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} customserver.GenericResponse{data=model.Product}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 404 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products/{id} [GET]
func (h *ProductHandler) GetProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	product, err := h.Service.GetProduct(ctx, id)
	if err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, customserver.GenerateSuccessGenericResponse(product))
}

// GetProducts handles GET /products.
//
// @Description Get every product
// @Tags Products
// @Produce json
// @ID GetProducts
// @Security BearerAuth
// @Success 200 {object} customserver.GenericResponse{data=[]model.Product}
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products [GET]
func (h *ProductHandler) GetProducts(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.Service.GetProducts(ctx)
	if err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, customserver.GenerateSuccessGenericResponse(products))
}

// UpdateProduct handles PUT /products/:id.
//
// @Description Update an existing product
// @Tags Products
// @Accept json
// @Produce json
// @ID UpdateProduct
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param ProductRequest body model.Product true "Product data"
// @Success 204 "Product updated successfully"
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 404 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products/{id} [PUT]
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	var product model.Product
	if err := c.Bind(&product); err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	if err := h.Service.UpdateProduct(ctx, id, &product); err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteProduct handles DELETE /products/:id.
//
// @Description Delete a product
// @Tags Products
// @Produce json
// @ID DeleteProduct
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 204 "Product deleted successfully"
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 404 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products/{id} [DELETE]
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return middlewares.HandleError(c, apperrors.ErrInvalidInput)
	}

	if err := h.Service.DeleteProduct(ctx, id); err != nil {
		return middlewares.HandleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
