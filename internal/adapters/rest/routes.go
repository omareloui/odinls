package rest

import (
	"github.com/omareloui/odinls/web/views"
)

func (a *Adapter) registerRoutes() {
	a.server.Get("/", respondWithTemplate(views.Homepage()))
	a.server.Get("/login", respondWithTemplate(views.Login()))
	a.server.Get("/register", respondWithTemplate(views.Register(views.RegisterFormErrors{})))

	a.server.Post("/login", a.loginHandler)
	a.server.Post("/register", a.registerHandler)
}
