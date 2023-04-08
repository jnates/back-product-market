package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type partialWriteErrorWriter struct {
	http.ResponseWriter
	n int
}

type errorWriter struct {
	http.ResponseWriter
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		data        interface{}
		expectedErr error
	}{
		{
			name:        "Should return no error when data is nil",
			statusCode:  http.StatusOK,
			data:        nil,
			expectedErr: nil,
		},
		{
			name:       "Should return no error when valid data is provided",
			statusCode: http.StatusOK,
			data: SuccessfullyMessage{
				Status:  http.StatusOK,
				Message: "Success",
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertion := assert.New(t)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			err := JSON(w, r, test.statusCode, test.data)

			if test.expectedErr != nil {
				assertion.IsType(test.expectedErr, err)
			} else {
				assertion.NoError(err)
			}
		})
	}
}

func TestJSONError(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		data        interface{}
		expectedErr error
		writer      http.ResponseWriter
	}{
		{
			name:        "Should return UnsupportedTypeError when invalid data is provided",
			statusCode:  http.StatusOK,
			data:        make(chan int),
			expectedErr: &json.UnsupportedTypeError{},
			writer:      httptest.NewRecorder(),
		},
		{
			name:        "Should return error when writing response fails",
			statusCode:  http.StatusOK,
			data:        "This is a test message",
			expectedErr: fmt.Errorf(" Error writing response body: wrote 0 bytes out of 21"),
			writer:      &errorWriter{ResponseWriter: httptest.NewRecorder()},
		},
		{
			name:        "Should return error when partial write occurs",
			statusCode:  http.StatusOK,
			data:        "This is a test message",
			expectedErr: fmt.Errorf(" Error writing response body: wrote 10 bytes out of 21"),
			writer:      &partialWriteErrorWriter{ResponseWriter: httptest.NewRecorder(), n: 10},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)

			err := JSON(test.writer, r, test.statusCode, test.data)

			if test.expectedErr != nil {
				assert.IsType(t, test.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (p *partialWriteErrorWriter) Write([]byte) (int, error) {
	return p.n, nil
}

func (e *errorWriter) Write([]byte) (int, error) {
	return 0, errors.New(" Error writing response ")
}

func TestHTTPError(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		error       string
		message     string
		expectedErr error
	}{
		{
			name:        "Should return error message with given status code and custom message",
			statusCode:  http.StatusBadRequest,
			error:       "Bad Request",
			message:     "Invalid data",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertion := assert.New(t)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)

			err := HTTPError(w, r, test.statusCode, test.error, test.message)

			if test.expectedErr != nil {
				assertion.Equal(test.expectedErr.Error(), err.Error())
			} else {
				assertion.NoError(err)
				var res ErrorMessage
				json.Unmarshal(w.Body.Bytes(), &res)
				assertion.Equal(test.statusCode, res.Status)
				assertion.Equal(test.error, res.Error)
				assertion.Equal(test.message, res.Message)
			}
		})
	}
}
