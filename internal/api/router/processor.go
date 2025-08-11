package router

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/logger"
	"go.uber.org/zap"
)

type process func(w http.ResponseWriter, r *http.Request) (templ.Component, error)

func (fn process) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := logger.FromCtx(r.Context())
	defer func() {
		if p := recover(); p != nil {
			l.Fatal("panicked", zap.Any("message", p))
			panic(p)
		}
	}()

	comp, err := fn(w, r)
	code := http.StatusOK

	if err != nil {
		code = getHTTPError(err).Code
	}

	w.WriteHeader(code)

	if err != nil && comp == nil {
		errorResponse(l, w, err)
		return
	}

	if comp != nil {
		w.Header().Set("Content-Type", "text/html")
		if err = comp.Render(r.Context(), w); err != nil {
			l.Error("rendering the component")
			errorResponse(l, w, err)
			return
		}
	}
}

func errorResponse(l *zap.Logger, w http.ResponseWriter, err error) {
	httpErr := getHTTPError(err)

	l.Error("server error", zap.Error(httpErr),
		zap.String("message", httpErr.Message),
		zap.Int("status_code", httpErr.Code),
	)

	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, httpErr.Message, httpErr.Code)
}

func getHTTPError(err error) *errs.RespError {
	var httperr *errs.RespError
	if !errors.As(err, &httperr) {
		httperr = errs.NewRespError(http.StatusInternalServerError, err.Error())
	}
	return httperr
}
