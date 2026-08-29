package handler

import "myapp/internal/transport/dto"

// MessageResponse is a generic success response with a message.
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

// ErrorResponse is the standard error response returned by the API.
type ErrorResponse struct {
	Message string `json:"message" example:"validation failed"`
}

// AuthResponse is returned on successful login and token refresh.
type AuthResponse struct {
	AccessToken string   `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Roles       []string `json:"roles" example:"user,seller"`
}

// StatusResponse is used by the health endpoint.
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// UserIDResponse is returned after user registration.
type UserIDResponse struct {
	UserID int `json:"user_id" example:"42"`
}

// CategoryIDResponse is returned after creating a category.
type CategoryIDResponse struct {
	CategoryID int `json:"category_id" example:"7"`
}

// ProductIDResponse is returned after creating or updating a product.
type ProductIDResponse struct {
	ProductID int `json:"product_id" example:"15"`
}

// OrderIDResponse is returned after creating an order.
type OrderIDResponse struct {
	OrderID int `json:"order_id" example:"101"`
}

// CartItemIDResponse is returned after adding an item to the cart.
type CartItemIDResponse struct {
	CartItemID int `json:"cart_item_id" example:"33"`
}

// QuantityResponse is returned after updating a cart item quantity.
type QuantityResponse struct {
	Quantity int `json:"quantity" example:"3"`
}

// ProductListResponse is returned by product listing endpoints.
type ProductListResponse struct {
	Products []dto.GetProductDTO `json:"products"`
	Page     int                 `json:"page" example:"1"`
	Limit    int                 `json:"limit" example:"20"`
}
