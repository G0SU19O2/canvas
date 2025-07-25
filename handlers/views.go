package handlers

import (
	"canvas/views"
	"net/http"

	"github.com/go-chi/chi"
)

func FrontPage(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_ = views.FrontPage().Render(w)
	})
}
