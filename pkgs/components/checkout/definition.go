package checkout

import (
	"bizbundl/internal/modules/cart/service"
	catalogService "bizbundl/internal/modules/catalog/service"
	"bizbundl/pkgs/components/registry"

	"github.com/a-h/templ"
)

func Register(cartSvc *service.CartService, catSvc *catalogService.CatalogService) {
	registry.Register(&registry.Component{
		Type:     "checkout_widget",
		Resolver: NewResolver(cartSvc, catSvc),
		Renderer: func(props map[string]interface{}) templ.Component {
			return View(props)
		},
	})
}
