package dto

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BaseResponse is a unified envelope for all API responses.
// - All timestamps are in UTC RFC3339Nano for consistent logging and searching.
// - Data is a pointer so that `omitempty` actually hides the field when nil.
type BaseResponse[T any] struct {
	TraceID      string        `json:"trace_id,omitempty"` // Trace identifier for correlation
	ResponseCode int           `json:"response_code"`      // Business/application-specific code (e.g., 200, 10001)
	Message      string        `json:"message"`            // Human-readable message
	Data         *T            `json:"data,omitempty"`     // Optional payload data
	Error        *ErrorPayload `json:"error,omitempty"`    // Error details if any

	RequestTime  time.Time `json:"request_time"`  // When the request was received
	ResponseTime time.Time `json:"response_time"` // When the response was sent
	ProcessingMs int64     `json:"processing_ms"` // Total processing time in milliseconds
}

// ErrorPayload contains structured error information for clients and log systems.
type ErrorPayload struct {
	Code    string       `json:"code,omitempty"`    // Application-level error code (e.g., "VALIDATION_FAILED")
	Details []FieldError `json:"details,omitempty"` // List of field-specific validation errors
}

// FieldError follows a common convention for validation errors.
type FieldError struct {
	Field   string `json:"field"`           // Field name
	Message string `json:"message"`         // Human-readable validation message
	Tag     string `json:"tag,omitempty"`   // Validation tag (e.g., "required", "min")
	Param   string `json:"param,omitempty"` // Parameter for the validation tag (e.g., "8" for min=8)
}

// EmptyData is used for success responses that have no payload.
type EmptyData struct{}

// ErrorResponse is the standard error envelope.
// @name ErrorResponse
type ErrorResponse = BaseResponse[EmptyData]

// ==== Constructors ====

// NewSuccess wraps a payload in a success response envelope.
// `code` is your business/application code (often 200 for OK).
func NewSuccess[T any](code int, message, traceID string, data T, reqTime time.Time) BaseResponse[T] {
	now := time.Now().UTC()
	return BaseResponse[T]{
		TraceID:      traceID,
		ResponseCode: code,
		Message:      message,
		Data:         ptr(data),
		RequestTime:  reqTime.UTC(),
		ResponseTime: now,
		ProcessingMs: now.Sub(reqTime).Milliseconds(),
	}
}

// NewSuccessEmpty creates a success response without payload data.
func NewSuccessEmpty(code int, message, traceID string, reqTime time.Time) BaseResponse[EmptyData] {
	now := time.Now().UTC()
	return BaseResponse[EmptyData]{
		TraceID:      traceID,
		ResponseCode: code,
		Message:      message,
		RequestTime:  reqTime.UTC(),
		ResponseTime: now,
		ProcessingMs: now.Sub(reqTime).Milliseconds(),
	}
}

// NewError creates an error response with structured error details.
// `err` can be a validator.ValidationErrors or a regular error.
func NewError(code int, appErrCode, message, traceID string, reqTime time.Time, err error) BaseResponse[EmptyData] {
	now := time.Now().UTC()
	return BaseResponse[EmptyData]{
		TraceID:      traceID,
		ResponseCode: code,
		Message:      message,
		Error:        &ErrorPayload{Code: appErrCode, Details: validationDetails(err)},
		RequestTime:  reqTime.UTC(),
		ResponseTime: now,
		ProcessingMs: now.Sub(reqTime).Milliseconds(),
	}
}

// ==== Gin helpers ====

// WriteJSON writes the provided response as JSON with the given HTTP status.
func WriteJSON(c *gin.Context, status int, resp any) {
	c.JSON(status, resp)
}

// Ok is a shorthand for writing a 200 OK success response.
func Ok[T any](c *gin.Context, resp BaseResponse[T]) {
	WriteJSON(c, http.StatusOK, resp)
}

// BadRequest is a shorthand for writing a 400 Bad Request error response.
func BadRequest(c *gin.Context, resp BaseResponse[EmptyData]) {
	WriteJSON(c, http.StatusBadRequest, resp)
}

// ==== Validation mapping ====

// validationDetails converts go-playground/validator errors into a list of FieldError.
// Non-validation errors will be wrapped into a single FieldError with the key "error".
func validationDetails(err error) []FieldError {
	if err == nil {
		return nil
	}
	if verrs, ok := err.(validator.ValidationErrors); ok {
		out := make([]FieldError, 0, len(verrs))
		for _, fe := range verrs {
			out = append(out, FieldError{
				Field:   fe.Field(),
				Message: translate(fe), // Replace with your i18n implementation if needed
				Tag:     fe.Tag(),
				Param:   fe.Param(),
			})
		}
		return out
	}
	return []FieldError{{
		Field:   "error",
		Message: err.Error(),
	}}
}

// translate maps validator tags to user-friendly English messages.
// This is a placeholder for i18n; replace with your translation mechanism.
func translate(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email address"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	default:
		return fe.Error()
	}
}

// ==== Utils ====

// ptr returns a pointer to the provided value.
func ptr[T any](v T) *T { return &v }
