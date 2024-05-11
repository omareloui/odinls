package rest

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/omareloui/odinls/internal/application/core/domain"
	"github.com/omareloui/odinls/internal/misc/app_errors"
	"github.com/omareloui/odinls/web/views"
)

func (a *Adapter) loginHandler(c fiber.Ctx) error {
	dto := domain.NewLogin(c.FormValue("email"), c.FormValue("password"))
	usr, err := a.api.Login(c.Context(), *dto)
	if err != nil {
		return err
	}
	return c.SendString(fmt.Sprintf("the existing user is: %+v\n", usr))
}

func (a *Adapter) registerHandler(c fiber.Ctx) error {
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	email := c.FormValue("email")
	password := c.FormValue("password")
	cpassword := c.FormValue("cpassword")

	dto := domain.NewRegister(domain.Name{First: firstName, Last: lastName}, email, password, cpassword)

	usr, err := a.api.Register(c.Context(), *dto)
	if err != nil {
		vErr := app_errors.NewValidationErr(err.Error())
		c.Status(int(vErr.Code))
		return renderToBody(c, views.Register(views.RegisterFormErrors{}))
	}

	return c.SendString(fmt.Sprintf("the created user is: %+v\n", usr))
}
