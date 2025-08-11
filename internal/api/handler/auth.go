package handler

import (
	"net/http"
	"seno-blackdragon/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"require,email"`
	Password string `json:"password" binding:"require"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Expires      int    `json:"expires"`
}

// Login godoc
// @Summary     Login
// @Description  Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      LoginRequest  true  "Info Login"
// @Success      200   {object}  LoginResponse
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login Fail"})
	}
	c.JSON(http.StatusOK, LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}
