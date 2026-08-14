// Package dto
package dto

type CreateUserDTO struct {
	Username string `json:"username" validate:"required,min=3,max=40"`
	Password string `json:"password" validate:"required,min=8,max=60"`
	Email    string `json:"email" validate:"required,email"`
}
