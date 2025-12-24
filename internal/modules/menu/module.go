package menu

import (
	menupage "portfolioed/internal/modules/menu/page"
	"portfolioed/internal/server"
	"portfolioed/util"

	"github.com/gofiber/fiber/v2"
)

func Init(server *server.Server) {
	router := server.GetRouter()

	// Menu page
	router.Get("/menu", func(c *fiber.Ctx) error {
		return util.Render(c, menupage.MenuPage())
	})
}
