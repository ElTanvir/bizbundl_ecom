package handler

import (
	"bizbundl/internal/modules/cart/service"
	"bizbundl/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(service *service.CartService) *CartHandler {
	return &CartHandler{service: service}
}

// RegisterRoutes sets up the API routes for Cart
func (h *CartHandler) RegisterRoutes(router fiber.Router) {
	cartGroup := router.Group("/cart")
	cartGroup.Get("/", h.GetCart)
	cartGroup.Post("/items", h.AddItem)
	// cartGroup.Delete("/items/:id", h.RemoveItem)
}

// Helper to get Session/User from Context (Middleware integration later)
// For MVP, we extract from Cookie or Header manually if middleware not fully set.
func (h *CartHandler) getContextIdentities(c *fiber.Ctx) (pgtype.UUID, pgtype.UUID) {
	// 1. Session Token from Cookie
	// NOTE: In real app, Middleware validates token and sets UserID in locals.
	// For now, let's assume we might have a header or just extraction logic validation.
	// If we use the Auth middleware, Locals("user_id") is set.

	// Simplification for MVP:
	// If Auth Middleware ran, we trust it.
	// If not, we look for a "guest_token" cookie?
	// Let's rely on a "session_token" for both.

	// TODO: We need real Session Middleware to resolve Token -> SessionID/UserID.
	// Making Cart API dependent on Session Middleware is correct.

	// MOCK for scaffolding:
	// We will assume `c.Locals("session_id")` and `c.Locals("user_id")` are set by middleware.
	// If not, we default to invalid.

	/*
	   For this Step, since Middleware isn't fully wired for Guest,
	   we can't fully implement this without that middleware.
	   I will draft the handlers expecting these values.
	*/

	sessID := pgtype.UUID{}
	userID := pgtype.UUID{}

	// Parse from Locals (Assuming middleware will put strings or UUIDs there)
	if v, ok := c.Locals("user_id").(string); ok && v != "" {
		_ = userID.Scan(v)
	}

	// Session ID is internal. The browser sends a Token.
	// The middleware resolves Token -> Session ID.
	if v, ok := c.Locals("session_id").(string); ok && v != "" {
		_ = sessID.Scan(v)
	} else {
		// Fallback: If no session ID in locals, maybe we generate one for a new guest?
		// That logic belongs in middleware.
	}

	return sessID, userID
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	sessID, userID := h.getContextIdentities(c)

	// Temporary Hack until Middleware: Accept headers for testing
	if c.Get("X-Test-Session-ID") != "" {
		_ = sessID.Scan(c.Get("X-Test-Session-ID"))
	}

	if !sessID.Valid && !userID.Valid {
		return util.APIError(c, fiber.StatusUnauthorized, fmt.Errorf("no session or user"))
	}

	cart, err := h.service.GetOrCreateCart(c.Context(), sessID, userID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	items, err := h.service.GetCartItems(c.Context(), cart.ID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	return util.JSON(c, fiber.StatusOK, fiber.Map{
		"cart":  cart,
		"items": items,
	}, "Cart retrieved")
}

type AddItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	var req AddItemRequest
	if err := c.BodyParser(&req); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, err)
	}

	sessID, userID := h.getContextIdentities(c)
	// Hack for test
	if c.Get("X-Test-Session-ID") != "" {
		_ = sessID.Scan(c.Get("X-Test-Session-ID"))
	} else if !userID.Valid {
		// Create new guest session ID logic if missing?
		// For API, we expect client to manage session token.
		// If testing without middleware, generate one.
		newId := uuid.New()
		_ = sessID.Scan(newId.String())
	}

	var pID pgtype.UUID
	if err := pID.Scan(req.ProductID); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, fmt.Errorf("invalid product id"))
	}

	item, err := h.service.AddToCart(c.Context(), sessID, userID, pID, pgtype.UUID{}, int32(req.Quantity))
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	return util.JSON(c, fiber.StatusOK, item, "Item added to cart")
}
