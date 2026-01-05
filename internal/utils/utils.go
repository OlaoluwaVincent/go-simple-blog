package utils

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Cost 12-14 is good (higher = slower/more secure, but ~100-300ms is acceptable)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetPathParam(r *http.Request, index int) string {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if index < 0 || index >= len(parts) {
		return ""
	}

	return parts[index]
}

func GetQueryParam(r *http.Request, key string) string {
	values := r.URL.Query()
	return values.Get(key)
}
