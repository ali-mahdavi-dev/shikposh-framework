package middleware

import (
	"time"

	"github.com/ali-mahdavi-dev/framework/infrastructure/logging"

	"github.com/gofiber/fiber/v3"
)

func DefaultStructuredLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Run next handler
		err := c.Next()

		// Capture response body
		respBody := c.Response().Body()

		// Collect metadata
		latency := time.Since(start)
		clientIP := c.IP()
		method := c.Method()
		status := c.Response().StatusCode()
		path := c.OriginalURL()
		requestID := GetRequestID(c)

		entry := logging.Info("HTTP Request").
			WithAny("path", path).
			WithAny("client_ip", clientIP).
			WithAny("method", method).
			WithAny("latency", latency).
			WithAny("status_code", status).
			WithAny("body_size", len(respBody))

		// Add request_id if available
		if requestID != "" {
			entry = entry.WithString("request_id", requestID)
		}

		if err != nil {
			entry = entry.WithError(err)
		}

		entry.Log()

		return err
	}
}
