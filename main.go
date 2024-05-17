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
	// account setup
	router.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))
	router.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))

	router.Group(func(auth chi.Router) {
		// make sure user is has an account
		auth.Use(handler.WithAccountSetup)
		// home page
		auth.Get("/", handler.Make(handler.HandleHomeIndex))
		// settings
		auth.Get("/settings", handler.Make(handler.HandleSettingsIndex))
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
