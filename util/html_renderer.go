package util

import (
	"bytes"
	"sync"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	err := component.Render(c.Context(), buf)
	if err != nil {
		return err
	}

	return c.Send(buf.Bytes())
}
