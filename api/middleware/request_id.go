package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the header name for request ID
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey is the key used to store request ID in Fiber context
	RequestIDKey = "request_id"
)

// RequestIDMiddleware generates or extracts request ID from headers
// and stores it in the context for logging and tracing purposes
func RequestIDMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Check if request ID already exists in header
		requestID := c.Get(RequestIDHeader)

		// If not present, generate a new UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store request ID in Fiber context
		c.Locals(RequestIDKey, requestID)

		// Add request ID to response header
		c.Set(RequestIDHeader, requestID)

		return c.Next()
	}
}

// GetRequestID extracts request ID from Fiber context
func GetRequestID(c fiber.Ctx) string {
	if requestID, ok := c.Locals(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
