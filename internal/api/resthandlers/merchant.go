package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetMerchants(w http.ResponseWriter, r *http.Request) error {
	merchants, err := h.app.MerchantService.GetMerchants()
	if err != nil {
		return err
	}

	accessClaims, _ := h.getAuthFromContext(r)
	return respondWithTemplate(w, r, http.StatusOK,
		views.MerchantPage(accessClaims, merchants, mapMerchantToFormData(&merchant.Merchant{}, &errs.ValidationError{})))
}

func (h *handler) GetMerchant(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		mer, err := h.app.MerchantService.GetMerchantByID(id)
		if err != nil {
			return err
		}
		return respondWithTemplate(w, r, http.StatusOK, views.Merchant(mer))
	}
}

func (h *handler) CreateMerchant(w http.ResponseWriter, r *http.Request) error {
	mer := mapFormToMerchant(r)
	err := h.app.MerchantService.CreateMerchant(mer)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			e := mapMerchantToFormData(mer, &valerr)
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateMerchantForm(e))
		}
		return err
	}

	if err := renderToBody(w, r, views.MerchantOOB(mer)); err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusCreated, views.CreateMerchantForm(mapMerchantToFormData(&merchant.Merchant{}, &errs.ValidationError{})))
}

func (h *handler) GetEditMerchant(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		mer, err := h.app.MerchantService.GetMerchantByID(id)
		if err != nil {
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK,
			views.EditMerchant(mer, mapMerchantToFormData(mer, &errs.ValidationError{})))
	}
}

func (h *handler) EditMerchant(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		mer := mapFormToMerchant(r)

		err := h.app.MerchantService.UpdateMerchantByID(id, mer)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditMerchant(mer, mapMerchantToFormData(mer, &valerr)))
			}
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.Merchant(mer))
	}
}

func mapMerchantToFormData(merchant *merchant.Merchant, valerr *errs.ValidationError) *views.MerchantFormData {
	return &views.MerchantFormData{
		Name: views.FormInputData{Value: merchant.Name, Error: valerr.Errors.MsgFor("Name")},
		Logo: views.FormInputData{Value: merchant.Logo, Error: valerr.Errors.MsgFor("Logo")},
	}
}

func mapFormToMerchant(r *http.Request) *merchant.Merchant {
	return &merchant.Merchant{
		Name: r.FormValue("name"),
		Logo: r.FormValue("logo"),
	}
}
