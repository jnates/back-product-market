// Package handler exposes the HTTP handlers for the products domain.
package handler

import (
	"net/http"
	"strconv"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/products/constants"
	"backend_crudgo/products/interfaces"
	"backend_crudgo/products/models"

	"github.com/jnates/go-toolkit/tools/customserver"
	"github.com/labstack/echo/v4"
)

type productHandler struct {
	service interfaces.ProductService
}

// NewProductHandler builds a ProductHandler backed by the given service.
func NewProductHandler(service interfaces.ProductService) interfaces.ProductHandler {
	return &productHandler{service: service}
}

// CreateProduct handles POST /products.
//
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @ID CreateProduct
// @Security BearerAuth
// @Param ProductRequest body models.Product true "Product data"
// @Success 201 {object} customserver.GenericResponse{data=models.Product}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products [POST]
func (h *productHandler) CreateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	created, err := h.service.CreateProduct(ctx, &product)
	if err != nil {
		return customErrors.HandleError(c, err)
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
// @Success 200 {object} customserver.GenericResponse{data=models.Product}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 404 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products/{id} [GET]
func (h *productHandler) GetProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	product, err := h.service.GetProduct(ctx, id)
	if err != nil {
		return customErrors.HandleError(c, err)
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
// @Success 200 {object} customserver.GenericResponse{data=[]models.Product}
// @Failure 401 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products [GET]
func (h *productHandler) GetProducts(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.GetProducts(ctx)
	if err != nil {
		return customErrors.HandleError(c, err)
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
// @Param ProductRequest body models.Product true "Product data"
// @Success 204 "Product updated successfully"
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 404 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products/{id} [PUT]
func (h *productHandler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	if err := h.service.UpdateProduct(ctx, id, &product); err != nil {
		return customErrors.HandleError(c, err)
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
func (h *productHandler) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	if err := h.service.DeleteProduct(ctx, id); err != nil {
		return customErrors.HandleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
