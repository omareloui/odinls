package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	claims, _ := h.getAuthFromContext(r)
	respondWithTemplate(w, r, http.StatusOK, views.OrdersPage(claims))
}
