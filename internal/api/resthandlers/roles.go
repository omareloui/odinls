package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.app.RoleService.GetRoles()
	if err != nil {
		respondWithInternalServerError(w, r)
		return
	}

	accessClaims, _ := h.getAuthFromContext(r)
	respondWithTemplate(w, r, http.StatusOK, views.RolesPage(accessClaims, roles))
}
