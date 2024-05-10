package main

import (
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/omareloui/odinls/config"
	"github.com/omareloui/odinls/web/views"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	app := fiber.New(fiber.Config{AppName: "Odin Leather Store"})

	app.Static("/", "./web/public/")

	router(app)

	log.Fatal(app.Listen(fmt.Sprintf(":%d", config.GetApplicationPort())))
}

func router(app *fiber.App) {
	app.Get("/", respondWithTemplate(views.Homepage()))
	app.Get("/login", respondWithTemplate(views.Login()))
	app.Post("/login", postLogin)
	app.Post("/register", postRegister)
	app.Get("/register", respondWithTemplate(views.Register()))
}

func postLogin(c fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	return c.SendString("email is: " + email + " and the password is: " + password)
}

func postRegister(c fiber.Ctx) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	cpassword := c.FormValue("cpassword")
	return c.SendString("username is: " + username + " email is: " + email + " and the password is: " + password + " and confirm password is " + cpassword)
}

func respondWithTemplate(template templ.Component) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Type(".html")
		return renderToBody(c, template)
	}
}

func renderToBody(c fiber.Ctx, template templ.Component) error {
	return template.Render(c.Context(), c.Response().BodyWriter())
}
