package adapter

import "sync"

// Entity is the base interface that all domain entities must implement.
// GetID returns the primary identifier of the entity as uint64.
// Each entity can have its own strong typedef for ID (e.g. UserID uint64)
// and should convert it to uint64 in GetID.
type Entity interface {
	GetID() uint64
	Event() []any
	AddEvent(event any)
}

type BaseEntity struct {
	Events   []any `json:",omitempty" gorm:"-"`
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
