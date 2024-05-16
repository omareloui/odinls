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

	a.server.Get("/merchant", adaptor.HTTPHandlerFunc(a.handler.GetMerchant))
	a.server.Post("/merchant", adaptor.HTTPHandlerFunc(a.handler.PostMerchant))

	log.Fatal(a.server.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}
