package handler

import (
	db "bizbundl/internal/db/sqlc"
	cartservice "bizbundl/internal/modules/cart/service"
	"bizbundl/internal/modules/catalog/service"
	pb_resolver "bizbundl/internal/modules/page_builder/resolver"
	pb "bizbundl/internal/modules/page_builder/service"
	"bizbundl/internal/views/frontend/pages"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
)

type FrontendHandler struct {
	catalogService *service.CatalogService
	cartService    *cartservice.CartService
	pbService      *pb.PageBuilderService
	pbResolver     *pb_resolver.PageResolver
}

func NewFrontendHandler(catalogService *service.CatalogService, cartService *cartservice.CartService, pbService *pb.PageBuilderService, pbResolver *pb_resolver.PageResolver) *FrontendHandler {
	return &FrontendHandler{
		catalogService: catalogService,
		cartService:    cartService,
		pbService:      pbService,
		pbResolver:     pbResolver,
	}
}

func (h *FrontendHandler) HomePage(c *fiber.Ctx) error {
	// 1. Fetch Page Config
	page, err := h.pbService.GetPage(c.Context(), "/")
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	// 2. Resolve Data (Enrichment)
	if err := h.pbResolver.Resolve(c.Context(), page); err != nil {
		// Log but proceed? Or error?
		// For now proceed, sections might differ slightly.
	}

	// 3. Render Dynamic Page
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return pages.DynamicPage(page).Render(c.Context(), c.Response().BodyWriter())
}

func (h *FrontendHandler) ProductPage(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return util.APIError(c, fiber.StatusBadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid slug"))
	}

	product, err := h.catalogService.GetProductBySlug(c.Context(), slug)
	if err != nil {
		// handle not found specifically if possible, else 500
		return util.APIError(c, fiber.StatusNotFound, err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return pages.Product(product).Render(c.Context(), c.Response().BodyWriter())
}

// -- Cart --

func (h *FrontendHandler) CartPage(c *fiber.Ctx) error {
	sessionID, userID := h.getIdentities(c)
	cart, err := h.cartService.GetOrCreateCart(c.Context(), sessionID, userID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	items, err := h.cartService.GetCartItems(c.Context(), cart.ID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return pages.Cart(cart, items).Render(c.Context(), c.Response().BodyWriter())
}

// Helpers
func (h *FrontendHandler) getIdentities(c *fiber.Ctx) (pgtype.UUID, pgtype.UUID) {
	// Logic copied from CartHandler or shared lib.
	// Ideally shared, but for now duplicate to avoid strict coupling to handler package.
	sessID := pgtype.UUID{}
	userID := pgtype.UUID{}

	idStr, ok := c.Locals("user_id").(string)
	if !ok || idStr == "" {
		return sessID, userID
	}
	role, _ := c.Locals("user_role").(string)

	if role == "guest" {
		_ = sessID.Scan(idStr)
	} else {
		_ = userID.Scan(idStr)
	}
	return sessID, userID
}

// HTMX Actions

type AddToCartRequest struct {
	ProductID string `json:"product_id" form:"product_id"`
	Quantity  int    `json:"quantity" form:"quantity"`
}

func (h *FrontendHandler) AddToCart(c *fiber.Ctx) error {
	var req AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, err)
	}

	sessionID, userID := h.getIdentities(c)
	var pID pgtype.UUID
	if err := pID.Scan(req.ProductID); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, fiber.NewError(fiber.StatusBadRequest, "invalid product id"))
	}

	// Default quantity 1 if missing
	qty := req.Quantity
	if qty <= 0 {
		qty = 1
	}

	_, err := h.cartService.AddToCart(c.Context(), sessionID, userID, pID, pgtype.UUID{}, int32(qty))
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	// For HTMX, we can return a toast or update cart count.
	// For now, simpler: Redirect to cart or re-render button with 'Added!'
	// Let's redirect to Cart Page for MVP flow.
	// Or better: header HX-Trigger to update cart count.
	c.Set("HX-Redirect", "/cart")
	return c.SendStatus(fiber.StatusOK)
}

func (h *FrontendHandler) UpdateCartItem(c *fiber.Ctx) error {
	itemIDStr := c.FormValue("item_id")
	qtyStr := c.FormValue("quantity")

	var itemID pgtype.UUID
	if err := itemID.Scan(itemIDStr); err != nil {
		return util.APIError(c, fiber.StatusBadRequest, err)
	}
	qty, _ := util.StringToInt32(qtyStr) // Need helper or strconv

	sessionID, userID := h.getIdentities(c)
	cart, err := h.cartService.GetOrCreateCart(c.Context(), sessionID, userID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	item, err := h.cartService.UpdateItemQuantity(c.Context(), itemID, cart.ID, qty)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	// Fetch full row details to re-render the row (need joins)
	// Or just return the row component.
	// Since GetCartItemsRow has Join data, we need to fetch it again properly.
	// Shortcut: Fetch all items and find this one, or dedicated GetCartItem.
	// For MVP: Re-render the whole cart or just the row if we can fetch it.
	// Let's re-render the specific row by fetching items and filtering (inefficient but safe).
	items, _ := h.cartService.GetCartItems(c.Context(), cart.ID)
	var targetItem db.GetCartItemsRow
	for _, it := range items {
		if it.ID == item.ID {
			targetItem = it
			break
		}
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return pages.CartItemRow(targetItem).Render(c.Context(), c.Response().BodyWriter())
}

func (h *FrontendHandler) RemoveCartItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	var itemID pgtype.UUID
	itemID.Scan(idStr)

	sessionID, userID := h.getIdentities(c)
	cart, err := h.cartService.GetOrCreateCart(c.Context(), sessionID, userID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	err = h.cartService.RemoveItem(c.Context(), itemID, cart.ID)
	if err != nil {
		return util.APIError(c, fiber.StatusInternalServerError, err)
	}

	// Return empty string to remove element
	return c.SendString("")
}
