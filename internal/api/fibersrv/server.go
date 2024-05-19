package fibersrv

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/api/resthandlers"
)

type APIAdapter struct {
	handler resthandlers.Handler
	port    int
	server  *fiber.App
}

func NewAdapter(handler resthandlers.Handler, port int) *APIAdapter {
	return &APIAdapter{handler: handler, port: port}
}

func (a *APIAdapter) Run() {
	a.server = fiber.New(fiber.Config{AppName: "Odin Leather Store"})
	a.server.Static("/", "./web/public/")

	a.server.Use(logger.New())

	a.server.Get("/", adaptor.HTTPHandlerFunc(a.handler.GetHomepage))

	a.server.Get("/login", adaptor.HTTPHandlerFunc(a.handler.GetLogin))
	a.server.Post("/login", adaptor.HTTPHandlerFunc(a.handler.PostLogin))
	a.server.Get("/register", adaptor.HTTPHandlerFunc(a.handler.GetRegister))
	a.server.Post("/register", adaptor.HTTPHandlerFunc(a.handler.PostRegister))

	a.server.Get("/merchants", adaptor.HTTPHandlerFunc(a.handler.GetMerchants))
	a.server.Post("/merchants", adaptor.HTTPHandlerFunc(a.handler.PostMerchant))
	a.server.Get("/merchants/:id", func(c fiber.Ctx) error {
		return adaptor.HTTPHandlerFunc(a.handler.GetMerchant(c.Params("id")))(c)
	})
	a.server.Patch("/merchants/:id", func(c fiber.Ctx) error {
		return adaptor.HTTPHandlerFunc(a.handler.EditMerchant(c.Params("id")))(c)
	})
	a.server.Get("/merchants/edit/:id", func(c fiber.Ctx) error {
		return adaptor.HTTPHandlerFunc(a.handler.GetEditMerchant(c.Params("id")))(c)
	})

	log.Fatal(a.server.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}
