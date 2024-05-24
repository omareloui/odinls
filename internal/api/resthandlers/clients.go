package resthandlers

import (
	"errors"
	"fmt"
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
		cli.ContactInfo.PhoneNumbers["default"] = phone
	}
	if email := r.FormValue("email"); email != "" {
		cli.ContactInfo.Emails["default"] = email
	}
	if link := r.FormValue("link"); link != "" {
		cli.ContactInfo.Links["default"] = link
	}
	if location := r.FormValue("location"); location != "" {
		cli.ContactInfo.Locations["default"] = location
	}

	err := h.app.ClientService.CreateClient(claims, cli)
	fmt.Println("err ==>", err)
	if err != nil {
		// TODO: handler validation and duplication errors
		// TODO: handler forbidden
		respondWithInternalServerError(w, r)
		return
	}

	_ = renderToBody(w, r, views.ClientOOB(cli))
	respondWithTemplate(w, r, http.StatusOK, views.EditClient(cli, &views.CreateClientFormData{}))
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
