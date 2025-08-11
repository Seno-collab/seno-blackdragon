package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ContextKeyTraceID = "trace_id"
	HeaderKeyTraceID  = "X-Trace-Id"
	maxLogBodyBytes   = 8 << 10 // 8KB: enough for debugging without being too noisy
	defaultSkipPaths  = "/healthz,/metrics"
)

var sensitiveKeys = map[string]struct{}{
	"password":      {},
	"pass":          {},
	"pwd":           {},
	"token":         {},
	"access_token":  {},
	"refresh_token": {},
	"authorization": {},
	"secret":        {},
	"cardNumber":    {},
	"cvv":           {},
	"pin":           {},
}

type bodyLogWriter struct {
	gin.ResponseWriter
	buf bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	// Copy to buffer for logging without blocking the actual response flow
	if len(b) > 0 {
		// Only write up to maxLogBodyBytes to avoid large logs and memory usage
		remain := maxLogBodyBytes - w.buf.Len()
		if remain > 0 {
			if len(b) > remain {
				w.buf.Write(b[:remain])
			} else {
				w.buf.Write(b)
			}
		}
	}
	return w.ResponseWriter.Write(b)
}

func (w *bodyLogWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

type TraceLogOptions struct {
	SkipPaths string // CSV format: "/healthz,/metrics"
}

func TraceAndLogFullMiddleware(logger *zap.Logger, opts *TraceLogOptions) gin.HandlerFunc {
	if opts == nil {
		opts = &TraceLogOptions{SkipPaths: defaultSkipPaths}
	}
	skip := toSet(opts.SkipPaths)

	return func(c *gin.Context) {
		// Skip noisy paths (prefer FullPath; fallback to URL.Path)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if _, ok := skip[path]; ok {
			c.Next()
			return
		}

		start := time.Now()
		traceID := c.GetHeader(HeaderKeyTraceID)
		if traceID == "" {
			traceID = uuid.New().String()
		}
		c.Set(ContextKeyTraceID, traceID)
		c.Writer.Header().Set(HeaderKeyTraceID, traceID)

		// Read request body (must restore it for the handler)
		var reqBody []byte
		if c.Request != nil && c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Wrap ResponseWriter to capture response body (with size limit)
		blw := &bodyLogWriter{ResponseWriter: c.Writer}
		c.Writer = blw

		// ===== Log Request (Info level for traceability) =====
		reqFields := []zap.Field{
			zap.String("trace_id", traceID),
			zap.String("client_ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Any("headers", pickHeaders(c.Request.Header, []string{
				"User-Agent", "Content-Type", "Accept", "Accept-Encoding",
			})),
		}
		var jsonBody map[string]interface{}
		if isTextLike(c.ContentType()) && len(reqBody) > 0 {
			var jsonBody any
			if err := json.Unmarshal(reqBody, &jsonBody); err == nil {
				// redact fields in-place
				redactJSON(jsonBody)

				// limit size: if too big, log truncated string; else log as object (no \" escaping)
				if b, err := json.Marshal(jsonBody); err == nil && len(b) > maxLogBodyBytes {
					reqFields = append(reqFields, zap.String("request_body", limit(string(b), maxLogBodyBytes)))
				} else {
					reqFields = append(reqFields, zap.Any("request_body", jsonBody))
				}
			} else {
				// not JSON → log raw (truncated)
				reqFields = append(reqFields, zap.String("request_body_raw", limit(string(reqBody), maxLogBodyBytes)))
			}
		}

		logger.Info("http_request", reqFields...)

		// Execute handler
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// ===== Log Response =====
		respFields := []zap.Field{
			zap.String("trace_id", traceID),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("bytes_out", c.Writer.Size()),
		}

		ct := c.Writer.Header().Get("Content-Type")
		if isTextLike(ct) && blw.buf.Len() > 0 {
			if json.Unmarshal(blw.buf.Bytes(), &jsonBody) == nil {
				respFields = append(respFields, zap.Any("response_body", jsonBody))
			} else {
				respFields = append(respFields, zap.String("response_body_raw", limit(blw.buf.String(), maxLogBodyBytes)))
			}
		}

		if len(c.Errors) > 0 {
			// Collect gin errors if any
			respFields = append(respFields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case status >= 500:
			logger.Error("http_response", respFields...)
		case status >= 400:
			logger.Warn("http_response", respFields...)
		default:
			logger.Info("http_response", respFields...)
		}
	}
}

// ===== Helpers =====

func toSet(csv string) map[string]struct{} {
	m := map[string]struct{}{}
	for _, p := range strings.Split(csv, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			m[p] = struct{}{}
		}
	}
	return m
}

func isTextLike(ct string) bool {
	ct = strings.ToLower(ct)
	return strings.Contains(ct, "json") ||
		strings.Contains(ct, "text/") ||
		strings.Contains(ct, "xml") ||
		strings.Contains(ct, "x-www-form-urlencoded")
}

// limitAndRedact: trim body and redact sensitive fields if JSON
func limitAndRedact(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	if json.Valid(b) {
		var v any
		if err := json.Unmarshal(b, &v); err == nil {
			redactJSON(v)
			b2, _ := json.Marshal(v)
			return limit(string(b2), maxLogBodyBytes)
		}
	}
	return limit(string(b), maxLogBodyBytes)
}

func limit(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "...(truncated)"
}

func redactJSON(v any) {
	switch t := v.(type) {
	case map[string]any:
		for k, val := range t {
			if _, bad := sensitiveKeys[strings.ToLower(k)]; bad {
				t[k] = "***REDACTED***"
				continue
			}
			redactJSON(val)
		}
	case []any:
		for i := range t {
			redactJSON(t[i])
		}
	}
}

func pickHeaders(h http.Header, keys []string) map[string]string {
	out := make(map[string]string, len(keys))
	for _, k := range keys {
		if v := h.Get(k); v != "" {
			out[k] = v
		}
	}
	return out
}
