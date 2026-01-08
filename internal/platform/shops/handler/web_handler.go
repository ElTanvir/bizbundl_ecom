package handler

import (
	"bizbundl/internal/platform/shops/service"
	"bizbundl/internal/views/platform"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type PlatformWebHandler struct {
	service *service.PlatformService
}

func NewPlatformWebHandler(s *service.PlatformService) *PlatformWebHandler {
	return &PlatformWebHandler{service: s}
}

// ShowDashboard renders the main list of shops
func (h *PlatformWebHandler) ShowDashboard(c *fiber.Ctx) error {
	// 1. Get User ID from Session/Context (Middleware must set this)
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		// Redirect to Login if not found
		return c.Redirect("/login")
	}

	var ownerID pgtype.UUID
	ownerID.Scan(userIDStr)

	// 2. Call Service (Pure Logic)
	shops, err := h.service.ListShops(c.Context(), ownerID)
	if err != nil {
		return c.Status(500).SendString("Failed to load shops")
	}

	// 3. Render Template (Web Layer)
	return util.Render(c, platform.Dashboard(shops))
}

func (h *PlatformWebHandler) ShowCreateShopForm(c *fiber.Ctx) error {
	return util.Render(c, platform.CreateShopForm())
}

func (h *PlatformWebHandler) HandleCreateShop(c *fiber.Ctx) error {
	// 1. Parse Form
	name := c.FormValue("name")

	userIDStr, _ := c.Locals("user_id").(string)
	var ownerID pgtype.UUID
	ownerID.Scan(userIDStr)

	// 2. Call Service
	_, err := h.service.CreateShop(c.Context(), ownerID, name)
	if err != nil {
		// Return Form with Error
		return c.SendString("Error: " + err.Error())
	}

	// 3. Redirect
	return c.Redirect("/dashboard")
}
