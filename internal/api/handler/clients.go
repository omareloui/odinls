package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetClients(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	clients, err := h.app.ClientService.GetClients(claims)
	if err != nil {
		return responder.Error(err)
	}

	comp := views.ClientsPage(claims, clients, &views.ClientFormData{})
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	cli := new(client.Client)
	err := former.Populate(r, cli)
	if err != nil {
		return responder.BadRequest()
	}

	cli, err = h.app.ClientService.CreateClient(claims, cli)
	if err != nil {
		var fd *views.ClientFormData
		h.fm.MapToForm(cli, err, fd)
		return responder.Error(err, responder.WithComponentIfValidationErr(views.CreateClientForm(fd)))
	}

	return responder.Created(responder.WithOOBComponent(w, r.Context(), views.ClientOOB(cli)),
		responder.WithComponent(views.CreateClientForm(new(views.ClientFormData))))
}

func (h *handler) GetClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())
	c, err := h.app.ClientService.GetClientByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}
	return responder.OK(responder.WithComponent(views.Client(c)))
}

func (h *handler) GetEditClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	cli, err := h.app.ClientService.GetClientByID(claims, id)
	if err != nil {
		return responder.Error(err)
	}

	var fd *views.ClientFormData
	h.fm.MapToForm(cli, nil, fd)
	return responder.OK(responder.WithComponent(views.EditClient(cli, fd)))
}

func (h *handler) EditClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	cli := new(client.Client)
	err := former.Populate(r, cli)
	if err != nil {
		return responder.BadRequest()
	}

	cli, err = h.app.ClientService.UpdateClientByID(claims, id, cli)
	if err != nil {
		var fd *views.ClientFormData
		h.fm.MapToForm(cli, err, fd)
		return responder.Error(err,
			responder.WithComponentIfValidationErr(views.EditClient(cli, fd)))
	}

	return responder.OK(responder.WithComponent(views.Client(cli)))
}
