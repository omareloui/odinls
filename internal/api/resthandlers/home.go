package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetHomepage(w http.ResponseWriter, r *http.Request) {
	respondWithTemplate(w, r, http.StatusOK, views.Homepage())
}
