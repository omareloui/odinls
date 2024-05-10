package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/omareloui/odinls/internal/application/core/domain"
	"github.com/omareloui/odinls/web/views"
)

func (a *Adapter) registerRoutes() {
	a.server.Get("/", respondWithTemplate(views.Homepage()))
	a.server.Get("/login", respondWithTemplate(views.Login()))
	a.server.Get("/register", respondWithTemplate(views.Register()))

	a.server.Post("/login", func(c fiber.Ctx) error {
		dto := domain.NewLogin(c.FormValue("email"), c.FormValue("password"))
		usr, err := a.api.Login(c.Context(), *dto)
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("the existing user is: %+v\n", usr))
	})

	a.server.Post("/register", func(c fiber.Ctx) error {
		firstName := c.FormValue("firstName")
		lastName := c.FormValue("lastName")
		email := c.FormValue("email")
		password := c.FormValue("password")
		cpassword := c.FormValue("cpassword")

		dto := domain.NewRegister(domain.Name{First: firstName, Last: lastName}, email, password, cpassword)

		usr, err := a.api.Register(c.Context(), *dto)
		if err != nil {
			return err
		}
		return c.SendString(fmt.Sprintf("the created user is: %+v\n", usr))
	})
}
