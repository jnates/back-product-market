package v1

import (
	"backend_crudgo/infrastructure/kit/enum"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/domain/products/domain/service"
	"backend_crudgo/domain/products/infrastructure/persistence"
	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/middlewares"

	"github.com/go-chi/chi"
)

// ProductRouter router
type ProductRouter struct {
	Service service.ProductService
}

// NewProductHandler Should initialize the dependencies for this service.
func NewProductHandler(db *database.DataDB) *ProductRouter {
	return &ProductRouter{
		Service: service.NewProductService(persistence.NewProductRepository(db)),
	}
}

// CreateProductHandler Created initialize handler product.
func (prod *ProductRouter) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	var ctx = r.Context()

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := prod.Service.CreateProduct(ctx, &product)
	if err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusConflict, "Conflict", err.Error())
		return
	}

	w.Header().Add(enum.Location, fmt.Sprintf("%s%s", r.URL.String(), result))
	_ = middlewares.JSON(w, r, http.StatusCreated, result)
}

// GetProductHandler Created initialize get product.
func (prod *ProductRouter) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, enum.ID)

	productResponse, err := prod.Service.GetProduct(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			_ = middlewares.HTTPError(w, r, http.StatusNotFound, "Product not found", err.Error())
			return
		}
		_ = middlewares.HTTPError(w, r, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	writeJSONResponseWithMarshalling(w, http.StatusOK, productResponse)
}

func (prod *ProductRouter) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	productResponse, err := prod.Service.GetProducts(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSONResponseWithMarshalling(w, http.StatusOK, productResponse)
}

// UpdateProductHandler is the HTTP handler for updating a product. It receives an HTTP request with a JSON body containing the updated product information. It verifies the product ID and updates the product information through the product service. If the update is successful, it returns an HTTP response with a status code of 204 (No Content). If there is an error processing the request, it returns an appropriate HTTP error response.
func (prod *ProductRouter) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, enum.ID)

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	response, err := prod.Service.UpdateProduct(ctx, id, &product)
	if err != nil {
		_ = middlewares.HTTPError(w, r, http.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	writeJSONResponseWithMarshalling(w, http.StatusOK, response)
}

// DeleteProductHandler is the HTTP handler for deleting a product. It receives an HTTP request with the ID of the product to delete. It verifies the product ID and deletes the product through the product service. If the delete is successful, it returns an HTTP response with a status code of 204 (No Content). If there is an error processing the request, it returns an appropriate HTTP error response.
func (prod *ProductRouter) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, enum.ID)

	response, err := prod.Service.DeleteProduct(ctx, id)
	if err != nil {
		writeJSONResponseWithMarshalling(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponseWithMarshalling(w, http.StatusOK, response)
}

// writeJSONResponseWithMarshalling is a helper function that writes a JSON response to the HTTP response writer.
// It takes the response writer, the HTTP status code to set in the response, and the data to be written as a JSON payload.
// If there is an error while marshalling the data to JSON, it returns an HTTP error response with a status code of 500 (Internal Server Error).
func writeJSONResponseWithMarshalling(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err = w.Write(jsonBytes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
