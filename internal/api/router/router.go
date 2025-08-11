package router

import (
	"fmt"
	"net/http"

	"github.com/omareloui/odinls/internal/api/handler"
	"github.com/omareloui/odinls/internal/api/middleware"
)

type Router interface {
	http.Handler
}

func New(h handler.Handler) Router {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", handlePub(h.GetHomepage))

	mux.Handle("GET /login", handlePub(h.GetLogin))
	mux.Handle("POST /login", handlePub(h.Login))
	mux.Handle("GET /register", handlePub(h.GetRegister))
	mux.Handle("POST /register", handlePub(h.Register))
	mux.Handle("POST /logout", handle(h.Logout))
	mux.Handle("GET /refresh-tokens", handlePub(h.RefreshTokens))

	mux.Handle("GET /users", handle(h.GetUsers))
	mux.Handle("GET /users/{id}", handle(h.GetUser))
	mux.Handle("GET /users/{id}/edit", handle(h.GetEditUser))
	mux.Handle("PUT /users/{id}", handle(h.EditUser))
	mux.Handle("PATCH /users/{id}/unset-craftsman", handle(h.UnsetCraftsman))
	mux.Handle("GET /users/craftsman-form", handle(h.GetCraftsmanForm))

	mux.Handle("GET /clients", handle(h.GetClients))
	mux.Handle("GET /clients/{id}", handle(h.GetClient))
	mux.Handle("GET /clients/{id}/edit", handle(h.GetEditClient))
	mux.Handle("PUT /clients/{id}", handle(h.EditClient))
	mux.Handle("POST /clients", handle(h.CreateClient))

	mux.Handle("GET /products", handle(h.GetProducts))
	mux.Handle("GET /products/{id}", handle(h.GetProduct))
	mux.Handle("GET /products/{id}/edit", handle(h.GetEditProduct))
	mux.Handle("PUT /products/{id}", handle(h.EditProduct))
	mux.Handle("POST /products", handle(h.CreateProduct))

	mux.Handle("GET /orders", handle(h.GetOrders))
	mux.Handle("GET /orders/{id}", handle(h.GetOrder))
	mux.Handle("GET /orders/{id}/edit", handle(h.GetEditOrder))
	mux.Handle("PUT /orders/{id}", handle(h.EditOrder))
	mux.Handle("POST /orders", handle(h.CreateOrder))

	mux.Handle("GET /unauthorized", handlePub(h.Unauthorized))

	static(mux, []string{"styles", "js", "images"}, "./web/public")
	mux.Handle("/", handle(h.NotFound))

	return mux
}

func static(mux *http.ServeMux, paths []string, fspath string) {
	fs := http.StripPrefix("/", http.FileServer(http.Dir(fspath)))
	for _, prefix := range paths {
		mux.Handle(fmt.Sprintf("GET /%s/{path...}", prefix), middleware.CorrelationID(middleware.RequestLogger(fs)))
	}
}

func handle(h handler.HandlerMethod, appendMiddlewares ...(func(http.Handler) http.Handler)) http.Handler {
	appendMiddlewares = append(appendMiddlewares, middleware.Protected)
	return handlePub(h, appendMiddlewares...)
}

func handlePub(h handler.HandlerMethod, appendMiddlewares ...(func(http.Handler) http.Handler)) http.Handler {
	var handler http.Handler = process(h)

	for i := len(appendMiddlewares) - 1; i >= 0; i-- {
		handler = appendMiddlewares[i](handler)
	}

	handler = middleware.Auth(handler)
	handler = middleware.RequestLogger(handler)
	handler = middleware.CorrelationID(handler)

	return handler
}
