package restfiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
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

func (h *handler) GetMerchant(c fiber.Ctx) error {
	merchants, err := h.app.MerchantService.GetMerchants()
	status := http.StatusOK
	if err != nil {
		status = http.StatusInternalServerError
	}
	return respondWithTemplate(c, status, views.MerchantPage(merchants, newCreateMerchantFormData(&merchant.Merchant{}, &errs.ValidationError{})))
}

func (h *handler) PostMerchant(c fiber.Ctx) error {
	merchantform := &merchant.Merchant{
		Name: c.FormValue("name"),
		Logo: c.FormValue("logo"),
	}

	err := h.app.MerchantService.CreateMerchant(merchantform)

	if err == nil {
		renderToBody(c, views.MerchantOOB(merchantform))
		return respondWithTemplate(c, http.StatusOK, views.CreateMerchantForm(newCreateMerchantFormData(&merchant.Merchant{}, &errs.ValidationError{})))
	}

	if valerr, ok := err.(errs.ValidationError); ok {
		defaultValueAndErrs := newCreateMerchantFormData(merchantform, &valerr)
		return respondWithTemplate(c, http.StatusUnprocessableEntity, views.CreateMerchantForm(defaultValueAndErrs))
	}

	return respondWithTemplate(c, http.StatusInternalServerError, views.CreateMerchantForm(newCreateMerchantFormData(merchantform, &errs.ValidationError{})))
}
