package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/service"
	"myapp/internal/transport/dto"
	"myapp/internal/transport/middleware"
	"myapp/internal/transport/validator"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type CartItemService interface {
	Create(ctx context.Context, input service.CreateCartItemInput) (int, error)
	List(ctx context.Context, userID int) ([]domain.CartItem, error)
	UpdateQuantity(ctx context.Context, userID, cartItemID, delta int) (int, error)
	Delete(ctx context.Context, userID, cartItemID int) error
}

type CartItemHandler struct {
	cartItemService CartItemService
}

func NewCartItemHandler(cartItemService CartItemService) CartItemHandler {
	return CartItemHandler{
		cartItemService: cartItemService,
	}
}

// Create godoc
//
//	@Summary		Add item to cart
//	@Description	Add a product to the authenticated user's cart.
//	@Tags			Cart
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		dto.CreateCartItemDTO	true	"Product to add"
//	@Success		200		{object}	CartItemIDResponse		"Created cart item ID"
//	@Failure		400		{object}	ErrorResponse			"Parse or validation error"
//	@Failure		401		{object}	ErrorResponse			"Unauthorized"
//	@Failure		404		{object}	ErrorResponse			"Product not found"
//	@Failure		500		{object}	ErrorResponse			"Internal server error"
//	@Router			/cart/items [post]
func (h CartItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var dto dto.CreateCartItemDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("create cart item dto: %w", domain.ErrParse))
		return
	}

	if err := validator.Struct(dto); err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("create cart item dto: %w", domain.ErrValidation))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	cartItemID, err := h.cartItemService.Create(ctx, service.CreateCartItemInput{
		ProductID: dto.ProductID,
		UserID:    user.UserID,
	})

	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"cart_item_id": cartItemID})
}

// List godoc
//
//	@Summary		List cart items
//	@Description	Get all items in the authenticated user's cart.
//	@Tags			Cart
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		dto.GetCartItemDTO	"List of cart items"
//	@Failure		401	{object}	ErrorResponse		"Unauthorized"
//	@Failure		500	{object}	ErrorResponse		"Internal server error"
//	@Router			/cart [get]
func (h CartItemHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	cartItems, err := h.cartItemService.List(ctx, user.UserID)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	dtos := make([]dto.GetCartItemDTO, len(cartItems))
	for i := range dtos {
		dtos[i] = dto.ToGetCartItemDTO(cartItems[i])
	}

	utils.WriteJSON(w, dtos)
}

// UpdateQuantity godoc
//
//	@Summary		Update cart item quantity
//	@Description	Change quantity of a cart item by delta (can be negative).
//	@Tags			Cart
//	@Produce		json
//	@Security		BearerAuth
//	@Param			cart_item_id	path		int	true	"Cart item ID"
//	@Param			delta			query		int	true	"Quantity change (positive or negative)"
//	@Success		200				{object}	QuantityResponse	"Updated quantity"
//	@Failure		400				{object}	ErrorResponse		"Invalid parameters"
//	@Failure		401				{object}	ErrorResponse		"Unauthorized"
//	@Failure		404				{object}	ErrorResponse		"Cart item not found"
//	@Failure		500				{object}	ErrorResponse		"Internal server error"
//	@Router			/cart/items/{cart_item_id} [patch]
func (h CartItemHandler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cartItemID, err := strconv.Atoi(r.PathValue("cart_item_id"))
	if err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("cart item id parse: %w", domain.ErrParse))
		return
	}
	delta, err := strconv.Atoi(r.URL.Query().Get("delta"))
	if err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("delta parse: %w", domain.ErrParse))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	quantity, err := h.cartItemService.UpdateQuantity(ctx, user.UserID, cartItemID, delta)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]int{"quantity": quantity})
}

// Delete godoc
//
//	@Summary		Remove item from cart
//	@Description	Delete a cart item by ID.
//	@Tags			Cart
//	@Produce		json
//	@Security		BearerAuth
//	@Param			cart_item_id	path		int	true	"Cart item ID"
//	@Success		200				{object}	MessageResponse	"Deletion confirmation"
//	@Failure		400				{object}	ErrorResponse		"Invalid cart_item_id"
//	@Failure		401				{object}	ErrorResponse		"Unauthorized"
//	@Failure		404				{object}	ErrorResponse		"Cart item not found"
//	@Failure		500				{object}	ErrorResponse		"Internal server error"
//	@Router			/cart/items/{cart_item_id} [delete]
func (h CartItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cartItemID, err := strconv.Atoi(r.PathValue("cart_item_id"))
	if err != nil {
		utils.WriteError(ctx, w, fmt.Errorf("cart item id parse: %w", domain.ErrParse))
		return
	}

	user, err := middleware.UserFromContext(ctx)
	if err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	if err := h.cartItemService.Delete(ctx, user.UserID, cartItemID); err != nil {
		utils.WriteError(ctx, w, err)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "delete successful"})
}
