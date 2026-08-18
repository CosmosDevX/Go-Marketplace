package utils

import (
	"errors"
	"myapp/internal/domain"
	"net/http"
)

func MapError(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrUniqueViolation):
		return http.StatusConflict
	case errors.Is(err, domain.ErrTimeout):
		return http.StatusGatewayTimeout
	case errors.Is(err, domain.ErrValidation), errors.Is(err, domain.ErrParse):
		return http.StatusBadRequest
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, domain.ErrTooManyRequests):
		return http.StatusTooManyRequests
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func PublicMessage(err error) string {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return "not found"
	case errors.Is(err, domain.ErrUniqueViolation):
		return "already exists"
	case errors.Is(err, domain.ErrTimeout):
		return "request timeout"
	case errors.Is(err, domain.ErrValidation):
		return "validation failed"
	case errors.Is(err, domain.ErrParse):
		return "invalid request"
	case errors.Is(err, domain.ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, domain.ErrTooManyRequests):
		return "too many requests"
	case errors.Is(err, domain.ErrForbidden):
		return "forbidden"
	default:
		return "internal error"
	}
}
