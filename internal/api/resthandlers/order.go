package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)
	ords, err := h.app.OrderService.GetOrders(claims)
	if err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.OrdersPage(claims, ords))
}
