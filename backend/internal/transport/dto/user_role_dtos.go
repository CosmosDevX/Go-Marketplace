package dto

type UserRoleDTO struct {
	Username string `json:"username" validate:"required"`
	Role     string `json:"role" validate:"required"`
}
