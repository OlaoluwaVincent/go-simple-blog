package main

import (
	"encoding/json"
	"net/http"

	"github.com/olaoluwavincent/full-course/internal/services"
	"github.com/olaoluwavincent/full-course/internal/utils"
)

type UserController struct {
	uc          *application
	userService *services.UserService
}

func NewUserController(uc *application) *UserController {
	userService := services.NewUserService(uc.store)
	return &UserController{uc: uc, userService: userService}
}

func (uc *UserController) registerHandler(w http.ResponseWriter, r *http.Request) {

	var req utils.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	uc.userService.RegisterUser(req, w, r)
}

func (uc *UserController) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req utils.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	uc.userService.LoginUser(req, w, r)
}

func (uc *UserController) updateHandler(w http.ResponseWriter, r *http.Request) {
	uc.userService.UpdateUser(w, r)
}

func (uc *UserController) userDetailsHandler(w http.ResponseWriter, r *http.Request) {
	uc.userService.GetUser(w, r)
}

func (uc *UserController) getMeHandler(w http.ResponseWriter, r *http.Request) {
	uc.userService.GetMe(w, r)
}
