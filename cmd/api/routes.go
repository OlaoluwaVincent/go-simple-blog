package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func appRoutes(r *chi.Mux, app *application) http.Handler {
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// r.route("")
	})

	return r
}
