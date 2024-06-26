package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tjfleming0101/dreampicai/db"
	"github.com/tjfleming0101/dreampicai/handler"
	"github.com/tjfleming0101/dreampicai/pkg/sb"
)

//go:embed static
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(handler.WithUser)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	// home page
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	// login
	router.Get("/login", handler.Make(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.Make(handler.HandleLoginWithGoogle))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	// signUp
	router.Get("/signup", handler.Make(handler.HandleSignupIndex))
	router.Post("/signup", handler.Make(handler.HandleSignupCreate))
	// callback
	router.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	// logout
	router.Post("/logout", handler.Make(handler.HandleLogoutCreate))

	// This is to make sure users are authenticated and then checks if they have setup an account before accessing the Settings page
	router.Group(func(auth chi.Router) {
		// make sure user is an account setup
		auth.Use(handler.WithAuth, handler.WithAccountSetup)
		// settings
		auth.Get("/settings", handler.Make(handler.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handler.Make(handler.HandleSettingsUsernameUpdate))

	})

	// This is the route for making sure users are authenticated but has no account setup
	router.Group(func(account chi.Router) {
		// make sure user has an account
		account.Use(handler.WithAuth)
		// account setup
		account.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))
		account.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))

	})

	port := os.Getenv("PORT")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	if err := db.Init(); err != nil {
		return err
	}
	return sb.Init()
}
