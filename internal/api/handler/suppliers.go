package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/supplier"
	"github.com/omareloui/odinls/internal/logger"
	"github.com/omareloui/odinls/web/views"
	"go.uber.org/zap"
)

func (h *handler) GetSuppliers(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	ctx := r.Context()
	claims := getClaims(ctx)
	l := logger.FromCtx(ctx)

	suppliers, err := h.app.SupplierService.GetSuppliers(claims)
	if err != nil {
		l.Error("failed to get suppliers", zap.Any("claims", claims), zap.Error(err))
		return responder.Error(err)
	}

	comp := views.SuppliersPage(claims, suppliers, &views.SupplierFormData{})
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) CreateSupplier(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	sup := new(supplier.Supplier)
	err := former.Populate(r, sup)
	if err != nil {
		return responder.BadRequest()
	}

	sup, err = h.app.SupplierService.CreateSupplier(claims, sup)
	if err != nil {
		var fd *views.SupplierFormData
		h.fm.MapToForm(sup, err, fd)
		return responder.Error(err, responder.WithComponentIfValidationErr(views.CreateSupplierForm(fd)))
	}

	return responder.Created(responder.WithOOBComponent(w, r.Context(), views.SupplierOOB(sup)),
		responder.WithComponent(views.CreateSupplierForm(new(views.SupplierFormData))))
}

func (h *handler) GetSupplier(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())
	c, err := h.app.SupplierService.GetSupplierByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}
	return responder.OK(responder.WithComponent(views.Supplier(c)))
}

func (h *handler) GetEditSupplier(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	sup, err := h.app.SupplierService.GetSupplierByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}

	fd := new(views.SupplierFormData)
	h.fm.MapToForm(sup, nil, fd)
	return responder.OK(responder.WithComponent(views.EditSupplier(sup, fd)))
}

func (h *handler) EditSupplier(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	sup := new(supplier.Supplier)
	err := former.Populate(r, sup)
	if err != nil {
		return responder.BadRequest()
	}

	sup, err = h.app.SupplierService.UpdateSupplierByID(claims, id, sup)
	if err != nil {
		fd := new(views.SupplierFormData)
		h.fm.MapToForm(sup, err, fd)
		return responder.Error(err,
			responder.WithComponentIfValidationErr(views.EditSupplier(sup, fd)))
	}

	return responder.OK(responder.WithComponent(views.Supplier(sup)))
}
