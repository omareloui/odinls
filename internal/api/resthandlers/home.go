package resthandlers

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetHomepage(w http.ResponseWriter, r *http.Request) error {
	accessClaims, err := h.getAuthFromContext(r)

	if errors.Is(err, ErrNoAccessCookie) {
		return respondWithTemplate(w, r, http.StatusOK, views.Homepage(accessClaims))
	}

	return respondWithTemplate(w, r, http.StatusOK, views.Homepage(accessClaims))
}
