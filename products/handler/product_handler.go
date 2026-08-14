// Package handler exposes the HTTP handlers for the products domain.
package handler

import (
	"net/http"
	"strconv"

	"backend_crudgo/pkg/kit/customErrors"
	"backend_crudgo/pkg/kit/validation"
	"backend_crudgo/products/constants"
	"backend_crudgo/products/interfaces"
	"backend_crudgo/products/models"

	"github.com/jnates/go-toolkit/tools/customserver"
	"github.com/jnates/go-toolkit/tools/jwttools"
	"github.com/labstack/echo/v4"
)

// productHandler implements interfaces.ProductHandler, delegating to the product service.
type productHandler struct {
	service interfaces.ProductService
}

// NewProductHandler builds a ProductHandler backed by the given service.
func NewProductHandler(service interfaces.ProductService) interfaces.ProductHandler {
	return &productHandler{service: service}
}

// CreateProductRequest is the request body accepted by POST /products.
// It deliberately excludes the audit fields (product_user_created/modify):
// those are derived from the authenticated JWT, never trusted from the client.
type CreateProductRequest struct {
	ProductName   string  `json:"product_name" validate:"required,min=2,max=200"`
	ProductAmount int     `json:"product_amount" validate:"gte=0"`
	ProductPrice  float64 `json:"product_price" validate:"gte=0"`
}

// UpdateProductRequest is the request body accepted by PUT /products/:id.
// Mirrors the fields ProductRepository.UpdateProduct actually persists.
type UpdateProductRequest struct {
	ProductName   string `json:"product_name" validate:"required,min=2,max=200"`
	ProductAmount int    `json:"product_amount" validate:"gte=0"`
}

// CreateProduct handles POST /products.
//
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
//
// @Description Create a new product
// @Tags Products
// @Accept json
// @Produce json
// @ID CreateProduct
// @Security BearerAuth
// @Param ProductRequest body CreateProductRequest true "Product data"
// @Success 201 {object} customserver.GenericResponse{data=models.Product}
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 422 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products [POST]
func (h *productHandler) CreateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customserver.GenerateErrorGenericResponse(
			http.StatusUnprocessableEntity, "validation failed", validation.FormatErrors(err)))
	}

	userID, err := jwttools.UserIDFromContext(ctx)
	if err != nil {
		return customErrors.HandleError(c, customErrors.ErrUnauthorized)
	}

	product := models.Product{
		ProductName:        req.ProductName,
		ProductAmount:      req.ProductAmount,
		ProductPrice:       req.ProductPrice,
		ProductUserCreated: int(userID),
	}

	created, err := h.service.CreateProduct(ctx, &product)
	if err != nil {
		return customErrors.HandleError(c, err)
	}

	return c.JSON(http.StatusCreated, customserver.GenerateSuccessGenericResponse(created))
}

// GetProduct handles GET /products/:id.
//
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
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
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
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
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
//
// @Description Update an existing product
// @Tags Products
// @Accept json
// @Produce json
// @ID UpdateProduct
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param ProductRequest body UpdateProductRequest true "Product data"
// @Success 204 "Product updated successfully"
// @Failure 400 {object} customserver.GenericResponse
// @Failure 401 {object} customserver.GenericResponse
// @Failure 404 {object} customserver.GenericResponse
// @Failure 422 {object} customserver.GenericResponse
// @Failure 500 {object} customserver.GenericResponse
// @Router /products/{id} [PUT]
func (h *productHandler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param(constants.ID), 10, 64)
	if err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	var req UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return customErrors.HandleError(c, customErrors.ErrInvalidInput)
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, customserver.GenerateErrorGenericResponse(
			http.StatusUnprocessableEntity, "validation failed", validation.FormatErrors(err)))
	}

	userID, err := jwttools.UserIDFromContext(ctx)
	if err != nil {
		return customErrors.HandleError(c, customErrors.ErrUnauthorized)
	}

	product := models.Product{
		ProductName:       req.ProductName,
		ProductAmount:     req.ProductAmount,
		ProductUserModify: int(userID),
	}

	if err := h.service.UpdateProduct(ctx, id, &product); err != nil {
		return customErrors.HandleError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteProduct handles DELETE /products/:id.
//
// Parameters:
//   - c: echo context containing the HTTP request and response
//
// Returns:
//   - error: HTTP error or nil on success
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
