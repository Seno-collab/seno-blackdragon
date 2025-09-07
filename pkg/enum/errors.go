package enum

import "errors"

// ===== Errors (sentinel) =====
var (
	// Token / AuthN
	ErrInvalidToken       = errors.New("invalid token")             // parse/verify fail, exp, nbf...
	ErrWrongType          = errors.New("wrong token type")          // access != refresh
	ErrWrongAlgorithm     = errors.New("unexpected signing method") // HS256 vs RS256...
	ErrInvalidCredentials = errors.New("invalid credentials")       // wrong email/password

	// Refresh flow
	ErrRefreshNotActive = errors.New("refresh token not active")                // not in allow-list
	ErrRefreshRevoked   = errors.New("refresh token already revoked (reuse)")   // reuse detected
	ErrFamilyBlocked    = errors.New("refresh family is blocked")               // block whole family
	ErrRotationRace     = errors.New("refresh rotation in progress, try again") // lock conflict

	// User / Business
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlready = errors.New("email already registered")
)

// ===== Error codes (machine-readable) =====
const (
	// Generic
	CodeBadRequest       = "BAD_REQUEST"
	CodeBusinessError    = "BUSINESS_ERROR"
	CodeInternalError    = "INTERNAL_ERROR"
	CodeValidationFailed = "VALIDATION_FAILED"
	CodeAuth             = "AUTH"

	// AuthN / Token
	CodeInvalidToken       = "INVALID_TOKEN"
	CodeWrongTokenType     = "WRONG_TOKEN_TYPE"
	CodeUnexpectedAlg      = "UNEXPECTED_SIGNING_METHOD"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"

	// Refresh
	CodeRefreshNotActive = "REFRESH_NOT_ACTIVE"
	CodeRefreshRevoked   = "REFRESH_REVOKED"
	CodeFamilyBlocked    = "FAMILY_BLOCKED"
	CodeRotationRace     = "ROTATION_RACE"

	// User
	CodeUserNotFound  = "USER_NOT_FOUND"
	CodeEmailConflict = "EMAIL_ALREADY_REGISTERED"
)
