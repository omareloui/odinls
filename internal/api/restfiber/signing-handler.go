package restfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/omareloui/odinls/web/views"
)

func newLoginFormData() *views.LoginFormData {
	return &views.LoginFormData{
		Email:    views.FormInputData{},
		Password: views.FormInputData{},
	}
}

func newRegisterFormData() *views.RegisterFormData {
	return &views.RegisterFormData{
		FristName:       views.FormInputData{},
		LastName:        views.FormInputData{},
		Email:           views.FormInputData{},
		Password:        views.FormInputData{},
		ConfirmPassword: views.FormInputData{},
	}
}

func (h *handler) GetLogin(c fiber.Ctx) error {
	return respondWithTemplate(c, http.StatusOK, views.Login(newLoginFormData()))
}

func (h *handler) GetRegister(c fiber.Ctx) error {
	return respondWithTemplate(c, http.StatusOK, views.Register(newRegisterFormData()))
}
