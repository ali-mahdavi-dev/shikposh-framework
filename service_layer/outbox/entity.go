package outbox

import (
	"time"

	"github.com/ali-mahdavi-dev/shikposh-framework/adapter"

	"gorm.io/gorm"
)

// OutboxEventStatus represents the status of an outbox event
type OutboxEventStatus string

const (
	OutboxStatusPending    OutboxEventStatus = "pending"
	OutboxStatusProcessing OutboxEventStatus = "processing"
	OutboxStatusCompleted  OutboxEventStatus = "completed"
	OutboxStatusFailed     OutboxEventStatus = "failed"
)

// OutboxEventID is the type for outbox event ID
type OutboxEventID uint64

// OutboxEvent represents a generic outbox event entity
type OutboxEvent struct {
	adapter.BaseEntity
	ID            OutboxEventID `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt         `gorm:"index"`
	EventType     string                 `json:"event_type" gorm:"event_type"`
	AggregateType string                 `json:"aggregate_type" gorm:"aggregate_type"`
	AggregateID   string                 `json:"aggregate_id" gorm:"aggregate_id"`
	Payload       map[string]interface{} `json:"payload" gorm:"type:jsonb"`
	Status        OutboxEventStatus      `json:"status" gorm:"status;default:'pending'"`
	RetryCount    int                    `json:"retry_count" gorm:"retry_count;default:0"`
	MaxRetries    int                    `json:"max_retries" gorm:"max_retries;default:5"`
	ErrorMessage  *string                `json:"error_message,omitempty" gorm:"error_message;type:text"`
	ProcessedAt   *time.Time             `json:"processed_at,omitempty" gorm:"processed_at"`
}

// TableName returns the table name for the outbox event
// This can be overridden by modules if they need different table names
func (o *OutboxEvent) TableName() string {
	return "outbox_events"
}
