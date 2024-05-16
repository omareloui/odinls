package resthandlers

import (
	"net/http"

	"github.com/omareloui/odinls/web/views"
)

func newLoginFormData() *views.LoginFormData {
	return &views.LoginFormData{
		Email:    views.FormInputData{},
		Password: views.FormInputData{},
	}
}

func newRegisterFormData() *views.RegisterFormData {
	return &views.RegisterFormData{
		FristName:       views.FormInputData{},
		LastName:        views.FormInputData{},
		Email:           views.FormInputData{},
		Password:        views.FormInputData{},
		ConfirmPassword: views.FormInputData{},
	}
}

func (h *handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	respondWithTemplate(w, r, http.StatusOK, views.Login(newLoginFormData()))
}

func (h *handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	respondWithTemplate(w, r, http.StatusOK, views.Register(newRegisterFormData()))
}
