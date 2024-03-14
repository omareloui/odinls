package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}
