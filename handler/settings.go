package handler

import (
	"net/http"

	"github.com/tjfleming0101/dreampicai/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(w, r, settings.Index(user))
}
