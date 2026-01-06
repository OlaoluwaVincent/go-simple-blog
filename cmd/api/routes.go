package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/olaoluwavincent/full-course/internal/middleware"
)

func appRoutes(r *chi.Mux, app *application) http.Handler {
	userController := NewUserController(app)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		// Public User routes (no auth required)
		r.Route("/users", func(r chi.Router) {
			// Public routes
			r.Post("/register", userController.registerHandler)
			r.Post("/login", userController.loginHandler)

			// Protected routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.AuthMiddleware) // JWT required
				r.Put("/update", userController.updateHandler)
				r.Get("/profile/me", userController.getMeHandler)
				r.Get("/profile/{id}", userController.userDetailsHandler)
			})
		})
	})

	return r
}
