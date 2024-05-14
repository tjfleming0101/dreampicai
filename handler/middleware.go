package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/tjfleming0101/dreampicai/types"
)

// find the user, see if the user is logged in or not
func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}
		user := types.AuthenticatedUser{}
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

