package apperrors

import (
	"errors"
	"net/http"
)

type appError struct {
	code       code
	httpStatus int
}

var (
	InvalidParameter    = &appError{code: InvalidParameterCode, httpStatus: http.StatusBadRequest}
	InternalServerError = &appError{code: InternalServerErrorCode, httpStatus: http.StatusInternalServerError}
	TodoNotFound        = &appError{code: TodoNotFoundCode, httpStatus: http.StatusNotFound}
)

func (e *appError) Error() string {
	return e.code.value()
}

func (e *appError) StatusCode() int {
	return e.httpStatus
}

func AsAppError(err error) *appError {
	var e *appError
	if errors.As(err, &e) {
		return e
	}

	// アサーションに失敗した場合InternalServerError
	return &appError{
		code:       InternalServerErrorCode,
		httpStatus: http.StatusInternalServerError,
	}
}
