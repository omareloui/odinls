package resthandlers

import "net/http"

func (h *handler) Unauthorized(w http.ResponseWriter, r *http.Request) {
	respondWithUnauthorized(w, r)
}
