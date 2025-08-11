package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/formmap"
	application "github.com/omareloui/odinls/internal/application/core"
	"github.com/omareloui/odinls/internal/interfaces"
)

type HandlerMethod func(w http.ResponseWriter, r *http.Request) (templ.Component, error)

type Handler interface {
	GetHomepage(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	GetLogin(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	Login(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetRegister(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	Register(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	Logout(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	RefreshTokens(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	GetUsers(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetUser(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetEditUser(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	EditUser(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	UnsetCraftsman(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetCraftsmanForm(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	GetClients(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	CreateClient(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetClient(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetEditClient(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	EditClient(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	GetProducts(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	CreateProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetEditProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	EditProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	GetOrders(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	CreateOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	GetEditOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	EditOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	// AttachAuthenticatedUserMiddleware(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	// ErrorHandler(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	// AuthGuard(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	// AlreadyAuthedGuard(w http.ResponseWriter, r *http.Request) (templ.Component, error)

	Unauthorized(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	NotFound(w http.ResponseWriter, r *http.Request) (templ.Component, error)
	// InternalServerError(w http.ResponseWriter, r *http.Request) (templ.Component, error)
}

type handler struct {
	app *application.Application
	fm  interfaces.FormMapper
}

func New(app *application.Application) Handler {
	mapper := formmap.NewMapper()
	return &handler{app: app, fm: mapper}
}
