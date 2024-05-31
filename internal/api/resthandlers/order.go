package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

// TODO(refactor): update the expected signature for the handlers to return the
// error and handle the errors in an centralized place.
// Research: how to pass the validation handler

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)
	ords, err := h.app.OrderService.GetOrders(claims)
	if err != nil {
		return err
	}
	respondWithTemplate(w, r, http.StatusOK, views.OrdersPage(claims, ords))
}
