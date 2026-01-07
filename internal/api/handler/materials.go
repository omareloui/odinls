package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/material"
	"github.com/omareloui/odinls/internal/logger"
	"github.com/omareloui/odinls/web/views"
	"go.uber.org/zap"
)

func (h *handler) GetMaterials(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	ctx := r.Context()
	claims := getClaims(ctx)

	materials, err := h.app.MaterialService.GetMaterials(claims)
	if err != nil {
		return responder.Error(err)
	}

	suppliers, err := h.app.SupplierService.GetSuppliers(claims)
	if err != nil {
		return responder.Error(err)
	}

	comp := views.MaterialsPage(claims, materials, suppliers, &views.MaterialFormData{})
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) CreateMaterial(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	ctx := r.Context()

	l := logger.FromCtx(ctx)
	claims := getClaims(ctx)

	mat := new(material.Material)
	err := former.Populate(r, mat)
	if err != nil {
		return responder.BadRequest()
	}

	l.Debug("creating a material", zap.Any("material", mat))
	mat, err = h.app.MaterialService.CreateMaterial(claims, mat)
	if err != nil {
		l.Error("creating a material", zap.Error(err), zap.Any("material", mat))

		suppliers, supErr := h.app.SupplierService.GetSuppliers(claims)
		if supErr != nil {
			return responder.Error(supErr)
		}

		fd := new(views.MaterialFormData)
		h.fm.MapToForm(mat, err, fd)
		return responder.Error(err, responder.WithComponentIfValidationErr(views.CreateMaterialForm(fd, suppliers)))
	}

	suppliers, err := h.app.SupplierService.GetSuppliers(claims)
	if err != nil {
		return responder.Error(err)
	}

	return responder.Created(responder.WithOOBComponent(w, ctx, views.MaterialOOB(mat)),
		responder.WithComponent(views.CreateMaterialForm(new(views.MaterialFormData), suppliers)))
}

func (h *handler) GetMaterial(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())
	c, err := h.app.MaterialService.GetMaterialByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}
	return responder.OK(responder.WithComponent(views.Material(c)))
}

func (h *handler) GetEditMaterial(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	mat, err := h.app.MaterialService.GetMaterialByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}

	suppliers, err := h.app.SupplierService.GetSuppliers(claims)
	if err != nil {
		return responder.Error(err)
	}

	fd := new(views.MaterialFormData)
	h.fm.MapToForm(mat, nil, fd)
	return responder.OK(responder.WithComponent(views.EditMaterial(mat, fd, suppliers)))
}

func (h *handler) EditMaterial(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	mat := new(material.Material)
	err := former.Populate(r, mat)
	if err != nil {
		return responder.BadRequest()
	}

	mat, err = h.app.MaterialService.UpdateMaterialByID(claims, id, mat)
	if err != nil {
		suppliers, supErr := h.app.SupplierService.GetSuppliers(claims)
		if supErr != nil {
			return responder.Error(supErr)
		}

		fd := new(views.MaterialFormData)
		h.fm.MapToForm(mat, err, fd)
		return responder.Error(err,
			responder.WithComponentIfValidationErr(views.EditMaterial(mat, fd, suppliers)))
	}

	return responder.OK(responder.WithComponent(views.Material(mat)))
}
