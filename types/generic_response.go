package types

// ProductCreateResponse to message for response handler products.
type ProductCreateResponse struct {
	Message string `json:"message,omitempty"`
}

// GenericResponse to message for response products.
type GenericResponse struct {
	Message string      `json:"message"`
	Product interface{} `json:"product,omitempty"`
	Error   string      `json:"error,omitempty"`
}
