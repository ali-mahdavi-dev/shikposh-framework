package outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/logging"

	"github.com/IBM/sarama"
)

// MessageConsumer defines the interface for consuming messages from a message broker
type MessageConsumer interface {
	ConsumeMessages(topic string, fn func(pc sarama.PartitionConsumer)) error
}

// EventHandler defines the interface for handling events consumed from the message broker
type EventHandler interface {
	HandleEvent(ctx context.Context, eventType string, payload map[string]interface{}) error
}

// Consumer handles consuming messages from a message broker and processing them
type Consumer struct {
	consumer MessageConsumer
	handler  EventHandler
	topic    string
	stopChan chan struct{}
}

// NewConsumer creates a new outbox consumer
func NewConsumer(consumer MessageConsumer, handler EventHandler, topic string) *Consumer {
	return &Consumer{
		consumer: consumer,
		handler:  handler,
		topic:    topic,
		stopChan: make(chan struct{}),
	}
}

// Start starts the Kafka consumer
func (c *Consumer) Start(ctx context.Context) error {
	handler := func(pc sarama.PartitionConsumer) {
		for {
			select {
			case <-ctx.Done():
				logging.Info("Outbox consumer stopped: context cancelled").Log()
				return
			case <-c.stopChan:
				logging.Info("Outbox consumer stopped: stop signal received").Log()
				return
			case message := <-pc.Messages():
				if err := c.handleMessage(ctx, message); err != nil {
					logging.Error("Failed to handle message").
						WithError(err).
						WithInt("partition", int(message.Partition)).
						WithInt64("offset", message.Offset).
						Log()
				}
			case err := <-pc.Errors():
				if err != nil {
					logging.Error("Consumer error").
						WithError(err.Err).
						Log()
				}
			}
		}
	}

	return c.consumer.ConsumeMessages(c.topic, handler)
}

// Stop stops the consumer
func (c *Consumer) Stop() {
	close(c.stopChan)
}

func (c *Consumer) handleMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	var kafkaMessage map[string]interface{}
	if err := json.Unmarshal(message.Value, &kafkaMessage); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	eventType, ok := kafkaMessage["event_type"].(string)
	if !ok {
		return fmt.Errorf("event_type is missing or invalid")
	}

	payload, ok := kafkaMessage["payload"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("payload is missing or invalid")
	}

	// Delegate to the event handler
	return c.handler.HandleEvent(ctx, eventType, payload)
}
