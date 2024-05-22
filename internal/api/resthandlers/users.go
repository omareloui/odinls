package resthandlers

import (
	"errors"
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.app.UserService.GetUsers()
	if err != nil {
		respondWithInternalServerError(w, r)
		return
	}
	accessClaims, _ := h.getAuthFromContext(r)

	respondWithTemplate(w, r, http.StatusOK, views.UserPage(accessClaims, users))
}

func (h *handler) GetUser(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, err := h.app.UserService.FindUser(id)
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
		usr, err := h.app.UserService.FindUser(id)
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

		respondWithTemplate(w, r, http.StatusOK, views.EditUser(usr, newEditUserFormData(usr, &errs.ValidationError{})))
	}
}

func (h *handler) EditUser(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		email := r.FormValue("email")
		username := r.FormValue("username")

		usr := &user.User{
			ID:       id,
			Name:     user.Name{First: firstName, Last: lastName},
			Email:    email,
			Username: username,
		}

		err := h.app.UserService.UpdateUserByID(id, usr)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditUser(usr, newEditUserFormData(usr, &valerr)))
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
				respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditUser(usr, e))
				return
			}

			respondWithInternalServerError(w, r)
			return
		}

		respondWithTemplate(w, r, http.StatusOK, views.User(usr))
	}
}

func newEditUserFormData(user *user.User, valerr *errs.ValidationError) *views.CreateUserFormData {
	return &views.CreateUserFormData{
		Name: views.NameFormData{
			First: views.FormInputData{Value: user.Name.First, Error: valerr.Errors.MsgFor("Name.First")},
			Last:  views.FormInputData{Value: user.Name.Last, Error: valerr.Errors.MsgFor("Name.Last")},
		},
		Email:    views.FormInputData{Value: user.Email, Error: valerr.Errors.MsgFor("Email")},
		Username: views.FormInputData{Value: user.Username, Error: valerr.Errors.MsgFor("Username")},
	}
}
