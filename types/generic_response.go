package types

//CreateResponse to message for response handler products.
type CreateResponse struct {
	Message string `json:"message,omitempty"`
}

// GenericResponse is a generic response type that can be used to return a message, data, and error in JSON format.
type GenericResponse struct {
	Message string      `json:"message"`
	Product interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GenericUserResponse is a generic response type that can be used to return a message, user data, and error in JSON format.
type GenericUserResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
