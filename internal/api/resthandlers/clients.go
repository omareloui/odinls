package resthandlers

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetClients(w http.ResponseWriter, r *http.Request) {
	claims, _ := h.getAuthFromContext(r)

	clients, err := h.app.ClientService.GetClients(claims)
	if err == nil {
		respondWithTemplate(w, r, http.StatusOK, views.ClientsPage(claims, clients, &views.CreateClientFormData{}))
		return
	}

	if errors.Is(client.ErrClientNotFound, err) {
		respondWithTemplate(w, r, http.StatusOK, views.ClientsPage(claims, []client.Client{}, &views.CreateClientFormData{}))
		return
	}

	if errors.Is(errs.ErrForbidden, err) {
		respondWithForbidden(w, r)
		return
	}
	respondWithInternalServerError(w, r)
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	claims, _ := h.getAuthFromContext(r)

	cli := &client.Client{
		Name:               r.FormValue("name"),
		Notes:              r.FormValue("notes"),
		WholesaleAsDefault: r.FormValue("wholesale_as_default") == "on",
	}

	if phone := r.FormValue("phone"); phone != "" {
		cli.ContactInfo.PhoneNumbers = make(map[string]string)
		cli.ContactInfo.PhoneNumbers["default"] = phone
	}
	if email := r.FormValue("email"); email != "" {
		cli.ContactInfo.Emails = make(map[string]string)
		cli.ContactInfo.Emails["default"] = email
	}
	if link := r.FormValue("link"); link != "" {
		cli.ContactInfo.Links = make(map[string]string)
		cli.ContactInfo.Links["default"] = link
	}
	if location := r.FormValue("location"); location != "" {
		cli.ContactInfo.Locations = make(map[string]string)
		cli.ContactInfo.Locations["default"] = location
	}

	err := h.app.ClientService.CreateClient(claims, cli)
	if err != nil {
		if errors.Is(errs.ErrForbidden, err) {
			respondWithForbidden(w, r)
			return
		}
		if valerr, ok := err.(errs.ValidationError); ok {
			e := newCreateClientFormData(cli, &valerr)
			respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateClientForm(e))
			return
		}
		if errors.Is(client.ErrClientExistsForMerchant, err) {
			formdata := newCreateClientFormData(cli, &errs.ValidationError{})
			// TODO(research): make the email unique, not the name?
			formdata.Name.Error = "You already have a client with this name"
			respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateClientForm(formdata))
			return
		}
		respondWithInternalServerError(w, r)
		return
	}

	_ = renderToBody(w, r, views.ClientOOB(cli))
	respondWithTemplate(w, r, http.StatusOK, views.CreateClientForm(&views.CreateClientFormData{}))
}

func (h *handler) GetClient(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (h *handler) GetEditClient(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (h *handler) EditClient(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func newCreateClientFormData(client *client.Client, valerr *errs.ValidationError) *views.CreateClientFormData {
	formData := &views.CreateClientFormData{
		Name:  views.FormInputData{Value: client.Name, Error: valerr.Errors.MsgFor("Name")},
		Notes: views.FormInputData{Value: client.Notes, Error: valerr.Errors.MsgFor("Notes")},
	}

	if client.ContactInfo.PhoneNumbers != nil && len(client.ContactInfo.PhoneNumbers) > 0 {
		formData.Phone = views.FormInputData{Value: client.ContactInfo.PhoneNumbers["default"], Error: valerr.Errors.MsgFor("ContactInfo.PhoneNumbers")}
	}
	if client.ContactInfo.Emails != nil && len(client.ContactInfo.Emails) > 0 {
		formData.Email = views.FormInputData{Value: client.ContactInfo.Emails["default"], Error: valerr.Errors.MsgFor("ContactInfo.Emails")}
	}
	if client.ContactInfo.Locations != nil && len(client.ContactInfo.Locations) > 0 {
		formData.Location = views.FormInputData{Value: client.ContactInfo.Locations["default"], Error: valerr.Errors.MsgFor("ContactInfo.Locations")}
	}
	if client.ContactInfo.Links != nil && len(client.ContactInfo.Links) > 0 {
		formData.Link = views.FormInputData{Value: client.ContactInfo.Links["default"], Error: valerr.Errors.MsgFor("ContactInfo.Links")}
	}

	return formData
}
