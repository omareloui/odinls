package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/api/middleware"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/logger"
	"github.com/omareloui/odinls/web/views"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type cookiePair struct {
	Access  *http.Cookie
	Refresh *http.Cookie
}

func (h *handler) GetLogin(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	comp := views.Login(mapLoginToFormData(&user.User{}, &errs.ValidationError{}))
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) GetRegister(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	comp := views.Register(mapRegisterToFormData(&user.User{}, &errs.ValidationError{}))
	return RespondOK(w, RespondWithComponent(comp))
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	l := logger.FromCtx(r.Context())

	usr := new(user.User)
	err := former.Populate(r, usr)
	if err != nil {
		return BadRequest()
	}

	usr, err = h.app.UserService.CreateUser(usr)
	if err != nil {
		l.Warn("Error creating user", zap.Error(err), zap.Any("user", usr))
		if valerr, ok := err.(errs.ValidationError); ok {
			l.Warn("Validation errors", zap.Any("valerr", valerr.Errors))
			comp := views.RegisterForm(mapRegisterToFormData(usr, &valerr))
			return UnprocessableEntity(RespondWithComponent(comp))
		}

		alreadyExists := errors.Is(err, errs.ErrDocumentAlreadyExists)

		if alreadyExists {
			l.Warn("Existing username or email", zap.Error(err))
			e := mapRegisterToFormData(usr, &errs.ValidationError{})
			e.Email.Error = "Email or Username already exists, try another one"
			e.Username.Error = "Email or Username already exists, try another one"
			comp := views.RegisterForm(e)
			return UnprocessableEntity(RespondWithComponent(comp))
		}

		return RespondError(err)
	}

	l.Info("Created the user", zap.Any("user", usr))

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return RespondError(err)
	}

	http.SetCookie(w, cookiesPair.Access)
	http.SetCookie(w, cookiesPair.Refresh)

	return RespondRedirectHX(w, RespondWithPath("/"))
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var err error

	emailOrUsername := r.FormValue("email_or_username")
	password := r.FormValue("password")
	inpUser := &user.User{Email: emailOrUsername, Password: password}

	usr, err := h.app.UserService.GetUserByEmailOrUsername(emailOrUsername)
	if err != nil {
		e := mapLoginToFormData(inpUser, &errs.ValidationError{})
		e.Email.Error = "Invalid email or username"
		comp := views.LoginForm(e)
		return UnprocessableEntity(RespondWithComponent(comp))
	}

	// TODO: move this to the business logic?
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		e := mapLoginToFormData(inpUser, &errs.ValidationError{})
		e.Password.Error = "Invalid password"
		comp := views.LoginForm(e)
		return UnprocessableEntity(RespondWithComponent(comp))
	}

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return RespondError(err)
	}

	http.SetCookie(w, cookiesPair.Refresh)
	http.SetCookie(w, cookiesPair.Access)

	return RespondRedirectHX(w, RespondWithPath("/"))
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	unsetCookie(w, middleware.AccessClaimsCookieName)
	unsetCookie(w, middleware.RefreshClaimsCookieName)

	return RespondRedirectHX(w, RespondWithPath("/"))
}

func (h *handler) newCookiesPairFromUser(usr *user.User) (*cookiePair, error) {
	tokens, err := jwtadapter.NewPair(usr)
	if err != nil {
		return nil, err
	}

	refreshCookie := newCookie(middleware.RefreshClaimsCookieName, tokens.Refresh.Encoded, tokens.Refresh.Expiration)
	accessCookie := newCookie(middleware.AccessClaimsCookieName, tokens.Access.Encoded, tokens.Access.Expiration)

	return &cookiePair{Refresh: refreshCookie, Access: accessCookie}, nil
}

func newCookie(name, value string, exp time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Expires:  exp,
		Path:     "/",
	}
}

func unsetCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
	})
}

func mapLoginToFormData(usr *user.User, valerr *errs.ValidationError) *views.LoginFormData {
	return &views.LoginFormData{
		Email:    views.FormInputData{Value: usr.Email, Error: valerr.Errors.MsgFor("Email")},
		Password: views.FormInputData{Value: usr.Password, Error: valerr.Errors.MsgFor("Password")},
	}
}

func mapRegisterToFormData(usr *user.User, valerr *errs.ValidationError) *views.RegisterFormData {
	return &views.RegisterFormData{
		FirstName:       views.FormInputData{Value: usr.Name.First, Error: valerr.Errors.MsgFor("Name.First")},
		LastName:        views.FormInputData{Value: usr.Name.Last, Error: valerr.Errors.MsgFor("Name.Last")},
		Username:        views.FormInputData{Value: usr.Username, Error: valerr.Errors.MsgFor("Username")},
		Email:           views.FormInputData{Value: usr.Email, Error: valerr.Errors.MsgFor("Email")},
		Password:        views.FormInputData{Value: usr.Password, Error: valerr.Errors.MsgFor("Password")},
		ConfirmPassword: views.FormInputData{Value: usr.Password, Error: valerr.Errors.MsgFor("ConfirmPassword")},
	}
}
