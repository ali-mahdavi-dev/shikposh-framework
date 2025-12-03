package service_host

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/logging"
)

// Service represents a service that can be started and stopped
type Service interface {
	// Start starts the service. It should block until the service is stopped or an error occurs.
	Start() error

	// Shutdown gracefully shuts down the service with the given context.
	Shutdown(ctx context.Context) error

	// Name returns the name of the service for logging purposes.
	Name() string
}

// ServiceHost orchestrates the lifecycle of multiple services
type ServiceHost struct {
	services        []Service
	shutdownTimeout time.Duration
}

// NewServiceHost creates a new ServiceHost with the given shutdown timeout.
// If timeout is 0, it defaults to 30 seconds.
func NewServiceHost(shutdownTimeout time.Duration) *ServiceHost {
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second
	}

	return &ServiceHost{
		services:        make([]Service, 0),
		shutdownTimeout: shutdownTimeout,
	}
}

// AddService adds a service to the host
func (h *ServiceHost) AddService(service Service) {
	if service == nil {
		return
	}
	h.services = append(h.services, service)
}

// AddServices adds multiple services to the host
func (h *ServiceHost) AddServices(services ...Service) {
	for _, service := range services {
		h.AddService(service)
	}
}

// Start starts all services and waits for shutdown signal
func (h *ServiceHost) Start() error {
	if len(h.services) == 0 {
		return fmt.Errorf("no services to start")
	}

	// Start all services in separate goroutines
	serviceErrs := make(chan error, len(h.services))
	for _, service := range h.services {
		go func(svc Service) {
			logging.Info("Starting service").
				WithString("service", svc.Name()).
				Log()

			if err := svc.Start(); err != nil {
				serviceErrs <- fmt.Errorf("service %s failed: %w", svc.Name(), err)
			}
		}(service)
	}

	// Wait for interrupt signal or service error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-serviceErrs:
		// One of the services failed, shutdown all
		logging.Error("Service error, initiating shutdown").
			WithError(err).
			Log()
		return h.shutdown(context.Background())
	case <-quit:
		logging.Info("Shutdown signal received").
			Log()
		return h.shutdown(context.Background())
	}
}

// shutdown gracefully shuts down all services
func (h *ServiceHost) shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, h.shutdownTimeout)
	defer cancel()

	logging.Info("Shutting down services").
		WithInt("count", len(h.services)).
		WithString("timeout", h.shutdownTimeout.String()).
		Log()

	var lastErr error
	for _, service := range h.services {
		logging.Info("Shutting down service").
			WithString("service", service.Name()).
			Log()

		if err := service.Shutdown(shutdownCtx); err != nil {
			logging.Warn("Service shutdown error").
				WithString("service", service.Name()).
				WithError(err).
				Log()
			lastErr = err
		} else {
			logging.Info("Service shut down successfully").
				WithString("service", service.Name()).
				Log()
		}
	}

	if lastErr != nil {
		return fmt.Errorf("some services failed to shutdown: %w", lastErr)
	}

	logging.Info("All services shut down successfully").
		Log()

	return nil
}

// GetServices returns all registered services
func (h *ServiceHost) GetServices() []Service {
	return h.services
}
