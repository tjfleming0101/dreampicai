package handler

import (
	"log/slog"
	"net/http"

	"github.com/nedpals/supabase-go"
	"github.com/tjfleming0101/dreampicai/pkg/sb"
	"github.com/tjfleming0101/dreampicai/pkg/validate"
	"github.com/tjfleming0101/dreampicai/view/auth"
)

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Signup())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email: r.FormValue("email"),
		Password: r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	errors := auth.SignupErrors{}

	if ok := validate.New(&params, validate.Fields{
		"Email": validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.Password), 
			validate.Message("Passwords do not match"),
		),		
	}).Validate(&errors); !ok {
		return render(w, r, auth.SignupForm(params, errors))
	}
	
	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email: params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return render(w, r, auth.SingupSuccess(user.Email))

}

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(w, r, auth.Login())
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email: r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(w, r, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))
	}

	setAuthCookie(w, resp.AccessToken)

	return hxRedirect(w, r, "/")
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider: "google",
		RedirectTo: "http://localhost:3000/auth/callback",
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	cookie := http.Cookie{
		Value: "",
		Name: "at",
		MaxAge: -1,
		HttpOnly: true,
		Path: "/",
		Secure: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(w, r, auth.CallbackScript())
	}
	setAuthCookie(w, accessToken)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Value: accessToken,
		Name: "at",
		Path: "/",
		HttpOnly: true,
		Secure: true,
	}
	http.SetCookie(w, cookie)
}