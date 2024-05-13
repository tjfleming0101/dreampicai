package handler

import (
	"net/http"

	"github.com/tjfleming0101/dreampicai/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}