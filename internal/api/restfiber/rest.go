package restfiber

import (
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	application "github.com/omareloui/odinls/internal/application/core"
)

type Handler interface {
	GetHomepage(fiber.Ctx) error

	GetLogin(fiber.Ctx) error
	GetRegister(fiber.Ctx) error

	GetMerchant(fiber.Ctx) error
	PostMerchant(fiber.Ctx) error
}

type handler struct {
	app *application.Application
}

func NewHandler(app *application.Application) Handler {
	return &handler{
		app: app,
	}
}

func respondWithTemplate(c fiber.Ctx, status int, template templ.Component) error {
	c.Status(status)
	c.Type(".html")
	return template.Render(c.Context(), c.Response().BodyWriter())
}
