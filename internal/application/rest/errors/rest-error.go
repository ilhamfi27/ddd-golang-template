package rest_error

import "net/http"

// Base custom error type that implements the error interface
type HttpError struct {
	Message string
	Code    int
}

func (r *HttpError) Error() string {
	return r.Message
}

func BadRequestError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func UnauthorizedError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

func ForbiddenError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusForbidden,
	}
}

func NotFoundError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func DuplicateError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusConflict,
	}
}

func InternalServerError(message string) *HttpError {
	return &HttpError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}
