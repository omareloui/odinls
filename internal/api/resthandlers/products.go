package resthandlers

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	claims, _ := h.getAuthFromContext(r)
	prods, err := h.app.ProductService.GetProducts(claims)
	if err != nil {
		if errors.Is(errs.ErrForbidden, err) {
			respondWithForbidden(w, r)
			return
		}
		respondWithInternalServerError(w, r)
		return
	}
	respondWithTemplate(w, r, http.StatusOK, views.ProductsPage(claims, prods))
}

func (h *handler) GetProductVariantForm(w http.ResponseWriter, r *http.Request) {
	// respondWithTemplate(w, r, http.StatusOK, views.ProductVariantFormBody())
}
