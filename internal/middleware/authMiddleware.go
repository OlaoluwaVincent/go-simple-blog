package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/olaoluwavincent/full-course/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
