package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.app.UserService.GetUsers()
	if err != nil {
		return err
	}

	claims, _ := h.getAuthFromContext(r)
	return respondWithTemplate(w, r, http.StatusOK, views.UserPage(claims, users))
}

func (h *handler) GetUser(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		usr, err := h.app.UserService.GetUserByID(id)
		if err != nil {
			return err
		}
		return respondWithTemplate(w, r, http.StatusOK, views.User(usr))
	}
}

func (h *handler) GetEditUser(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		usr, err := h.app.UserService.GetUserByID(id)
		if err != nil {
			return err
		}

		return h.responseWithEditUser(w, r, http.StatusOK, usr, mapUserToFormData(usr, &errs.ValidationError{}))
	}
}

func (h *handler) EditUser(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		usr, err := mapEditUserFormToUser(id, r)
		if err != nil {
			if errors.Is(err, errs.ErrInvalidFloat) {
				formdata := mapUserToFormData(usr, &errs.ValidationError{})
				formdata.HourlyRate.Error = "Invalid number"
				return h.responseWithEditUser(w, r, http.StatusUnprocessableEntity, usr, formdata)
			}
			return err
		}

		err = h.app.UserService.UpdateUserByID(id, usr)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				return h.responseWithEditUser(w, r, http.StatusUnprocessableEntity, usr, mapUserToFormData(usr, &valerr))
			}

			emailExists := errors.Is(err, user.ErrEmailAlreadyExists)
			usernameExists := errors.Is(err, user.ErrUsernameAlreadyExists)

			if emailExists || usernameExists {
				formdata := mapUserToFormData(usr, &errs.ValidationError{})
				if emailExists {
					formdata.Email.Error = "Email already exists, try another one"
				}
				if usernameExists {
					formdata.Username.Error = "Username already exists, try another one"
				}
				return h.responseWithEditUser(w, r, http.StatusUnprocessableEntity, usr, formdata)
			}

			return err
		}

		// TODO: if the updated user is the current user, update the cookie

		return h.GetUser(id)(w, r)
	}
}

func (h *handler) UnsetCraftsman(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := h.app.UserService.UnsetCraftsmanByID(id)
		if err != nil {
			return err
		}

		// TODO: if the updated user is the current user, update the cookie

		return h.GetUser(id)(w, r)
	}
}

func (h *handler) GetCraftsmanForm(w http.ResponseWriter, r *http.Request) error {
	return respondWithTemplate(w, r, http.StatusOK, views.CraftsmanForm(&views.UserFormData{}))
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

func mapEditUserFormToUser(id string, r *http.Request) (*user.User, error) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	username := r.FormValue("username")
	role := r.FormValue("role")
	hourlyRate := r.FormValue("hourly_rate")

	usr := &user.User{
		ID:       id,
		Name:     user.Name{First: firstName, Last: lastName},
		Email:    email,
		Username: username,
		Role:     user.RoleFromString(role),
	}

	isCraftsman := hourlyRate != ""
	if isCraftsman {
		usr.Craftsman = &user.Craftsman{}
		hourlyRate, err := strconv.ParseFloat(hourlyRate, 64)
		if err != nil {
			return usr, errs.ErrInvalidFloat
		}
		usr.Craftsman.HourlyRate = hourlyRate
	}

	return usr, nil
}

func (h *handler) responseWithEditUser(w http.ResponseWriter, r *http.Request, status int, usr *user.User, formdata *views.UserFormData) error {
	return respondWithTemplate(w, r, status, views.EditUser(usr, formdata))
}
