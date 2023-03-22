package types

import "backend_crudgo/domain/products/domain/model"

//ProductCreateResponse to message for response handler products.
type ProductCreateResponse struct {
	Message string `json:"message,omitempty"`
}

type ProductResponse struct {
	Message string         `json:"message"`
	Product *model.Product `json:"product,omitempty"`
	Error   string         `json:"error,omitempty"`
}
