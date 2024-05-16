package restfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetHomepage(c fiber.Ctx) error {
	return respondWithTemplate(c, http.StatusOK, views.Homepage())
}
