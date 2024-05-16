package rest

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/internal/ports"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *fiber.App
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a *Adapter) Run() {
	a.server = fiber.New(fiber.Config{AppName: "Odin Leather Store"})
	a.server.Static("/", "./web/public/")

	a.registerRoutes()

	log.Fatal(a.server.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}
