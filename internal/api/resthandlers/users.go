package resthandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.app.UserService.GetUsers(user.WithPopulatedRole)
	if err != nil {
		return err
	}

	claims, _ := h.getAuthFromContext(r)
	return respondWithTemplate(w, r, http.StatusOK, views.UserPage(claims, users))
}

func (h *handler) GetUser(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		usr, err := h.app.UserService.FindUser(id, user.WithPopulatedRole)
		if err != nil {
			return err
		}
		return respondWithTemplate(w, r, http.StatusOK, views.User(usr))
	}
}

func (h *handler) GetEditUser(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		usr, err := h.app.UserService.FindUser(id, user.WithPopulatedRole)
		if err != nil {
			return err
		}

		return h.responseWithEditUser(w, r, http.StatusOK, usr, mapUserToFormData(usr, &errs.ValidationError{}), usr.IsCraftsman())
	}
}

func (h *handler) EditUser(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		usr, err := mapEditUserFormToUser(id, r)
		if err != nil {
			if errors.Is(errs.ErrInvalidFloat, err) {
				formdata := mapUserToFormData(usr, &errs.ValidationError{})
				formdata.HourlyRate.Error = "Invalid number"
				return h.responseWithEditUser(w, r, http.StatusUnprocessableEntity, usr, mapUserToFormData(usr, &errs.ValidationError{}), true)
			}
			return err
		}

		err = h.app.UserService.UpdateUserByID(id, usr, user.WithPopulatedRole)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				isCraftsman := usr.Craftsman.MerchantID != "" || usr.Craftsman.HourlyRate != 0
				return h.responseWithEditUser(w, r, http.StatusUnprocessableEntity, usr, mapUserToFormData(usr, &valerr), isCraftsman)
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
				return h.responseWithEditUser(w, r, http.StatusUnprocessableEntity, usr, formdata, usr.IsCraftsman())
			}

			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.User(usr))
	}
}

func (h *handler) UnsetCraftsman(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := h.app.UserService.UnsetCraftsmanByID(id)
		if err != nil {
			return err
		}

		user, err := h.app.UserService.FindUser(id, user.WithPopulatedRole)
		if err != nil {
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.User(user))
	}
}

func (h *handler) GetCraftsmanForm(w http.ResponseWriter, r *http.Request) error {
	merchants, err := h.app.MerchantService.GetMerchants()
	if err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.CraftsmanForm(merchants, &views.CreateUserFormData{}))
}

func mapUserToFormData(user *user.User, valerr *errs.ValidationError) *views.CreateUserFormData {
	var hourlyRateStr string
	var merId string

	if user.Craftsman != nil {
		if user.Craftsman.HourlyRate != 0.00 {
			hourlyRateStr = strconv.FormatFloat(user.Craftsman.HourlyRate, 'f', -1, 64)
		}
		merId = user.Craftsman.MerchantID
	}

	return &views.CreateUserFormData{
		Name: views.NameFormData{
			First: views.FormInputData{Value: user.Name.First, Error: valerr.Errors.MsgFor("Name.First")},
			Last:  views.FormInputData{Value: user.Name.Last, Error: valerr.Errors.MsgFor("Name.Last")},
		},
		Email:      views.FormInputData{Value: user.Email, Error: valerr.Errors.MsgFor("Email")},
		Username:   views.FormInputData{Value: user.Username, Error: valerr.Errors.MsgFor("Username")},
		Role:       views.FormInputData{Value: user.RoleID, Error: valerr.Errors.MsgFor("RoleID")},
		HourlyRate: views.FormInputData{Value: hourlyRateStr, Error: valerr.Errors.MsgFor("HourlyRate")},
		MerchantID: views.FormInputData{Value: merId, Error: valerr.Errors.MsgFor("MerchantID")},
	}
}

func mapEditUserFormToUser(id string, r *http.Request) (*user.User, error) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	username := r.FormValue("username")
	role := r.FormValue("role")
	merId := r.FormValue("merchant")
	hourlyRate := r.FormValue("hourly_rate")

	usr := &user.User{
		ID:       id,
		Name:     user.Name{First: firstName, Last: lastName},
		Email:    email,
		Username: username,
		RoleID:   role,
	}

	isCraftsman := merId != "" || hourlyRate != ""
	if isCraftsman {
		hourlyRate, err := strconv.ParseFloat(hourlyRate, 64)
		if err != nil {
			return nil, errs.ErrInvalidFloat
		}

		usr.Craftsman = &user.Craftsman{
			MerchantID: merId,
			HourlyRate: hourlyRate,
		}
	}

	return usr, nil
}

func (h *handler) responseWithEditUser(w http.ResponseWriter, r *http.Request, status int, usr *user.User, formdata *views.CreateUserFormData, withCraftsmanInfo bool) error {
	opts := &views.EditUserOpts{WithCraftsmanInfo: withCraftsmanInfo}
	roles, err := h.app.RoleService.GetRoles()
	if err != nil {
		return err
	}
	if withCraftsmanInfo {
		merchants, err := h.app.MerchantService.GetMerchants()
		if err != nil {
			return err
		}
		opts.Merchants = merchants
	}
	return respondWithTemplate(w, r, status, views.EditUser(usr, roles, formdata, opts))
}
