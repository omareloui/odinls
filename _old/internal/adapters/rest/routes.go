package rest

import (
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/domain"
	"github.com/omareloui/odinls/web/views"
)

func (a *Adapter) registerRoutes() {
	a.server.Get("/", respondWithTemplate(http.StatusOK, views.Homepage()))
	a.server.Get("/login", respondWithTemplate(http.StatusOK, views.Login()))
	a.server.Get("/register", respondWithTemplate(http.StatusOK, views.Register(&domain.Register{}, nil)))

	a.server.Post("/login", a.loginHandler)
	a.server.Post("/register", a.registerHandler)
}
