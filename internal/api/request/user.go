package request

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Bio      string `json:"bio"`
	Password string `json:"password"  binding:"required"`
}
