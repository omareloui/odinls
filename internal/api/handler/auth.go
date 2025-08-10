package handler

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/omareloui/former"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/api/middleware"
	"github.com/omareloui/odinls/internal/api/responder"
	"github.com/omareloui/odinls/internal/application/core/user"
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
	fd := new(views.LoginFormData)
	h.fm.MapToForm(new(user.User), nil, fd)
	comp := views.Login(fd)
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) GetRegister(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	fd := new(views.RegisterFormData)
	h.fm.MapToForm(new(user.User), nil, fd)
	comp := views.Register(fd)
	return responder.OK(responder.WithComponent(comp))
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	l := logger.FromCtx(r.Context())

	usr := new(user.User)
	err := former.Populate(r, usr)
	if err != nil {
		return responder.BadRequest()
	}

	usr, err = h.app.UserService.CreateUser(usr)
	if err != nil {
		vfd := new(views.RegisterFormData)
		h.fm.MapToForm(usr, err, vfd)
		compIfVarErr := views.RegisterForm(vfd)

		cfd := new(views.RegisterFormData)
		h.fm.MapToForm(usr, nil, cfd)
		cfd.Email.Error = "Email or Username already exists, try another one"
		cfd.Username.Error = "Email or Username already exists, try another one"
		compIfAlreadyExists := views.RegisterForm(cfd)

		return responder.Error(err,
			responder.WithComponentIfValidationErr(compIfVarErr),
			responder.WithComponentIfErrIs(err, compIfAlreadyExists),
		)
	}

	l.Info("Created the user", zap.Any("user", usr))

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return responder.Error(err)
	}

	http.SetCookie(w, cookiesPair.Access)
	http.SetCookie(w, cookiesPair.Refresh)

	return responder.RedirectHX(w, responder.WithPath("/"))
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	var err error

	emailOrUsername := r.FormValue("email_or_username")
	password := r.FormValue("password")
	inpUser := &user.User{Email: emailOrUsername, Password: password}

	usr, err := h.app.UserService.GetUserByEmailOrUsername(emailOrUsername)
	if err != nil {
		fd := new(views.LoginFormData)
		h.fm.MapToForm(inpUser, nil, fd)
		fd.Email.Error = "Invalid email or username"
		comp := views.LoginForm(fd)
		return responder.UnprocessableEntity(responder.WithComponent(comp))
	}

	// TODO: move this to the business logic
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		fd := new(views.LoginFormData)
		h.fm.MapToForm(inpUser, nil, fd)
		fd.Password.Error = "Invalid password"
		comp := views.LoginForm(fd)
		return responder.UnprocessableEntity(responder.WithComponent(comp))
	}

	cookiesPair, err := h.newCookiesPairFromUser(usr)
	if err != nil {
		return responder.Error(err)
	}

	http.SetCookie(w, cookiesPair.Refresh)
	http.SetCookie(w, cookiesPair.Access)

	return responder.RedirectHX(w, responder.WithPath("/"))
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	unsetCookie(w, middleware.AccessClaimsCookieName)
	unsetCookie(w, middleware.RefreshClaimsCookieName)

	return responder.RedirectHX(w, responder.WithPath("/"))
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
