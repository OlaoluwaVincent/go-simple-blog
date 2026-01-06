package main

import "github.com/olaoluwavincent/full-course/internal/services"

type PostController struct {
	pc          *application
	userService *services.PostService
}

func NewPostController(uc *application) *PostController {
	userService := services.NewPostService(uc.store)
	return &PostController{pc: uc, userService: userService}
}
