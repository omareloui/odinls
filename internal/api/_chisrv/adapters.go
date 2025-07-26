package chisrv

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/omareloui/odinls/internal/api/resthandlers"
)

func (a *APIAdapter) errorHandlerAdapter(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err == nil {
			return
		}

		a.handler.ErrorHandler(w, r, err)
	}
}

func (a *APIAdapter) passParam(key string, handler func(string) resthandlers.HandlerFunc) resthandlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return handler(chi.URLParam(r, key))(w, r)
	}
}

func (a *APIAdapter) route(method, pattern string, handlerFunc resthandlers.HandlerFunc, options ...optsFunc) {
	opts := parseOpts(options...)

	var fn func(pattern string, handlerFunc http.HandlerFunc)

	switch method {
	case "PUT":
		fn = a.router.Put
	case "POST":
		fn = a.router.Post
	case "PATCH":
		fn = a.router.Patch
	case "DELETE":
		fn = a.router.Delete
	default:
		fn = a.router.Get
	}

	if opts.protected {
		handlerFunc = a.handler.AuthGuard(handlerFunc)
	}

	if opts.hasToNotBeSigned {
		handlerFunc = a.handler.AlreadyAuthedGuard(handlerFunc)
	}

	fn(pattern, a.errorHandlerAdapter(handlerFunc))
}

func (a *APIAdapter) Get(pattern string, handlerFunc resthandlers.HandlerFunc, opts ...optsFunc) {
	a.route("GET", pattern, handlerFunc, opts...)
}

func (a *APIAdapter) Post(pattern string, handlerFunc resthandlers.HandlerFunc, opts ...optsFunc) {
	a.route("POST", pattern, handlerFunc, opts...)
}

func (a *APIAdapter) Put(pattern string, handlerFunc resthandlers.HandlerFunc, opts ...optsFunc) {
	a.route("PUT", pattern, handlerFunc, opts...)
}

func (a *APIAdapter) Patch(pattern string, handlerFunc resthandlers.HandlerFunc, opts ...optsFunc) {
	a.route("PATCH", pattern, handlerFunc, opts...)
}

func (a *APIAdapter) Delete(pattern string, handlerFunc resthandlers.HandlerFunc, opts ...optsFunc) {
	a.route("DELETE", pattern, handlerFunc, opts...)
}
