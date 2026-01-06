package checkout

import (
	"bizbundl/internal/modules/cart/service"
	catalogService "bizbundl/internal/modules/catalog/service"
	"bizbundl/pkgs/components/registry"
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Resolver struct {
	cartSvc    *service.CartService
	catalogSvc *catalogService.CatalogService
}

func NewResolver(cartSvc *service.CartService, catalogSvc *catalogService.CatalogService) *Resolver {
	return &Resolver{
		cartSvc:    cartSvc,
		catalogSvc: catalogSvc,
	}
}

func (r *Resolver) Resolve(ctx context.Context, section *registry.Section) error {
	// Attempt to get Session/User from Context (Requires Middleware to set these on context.Context)
	// For now, assuming they are available or we handle "Empty" state gracefully.
	// NOTE: In the current architecture, Fiber Locals aren't auto-propagated to ctx.
	// We might just render the "Empty" structure and let HTMX load the cart?
	// OR we assume we can get it.

	// Let's implement a "Lazy" approach or assume the Handler injects "session_id" into ctx.
	// If context doesn't have it, we return empty props, and the widget serves as a placeholder
	// or relies on client-side/htmx to fetch real state?
	// NO, the widget is Server Side. We need the data.

	// CRITICAL: We need SessionID.
	// Temporary Hack: We will look for a convention or assume standard Context usage.
	// If fails, we default to empty.

	sessionID, _ := ctx.Value("session_id").(pgtype.UUID)
	userID, _ := ctx.Value("user_id").(pgtype.UUID)

	// Propulate the section props
	props := make(map[string]interface{})

	// 1. Get Cart
	// If Valid Session
	if sessionID.Valid {
		cart, err := r.cartSvc.GetOrCreateCart(ctx, sessionID, userID)
		if err == nil {
			items, _ := r.cartSvc.GetCartItems(ctx, cart.ID)
			props["Cart"] = cart
			props["Items"] = items

			hasPhysical, _ := r.cartSvc.HasPhysicalItems(ctx, cart.ID)
			props["IsPhysical"] = hasPhysical
		}
	}

	// 2. Direct Checkout (Check Props or Context for overrides?)
	// If this component is placed on a page, it usually acts as "Cart Checkout".
	// "Direct Checkout" is usually dynamic via URL.
	// If URL params are in context, we could read them.
	// For now, we focus on Cart Checkout.

	section.Props = props
	return nil
}
