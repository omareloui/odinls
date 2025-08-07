package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	users, err := h.app.UserService.GetUsers()
	if err != nil {
		return RespondError(err)
	}

	claims := getClaims(r.Context())
	comp := views.UserPage(claims, users)
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	usr, err := h.app.UserService.GetUserByID(id)
	if err != nil {
		return RespondError(err)
	}
	comp := views.User(usr)
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) GetEditUser(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")
	usr, err := h.app.UserService.GetUserByID(id)
	if err != nil {
		return RespondError(err)
	}

	fd := mapUserToFormData(usr, &errs.ValidationError{})
	comp := views.EditUser(usr, fd)
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) EditUser(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")

	usr := new(user.User)
	err := former.Populate(r, usr)
	if err != nil {
		return RespondError(err)
	}

	usr.ID = id

	usr, err = h.app.UserService.UpdateUserByID(id, usr)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			fd := mapUserToFormData(usr, &valerr)
			comp := views.EditUser(usr, fd)
			return UnprocessableEntity(RespondWithComponent(comp))
		}

		alreadyExists := errors.Is(err, errs.ErrDocumentAlreadyExists)

		if alreadyExists {
			formdata := mapUserToFormData(usr, &errs.ValidationError{})
			formdata.Email.Error = "Email or Username already exists, try another one"
			formdata.Username.Error = "Email or Username already exists, try another one"
			comp := views.EditUser(usr, formdata)
			return UnprocessableEntity(RespondWithComponent(comp))
		}

		return RespondError(err)
	}

	// TODO: if the updated user is the current user, update the cookie

	comp := views.User(usr)
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) UnsetCraftsman(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	id := r.PathValue("id")

	usr, err := h.app.UserService.UnsetCraftsmanByID(id)
	if err != nil {
		return RespondError(err)
	}

	// TODO: if the updated user is the current user, update the cookie

	comp := views.User(usr)
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) GetCraftsmanForm(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	comp := views.CraftsmanForm(&views.UserFormData{})
	return RespondOK(w, RespondWithComponent(comp))
}

func mapUserToFormData(user *user.User, valerr *errs.ValidationError) *views.UserFormData {
	var hourlyRateStr string

	if user.Craftsman != nil {
		if user.Craftsman.HourlyRate != 0.00 {
			hourlyRateStr = strconv.FormatFloat(user.Craftsman.HourlyRate, 'f', -1, 64)
		}
	}

	return &views.UserFormData{
		Name: views.NameFormData{
			First: views.FormInputData{Value: user.Name.First, Error: valerr.Errors.MsgFor("Name.First")},
			Last:  views.FormInputData{Value: user.Name.Last, Error: valerr.Errors.MsgFor("Name.Last")},
		},
		Email:      views.FormInputData{Value: user.Email, Error: valerr.Errors.MsgFor("Email")},
		Username:   views.FormInputData{Value: user.Username, Error: valerr.Errors.MsgFor("Username")},
		Role:       views.FormInputData{Value: user.Role.String(), Error: valerr.Errors.MsgFor("Role")},
		HourlyRate: views.FormInputData{Value: hourlyRateStr, Error: valerr.Errors.MsgFor("HourlyRate")},
	}
}
