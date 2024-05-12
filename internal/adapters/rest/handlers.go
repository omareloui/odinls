package rest

import (
	"fmt"

	"github.com/go-playground/validator/v10"
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

	dto := domain.NewRegister(domain.AuthName{First: firstName, Last: lastName}, email, password, cpassword)

	usr, err := a.api.Register(c.Context(), *dto)
	if err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			vErr := app_errors.NewValidationErr(&validationErrs)
			return respondWithTemplate(vErr.Code, views.RegisterForm(dto, vErr.Errors))(c)
		}
		if existingEmailErr, ok := err.(app_errors.EmailAlreadyInUse); ok {
			return respondWithTemplate(existingEmailErr.Code, views.RegisterForm(dto, existingEmailErr.Errors))(c)
		}
	}

	return c.SendString(fmt.Sprintf("the created user is: %+v\n", usr))
}
