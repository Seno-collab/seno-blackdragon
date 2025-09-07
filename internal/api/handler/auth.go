package handler

import (
	"net/http"
	"time"

	"seno-blackdragon/internal/model"
	"seno-blackdragon/internal/service"
	"seno-blackdragon/pkg/dto"
	"seno-blackdragon/pkg/enum"
	"seno-blackdragon/pkg/middleware"
	"seno-blackdragon/pkg/pass"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email      string            `json:"email" binding:"required,email"`
	Password   string            `json:"password" binding:"required"`
	DeviceID   string            `json:"device_id"`
	DeviceMeta map[string]string `json:"device_meta"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"` // e.g. "Bearer"
	Expires      int64  `json:"expires"`    // seconds until access token expires
}

type LoginSuccess = dto.BaseResponse[LoginResponse]

// @BasePath /api/v1
// Login godoc
// @Summary      Login
// @Description  Login user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      LoginRequest  true  "Request Login"
// @Success      200   {object}  LoginSuccess
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	reqTime := time.Now().UTC()
	traceID := c.GetString(middleware.ContextKeyTraceID)
	if traceID == "" {
		traceID = c.GetHeader(middleware.HeaderKeyTraceID)
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeInvalidToken,
			"Invalid login payload", traceID, reqTime, err))
		return
	}
	// TODO: write code verify password
	if err := pass.VerifyPassword(req.Password); err != nil {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeBadRequest, "Invalid login payload", traceID, reqTime, err))
		return
	}
	cmd := model.LoginCmd{
		Email:    req.Email,
		Password: req.Password,
		DeviceID: req.DeviceID,
		IP:       c.ClientIP(),
		UA:       c.GetHeader("User-Agent"),
	}
	token, err := h.authService.Login(c.Request.Context(), cmd)
	if err != nil {
		dto.WriteJSON(c, http.StatusUnauthorized, dto.NewError(http.StatusUnauthorized, enum.CodeAuth,
			"Invalid email or password", traceID, reqTime, err))
		return
	}
	resp := LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		Expires:      token.Expired,
	}
	dto.Ok(c, dto.NewSuccess(http.StatusOK, "Login success", traceID, resp, reqTime))
}

type RegisterRequest struct {
	FullName string `json:"full_name" binding:"omitempty"`
	Bio      string `json:"bio" binding:"omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterSuccess = dto.BaseResponse[dto.EmptyData]

// @BasePath /api/v1
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
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeInvalidToken,
			"Invalid register payload", traceID, reqTime, err))
		return
	}
	// TODO: write code verify password
	if err := pass.VerifyPassword(req.Password); err != nil {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeBadRequest, "Invalid login payload", traceID, reqTime, err))
		return
	}
	_, err := h.authService.Register(c.Request.Context(), req.FullName, req.Bio, req.Email, req.Password)
	if err != nil {
		dto.WriteJSON(c, http.StatusBadRequest, dto.NewError(http.StatusBadRequest, enum.CodeBusinessError,
			"Register user failed", traceID, reqTime, err))
		return
	}

	dto.Ok(c, dto.NewSuccessEmpty(http.StatusOK, "Create user success", traceID, reqTime))
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"token_type"`
}

type RefreshTokenSuccess = dto.BaseResponse[RefreshTokenRequest]

// @BasePath /api/v1
// Register godoc
// @Summary      Refresh token
// @Description  Refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data  body      RefreshTokenRequest  true  "Register User"
// @Success      200   {object}  RefreshTokenSuccess
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	reqTime := time.Now().UTC()
	var req RefreshTokenRequest
	traceID := c.GetString(middleware.ContextKeyTraceID)
	if traceID == "" {
		traceID = c.GetHeader(middleware.HeaderKeyTraceID)
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeBadRequest, "Invalid refresh token request", traceID, reqTime, err))
		return
	}
	if req.RefreshToken == "" {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeBadRequest, "Invalid refresh token", traceID, reqTime, enum.ErrWrongType))
		return
	}
	token, err := h.authService.Refresh(c, req.RefreshToken)
	if err != nil {
		dto.BadRequest(c, dto.NewError(http.StatusBadRequest, enum.CodeBadRequest, "error", traceID, reqTime, err))
	}
	resp := LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    "Bearer",
		Expires:      token.Expired,
	}
	dto.Ok(c, dto.NewSuccess(http.StatusOK, "Refresh token success", traceID, resp, reqTime))
	return
}
