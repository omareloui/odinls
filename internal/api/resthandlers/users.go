package resthandlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.app.UserService.GetUsers(user.WithPopulatedRole)
	if err != nil {
		respondWithInternalServerError(w, r)
		return
	}
	accessClaims, _ := h.getAuthFromContext(r)

	respondWithTemplate(w, r, http.StatusOK, views.UserPage(accessClaims, users))
}

func (h *handler) GetUser(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := h.app.UserService.FindUser(id, user.WithPopulatedRole)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, err = w.Write([]byte(err.Error()))
				if err != nil {
					respondWithInternalServerError(w, r)
					return
				}
			}
			respondWithInternalServerError(w, r)
			return
		}

		respondWithTemplate(w, r, http.StatusOK, views.User(usr))
	}
}

func (h *handler) GetEditUser(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := h.app.UserService.FindUser(id, user.WithPopulatedRole)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			respondWithInternalServerError(w, r)
			return
		}

		roles, err := h.app.RoleService.GetRoles()
		if err != nil {
			respondWithInternalServerError(w, r)
			return
		}

		opts := &views.EditUserOpts{}
		if usr.IsCraftsman() {
			opts.WithCraftsmanInfo = true
			merchants, _ := h.app.MerchantService.GetMerchants()
			opts.Merchants = merchants
		}

		respondWithTemplate(w, r, http.StatusOK, views.EditUser(usr, roles,
			newEditUserFormData(usr, &errs.ValidationError{}), opts))
	}
}

func (h *handler) EditUser(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				merchants, _ := h.app.MerchantService.GetMerchants()
				roles, _ := h.app.RoleService.GetRoles()
				data := newEditUserFormData(usr, &errs.ValidationError{})
				data.HourlyRate.Error = "Invalid number"
				respondWithTemplate(w, r, http.StatusUnprocessableEntity,
					views.EditUser(usr, roles, data, &views.EditUserOpts{WithCraftsmanInfo: true, Merchants: merchants}))
				return
			}

			usr.Craftsman = &user.Craftsman{
				MerchantID: merId,
				HourlyRate: hourlyRate,
			}
		}

		err := h.app.UserService.UpdateUserByID(id, usr, user.WithPopulatedRole)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				roles, _ := h.app.RoleService.GetRoles()
				if isCraftsman {
					merchants, _ := h.app.MerchantService.GetMerchants()
					respondWithTemplate(w, r, http.StatusUnprocessableEntity,
						views.EditUser(usr, roles, newEditUserFormData(usr, &valerr),
							&views.EditUserOpts{WithCraftsmanInfo: true, Merchants: merchants}))
				} else {
					respondWithTemplate(w, r, http.StatusUnprocessableEntity,
						views.EditUser(usr, roles, newEditUserFormData(usr, &valerr)))
				}
				return
			}

			emailExists := errors.Is(err, user.ErrEmailAlreadyExists)
			usernameExists := errors.Is(err, user.ErrUsernameAlreadyExists)

			if emailExists || usernameExists {
				e := newEditUserFormData(usr, &errs.ValidationError{})
				if emailExists {
					e.Email.Error = "Email already exists, try another one"
				}
				if usernameExists {
					e.Username.Error = "Username already exists, try another one"
				}
				roles, _ := h.app.RoleService.GetRoles()
				respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditUser(usr, roles, e))
				return
			}

			respondWithInternalServerError(w, r)
			return
		}

		respondWithTemplate(w, r, http.StatusOK, views.User(usr))
	}
}

func (h *handler) UnsetCraftsman(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.app.UserService.UnsetCraftsmanByID(id)
		if err != nil {
			if errors.Is(user.ErrUserNotFound, err) {
				respondWithNotFound(w, r)
				return
			}
			respondWithInternalServerError(w, r)
			return
		}

		user, err := h.app.UserService.FindUser(id, user.WithPopulatedRole)
		if err != nil {
			respondWithInternalServerError(w, r)
			return
		}
		respondWithTemplate(w, r, http.StatusOK, views.User(user))
	}
}

func (h *handler) GetCraftsmanForm(w http.ResponseWriter, r *http.Request) {
	merchants, err := h.app.MerchantService.GetMerchants()
	if err != nil {
		respondWithInternalServerError(w, r)
		return
	}
	respondWithTemplate(w, r, http.StatusOK, views.CraftsmanForm(merchants, &views.CreateUserFormData{}))
}

func newEditUserFormData(user *user.User, valerr *errs.ValidationError) *views.CreateUserFormData {
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
