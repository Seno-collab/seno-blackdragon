package enum

import "errors"

// Parse & validate helpers
var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrWrongType          = errors.New("wrong token type")
	ErrWrongAlgorithm     = errors.New("unexpected signing method")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlready       = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrRefreshToken       = errors.New("refresh session not found (revoked or expired)")
)

var (
	VALIDATION_FAILED = "VALIDATION_FAILED"
	BUSINESS_ERROR    = "BUSINESS_ERROR"
	AUTH_FAILED       = "AUTH_FAILED"
)
