package resthandlers

import (
	"net/http"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	application "github.com/omareloui/odinls/internal/application/core"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

type Handler interface {
	AttachAuthenticatedUserMiddleware(next http.Handler) http.Handler

	ErrorHandler(w http.ResponseWriter, r *http.Request, err error)

	AuthGuard(next HandlerFunc) HandlerFunc
	AlreadyAuthedGuard(next HandlerFunc) HandlerFunc

	Unauthorized(w http.ResponseWriter, r *http.Request) error

	GetHomepage(w http.ResponseWriter, r *http.Request) error

	GetLogin(w http.ResponseWriter, r *http.Request) error
	Login(w http.ResponseWriter, r *http.Request) error
	GetRegister(w http.ResponseWriter, r *http.Request) error
	Register(w http.ResponseWriter, r *http.Request) error
	Logout(w http.ResponseWriter, r *http.Request) error

	GetUsers(w http.ResponseWriter, r *http.Request) error
	GetUser(id string) HandlerFunc
	GetEditUser(id string) HandlerFunc
	EditUser(id string) HandlerFunc
	UnsetCraftsman(id string) HandlerFunc
	GetCraftsmanForm(w http.ResponseWriter, r *http.Request) error

	GetClients(w http.ResponseWriter, r *http.Request) error
	CreateClient(w http.ResponseWriter, r *http.Request) error
	GetClient(id string) HandlerFunc
	GetEditClient(id string) HandlerFunc
	EditClient(id string) HandlerFunc

	GetProducts(w http.ResponseWriter, r *http.Request) error
	CreateProduct(w http.ResponseWriter, r *http.Request) error
	GetProduct(id string) HandlerFunc
	GetEditProduct(id string) HandlerFunc
	EditProduct(id string) HandlerFunc

	GetOrders(w http.ResponseWriter, r *http.Request) error
	CreateOrder(w http.ResponseWriter, r *http.Request) error
	GetOrder(id string) HandlerFunc
	GetEditOrder(id string) HandlerFunc
	EditOrder(id string) HandlerFunc
}

type handler struct {
	app        *application.Application
	jwtAdapter jwtadapter.JwtAdapter
}

func NewHandler(app *application.Application, jwtAdapter jwtadapter.JwtAdapter) Handler {
	return &handler{app: app, jwtAdapter: jwtAdapter}
}
