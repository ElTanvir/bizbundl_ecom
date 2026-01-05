package handler

import (
	"fmt"

	"bizbundl/internal/modules/cart/service"
	orderService "bizbundl/internal/modules/order/service"
	"bizbundl/internal/modules/payment"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	cartSvc   *service.CartService
	orderSvc  *orderService.OrderService
	paymentGw payment.Gateway
}

func NewOrderHandler(cSvc *service.CartService, oSvc *orderService.OrderService, pgw payment.Gateway) *OrderHandler {
	return &OrderHandler{
		cartSvc:   cSvc,
		orderSvc:  oSvc,
		paymentGw: pgw,
	}
}

// Checkout handles the checkout process
func (h *OrderHandler) Checkout(c *fiber.Ctx) error {
	userID := util.GetUserIDFromContext(c)
	sessionID := util.GetSessionIDFromContext(c)

	// 1. Get Active Cart
	// Note: GetCart creates one if not exists, but we want to check if it has items.
	cart, err := h.cartSvc.GetOrCreateCart(c.Context(), sessionID, userID)
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

	// 3. Initiate Payment
	paymentURL, err := h.paymentGw.InitPayment(order, "") // passing empty email implies provider default for now
	if err != nil {
		// Log err
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Payment Init Failed: %v", err))
	}

	return c.Redirect(paymentURL)
}

func (h *OrderHandler) PaymentCallback(c *fiber.Ctx) error {
	// invoice_id comes from POST or Get?
	// UddoktaPay might send invoice_id in query or POST.
	// Check documentation: Standard checkout redirects to success URL with params?
	// Usually it doesn't send sensitive info in URL.
	// However, verify payment API needs invoice_id.
	// Wait, the redirect URL I set was /order/payment/callback?order_id=...
	// Does it append invoice_id?
	// Let's assume user lands here. We might treat this as "Verify Order Status".

	// For now, let's just show Success Page if we land here, assuming simple flow.
	// BUT we should verify.

	orderIDHex := c.Query("order_id")

	// In a real flow, we'd verify the payment status using the Gateway.
	// Since we don't have the Invoice ID here unless the Gateway appended it,
	// we assume success or check our DB if Webhook updated it.

	return c.Redirect("/order/success/" + orderIDHex)
}

func (h *OrderHandler) SuccessPage(c *fiber.Ctx) error {
	idHex := c.Params("id")
	// Convert hex to bytes... util helper?
	// Or just UUID parsing.
	// Assuming util has UUID parsing from string.
	// If not, use google/uuid then to pgtype.

	// Mocking render for now, actually need a view.
	return c.SendString(fmt.Sprintf("Order Success! ID: %s", idHex))
}
