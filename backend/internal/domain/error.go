package domain

import "myapp/internal/constants"

type DomainError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

func NewValidationError(validationMessage string) *DomainError {
	return &DomainError{Code: constants.ValidationError, Message: validationMessage}
}

func NewDeserializingError(message string) *DomainError {
	return &DomainError{Code: constants.DeserializingError, Message: message}
}
