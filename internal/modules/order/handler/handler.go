package handler

import (
	"fmt"

	"bizbundl/internal/modules/cart/service"
	orderService "bizbundl/internal/modules/order/service"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	cartSvc  *service.CartService
	orderSvc *orderService.OrderService
}

func NewOrderHandler(cSvc *service.CartService, oSvc *orderService.OrderService) *OrderHandler {
	return &OrderHandler{
		cartSvc:  cSvc,
		orderSvc: oSvc,
	}
}

// Checkout handles the checkout process
func (h *OrderHandler) Checkout(c *fiber.Ctx) error {
	userID := util.GetUserIDFromContext(c)
	sessionID := util.GetSessionIDFromContext(c)

	// 1. Get Active Cart
	// Note: GetCart creates one if not exists, but we want to check if it has items.
	cart, err := h.cartSvc.GetCart(c.Context(), userID, sessionID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get cart")
	}

	// Check items? Service CreateOrder checks items.

	// 2. Create Order
	// UserID logic: If Guest, we might have a placeholder UserID or allow NULL in DB (schema says NOT NULL).
	// PRD says "Guest Checkout".
	// Schema `users` table: id is PK.
	// If guest, do we create a shadow user? Or does CartService handle "Guest User"?
	// Currently `CreateOrderFromCart` expects `pgtype.UUID` for UserID.
	// WORKAROUND: For MVP, if userID is nil, we assume a "Guest User" exists or we fail.
	// But `GetUserIDFromContext` returns UUID. If invalid/empty, checking.

	if !userID.Valid {
		// Guest checkout logic needed.
		// For now, let's assume strict auth or fail, OR create a transient user.
		// Actually, schema `sessions` links to `user_id`.
		// Let's rely on what we have. If invalid, we might block checkout for now (MVP 1.0)
		// OR better: The user must be logged in OR we have a guest logic.
		// Let's check `cart.UserID`.
		if cart.UserID.Valid {
			userID = cart.UserID
		} else {
			// It's a true guest. CreateOrder needs a UserID.
			// We'll TODO this: "Guest Checkout User Creation".
			// For now, return Error "Please Login".
			return fiber.NewError(fiber.StatusUnauthorized, "Please login to checkout")
		}
	}

	order, err := h.orderSvc.CreateOrderFromCart(c.Context(), userID, cart.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Checkout failed: %v", err))
	}

	return c.Redirect(fmt.Sprintf("/order/success/%x", order.ID.Bytes))
}

func (h *OrderHandler) SuccessPage(c *fiber.Ctx) error {
	idHex := c.Params("id")
	var idBytes [16]byte
	// Convert hex to bytes... util helper?
	// Or just UUID parsing.
	// Assuming util has UUID parsing from string.
	// If not, use google/uuid then to pgtype.

	// Mocking render for now, actually need a view.
	return c.SendString(fmt.Sprintf("Order Success! ID: %s", idHex))
}
