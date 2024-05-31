package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetRoles(w http.ResponseWriter, r *http.Request) error {
	roles, err := h.app.RoleService.GetRoles()
	if err != nil {
		return err
	}

	claims, _ := h.getAuthFromContext(r)
	return respondWithTemplate(w, r, http.StatusOK, views.RolesPage(claims, roles))
}
