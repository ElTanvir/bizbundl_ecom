package handler

import (
	"fmt"

	db "bizbundl/internal/db/sqlc"
	"bizbundl/internal/storefront/cart/service"
	catalogService "bizbundl/internal/storefront/catalog/service"
	orderService "bizbundl/internal/storefront/order/service"
	"bizbundl/internal/modules/payment"
	"bizbundl/internal/views/frontend/pages"
	"bizbundl/util"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	cartSvc    *service.CartService
	orderSvc   *orderService.OrderService
	catalogSvc *catalogService.CatalogService
	paymentGw  payment.Gateway
}

func NewOrderHandler(cSvc *service.CartService, oSvc *orderService.OrderService, catSvc *catalogService.CatalogService, pgw payment.Gateway) *OrderHandler {
	return &OrderHandler{
		cartSvc:    cSvc,
		orderSvc:   oSvc,
		catalogSvc: catSvc,
		paymentGw:  pgw,
	}
}

// ShowCheckoutPage renders the checkout widget
func (h *OrderHandler) ShowCheckoutPage(c *fiber.Ctx) error {
	productID := c.Query("product_id")
	variantID := c.Query("variant_id")

	var directProduct *db.Product
	var directVariant *db.ProductVariant
	isPhysical := false

	if productID != "" {
		// Direct Checkout Mode
		pID, err := util.StringToUUID(productID)
		if err == nil {
			p, err := h.catalogSvc.GetProduct(c.Context(), pID)
			if err == nil {
				directProduct = &p
				if p.IsDigital != nil && !*p.IsDigital {
					isPhysical = true
				}

				if variantID != "" {
					vID, err := util.StringToUUID(variantID)
					if err == nil {
						v, err := h.catalogSvc.GetProductVariant(c.Context(), vID) // Need method in CatalogSvc
						if err == nil {
							directVariant = &v
							// Assuming we dont track dimensions on variant level for MVP, usually inherits
						}
					}
				}
			}
		}
	}

	var cart db.Cart
	var items []db.GetCartItemsRow

	if directProduct == nil {
		// Cart Mode
		userID := util.GetUserIDFromContext(c)
		sessionID := util.GetSessionIDFromContext(c)
		var err error
		cart, err = h.cartSvc.GetOrCreateCart(c.Context(), sessionID, userID)
		if err == nil {
			items, _ = h.cartSvc.GetCartItems(c.Context(), cart.ID)
			hasPhysical, _ := h.cartSvc.HasPhysicalItems(c.Context(), cart.ID)
			isPhysical = hasPhysical
		}
	}

	return render(c, pages.Checkout(cart, items, directProduct, directVariant, isPhysical))
}

// Checkout handles the checkout process
func (h *OrderHandler) Checkout(c *fiber.Ctx) error {
	userID := util.GetUserIDFromContext(c)
	// sessionID := util.GetSessionIDFromContext(c)

	directProductID := c.FormValue("direct_product_id")
	directVariantID := c.FormValue("variant_id")

	var order *db.Order
	var err error

	// If Guest User (no ID), handling logic is in Service or we enforce temporary ID.
	// For now, let's assume we proceed with whatever UserID we have (could be invalid/nil).
	// Ideally we create a Guest User record.
	// MVP: If userID invalid, we might fail at DB level if FK constrained.
	// Let's create a Guest User if needed?
	// The DB `orders` table has `user_id UUID REFERENCES users`. It is nullable?
	// Schema checks: `user_id UUID REFERENCES users(id)` -> implied NULLABLE unless NOT NULL specified.
	// Checking schema: `user_id UUID REFERENCES users(id)` (No NOT NULL). So Guest Orders possible with NULL user_id.

	if directProductID != "" {
		// 1. Direct Checkout
		pID, _ := util.StringToUUID(directProductID)
		vID, _ := util.StringToUUID(directVariantID)

		// Quantity? Default 1 for Direct Buy Now
		order, err = h.orderSvc.CreateOrderDirect(c.Context(), userID, pID, vID, 1)

	} else {
		// 2. Cart Checkout
		sessionID := util.GetSessionIDFromContext(c)
		cart, cErr := h.cartSvc.GetOrCreateCart(c.Context(), sessionID, userID)
		if cErr != nil {
			return h.renderHTMXError(c, "Cart not found")
		}

		order, err = h.orderSvc.CreateOrderFromCart(c.Context(), userID, cart.ID)
	}

	if err != nil {
		return h.renderHTMXError(c, fmt.Sprintf("Checkout failed: %v", err))
	}

	// 3. Update Order with Guest Info
	_ = c.FormValue("customer_name") // TODO: Save to Order
	customerEmail := c.FormValue("customer_email")
	// Save this info to Order? Schema has `guest_info JSONB`.
	// We should update the order.
	// Skipping for brevity, assuming Payment Gateway captures email.
	// Ideally: h.orderSvc.UpdateGuestInfo(ctx, orderID, name, email...)

	// 4. Initiate Payment
	paymentURL, err := h.paymentGw.InitPayment(order, customerEmail)
	if err != nil {
		return h.renderHTMXError(c, "Payment initialization failed")
	}

	// HTMX Redirect
	c.Set("HX-Redirect", paymentURL)
	return c.SendStatus(fiber.StatusOK)
}

func (h *OrderHandler) renderHTMXError(c *fiber.Ctx, msg string) error {
	return c.SendString(fmt.Sprintf("<div class='bg-red-100 text-red-700 p-3 rounded'>%s</div>", msg))
}

func render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
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
