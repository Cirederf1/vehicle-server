package httputil

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

var errTrailingGarbage = &APIError{
	Code:    ErrCodeRequestBodyTrailingGarbage,
	Message: "Unexpected garbage at the end on the request body",
}

func DecodeRequestAsJSON(r *http.Request, v any) error {
	if ct := r.Header.Get("Content-Type"); !strings.EqualFold(ct, "application/json") {
		return unexpectedRequestContentTypeError(ct)
	}

	return DecodeJSON(r.Body, v)
}

func DecodeJSON(body io.ReadCloser, v any) error {
	defer body.Close()

	dec := json.NewDecoder(body)
	if err := dec.Decode(v); err != nil {
		return err
	}

	if dec.More() {
		return errTrailingGarbage
	}

	return nil
}

func ServeJSON(rw http.ResponseWriter, statusCode int, payload any) {
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(statusCode)
	_ = json.NewEncoder(rw).Encode(payload)
}

func unexpectedRequestContentTypeError(got string) error {
	return &APIError{
		Code:    ErrCodeRequestUnexpectedContentType,
		Message: "Unexpected request content type",
		Details: map[string]string{
			"expected": "application/json",
			"got":      got,
		},
	}
}
