package dto

type LoginDTO struct {
	Username string `json:"username" validate:"required,min=3,max=40"`
	Password string `json:"password" validate:"required,min=8,max=60"`
}
