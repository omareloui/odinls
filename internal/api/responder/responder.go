package responder

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/formmap"
	"github.com/omareloui/odinls/internal/errs"
)

type responseOpts struct {
	context                  context.Context
	writer                   http.ResponseWriter
	message                  string
	path                     string
	component                templ.Component
	componentsIfErrIs        map[error]templ.Component
	componentIfValidationErr *templ.Component
	oobComponents            []templ.Component
}

type responseOptsFunc func(opts *responseOpts)

func WithPath(path string) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.path = path
	}
}

func WithMessage(msg string) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.message = msg
	}
}

func WithComponent(component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.component = component
	}
}

func WithComponentIfErrIs(err error, component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.componentsIfErrIs[err] = component
	}
}

func WithComponentIfValidationErr(component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.componentIfValidationErr = &component
	}
}

func WithOOBComponent(w http.ResponseWriter, context context.Context, component templ.Component) responseOptsFunc {
	return func(opts *responseOpts) {
		opts.oobComponents = append(opts.oobComponents, component)
		opts.context = context
		opts.writer = w
	}
}

func parseResponseOpts(opts ...responseOptsFunc) *responseOpts {
	_opts := &responseOpts{
		componentsIfErrIs: map[error]templ.Component{},
	}
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

func populateComponentIfErrorIsValidationError(opts *responseOpts, err error) {
	if _, ok := err.(*formmap.ValidationError); ok &&
		opts.componentIfValidationErr != nil {
		opts.component = *opts.componentIfValidationErr
	}
}

func OK(w http.ResponseWriter, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	w.WriteHeader(http.StatusOK)
	preRespond(_opts)
	return _opts.component, nil
}

func Created(w http.ResponseWriter, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	w.WriteHeader(http.StatusCreated)
	preRespond(_opts)
	return _opts.component, nil
}

func RedirectHX(w http.ResponseWriter, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	if _opts.path == "" {
		_opts.path = "/"
	}
	w.Header().Set("HX-Location", _opts.path)
	preRespond(_opts)
	return _opts.component, nil
}

func NotFound(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return notFound(_opts)
}

func notFound(opts *responseOpts) (templ.Component, error) {
	preRespond(opts)
	return opts.component, errs.NewRespError(http.StatusNotFound, opts.message)
}

func InternalServerError(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	preRespond(_opts)
	return _opts.component, errs.NewRespError(http.StatusInternalServerError, _opts.message)
}

func Unauthorized(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	preRespond(_opts)
	return _opts.component, errs.NewRespError(http.StatusUnauthorized, _opts.message)
}

func Forbidden(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	preRespond(_opts)
	return _opts.component, errs.NewRespError(http.StatusForbidden, _opts.message)
}

func BadRequest(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return badRequest(_opts)
}

func badRequest(opts *responseOpts) (templ.Component, error) {
	preRespond(opts)
	return opts.component, errs.NewRespError(http.StatusBadRequest, opts.message)
}

func UnprocessableEntity(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return unprocessableEntity(_opts)
}

func unprocessableEntity(opts *responseOpts) (templ.Component, error) {
	preRespond(opts)
	return opts.component, errs.NewRespError(http.StatusUnprocessableEntity, opts.message)
}

func Conflict(opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	return conflict(_opts)
}

func conflict(opts *responseOpts) (templ.Component, error) {
	preRespond(opts)
	return opts.component, errs.NewRespError(http.StatusConflict, opts.message)
}

func Error(err error, opts ...responseOptsFunc) (templ.Component, error) {
	_opts := parseResponseOpts(opts...)
	if errors.Is(err, errs.ErrDocumentNotFound) {
		populateComponentIfErrorIs(_opts, err, errs.ErrDocumentNotFound)
		return notFound(_opts)
	}

	_, isValerr := err.(*formmap.ValidationError)
	if isValerr ||
		errors.Is(err, errs.ErrInvalidID) ||
		errors.Is(err, errs.ErrInvalidFloat) ||
		errors.Is(err, errs.ErrInvalidNumber) ||
		errors.Is(err, errs.ErrInvalidDate) {
		populateComponentIfErrorIs(_opts, err,
			errs.ErrInvalidID, errs.ErrInvalidFloat,
			errs.ErrInvalidNumber, errs.ErrInvalidDate)
		populateComponentIfErrorIsValidationError(_opts, err)
		return unprocessableEntity(_opts)
	}

	if errors.Is(err, errs.ErrDocumentAlreadyExists) {
		populateComponentIfErrorIs(_opts, err, errs.ErrDocumentAlreadyExists)
		return conflict(_opts)
	}

	preRespond(_opts)
	return InternalServerError(opts...)
}

func preRespond(opts *responseOpts) {
	if len(opts.oobComponents) > 0 {
		for _, oobComp := range opts.oobComponents {
			if err := oobComp.Render(opts.context, opts.writer); err != nil {
				log.Panicf("error rendering template: %v", err)
			}
		}
	}
}
