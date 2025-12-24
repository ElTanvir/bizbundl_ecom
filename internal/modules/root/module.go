package root

import (
	homepage "portfolioed/internal/modules/root/page"
	"portfolioed/internal/server"
	"portfolioed/util"

	"github.com/gofiber/fiber/v2"
)

func Init(server *server.Server) {
	router := server.GetRouter()

	// Home page
	router.Get("/", func(c *fiber.Ctx) error {
		return util.Render(c, homepage.HomePage())
	})
}
