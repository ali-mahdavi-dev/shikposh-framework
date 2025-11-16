package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
)

type gormRepository[E Entity] struct {
	seen   []Entity
	seenMu sync.Mutex
	db     *gorm.DB
}

func NewGormRepository[E Entity](db *gorm.DB) BaseRepository[E] {
	return &gormRepository[E]{db: db}
}

func (c *gormRepository[E]) FindByID(ctx context.Context, id uint64) (E, error) {
	model, err := c.FindByField(ctx, "id", id)

	c.SetSeen(model)
	return model, err
}

func (c *gormRepository[E]) FindByField(ctx context.Context, field string, value interface{}) (E, error) {
	var e E
	err := c.db.WithContext(ctx).Model(e).Where(field+"=?", value).First(&e).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return e, ErrEntityNotFound
	}

	c.SetSeen(e)
	return e, err
}

func (c *gormRepository[E]) Remove(ctx context.Context, model E, softDelete bool) error {
	c.SetSeen(model)
	if softDelete {
		now := time.Now()
		var e E
		return c.db.WithContext(ctx).Model(e).Update("deleted_at", &now).Error
	}

	return c.db.WithContext(ctx).Delete(model).Error
}

func (c *gormRepository[E]) Save(ctx context.Context, model E) error {
	err := c.db.WithContext(ctx).Save(model).Error
	c.SetSeen(model)
	return err
}

func (c *gormRepository[E]) Modify(ctx context.Context, model E) error {
	m, err := structToMap(model)
	if err != nil {
		return fmt.Errorf("gormRepository.Modify fail convert data to map: %w", err)
	}
	err = c.db.WithContext(ctx).Updates(m).Error
	c.SetSeen(model)
	return err
}

func (c *gormRepository[E]) Seen() []Entity {
	c.seenMu.Lock()
	defer c.seenMu.Unlock()

	seen := make([]Entity, len(c.seen))
	copy(seen, c.seen)
	c.seen = []Entity{}
	return seen
}

func (c *gormRepository[E]) SetSeen(model Entity) {
	c.seenMu.Lock()
	defer c.seenMu.Unlock()
	c.seen = append(c.seen, model)
}

func structToMap(model any) (map[string]interface{}, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
