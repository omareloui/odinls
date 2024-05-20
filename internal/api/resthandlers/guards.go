package resthandlers

import (
	"fmt"
	"net/http"
)

func (h *handler) AuthGuard(next http.HandlerFunc) http.HandlerFunc {
	// TODO(auth): try to refresh here?
	return func(w http.ResponseWriter, r *http.Request) {
		access := r.Context().Value(authContextKey)
		fmt.Println(access)
		if access == nil {
			hxRespondWithRedirect(w, "/unauthorized")
			return
		}
		next(w, r)
	}
}
