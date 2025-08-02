package request

type LoginRequest struct {
	Email    string `json:"email" binding:"require,email"`
	Password string `json:"password" binding:"require"`
}
