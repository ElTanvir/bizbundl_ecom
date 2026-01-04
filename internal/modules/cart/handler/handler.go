package handler

import (
	"bizbundl/internal/modules/cart/service"
	"bizbundl/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
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
	sessID := pgtype.UUID{}
	userID := pgtype.UUID{}

	// Middleware GUARANTEES "user_id" and "user_role" are set now (Guest or User).
	// Extract from Locals.
	idStr, ok := c.Locals("user_id").(string)
	if !ok || idStr == "" {
		// Should not happen if middleware is active
		return sessID, userID
	}

	role, _ := c.Locals("user_role").(string)

	isGuest := role == "guest"
	if isGuest {
		_ = sessID.Scan(idStr)
	} else {
		_ = userID.Scan(idStr)
	}

	return sessID, userID
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	sessID, userID := h.getContextIdentities(c)

	// Middleware handles identity extraction

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
