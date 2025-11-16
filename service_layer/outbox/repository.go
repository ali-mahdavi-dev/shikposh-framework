package outbox

import (
	"context"
	"time"

	"github.com/ali-mahdavi-dev/framework/adapter"

	"gorm.io/gorm"
)

// Repository defines the interface for outbox event repository operations
type Repository interface {
	adapter.BaseRepository[*OutboxEvent]
	Model(ctx context.Context) *gorm.DB
	Create(ctx context.Context, event *OutboxEvent) error
	GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error)
	MarkAsProcessing(ctx context.Context, id OutboxEventID) error
	MarkAsCompleted(ctx context.Context, id OutboxEventID) error
	MarkAsFailed(ctx context.Context, id OutboxEventID, errorMsg string) error
	IncrementRetry(ctx context.Context, id OutboxEventID) error
}

// GormRepository is a GORM implementation of the outbox repository
type GormRepository struct {
	adapter.BaseRepository[*OutboxEvent]
	db        *gorm.DB
	tableName string
}

// NewGormRepository creates a new GORM-based outbox repository
// If tableName is empty, it will use the default table name from OutboxEvent.TableName()
func NewGormRepository(db *gorm.DB, tableName string) Repository {
	return &GormRepository{
		BaseRepository: adapter.NewGormRepository[*OutboxEvent](db),
		db:             db,
		tableName:      tableName,
	}
}

func (r *GormRepository) Model(ctx context.Context) *gorm.DB {
	model := r.db.WithContext(ctx).Model(&OutboxEvent{})
	if r.tableName != "" {
		model = model.Table(r.tableName)
	}
	return model
}

func (r *GormRepository) Create(ctx context.Context, event *OutboxEvent) error {
	return r.Model(ctx).Create(event).Error
}

func (r *GormRepository) GetPendingEvents(ctx context.Context, limit int) ([]*OutboxEvent, error) {
	var events []*OutboxEvent
	err := r.Model(ctx).
		Where("status = ?", OutboxStatusPending).
		Where("retry_count < max_retries").
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error
	return events, err
}

func (r *GormRepository) MarkAsProcessing(ctx context.Context, id OutboxEventID) error {
	return r.Model(ctx).
		Where("id = ?", uint64(id)).
		Updates(map[string]interface{}{
			"status":     OutboxStatusProcessing,
			"updated_at": time.Now(),
		}).Error
}

func (r *GormRepository) MarkAsCompleted(ctx context.Context, id OutboxEventID) error {
	now := time.Now()
	return r.Model(ctx).
		Where("id = ?", uint64(id)).
		Updates(map[string]interface{}{
			"status":       OutboxStatusCompleted,
			"processed_at": now,
			"updated_at":   now,
		}).Error
}

func (r *GormRepository) MarkAsFailed(ctx context.Context, id OutboxEventID, errorMsg string) error {
	return r.Model(ctx).
		Where("id = ?", uint64(id)).
		Updates(map[string]interface{}{
			"status":        OutboxStatusFailed,
			"error_message": errorMsg,
			"updated_at":    time.Now(),
		}).Error
}

func (r *GormRepository) IncrementRetry(ctx context.Context, id OutboxEventID) error {
	return r.Model(ctx).
		Where("id = ?", uint64(id)).
		Update("retry_count", gorm.Expr("retry_count + 1")).Error
}
