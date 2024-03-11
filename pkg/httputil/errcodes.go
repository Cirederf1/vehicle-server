package httputil

type ErrCode int64

const (
	// 1: catch all errors.
	ErrCodeInternalServerError = iota + 1
	// [2 - 999]: low level technical errors.
	ErrCodeRequestBodyTrailingGarbage
	ErrCodeRequestUnexpectedContentType

	// [1000 - 1999]: application level errors
	ErrCodeInvalidRequestPayload = iota + 1000
	ErrCodeResourceNotFound
)
