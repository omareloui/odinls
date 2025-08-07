package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetClients(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	clients, err := h.app.ClientService.GetClients(claims)
	if err != nil {
		return RespondError(err)
	}

	comp := views.ClientsPage(claims, clients, &views.ClientFormData{})
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims := getClaims(r.Context())

	cli := new(client.Client)
	err := former.Populate(r, cli)
	if err != nil {
		return BadRequest()
	}

	cli, err = h.app.ClientService.CreateClient(claims, cli)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			return UnprocessableEntity(RespondWithComponent(views.CreateClientForm(mapClientToFormData(cli, &valerr))))
		}
		return RespondError(err)
	}

	if err := renderToBody(w, r, views.ClientOOB(cli)); err != nil {
		return RespondError(err)
	}

	return RespondCreated(w, RespondWithComponent(views.CreateClientForm(&views.ClientFormData{})))
}

func (h *handler) GetClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())
	c, err := h.app.ClientService.GetClientByID(claims, id)
	if err != nil {
		return RespondError(err)
	}
	return RespondOK(w, RespondWithComponent(views.Client(c)))
}

func (h *handler) GetEditClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	c, err := h.app.ClientService.GetClientByID(claims, id)
	if err != nil {
		return RespondError(err)
	}

	return RespondOK(w,
		RespondWithComponent(views.EditClient(c, mapClientToFormData(c, &errs.ValidationError{}))))
}

func (h *handler) EditClient(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	claims := getClaims(r.Context())

	cli := new(client.Client)
	err := former.Populate(r, cli)
	if err != nil {
		return BadRequest()
	}

	cli, err = h.app.ClientService.UpdateClientByID(claims, id, cli)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			return UnprocessableEntity(RespondWithComponent(views.EditClient(cli, mapClientToFormData(cli, &valerr))))
		}
		return RespondError(err)
	}

	return RespondOK(w, RespondWithComponent(views.Client(cli)))
}

func mapClientToFormData(client *client.Client, valerr *errs.ValidationError) *views.ClientFormData {
	formData := &views.ClientFormData{
		Name:  views.FormInputData{Value: client.Name, Error: valerr.Errors.MsgFor("Name")},
		Notes: views.FormInputData{Value: client.Notes, Error: valerr.Errors.MsgFor("Notes")},
	}

	if client.WholesaleAsDefault {
		formData.WholesaleAsDefault = views.FormInputData{Value: "on", Error: valerr.Errors.MsgFor("WholesaleAsDefault")}
	} else {
		formData.WholesaleAsDefault = views.FormInputData{Value: "", Error: valerr.Errors.MsgFor("WholesaleAsDefault")}
	}

	if client.ContactInfo.PhoneNumbers != nil {
		formData.Phone = views.FormInputData{Value: client.ContactInfo.PhoneNumbers["default"], Error: valerr.Errors.MsgFor("ContactInfo.PhoneNumbers")}
	}
	if client.ContactInfo.Emails != nil {
		formData.Email = views.FormInputData{Value: client.ContactInfo.Emails["default"], Error: valerr.Errors.MsgFor("ContactInfo.Emails")}
	}
	if client.ContactInfo.Locations != nil {
		formData.Location = views.FormInputData{Value: client.ContactInfo.Locations["default"], Error: valerr.Errors.MsgFor("ContactInfo.Locations")}
	}
	if client.ContactInfo.Links != nil {
		formData.Link = views.FormInputData{Value: client.ContactInfo.Links["default"], Error: valerr.Errors.MsgFor("ContactInfo.Links")}
	}

	return formData
}
