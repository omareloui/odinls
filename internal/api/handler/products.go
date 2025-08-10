package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	prods, err := h.app.ProductService.GetProducts(claims)
	if err != nil {
		return responder.Error(err)
	}

	comp := views.ProductsPage(claims, prods)
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	prod := new(product.Product)
	err := former.Populate(r, prod)
	if err != nil {
		return responder.BadRequest()
	}

	prod, err = h.app.ProductService.CreateProduct(claims, prod)
	if err != nil {
		fd := new(views.ProductFormData)
		h.fm.MapToForm(prod, err, fd)
		comp := views.CreateProductForm(prod, fd, claims.Craftsman.HourlyRate)
		return responder.Error(err, responder.WithComponentIfValidationErr(comp))
	}

	oobComp := views.ProductOOB(prod, claims.Craftsman.HourlyRate)
	comp := views.CreateProductForm(&product.Product{},
		&views.ProductFormData{Variants: []views.ProductVariantFormData{{}}},
		claims.Craftsman.HourlyRate)
	return responder.OK(responder.WithOOBComponent(w, r.Context(), oobComp),
		responder.WithComponent(comp))
}

func (h *handler) GetProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	prod, err := h.app.ProductService.GetProductByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}
	comp := views.Product(prod, claims.Craftsman.HourlyRate)
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) GetEditProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	prod, err := h.app.ProductService.GetProductByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}
	fd := new(views.ProductFormData)
	h.fm.MapToForm(prod, nil, fd)
	comp := views.EditProduct(prod, fd, claims.Craftsman.HourlyRate)
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) EditProduct(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	prod := new(product.Product)
	err := former.Populate(r, prod)
	if err != nil {
		return responder.BadRequest()
	}

	prod, err = h.app.ProductService.UpdateProductByID(claims, id, prod)
	if err != nil {
		fd := new(views.ProductFormData)
		h.fm.MapToForm(prod, err, fd)
		comp := views.EditProduct(prod, fd, claims.Craftsman.HourlyRate)
		return responder.Error(err, responder.WithComponentIfValidationErr(comp))
	}

	return responder.OK(responder.WithComponent(views.Product(prod, claims.Craftsman.HourlyRate)))
}
