package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/tjfleming0101/dreampicai/pkg/sb"
	"github.com/tjfleming0101/dreampicai/types"
)

// find the user, see if the user is logged in or not
func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}
		accessToken, err := r.Cookie("at")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}		
		resp, err := sb.Client.Auth.User(r.Context(), accessToken.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user := types.AuthenticatedUser{
			Email: resp.Email,
			LoggedIn: true,
		}		
		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

