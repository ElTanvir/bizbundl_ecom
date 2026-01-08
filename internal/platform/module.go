package platform

import (

	// "bizbundl/internal/middleware" // Global?
	"bizbundl/internal/platform/auth"

	// cartservice "bizbundl/internal/storefront/cart/service"
	"bizbundl/internal/server"
)

func Init(app *server.Server) {
	auth.Init(app)
}
