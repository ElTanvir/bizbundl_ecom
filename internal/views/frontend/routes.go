package frontend

import (
	homepage "bizbundl/internal/views/frontend/homepage"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
)

type ViewRoute struct {
	Path    string
	Method  string
	Handler fiber.Handler
}

var Routes = []ViewRoute{
	{
		Path:   "/",
		Method: "GET",
		Handler: func(c *fiber.Ctx) error {
			return util.Render(c, homepage.HomePage())
		},
	},
}
