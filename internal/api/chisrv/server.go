package chisrv

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/omareloui/odinls/internal/api/resthandlers"
)

type APIAdapter struct {
	handler resthandlers.Handler
	port    int
	router  *chi.Mux
}

func NewAdapter(handler resthandlers.Handler, port int) *APIAdapter {
	return &APIAdapter{handler: handler, port: port}
}

func (a *APIAdapter) Run() {
	a.router = chi.NewRouter()

	a.router.Use(middleware.Logger)
	a.router.Use(a.handler.AttachAuthenticatedUserMiddleware)

	a.Get("/", a.handler.GetHomepage)

	a.Get("/login", a.handler.GetLogin, withHasToNotBeSigned)
	a.Post("/login", a.handler.PostLogin, withHasToNotBeSigned)
	a.Get("/register", a.handler.GetRegister, withHasToNotBeSigned)
	a.Post("/register", a.handler.PostRegister, withHasToNotBeSigned)
	a.Post("/logout", a.handler.Logout, withHasToNotBeSigned)

	a.Get("/users", a.handler.GetUsers, withProtection)
	a.Get("/users/{id}", a.passParam("id", a.handler.GetUser))
	a.Get("/users/{id}/edit", a.passParam("id", a.handler.GetEditUser))
	a.Put("/users/{id}", a.passParam("id", a.handler.EditUser))
	a.Patch("/users/{id}/unset-craftsman", a.passParam("id", a.handler.UnsetCraftsman))
	a.Get("/users/craftsman-form", a.handler.GetCraftsmanForm, withProtection)

	a.Get("/roles", a.handler.GetRoles, withProtection)

	a.Get("/merchants", a.handler.GetMerchants, withProtection)
	a.Post("/merchants", a.handler.CreateMerchant, withProtection)

	a.Get("/merchants/{id}", a.passParam("id", a.handler.GetMerchant), withProtection)
	a.Put("/merchants/{id}", a.passParam("id", a.handler.EditMerchant), withProtection)
	a.Get("/merchants/{id}/edit", a.passParam("id", a.handler.GetEditMerchant), withProtection)

	a.Get("/clients", a.handler.GetClients, withProtection)
	a.Post("/clients", a.handler.CreateClient, withProtection)
	a.Get("/clients/{id}", a.passParam("id", a.handler.GetClient), withProtection)
	a.Get("/clients/{id}/edit", a.passParam("id", a.handler.GetEditClient), withProtection)
	a.Put("/clients/{id}", a.passParam("id", a.handler.EditClient), withProtection)

	a.Get("/products", a.handler.GetProducts, withProtection)
	a.Post("/products", a.handler.CreateProduct, withProtection)
	a.Get("/products/{id}", a.passParam("id", a.handler.GetProduct), withProtection)
	a.Get("/products/{id}/edit", a.passParam("id", a.handler.GetEditProduct), withProtection)
	a.Put("/products/{id}", a.passParam("id", a.handler.EditProduct), withProtection)

	a.Get("/orders", a.handler.GetOrders, withProtection)

	a.Get("/unauthorized", a.handler.Unauthorized)

	static(a.router, []string{"styles", "js"}, "./web/public")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.port), a.router))
}

func static(r *chi.Mux, paths []string, fspath string) {
	fs := http.FileServer(http.Dir(fspath))
	for _, prefix := range paths {
		r.Handle(fmt.Sprintf("/%s/*", prefix), http.StripPrefix("/", fs))
	}
}
