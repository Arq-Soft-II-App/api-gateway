package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	HTTPStatusCode int    `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Constructor genérico
func NewError(code, message string, httpStatusCode int) *Error {
	return &Error{
		Code:           code,
		Message:        message,
		HTTPStatusCode: httpStatusCode,
	}
}

// Errores específicos
func NewBadRequestError(message string) *Error {
	return &Error{
		Code:           "BAD_REQUEST",
		Message:        message,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *Error {
	return &Error{
		Code:           "UNAUTHORIZED",
		Message:        message,
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

func NewNotFoundError(message string) *Error {
	return &Error{
		Code:           "NOT_FOUND",
		Message:        message,
		HTTPStatusCode: http.StatusNotFound,
	}
}

func NewInternalServerError(message string) *Error {
	return &Error{
		Code:           "INTERNAL_SERVER_ERROR",
		Message:        message,
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

// Función helper para obtener el código HTTP
func GetStatusCode(err error) int {
	if customErr, ok := err.(*Error); ok {
		return customErr.HTTPStatusCode
	}
	return http.StatusInternalServerError
}
