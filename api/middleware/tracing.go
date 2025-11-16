package middleware

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v3"
)

// fiberHeaderCarrier adapts Fiber headers to OpenTelemetry propagation
type fiberHeaderCarrier struct {
	ctx fiber.Ctx
}

func (c fiberHeaderCarrier) Get(key string) string {
	return c.ctx.Get(key)
}

func (c fiberHeaderCarrier) Set(key, value string) {
	c.ctx.Set(key, value)
}

func (c fiberHeaderCarrier) Keys() []string {
	// Fiber v3 doesn't expose header keys directly, so we return empty
	// This is acceptable as the propagator will call Get() for known keys
	return []string{}
}

// TracingMiddleware creates a middleware for distributed tracing
func TracingMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		// Extract trace context from headers
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, fiberHeaderCarrier{ctx: c})

		// Start span
		tracer := otel.Tracer("shikposh-backend")
		ctx, span := tracer.Start(ctx, c.Method()+" "+c.Path())
		defer span.End()

		// Set span attributes
		span.SetAttributes(
			attribute.String("http.method", c.Method()),
			attribute.String("http.path", c.Path()),
			attribute.String("http.url", c.OriginalURL()),
			attribute.String("http.scheme", c.Protocol()),
			attribute.String("http.host", c.Hostname()),
			attribute.String("http.client_ip", c.IP()),
			attribute.String("http.user_agent", c.Get("User-Agent")),
		)

		// Add request ID to span if available
		if requestID := GetRequestID(c); requestID != "" {
			span.SetAttributes(attribute.String("request.id", requestID))
		}

		// Store span in locals for potential use in handlers
		c.Locals("span", span)
		c.Locals("trace_ctx", ctx)

		// Record start time
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Set response attributes
		span.SetAttributes(
			attribute.Int("http.status_code", c.Response().StatusCode()),
			attribute.Int("http.response.size", len(c.Response().Body())),
			attribute.Int64("http.duration_ms", duration.Milliseconds()),
		)

		// Set span status
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else if c.Response().StatusCode() >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Inject trace context into response headers
		propagator.Inject(ctx, fiberHeaderCarrier{ctx: c})

		return err
	}
}

// GetSpan retrieves the current span from the context
func GetSpan(c fiber.Ctx) trace.Span {
	if span, ok := c.Locals("span").(trace.Span); ok {
		return span
	}
	if ctx, ok := c.Locals("trace_ctx").(context.Context); ok {
		return trace.SpanFromContext(ctx)
	}
	return trace.SpanFromContext(c.Context())
}
