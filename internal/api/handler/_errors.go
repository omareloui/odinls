package handler

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
)

func (h *handler) Unauthorized(w http.ResponseWriter, r *http.Request) error {
	return respondWithUnauthorized(w, r)
}

func (h *handler) ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(errs.ErrForbidden, err) {
		_ = respondWithForbidden(w, r)
		return
	}

	if errors.Is(errs.ErrInvalidID, err) || errors.Is(errs.ErrInvalidFloat, err) {
		_ = respondWithString(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if errors.Is(order.ErrOrderNotFound, err) ||
		errors.Is(user.ErrUserNotFound, err) ||
		errors.Is(product.ErrProductNotFound, err) ||
		errors.Is(product.ErrVariantNotFound, err) ||
		errors.Is(client.ErrClientNotFound, err) {
		_ = respondWithString(w, r, http.StatusNotFound, err.Error())
		return
	}

	_ = respondWithInternalServerError(w, r)
}
