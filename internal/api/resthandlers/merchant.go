package resthandlers

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

// TODO(refactor): add a page for the handler to show not found and 500 pages

func (h *handler) GetMerchants(w http.ResponseWriter, r *http.Request) {
	merchants, err := h.app.MerchantService.GetMerchants()
	if err != nil {
		respondWithInternalServerError(w, r)
		return
	}
	accessClaims, _ := h.getAuthFromContext(r)
	respondWithTemplate(w, r, http.StatusOK, views.MerchantPage(accessClaims, merchants, newCreateMerchantFormData(&merchant.Merchant{}, &errs.ValidationError{})))
}

func (h *handler) GetMerchant(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := h.app.MerchantService.FindMerchant(id)
		if ok := errors.Is(err, merchant.ErrMerchantNotFound); ok {
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				respondWithInternalServerError(w, r)
				return
			}
			return
		}

		respondWithTemplate(w, r, http.StatusOK, views.Merchant(m))
	}
}

func (h *handler) CreateMerchant(w http.ResponseWriter, r *http.Request) {
	merchantform := &merchant.Merchant{
		Name: r.FormValue("name"),
		Logo: r.FormValue("logo"),
	}

	err := h.app.MerchantService.CreateMerchant(merchantform)

	if err == nil {
		if err := renderToBody(w, r, views.MerchantOOB(merchantform)); err != nil {
			respondWithInternalServerError(w, r)
			return
		}
		respondWithTemplate(w, r, http.StatusCreated, views.CreateMerchantForm(newCreateMerchantFormData(&merchant.Merchant{}, &errs.ValidationError{})))
		return
	}

	if valerr, ok := err.(errs.ValidationError); ok {
		e := newCreateMerchantFormData(merchantform, &valerr)
		respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateMerchantForm(e))
		return
	}

	respondWithInternalServerError(w, r)
}

func (h *handler) GetEditMerchant(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := h.app.MerchantService.FindMerchant(id)
		if err != nil {
			respondWithInternalServerError(w, r)
			return
		}

		respondWithTemplate(w, r, http.StatusOK,
			views.EditMerchant(m, newCreateMerchantFormData(m, &errs.ValidationError{})))
	}
}

func (h *handler) EditMerchant(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		logo := r.FormValue("logo")

		m := &merchant.Merchant{ID: id, Name: name, Logo: logo}

		err := h.app.MerchantService.UpdateMerchantByID(id, m)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditMerchant(m, newCreateMerchantFormData(m, &valerr)))
				return
			}
			respondWithInternalServerError(w, r)
			return
		}

		respondWithTemplate(w, r, http.StatusOK, views.Merchant(m))
	}
}

func newCreateMerchantFormData(merchant *merchant.Merchant, valerr *errs.ValidationError) *views.CreateMerchantFormData {
	return &views.CreateMerchantFormData{
		Name: views.FormInputData{Value: merchant.Name, Error: valerr.Errors.MsgFor("Name")},
		Logo: views.FormInputData{Value: merchant.Logo, Error: valerr.Errors.MsgFor("Logo")},
	}
}
