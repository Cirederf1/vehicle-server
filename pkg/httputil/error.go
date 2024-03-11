package httputil

import (
	"errors"
	"fmt"
	"net/http"
)

type APIError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
	Details any     `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%d] %s (details: %+v)", e.Code, e.Message, e.Details)
}

func ServeError(rw http.ResponseWriter, statusCode int, err error) {
	if err == nil {
		return
	}

	// This is consumed by errors.As, so we need to reference a pointer here.
	apiError := &APIError{}

	if !errors.As(err, &apiError) {
		ServeJSON(rw, statusCode, &APIError{Code: ErrCodeInternalServerError, Message: "Unexpected error"})
		return
	}

	ServeJSON(rw, statusCode, apiError)
}
