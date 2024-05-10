package rest

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func respondWithTemplate(template templ.Component) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Type(".html")
		return renderToBody(c, template)
	}
}

func renderToBody(c fiber.Ctx, template templ.Component) error {
	return template.Render(c.Context(), c.Response().BodyWriter())
}
