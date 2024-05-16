package restfiber

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
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

	a.server.Get("/", a.handler.GetHomepage)
	a.server.Get("/login", a.handler.GetLogin)
	a.server.Get("/register", a.handler.GetRegister)

	log.Fatal(a.server.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}
