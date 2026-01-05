package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/olaoluwavincent/full-course/internal/store"
	"github.com/olaoluwavincent/full-course/internal/utils"
)

type UserService struct {
	store store.Storage
}

// NewUserService creates a new instance of UserService
func NewUserService(store store.Storage) *UserService {
	return &UserService{
		store: store,
	}
}

func (us *UserService) RegisterUser(req utils.RegisterRequest, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, _ := us.store.Users.GetByIDUsernameOrEmail(ctx, nil, nil, &req.Email)
	if user != nil {
		http.Error(w, "Email already exists", http.StatusBadRequest)
		return
	}

	user, _ = us.store.Users.GetByIDUsernameOrEmail(ctx, nil, &req.Username, nil)
	if user != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hashed_password, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Account creation failed, please try again later", http.StatusBadRequest)
		return
	}

	user, err = us.store.Users.Create(ctx, &store.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed_password,
	})

	if err != nil {
		http.Error(w, "Account creation failed, please try again later", http.StatusBadRequest)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Response
	response := map[string]any{
		"token": token,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (us *UserService) LoginUser(req utils.LoginRequest, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := us.store.Users.GetByIDUsernameOrEmail(ctx, nil, nil, &req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"token": token,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (us *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	user, err := us.store.Users.GetByIDUsernameOrEmail(ctx, &userID, nil, nil)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (us *UserService) UpdateUser(ctx context.Context, user *store.User) error {
	return us.store.Users.Update(ctx, user)
}
