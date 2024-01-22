package tool

import (
	"encoding/json"
	"net/http"
)

// WriteJSONResponseWithMarshalling is a helper function that writes a JSON response to the HTTP response writer.
// It takes the response writer, the HTTP status code to set in the response, and the data to be written as a JSON payload.
// If there is an error while marshall the data to JSON, it returns an HTTP error response with a status code.
// of 500 (Internal Server Error).
func WriteJSONResponseWithMarshalling(w http.ResponseWriter, statusCode int, data interface{}) {
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

