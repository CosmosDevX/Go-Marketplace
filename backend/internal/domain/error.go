package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrUniqueViolation = errors.New("unique violation")
	ErrTimeout         = errors.New("timeout")
	ErrParse           = errors.New("parse error")
	ErrValidation      = errors.New("validation error")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrTooManyRequests = errors.New("too many requests")
)
