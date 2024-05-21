package resthandlers

import (
	"errors"
	"net/http"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
	"golang.org/x/crypto/bcrypt"
)

func newLoginFormData(usr *user.User, valerr *errs.ValidationError) *views.LoginFormData {
	return &views.LoginFormData{
		Email:    views.FormInputData{Value: usr.Email, Error: valerr.Errors.MsgFor("Email")},
		Password: views.FormInputData{Value: usr.Password, Error: valerr.Errors.MsgFor("Password")},
	}
}

func newRegisterFormData(usr *user.User, valerr *errs.ValidationError) *views.RegisterFormData {
	return &views.RegisterFormData{
		FirstName:       views.FormInputData{Value: usr.Name.First, Error: valerr.Errors.MsgFor("Name.First")},
		LastName:        views.FormInputData{Value: usr.Name.Last, Error: valerr.Errors.MsgFor("Name.Last")},
		Username:        views.FormInputData{Value: usr.Username, Error: valerr.Errors.MsgFor("Username")},
		Email:           views.FormInputData{Value: usr.Email, Error: valerr.Errors.MsgFor("Email")},
		Password:        views.FormInputData{Value: usr.Password, Error: valerr.Errors.MsgFor("Password")},
		ConfirmPassword: views.FormInputData{Value: usr.Password, Error: valerr.Errors.MsgFor("ConfirmPassword")},
	}
}

func (h *handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	respondWithTemplate(w, r, http.StatusOK, views.Login(newLoginFormData(&user.User{}, &errs.ValidationError{})))
}

func (h *handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	respondWithTemplate(w, r, http.StatusOK, views.Register(newRegisterFormData(&user.User{}, &errs.ValidationError{})))
}

func (h *handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	usrform := &user.User{
		Name:            user.Name{First: r.FormValue("first_name"), Last: r.FormValue("last_name")},
		Username:        r.FormValue("username"),
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("cpassword"),
	}

	err := h.app.UserService.CreateUser(usrform)
	if err == nil {
		cookiesPair, err := h.jwtAdapter.GenTokenPairInCookie(usrform)
		if err != nil {
			respondWithInternalServerError(w, r)
			return
		}

		http.SetCookie(w, cookiesPair.Access)
		http.SetCookie(w, cookiesPair.Refresh)

		hxRespondWithRedirect(w, "/")
		return
	}

	if valerr, ok := err.(errs.ValidationError); ok {
		e := newRegisterFormData(usrform, &valerr)
		respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.RegisterForm(e))
		return
	}

	emailExists := errors.Is(err, user.ErrEmailAlreadyExists)
	usernameExists := errors.Is(err, user.ErrUsernameAlreadyExists)

	if emailExists || usernameExists {
		e := newRegisterFormData(usrform, &errs.ValidationError{})
		if emailExists {
			e.Email.Error = "Email already exists, try another one"
		}
		if usernameExists {
			e.Username.Error = "Username already exists, try another one"
		}
		respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.RegisterForm(e))
		return
	}

	respondWithTemplate(w, r, http.StatusInternalServerError, views.RegisterForm(newRegisterFormData(usrform, &errs.ValidationError{})))
}

func (h *handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	emailOrPassword := r.FormValue("email_or_username")
	usrform := &user.User{
		Email:    emailOrPassword,
		Password: r.FormValue("password"),
	}

	usr, err := h.app.UserService.FindUserByEmailOrUsername(emailOrPassword)
	if err != nil {
		e := newLoginFormData(usrform, &errs.ValidationError{})
		e.Email.Error = "Invalid email or username"
		respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.LoginForm(e))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(usrform.Password))
	if err != nil {
		e := newLoginFormData(usrform, &errs.ValidationError{})
		e.Password.Error = "Invalid password"
		respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.LoginForm(e))
		return
	}

	cookiesPair, err := h.jwtAdapter.GenTokenPairInCookie(usr)
	if err != nil {
		respondWithInternalServerError(w, r)
		return
	}
	http.SetCookie(w, cookiesPair.Access)
	http.SetCookie(w, cookiesPair.Refresh)

	hxRespondWithRedirect(w, "/")
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	emptyAccessCookie := h.jwtAdapter.NewCookie(jwtadapter.AccessToken, "", time.Unix(0, 0))
	emptyRefreshCookie := h.jwtAdapter.NewCookie(jwtadapter.RefreshToken, "", time.Unix(0, 0))

	http.SetCookie(w, emptyAccessCookie)
	http.SetCookie(w, emptyRefreshCookie)

	hxRespondWithRedirect(w, "/")
}

func (h *handler) refreshTokens(w http.ResponseWriter, r *http.Request) *jwtadapter.CookiePair {
	cookie, err := r.Cookie(jwtadapter.RefreshTokenCookieName)
	if err != nil || cookie.Value == "" {
		return nil
	}

	parsed, err := h.jwtAdapter.ParseRefreshClaims(cookie.Value)
	if err != nil {
		return nil
	}

	usr, err := h.app.UserService.FindUser(parsed.ID)
	if err != nil {
		return nil
	}

	cookiesPair, err := h.jwtAdapter.GenTokenPairInCookie(usr)
	if err != nil {
		return nil
	}

	http.SetCookie(w, cookiesPair.Access)
	http.SetCookie(w, cookiesPair.Refresh)

	return cookiesPair
}

func (h *handler) getMe(r *http.Request) (*user.User, error) {
	access, err := h.getAuthFromContext(r)
	if errors.Is(err, ErrNoAccessCookie) {
		return nil, err
	}

	me, err := h.app.UserService.FindUser(access.ID)
	if errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	return me, nil
}
