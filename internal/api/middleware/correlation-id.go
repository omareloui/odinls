package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid"
)

type correlationIDCtxKey struct{}

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
		correlationID := ulid.MustNew(ulid.Now(), entropy).String()

		ctx := context.WithValue(r.Context(), correlationIDCtxKey{}, correlationID)
		w.Header().Add("X-Correlation-ID", correlationID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
