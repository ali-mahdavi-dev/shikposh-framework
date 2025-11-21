package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/logging"

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

		// Determine if this is an error response
		isError := err != nil || status >= 400

		// Create log entry based on error status
		var entry logging.Logger
		if isError {
			entry = logging.Error("HTTP Request Error").
				WithAny("path", path).
				WithAny("client_ip", clientIP).
				WithAny("method", method).
				WithAny("latency", latency).
				WithAny("status_code", status).
				WithAny("body_size", len(respBody))
		} else {
			entry = logging.Info("HTTP Request").
				WithAny("path", path).
				WithAny("client_ip", clientIP).
				WithAny("method", method).
				WithAny("latency", latency).
				WithAny("status_code", status).
				WithAny("body_size", len(respBody))
		}

		// Add request_id if available
		if requestID != "" {
			entry = entry.WithString("request_id", requestID)
		}

		// Add error details if error exists
		if err != nil {
			entry = entry.WithError(err).
				WithString("error_message", err.Error())
		}

		// Add error status code details and parse error from response body
		if status >= 400 {
			statusText := http.StatusText(status)
			if statusText != "" {
				entry = entry.WithString("error_status", statusText)
			}

			// Try to parse error details from response body
			if len(respBody) > 0 {
				var responseResult struct {
					Success bool `json:"success"`
					Error   *struct {
						Code      string `json:"code"`
						Message   string `json:"message"`
						Detail    string `json:"detail"`
						Status    string `json:"status"`
						RequestID string `json:"request_id,omitempty"`
					} `json:"error,omitempty"`
				}

				if err := json.Unmarshal(respBody, &responseResult); err == nil {
					if responseResult.Error != nil {
						entry = entry.
							WithString("error_code", responseResult.Error.Code).
							WithString("error_message", responseResult.Error.Message).
							WithString("error_detail", responseResult.Error.Detail)
					}
				}
			}
		}

		entry.Log()

		return err
	}
}
