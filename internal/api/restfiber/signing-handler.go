package restfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetLogin(c fiber.Ctx) error {
	return respondWithTemplate(c, http.StatusOK, views.Login())
}

func (h *handler) GetRegister(c fiber.Ctx) error {
	return respondWithTemplate(c, http.StatusOK, views.Register())
}
