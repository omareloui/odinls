package restfiber

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/omareloui/odinls/config"
)

type APIAdapter struct {
	handler Handler
	port    int
	server  *fiber.App
}

func NewAdapter(handler Handler, port int) *APIAdapter {
	return &APIAdapter{handler: handler, port: port}
}

func (a *APIAdapter) Run() {
	a.server = fiber.New(fiber.Config{AppName: "Odin Leather Store"})
	a.server.Static("/", "./web/public/")

	a.server.Use(logger.New())

	a.server.Get("/", adaptor.HTTPHandlerFunc(a.handler.GetHomepage))

	a.server.Get("/login", adaptor.HTTPHandlerFunc(a.handler.GetLogin))
	a.server.Get("/register", adaptor.HTTPHandlerFunc(a.handler.GetRegister))

	a.server.Get("/merchant", adaptor.HTTPHandlerFunc(a.handler.GetMerchant))
	a.server.Post("/merchant", adaptor.HTTPHandlerFunc(a.handler.PostMerchant))

	log.Fatal(a.server.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}
