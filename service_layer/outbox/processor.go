package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/ali-mahdavi-dev/framework/infrastructure/logging"
)

// MessagePublisher defines the interface for publishing messages to a message broker
type MessagePublisher interface {
	SendMessage(topic string, message interface{}) error
}

// ProcessorConfig holds configuration for the outbox processor
type ProcessorConfig struct {
	BatchSize    int           // Number of events to process in each batch
	PollInterval time.Duration // Interval between polling cycles
	Topic        string        // Kafka topic to publish events to
}

// DefaultProcessorConfig returns default configuration for the processor
func DefaultProcessorConfig(topic string) ProcessorConfig {
	return ProcessorConfig{
		BatchSize:    10,
		PollInterval: 5 * time.Second,
		Topic:        topic,
	}
}

// Processor handles processing of outbox events and publishing them to a message broker
type Processor struct {
	repo     Repository
	pub      MessagePublisher
	config   ProcessorConfig
	stopChan chan struct{}
}

// NewProcessor creates a new outbox processor
func NewProcessor(repo Repository, publisher MessagePublisher, config ProcessorConfig) *Processor {
	return &Processor{
		repo:     repo,
		pub:      publisher,
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// Start starts the outbox processor in a background goroutine
func (p *Processor) Start(ctx context.Context) {
	go p.run(ctx)
}

// Stop stops the outbox processor
func (p *Processor) Stop() {
	close(p.stopChan)
}

func (p *Processor) run(ctx context.Context) {
	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logging.Info("Outbox processor stopped: context cancelled").Log()
			return
		case <-p.stopChan:
			logging.Info("Outbox processor stopped: stop signal received").Log()
			return
		case <-ticker.C:
			if err := p.processBatch(ctx); err != nil {
				logging.Error("Failed to process outbox batch").
					WithError(err).
					Log()
			}
		}
	}
}

func (p *Processor) processBatch(ctx context.Context) error {
	// Get pending events
	events, err := p.repo.GetPendingEvents(ctx, p.config.BatchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending events: %w", err)
	}

	if len(events) == 0 {
		return nil // No events to process
	}

	logging.Info("Processing outbox batch").
		WithInt("count", len(events)).
		Log()

	for _, event := range events {
		if err := p.processEvent(ctx, event); err != nil {
			logging.Error("Failed to process outbox event").
				WithInt64("event_id", int64(event.ID)).
				WithString("event_type", event.EventType).
				WithError(err).
				Log()

			// Increment retry count
			if incErr := p.repo.IncrementRetry(ctx, event.ID); incErr != nil {
				logging.Error("Failed to increment retry count").
					WithInt64("event_id", int64(event.ID)).
					WithError(incErr).
					Log()
			}

			// Mark as failed if max retries reached
			if event.RetryCount+1 >= event.MaxRetries {
				if markErr := p.repo.MarkAsFailed(ctx, event.ID, err.Error()); markErr != nil {
					logging.Error("Failed to mark event as failed").
						WithInt64("event_id", int64(event.ID)).
						WithError(markErr).
						Log()
				}
			}

			continue // Continue with next event
		}
	}

	return nil
}

func (p *Processor) processEvent(ctx context.Context, event *OutboxEvent) error {
	// Mark as processing
	if err := p.repo.MarkAsProcessing(ctx, event.ID); err != nil {
		return fmt.Errorf("failed to mark event as processing: %w", err)
	}

	// Prepare message
	message := map[string]interface{}{
		"event_id":       event.ID,
		"event_type":     event.EventType,
		"aggregate_type": event.AggregateType,
		"aggregate_id":   event.AggregateID,
		"payload":        event.Payload,
		"created_at":     event.CreatedAt,
	}

	// Send to message broker
	if err := p.pub.SendMessage(p.config.Topic, message); err != nil {
		return fmt.Errorf("failed to send message to broker: %w", err)
	}

	// Mark as completed
	if err := p.repo.MarkAsCompleted(ctx, event.ID); err != nil {
		return fmt.Errorf("failed to mark event as completed: %w", err)
	}

	logging.Info("Outbox event sent successfully").
		WithInt64("event_id", int64(event.ID)).
		WithString("event_type", event.EventType).
		WithString("aggregate_id", event.AggregateID).
		Log()

	return nil
}
