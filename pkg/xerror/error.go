package xerror

import (
	"net/http"

	"github.com/pkg/errors"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}
func (e *AppError) Cause() error {
	return e.Reason
}

// New error
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Reason:  errors.New(message),
	}
}

func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Reason:  errors.WithStack(err),
	}
}

// 400
func BadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

// 401
func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}

// 403
func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, message)
}

// 404
func NotFound(message string) *AppError {
	return New(http.StatusNotFound, message)
}

// 500
func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}
