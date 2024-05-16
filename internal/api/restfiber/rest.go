package restfiber

import (
	"net/http"

	"github.com/a-h/templ"
	application "github.com/omareloui/odinls/internal/application/core"
)

type Handler interface {
	GetHomepage(w http.ResponseWriter, r *http.Request)

	GetLogin(w http.ResponseWriter, r *http.Request)
	GetRegister(w http.ResponseWriter, r *http.Request)

	GetMerchant(w http.ResponseWriter, r *http.Request)
	PostMerchant(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	app *application.Application
}

func NewHandler(app *application.Application) Handler {
	return &handler{app: app}
}

func respondWithTemplate(w http.ResponseWriter, r *http.Request, status int, template templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if err := renderToBody(w, r, template); err != nil {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func renderToBody(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	return template.Render(r.Context(), w)
}
