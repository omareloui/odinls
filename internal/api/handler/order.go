package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	ords, err := h.app.OrderService.GetOrders(claims)
	if err != nil {
		return responder.Error(err)
	}

	prods, clients, err := h.getProdsAndClients(claims)
	if err != nil {
		return responder.Error(err)
	}
	return responder.OK(responder.WithComponent(views.OrdersPage(claims, prods, clients, ords)))
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	ord := new(order.Order)
	err := former.Populate(r, ord)
	if err != nil {
		return responder.Error(err)
	}

	prods, clients, err := h.getProdsAndClients(claims)
	if err != nil {
		return responder.Error(err)
	}

	ord, err = h.app.OrderService.CreateOrder(claims, ord)
	if err != nil {
		fd := new(views.OrderFormData)
		h.fm.MapToForm(ord, err, fd)
		comp := views.CreateOrderForm(ord, prods, clients, fd)
		return responder.Error(err, responder.WithComponentIfValidationErr(comp))
	}

	return responder.OK(responder.WithOOBComponent(w, r.Context(), views.OrderOOB(ord)),
		responder.WithComponent(views.CreateOrderForm(new(order.Order), prods, clients,
			views.NewDefaultOrderFormData())))
}

func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	ord, err := h.app.OrderService.GetOrderByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}

	return responder.OK(responder.WithComponent(views.Order(ord)))
}

func (h *handler) GetEditOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	ord, err := h.app.OrderService.GetOrderByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}

	prods, clients, err := h.getProdsAndClients(claims)
	if err != nil {
		return responder.Error(err)
	}

	fd := new(views.OrderFormData)
	h.fm.MapToForm(ord, nil, fd)

	return responder.OK(responder.WithComponent(views.EditOrder(ord, prods, clients, fd)))
}

func (h *handler) EditOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	ord := new(order.Order)
	err := former.Populate(r, ord)
	if err != nil {
		return responder.BadRequest()
	}

	ord, err = h.app.OrderService.UpdateOrderByID(claims, id, ord)
	if err != nil {
		prods, clients, err := h.getProdsAndClients(claims)
		if err != nil {
			return responder.Error(err)
		}
		fd := new(views.OrderFormData)
		h.fm.MapToForm(ord, err, fd)
		comp := views.EditOrder(ord, prods, clients, fd)
		return responder.Error(err, responder.WithComponentIfValidationErr(comp))
	}

	return responder.OK(responder.WithComponent(views.Order(ord)))
}

func (h *handler) getProdsAndClients(claims *jwtadapter.AccessClaims) ([]product.Product, []client.Client, error) {
	prods, err := h.app.ProductService.GetProducts(claims)
	if err != nil {
		return nil, nil, err
	}
	clients, err := h.app.ClientService.GetClients(claims)
	if err != nil {
		return nil, nil, err
	}

	return prods, clients, nil
}
