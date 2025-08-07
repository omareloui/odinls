package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/omareloui/odinls/internal/errs"
	"go.mongodb.org/mongo-driver/mongo"
)

type responseOpts struct {
	message            string
	path               string
	component          templ.Component
	componentsIfErrIs  map[error]templ.Component
	componentsIfErrMsg map[string]templ.Component
}

type responseOptsFunc func(opts *responseOpts)

func RespondWithPath(path string) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.path = path
	}
}

func RespondWithMessage(msg string) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.message = msg
	}
}

func RespondWithComponent(component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.component = component
	}
}

func RespondWithComponentIfErrIs(err error, component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.componentsIfErrIs[err] = component
	}
}

func RespondWithComponentIfErrMessageContains(err error, component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.componentsIfErrMsg[err.Error()] = component
	}
}

func parseResponseOpts(opts ...responseOptsFunc) *responseOpts {
	_opts := &responseOpts{}
	for _, opt := range opts {
		opt(_opts)
	}
	return _opts
}

func populateComponentIfErrorIs(opts *responseOpts, err error, errIs error, orErrIs ...error) {
	if len(opts.componentsIfErrIs) == 0 {
		return
	}
	errs := append([]error{errIs}, orErrIs...)
	for _, _errIs := range errs {
		if errors.Is(err, _errIs) {
			for e, comp := range opts.componentsIfErrIs {
				if errors.Is(err, e) {
					opts.component = comp
				}
			}
		}
	}
}

func populateComponentIfErrorMessageContains(opts *responseOpts, err error, msg string, orMsg ...string) {
	if len(opts.componentsIfErrMsg) == 0 {
		return
	}
	msgs := append([]string{msg}, orMsg...)
	for _, _msg := range msgs {
		if err.Error() == _msg {
			for mmsg, comp := range opts.componentsIfErrMsg {
				if strings.Contains(err.Error(), mmsg) {
					opts.component = comp
				}
			}
		}
	}
}

func RespondOK(w http.ResponseWriter, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	w.WriteHeader(http.StatusOK)
	return _opts.component, nil
}

func RespondCreated(w http.ResponseWriter, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	w.WriteHeader(http.StatusCreated)
	return _opts.component, nil
}

func RespondRedirectHX(w http.ResponseWriter, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	if _opts.path == "" {
		_opts.path = "/"
	}
	w.Header().Set("HX-Location", _opts.path)
	return _opts.component, nil
}

func NotFound(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusNotFound, _opts.message)
}

func InternalServerError(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusInternalServerError, _opts.message)
}

func Unauthorized(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusUnauthorized, _opts.message)
}

func Forbidden(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusForbidden, _opts.message)
}

func BadRequest(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusBadRequest, _opts.message)
}

func UnprocessableEntity(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusUnprocessableEntity, _opts.message)
}

func Conflict(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return _opts.component, errs.NewRespError(http.StatusConflict, _opts.message)
}

func RespondError(err error, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	if errors.Is(err, mongo.ErrNoDocuments) {
		populateComponentIfErrorIs(_opts, err, mongo.ErrNoDocuments)
		return NotFound(opts...)
	}

	valerr := &errs.ValidationError{}
	if valerr.Error() == err.Error() ||
		errors.Is(err, errs.ErrInvalidID) ||
		errors.Is(err, errs.ErrInvalidFloat) ||
		errors.Is(err, errs.ErrInvalidNumber) ||
		errors.Is(err, errs.ErrInvalidDate) {
		populateComponentIfErrorIs(_opts, err,
			errs.ErrInvalidID, errs.ErrInvalidFloat,
			errs.ErrInvalidNumber, errs.ErrInvalidDate)
		populateComponentIfErrorMessageContains(_opts, err, valerr.Error())
		return UnprocessableEntity(opts...)
	}

	if errors.Is(err, errs.ErrDocumentAlreadyExists) {
		populateComponentIfErrorIs(_opts, err, errs.ErrDocumentAlreadyExists)
		return Conflict(opts...)
	}

	return InternalServerError(opts...)
}
