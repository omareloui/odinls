package resthandlers

import (
	"net/http"

	"github.com/a-h/templ"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	application "github.com/omareloui/odinls/internal/application/core"
	"github.com/omareloui/odinls/web/views"
)

type Handler interface {
	AttachAuthenticatedUserMiddleware(next http.Handler) http.Handler

	AuthGuard(next http.HandlerFunc) http.HandlerFunc
	AlreadyAuthedGuard(next http.HandlerFunc) http.HandlerFunc

	Unauthorized(w http.ResponseWriter, r *http.Request)

	GetHomepage(w http.ResponseWriter, r *http.Request)

	GetLogin(w http.ResponseWriter, r *http.Request)
	PostLogin(w http.ResponseWriter, r *http.Request)
	GetRegister(w http.ResponseWriter, r *http.Request)
	PostRegister(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUser(id string) http.HandlerFunc
	GetEditUser(id string) http.HandlerFunc
	EditUser(id string) http.HandlerFunc
	UnsetCraftsman(id string) http.HandlerFunc
	GetCraftsmanForm(w http.ResponseWriter, r *http.Request)

	GetRoles(w http.ResponseWriter, r *http.Request)

	GetMerchants(w http.ResponseWriter, r *http.Request)
	CreateMerchant(w http.ResponseWriter, r *http.Request)
	GetMerchant(id string) http.HandlerFunc
	GetEditMerchant(id string) http.HandlerFunc
	EditMerchant(id string) http.HandlerFunc

	GetClients(w http.ResponseWriter, r *http.Request)
	CreateClient(w http.ResponseWriter, r *http.Request)
	GetClient(id string) http.HandlerFunc
	GetEditClient(id string) http.HandlerFunc
	EditClient(id string) http.HandlerFunc

	GetProducts(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProduct(id string) http.HandlerFunc
	GetEditProduct(id string) http.HandlerFunc
	EditProduct(id string) http.HandlerFunc

	GetOrders(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	app        *application.Application
	jwtAdapter jwtadapter.JwtAdapter
}

func NewHandler(app *application.Application, jwtAdapter jwtadapter.JwtAdapter) Handler {
	return &handler{app: app, jwtAdapter: jwtAdapter}
}

func respondWithTemplate(w http.ResponseWriter, r *http.Request, status int, template templ.Component) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if err := renderToBody(w, r, template); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func respondWithInternalServerError(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusInternalServerError)
}

func respondWithUnauthorized(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusUnauthorized)
}

func respondWithForbidden(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusForbidden)
}

func respondWithNotFound(w http.ResponseWriter, r *http.Request) {
	respondWithErrorPage(w, r, http.StatusNotFound)
}

func respondWithErrorPage(w http.ResponseWriter, r *http.Request, status int) {
	auth := r.Context().Value(authContextKey)
	respondWithTemplate(w, r, status, views.ErrorPage(auth.(*jwtadapter.JwtAccessClaims), http.StatusText(status), status))
}

func renderToBody(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	return template.Render(r.Context(), w)
}

func hxRespondWithRedirect(w http.ResponseWriter, path string) {
	w.Header().Set("HX-Location", path)
}
