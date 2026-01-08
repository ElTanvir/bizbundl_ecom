package handler

import (
	"bizbundl/internal/storefront/catalog/service"
	"bizbundl/util"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type CatalogHandler struct {
	service *service.CatalogService
}

func NewCatalogHandler(service *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

// RegisterRoutes sets up the API routes for Catalog
func (h *CatalogHandler) RegisterRoutes(router fiber.Router) {
	catalogGroup := router.Group("/catalog")
	catalogGroup.Get("/categories", h.ListCategories)
	catalogGroup.Get("/products", h.ListProducts)
	catalogGroup.Get("/products/:id", h.GetProduct)
}

func (h *CatalogHandler) ListCategories(c *fiber.Ctx) error {
	cats, err := h.service.ListCategories(c.Context())
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}
	return util.JSON(c, fiber.StatusOK, cats, "Categories retrieved")
}

func (h *CatalogHandler) ListProducts(c *fiber.Ctx) error {
	products, err := h.service.ListProducts(c.Context())
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}
	return util.JSON(c, fiber.StatusOK, products, "Products retrieved")
}

func (h *CatalogHandler) GetProduct(c *fiber.Ctx) error {
	idHex := c.Params("id")
	var uuid pgtype.UUID
	err := uuid.Scan(idHex)
	if err != nil {
		return util.APIError(c, fiber.StatusBadRequest, fmt.Errorf("invalid product ID"))
	}

	product, err := h.service.GetProduct(c.Context(), uuid)
	if err != nil {
		return util.APIError(c, fiber.StatusNotFound, err)
	}
	return util.JSON(c, fiber.StatusOK, product, "Product retrieved")
}
