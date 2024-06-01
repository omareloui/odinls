package resthandlers

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/a-h/templ"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

type ValidationErrorResponseFunc func(valerr *errs.ValidationError) templ.Component

var nonDigitRegexp *regexp.Regexp = regexp.MustCompile(`\D+`)

func respondWithTemplate(w http.ResponseWriter, r *http.Request, status int, template templ.Component) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	return renderToBody(w, r, template)
}

func respondWithString(w http.ResponseWriter, r *http.Request, status int, str string) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	if _, err := w.Write([]byte(str)); err != nil {
		return respondWithInternalServerError(w, r)
	}
	return nil
}

func respondWithInternalServerError(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusInternalServerError)
}

func respondWithUnauthorized(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusUnauthorized)
}

func respondWithForbidden(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusForbidden)
}

func respondWithNotFound(w http.ResponseWriter, r *http.Request) error {
	return respondWithErrorPage(w, r, http.StatusNotFound)
}

func respondWithErrorPage(w http.ResponseWriter, r *http.Request, status int) error {
	auth := r.Context().Value(authContextKey)
	return respondWithTemplate(w, r, status, views.ErrorPage(auth.(*jwtadapter.JwtAccessClaims), http.StatusText(status), status))
}

func renderToBody(w http.ResponseWriter, r *http.Request, template templ.Component) error {
	return template.Render(r.Context(), w)
}

func hxRespondWithRedirect(w http.ResponseWriter, path string) error {
	w.Header().Set("HX-Location", path)
	return nil
}

func parseIntIfExists(str string) (int, error) {
	if str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return 0, errs.ErrInvalidFloat
		}
		return num, nil
	}
	return 0, nil
}

func parseFloatIfExists(str string) (float64, error) {
	if str != "" {
		str = nonDigitRegexp.ReplaceAllLiteralString(str, "")
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0, errs.ErrInvalidFloat
		}
		return num, nil
	}
	return 0, nil
}

func parseDateOnlyIfExists(str string) (time.Time, error) {
	if str != "" {
		date, err := time.Parse(time.DateOnly, str)
		if err != nil {
			return time.Time{}, errs.ErrInvalidFloat
		}
		return date, nil
	}
	return time.Time{}, nil
}

func formatDateOnlyIfNonZero(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format(time.DateOnly)
}

func formatFloatIfNonZero(f float64) string {
	if f == 0 {
		return ""
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func formatBooleanIfNonZero(b bool) string {
	if !b {
		return ""
	}
	return "on"
}
