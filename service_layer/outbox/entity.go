package outbox

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
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

// JSONBMap is a custom type for JSONB map fields that implements driver.Valuer and sql.Scanner
type JSONBMap map[string]interface{}

// Value implements driver.Valuer interface
func (j JSONBMap) Value() (driver.Value, error) {
	if j == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(j)
}

// Scan implements sql.Scanner interface
func (j *JSONBMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONBMap)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = JSONBMap(result)
	return nil
}

// OutboxEvent represents a generic outbox event entity
type OutboxEvent struct {
	adapter.BaseEntity
	ID            OutboxEventID `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt    `gorm:"index"`
	EventType     string            `json:"event_type" gorm:"event_type"`
	AggregateType string            `json:"aggregate_type" gorm:"aggregate_type"`
	AggregateID   string            `json:"aggregate_id" gorm:"aggregate_id"`
	Payload       JSONBMap          `json:"payload" gorm:"type:jsonb"`
	Status        OutboxEventStatus `json:"status" gorm:"status;default:'pending'"`
	RetryCount    int               `json:"retry_count" gorm:"retry_count;default:0"`
	MaxRetries    int               `json:"max_retries" gorm:"max_retries;default:5"`
	ErrorMessage  *string           `json:"error_message,omitempty" gorm:"error_message;type:text"`
	ProcessedAt   *time.Time        `json:"processed_at,omitempty" gorm:"processed_at"`
}

// TableName returns the table name for the outbox event
// This can be overridden by modules if they need different table names
func (o *OutboxEvent) TableName() string {
	return "outbox_events"
}
