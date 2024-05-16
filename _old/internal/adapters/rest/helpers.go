package rest

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
)

func respondWithTemplate(c fiber.Ctx, status int, template templ.Component) {
	c.Status(status)
	c.Type(".html")
	return renderToResponseBody(c, template)
}

func renderToResponseBody(c fiber.Ctx, template templ.Component) error {
	return template.Render(c.Context(), c.Response().BodyWriter())
}
