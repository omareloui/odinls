package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	users, err := h.app.UserService.GetUsers()
	if err != nil {
		return responder.Error(err)
	}

	claims := getClaims(r.Context())
	comp := views.UserPage(claims, users)
	return responder.OK(w, responder.WithComponent(comp))
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	usr, err := h.app.UserService.GetUserByID(id)
	if err != nil {
		return responder.Error(err)
	}
	comp := views.User(usr)
	return responder.OK(w, responder.WithComponent(comp))
}

func (h *handler) GetEditUser(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	usr, err := h.app.UserService.GetUserByID(id)
	if err != nil {
		return responder.Error(err)
	}

	fd := new(views.UserFormData)
	h.fm.MapToForm(usr, nil, fd)
	comp := views.EditUser(usr, fd)
	return responder.OK(w, responder.WithComponent(comp))
}

func (h *handler) EditUser(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")

	usr := new(user.User)
	err := former.Populate(r, usr)
	if err != nil {
		return responder.Error(err)
	}

	usr.ID = id

	usr, err = h.app.UserService.UpdateUserByID(id, usr)
	if err != nil {
		vfd := new(views.UserFormData)
		h.fm.MapToForm(usr, err, vfd)
		vcomp := views.EditUser(usr, vfd)

		cfd := new(views.UserFormData)
		h.fm.MapToForm(usr, nil, cfd)
		cfd.Email.Error = "Email or Username already exists, try another one"
		cfd.Username.Error = "Email or Username already exists, try another one"
		ccomp := views.EditUser(usr, cfd)

		return responder.Error(err,
			responder.WithComponentIfValidationErr(vcomp),
			responder.WithComponentIfErrIs(err, ccomp))
	}

	claims := getClaims(r.Context())
	if claims.ID == usr.ID {
		cookiesPair, err := h.newCookiesPairFromUser(usr)
		if err != nil {
			return responder.Error(err)
		}

		http.SetCookie(w, cookiesPair.Access)
		http.SetCookie(w, cookiesPair.Refresh)
	}

	comp := views.User(usr)
	return responder.OK(w, responder.WithComponent(comp))
}

func (h *handler) UnsetCraftsman(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")

	usr, err := h.app.UserService.UnsetCraftsmanByID(id)
	if err != nil {
		return responder.Error(err)
	}

	claims := getClaims(r.Context())
	if claims.ID == usr.ID {
		cookiesPair, err := h.newCookiesPairFromUser(usr)
		if err != nil {
			return responder.Error(err)
		}

		http.SetCookie(w, cookiesPair.Access)
		http.SetCookie(w, cookiesPair.Refresh)
	}

	comp := views.User(usr)
	return responder.OK(w, responder.WithComponent(comp))
}

func (h *handler) GetCraftsmanForm(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	comp := views.CraftsmanForm(&views.UserFormData{})
	return responder.OK(w, responder.WithComponent(comp))
}
