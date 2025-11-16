package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/logging"
)

type Config struct {
	ServiceName  string
	OTLPEndpoint string // e.g., "http://localhost:4318" for HTTP OTLP endpoint
	Environment  string
	SamplingRate float64 // 0.0 to 1.0
	Enabled      bool
}

type Tracer struct {
	tracer trace.Tracer
	tp     *tracesdk.TracerProvider
}

// New initializes Jaeger tracing
func New(cfg Config) (*Tracer, error) {
	if !cfg.Enabled {
		return &Tracer{}, nil
	}

	// Create OTLP HTTP exporter
	otlpEndpoint := cfg.OTLPEndpoint
	if otlpEndpoint == "" {
		// Default to localhost OTLP HTTP endpoint (Jaeger accepts OTLP on port 4318)
		otlpEndpoint = "localhost:4318"
	}

	// Parse endpoint URL - OTLP expects host:port format
	// Remove http:// or https:// prefix if present
	endpoint := otlpEndpoint
	if len(endpoint) > 7 && endpoint[:7] == "http://" {
		endpoint = endpoint[7:]
	} else if len(endpoint) > 8 && endpoint[:8] == "https://" {
		endpoint = endpoint[8:]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), // Use insecure for local development
	)
	if err != nil {
		logging.Error("Failed to create OTLP exporter").
			WithString("endpoint", otlpEndpoint).
			WithError(err).
			Log()
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Set sampling rate
	samplingRate := cfg.SamplingRate
	if samplingRate < 0.0 || samplingRate > 1.0 {
		samplingRate = 1.0
	}

	// Create tracer provider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(samplingRate)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer := tp.Tracer(cfg.ServiceName)

	logging.Info("Jaeger tracing initialized via OTLP").
		WithString("service_name", cfg.ServiceName).
		WithString("otlp_endpoint", otlpEndpoint).
		WithAny("sampling_rate", samplingRate).
		Log()

	return &Tracer{
		tracer: tracer,
		tp:     tp,
	}, nil
}

// GetTracer returns the OpenTelemetry tracer
func (t *Tracer) GetTracer() trace.Tracer {
	return t.tracer
}

// Shutdown gracefully shuts down the tracer provider
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.tp != nil {
		return t.tp.Shutdown(ctx)
	}
	return nil
}
