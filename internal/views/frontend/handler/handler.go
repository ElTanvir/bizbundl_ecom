package handler

import (
	"bizbundl/internal/modules/catalog/service"
	"bizbundl/internal/views/frontend/pages"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
)

type FrontendHandler struct {
	catalogService *service.CatalogService
}

func NewFrontendHandler(catalogService *service.CatalogService) *FrontendHandler {
	return &FrontendHandler{catalogService: catalogService}
}

func (h *FrontendHandler) HomePage(c *fiber.Ctx) error {
	// 1. Fetch products (Limit 20 for MVP)
	// TODO: Create ListProductsParams if needed, for now assuming ListProducts returns all or top.
	products, err := h.catalogService.ListProducts(c.Context())
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	// 2. Render Page
	// We need an adapter to converting Templ component to Fiber Handler, or use a util.
	// We'll use a local helper or `templ.Handler` logic adapted for Fiber.
	// Fiber's adaptor: `adaptor.HTTPHandler(templ.Handler(component))` works but ignores Fiber context somewhat.
	// Better: Set Content-Type and write body.

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return pages.Home(products).Render(c.Context(), c.Response().BodyWriter())
}
