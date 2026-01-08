package handler

import (
	root "bizbundl/internal/platform/root/view"
	"bizbundl/util"

	"github.com/gofiber/fiber/v2"
)

type RootHandler struct{}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (h *RootHandler) ShowHome(c *fiber.Ctx) error {
	return util.Render(c, root.Home())
}
