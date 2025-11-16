package adapter

import (
	"sync"
)

type Entity interface {
	Event() []any
	AddEvent(event any)
}

type BaseEntity struct {
	Events   []any `gorm:"-"`
	eventsMu sync.Mutex
}

// Event returns all events and clears them atomically
func (u *BaseEntity) Event() []any {
	u.eventsMu.Lock()
	defer u.eventsMu.Unlock()

	events := append([]any(nil), u.Events...)
	u.Events = nil
	return events
}

// AddEvent adds an event to the entity in a thread-safe manner
func (u *BaseEntity) AddEvent(event any) {
	u.eventsMu.Lock()
	defer u.eventsMu.Unlock()
	u.Events = append(u.Events, event)
}
