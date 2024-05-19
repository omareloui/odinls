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

	a.router.Get("/", a.handler.GetHomepage)

	a.router.Get("/login", a.handler.GetLogin)
	a.router.Post("/login", a.handler.PostLogin)
	a.router.Get("/register", a.handler.GetRegister)
	a.router.Post("/register", a.handler.PostRegister)

	a.router.Get("/merchants", a.handler.GetMerchants)
	a.router.Post("/merchants", a.handler.PostMerchant)

	a.router.Get("/merchants/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		a.handler.GetMerchant(id)(w, r)
	})
	a.router.Patch("/merchants/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		a.handler.EditMerchant(id)(w, r)
	})
	a.router.Get("/merchants/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		a.handler.GetEditMerchant(id)(w, r)
	})

	static(a.router, []string{"styles"}, "./web/public")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.port), a.router))
}

func static(r *chi.Mux, paths []string, fspath string) {
	fs := http.FileServer(http.Dir(fspath))
	for _, prefix := range paths {
		r.Handle(fmt.Sprintf("/%s/*", prefix), http.StripPrefix("/", fs))
	}
}
