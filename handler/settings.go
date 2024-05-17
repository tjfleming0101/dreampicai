package handler

import (
	"net/http"

	"github.com/tjfleming0101/dreampicai/db"
	"github.com/tjfleming0101/dreampicai/pkg/validate"
	"github.com/tjfleming0101/dreampicai/view/settings"
)

func HandleSettingsIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	return render(w, r, settings.Profile(user))
}

func HandleSettingsUsernameUpdate(w http.ResponseWriter, r *http.Request) error {
	params := settings.ProfileParams{
		Username: r.FormValue("username"),
	}

	errors := settings.ProfileErrors{}

	if ok := validate.New(&params, validate.Fields{
		"Username": validate.Rules(validate.Min(2), validate.Max(25)),
	}).Validate(&errors); !ok {
		return render(w, r, settings.ProfileForm(params, errors))
	}

	user := getAuthenticatedUser(r)
	user.Account.Username = params.Username
	if err := db.UpdateAccount(&user.Account); err != nil {
		return err
	}
	params.Success = true
	return render(w, r, settings.ProfileForm(params, settings.ProfileErrors{}))
}
