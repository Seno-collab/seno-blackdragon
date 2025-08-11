package handler

import (
	"net/http"
	"time"

	"seno-blackdragon/internal/service"
	"seno-blackdragon/pkg/dto"
	"seno-blackdragon/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // e.g. "Bearer"
	Expires      int    `json:"expires"`    // seconds until access token expires
}

type LoginSuccess = dto.BaseResponse[LoginResponse]

// // Login godoc
// // @Summary      Login
// // @Description  Login user
// // @Tags         auth
// // @Accept       json
// // @Produce      json
// // @Param        data  body      LoginRequest  true  "Request Login"
// // @Success      200   {object}  LoginSuccess
// // @Failure      400   {object}  dto.ErrorResponse
// // @Failure      401   {object}  dto.ErrorResponse
// // @Router       /api/v1/auth/login [post]
// func (h *AuthHandler) Login(c *gin.Context) {
// 	reqTime := time.Now().UTC()
// 	traceID := c.GetString(middleware.ContextKeyTraceID)
// 	if traceID == "" {
// 		traceID = c.GetHeader(middleware.HeaderKeyTraceID)
// 	}

// 	var req LoginRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, "VALIDATION_FAILED",
// 			"Invalid login payload", traceID, reqTime, err))
// 		return
// 	}

// 	accessToken, refreshToken, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
// 	if err != nil {
// 		dto.WriteJSON(c, http.StatusUnauthorized, dto.NewError(http.StatusUnauthorized, "AUTH_FAILED",
// 			"Invalid email or password", traceID, reqTime, err))
// 		return
// 	}

// 	// Nếu có TTL từ service, set vào Expires; nếu chưa có, tạm để 0.
// 	resp := LoginResponse{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 		TokenType:    "Bearer",
// 		Expires:      0,
// 	}
// 	dto.Ok(c, dto.NewSuccess(http.StatusOK, "Login success", traceID, resp, reqTime))
// }

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"omitempty"`
	Bio      string `json:"bio" binding:"omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterSuccess = dto.BaseResponse[dto.EmptyData]

// Register godoc
// @Summary      Register
// @Description  Register user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      RegisterRequest  true  "Register User"
// @Success      200   {object}  RegisterSuccess
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	reqTime := time.Now().UTC()
	traceID := c.GetString(middleware.ContextKeyTraceID)
	if traceID == "" {
		traceID = c.GetHeader(middleware.HeaderKeyTraceID)
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, "VALIDATION_FAILED",
			"Invalid register payload", traceID, reqTime, err))
		return
	}
	_, err := h.authService.Register(c.Request.Context(), req.FullName, req.Bio, req.Email, req.Password)
	if err != nil {
		dto.WriteJSON(c, http.StatusBadRequest, dto.NewError(http.StatusBadRequest, "BUSINESS_ERROR",
			"Register user failed", traceID, reqTime, err))
		return
	}

	dto.Ok(c, dto.NewSuccessEmpty(http.StatusOK, "Create user success", traceID, reqTime))
}
