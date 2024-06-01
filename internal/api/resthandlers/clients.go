package resthandlers

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetClients(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)

	clients, err := h.app.ClientService.GetClients(claims)
	if err != nil {
		return err
	}

	return respondWithTemplate(w, r, http.StatusOK, views.ClientsPage(claims, clients, &views.ClientFormData{}))
}

func (h *handler) CreateClient(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)
	cli := mapFormToClient(r)

	err := h.app.ClientService.CreateClient(claims, cli)
	if err != nil {
		if errors.Is(client.ErrClientExistsForMerchant, err) {
			formdata := mapClientToFormData(cli, &errs.ValidationError{})
			// TODO(research): make the email unique if it exists, not the name?
			formdata.Name.Error = "You already have a client with this name"
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateClientForm(formdata))
		}
		if valerr, ok := err.(errs.ValidationError); ok {
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateClientForm(mapClientToFormData(cli, &valerr)))
		}
		return err
	}

	if err := renderToBody(w, r, views.ClientOOB(cli)); err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.CreateClientForm(&views.ClientFormData{}))
}

func (h *handler) GetClient(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, _ := h.getAuthFromContext(r)
		c, err := h.app.ClientService.GetClientByID(claims, id)
		if err != nil {
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.Client(c))
	}
}

func (h *handler) GetEditClient(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, _ := h.getAuthFromContext(r)
		c, err := h.app.ClientService.GetClientByID(claims, id)
		if err != nil {
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.EditClient(c, mapClientToFormData(c, &errs.ValidationError{})))
	}
}

func (h *handler) EditClient(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, _ := h.getAuthFromContext(r)

		cli := mapFormToClient(r)

		err := h.app.ClientService.UpdateClientByID(claims, id, cli)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditClient(cli, mapClientToFormData(cli, &valerr)))
			}
			if errors.Is(client.ErrClientExistsForMerchant, err) {
				formdata := mapClientToFormData(cli, &errs.ValidationError{})
				formdata.Name.Error = "You already have a client with this name"
				return respondWithTemplate(w, r, http.StatusConflict, views.EditClient(cli, formdata))
			}
			return err
		}

		if err := renderToBody(w, r, views.ClientOOB(cli)); err != nil {
			return err
		}
		return respondWithTemplate(w, r, http.StatusOK, views.Client(cli))
	}
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

func mapFormToClient(r *http.Request) *client.Client {
	cli := &client.Client{
		Name:               r.FormValue("name"),
		Notes:              r.FormValue("notes"),
		WholesaleAsDefault: r.FormValue("wholesale_as_default") == "on",
		ContactInfo:        client.ContactInfo{},
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
	return cli
}
