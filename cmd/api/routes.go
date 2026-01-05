package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func appRoutes(r *chi.Mux, app *application) http.Handler {
	userController := NewUserController(app)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// User routes - grouped nicely
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", userController.registerHandler)
			r.Post("/login", userController.loginHandler)
			r.Put("/update", userController.updateHandler)
			r.Get("/profile/{id}", userController.userDetailsHandler)
		})
	})

	return r
}
