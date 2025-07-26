package handler

import (
	"errors"
	"net/http"
	"time"

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

func (h *handler) GetLogin(w http.ResponseWriter, r *http.Request) error {
	return respondWithTemplate(w, r, http.StatusOK, views.Login(mapLoginToFormData(&user.User{}, &errs.ValidationError{})))
}

func (h *handler) GetRegister(w http.ResponseWriter, r *http.Request) error {
	return respondWithTemplate(w, r, http.StatusOK, views.Register(mapRegisterToFormData(&user.User{}, &errs.ValidationError{})))
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) error {
	l := logger.FromCtx(r.Context())

	usr := mapFormToUser(r)

	err := h.app.UserService.CreateUser(usr)
	if err != nil {
		l.Warn("Error creating user", zap.Error(err), zap.Any("user", usr))
		if valerr, ok := err.(errs.ValidationError); ok {
			l.Warn("Validation errors", zap.Any("valerr", valerr.Errors))
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.RegisterForm(mapRegisterToFormData(usr, &valerr)))
		}

		emailExists := errors.Is(err, user.ErrEmailAlreadyExists)
		usernameExists := errors.Is(err, user.ErrUsernameAlreadyExists)

		if emailExists || usernameExists {
			l.Warn("Existing username or email", zap.Error(err))
			e := mapRegisterToFormData(usr, &errs.ValidationError{})
			if emailExists {
				e.Email.Error = "Email already exists, try another one"
			}
			if usernameExists {
				e.Username.Error = "Username already exists, try another one"
			}
			return respondWithTemplate(w, r, http.StatusConflict, views.RegisterForm(e))
		}

		return err
	}

	l.Info("Created the user", zap.Any("user", usr))

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return err
	}

	http.SetCookie(w, cookiesPair.Access)
	http.SetCookie(w, cookiesPair.Refresh)

	return hxRespondWithRedirect(w, "/")
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) error {
	var err error

	emailOrUsername := r.FormValue("email_or_username")
	password := r.FormValue("password")
	inpUser := &user.User{Email: emailOrUsername, Password: password}

	usr, err := h.app.UserService.GetUserByEmailOrUsername(emailOrUsername)
	if err != nil {
		e := mapLoginToFormData(inpUser, &errs.ValidationError{})
		e.Email.Error = "Invalid email or username"
		return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.LoginForm(e))
	}

	// TODO: move this to the business logic
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		e := mapLoginToFormData(inpUser, &errs.ValidationError{})
		e.Password.Error = "Invalid password"
		return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.LoginForm(e))
	}

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return err
	}

	http.SetCookie(w, cookiesPair.Refresh)
	http.SetCookie(w, cookiesPair.Access)

	return hxRespondWithRedirect(w, "/")
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) error {
	unsetCookie(w, refreshTokenCookieName)
	unsetCookie(w, accessTokenCookieName)

	return hxRespondWithRedirect(w, "/")
}

func (h *handler) newCookiesPairFromUser(usr *user.User) (*cookiePair, error) {
	tokens, err := h.jwtAdapter.NewPair(usr)
	if err != nil {
		return nil, err
	}

	refreshCookie := newCookie(refreshTokenCookieName, tokens.Refresh.Encoded, tokens.Refresh.Expiration)
	accessCookie := newCookie(accessTokenCookieName, tokens.Access.Encoded, tokens.Access.Expiration)

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

func mapFormToUser(r *http.Request) *user.User {
	return &user.User{
		Name:            user.Name{First: r.FormValue("first_name"), Last: r.FormValue("last_name")},
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("cpassword"),
	}
}
