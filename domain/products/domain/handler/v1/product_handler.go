package v1

import (
	"backend_crudgo/infrastructure/kit/enum"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"backend_crudgo/domain/products/domain/model"
	"backend_crudgo/domain/products/domain/service"
	"backend_crudgo/domain/products/infrastructure/persistence"
	"backend_crudgo/infrastructure/database"
	"backend_crudgo/infrastructure/middleware"
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
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, "Bad request", err.Error())
		return
	}

	result, err := prod.Service.CreateProductHandler(ctx, &product)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, "Conflict", err.Error())
		return
	}

	w.Header().Add(enum.Location, fmt.Sprintf("%s%s", r.URL.String(), result))
	_ = middleware.JSON(w, r, http.StatusCreated, result)
}

// GetProductHandler Created initialize get product.
func (prod *ProductRouter) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var id = chi.URLParam(r, enum.Id)

	productResponse, err := prod.Service.GetProductHandler(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(productResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (prod *ProductRouter) GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	productResponse, err := prod.Service.GetProductsHandler(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonBytes, err := json.Marshal(productResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonBytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
