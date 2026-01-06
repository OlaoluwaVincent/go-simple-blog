package utils

type RegisterRequest struct {
	Username string `json:"username" form:"username" bind:"required"`
	Email    string `json:"email" form:"email" bind:"required"`
	Password string `json:"password" form:"password" bind:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" bind:"required"`
	Password string `json:"password" form:"password" bind:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
