package restfiber

import (
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/merchant"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func newCreateMerchantFormData(merchant *merchant.Merchant, valerr *errs.ValidationError) *views.CreateMerchantFormData {
	return &views.CreateMerchantFormData{
		Name: views.FormInputData{Value: merchant.Name, Error: valerr.Errors.MsgFor("Name")},
		Logo: views.FormInputData{Value: merchant.Logo, Error: valerr.Errors.MsgFor("Logo")},
	}
}

// TODO: add a page for the handler to show not found and 500 pages

func (h *handler) GetMerchant(w http.ResponseWriter, r *http.Request) {
	merchants, err := h.app.MerchantService.GetMerchants()
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}
	respondWithTemplate(w, r, status, views.MerchantPage(merchants, newCreateMerchantFormData(&merchant.Merchant{}, &errs.ValidationError{})))
}

func (h *handler) PostMerchant(w http.ResponseWriter, r *http.Request) {
	merchantform := &merchant.Merchant{
		Name: r.FormValue("name"),
		Logo: r.FormValue("logo"),
	}

	err := h.app.MerchantService.CreateMerchant(merchantform)

	if err == nil {
		renderToBody(w, r, views.MerchantOOB(merchantform))
		respondWithTemplate(w, r, http.StatusOK, views.CreateMerchantForm(newCreateMerchantFormData(&merchant.Merchant{}, &errs.ValidationError{})))
	}

	if valerr, ok := err.(errs.ValidationError); ok {
		defaultValueAndErrs := newCreateMerchantFormData(merchantform, &valerr)
		respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateMerchantForm(defaultValueAndErrs))
	}

	respondWithTemplate(w, r, http.StatusInternalServerError, views.CreateMerchantForm(newCreateMerchantFormData(merchantform, &errs.ValidationError{})))
}
