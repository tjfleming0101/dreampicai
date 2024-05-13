package handler

import (
	"fmt"
	"net/http"
	"strings"
)

// find the user, see if the user is logged in or not
func WithUser(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("from the with user middleware")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}