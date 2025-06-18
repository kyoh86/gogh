package script

import (
	"io"
	"time"

	"github.com/google/uuid"
)

type Entry struct {
	Name    string
	Content io.Reader
}

type Script interface {
	ID() string
	UUID() uuid.UUID
	Name() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// NewScript creates a new Script with the given entry.
func NewScript(entry Entry) Script {
	now := time.Now()
	return scriptElement{
		id:        uuid.Must(uuid.NewRandom()),
		name:      entry.Name,
		createdAt: now,
		updatedAt: now,
	}
}

// ConcreteScript creates a Script with the given parameters.
func ConcreteScript(id uuid.UUID, name string, createdAt, updatedAt time.Time) Script {
	return &scriptElement{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

type scriptElement struct {
	id   uuid.UUID
	name string

	createdAt time.Time
	updatedAt time.Time
}

func (h scriptElement) ID() string {
	return h.id.String()
}

func (h scriptElement) UUID() uuid.UUID {
	return h.id
}

func (h scriptElement) Name() string {
	return h.name
}

func (h scriptElement) CreatedAt() time.Time {
	return h.createdAt
}

func (h scriptElement) UpdatedAt() time.Time {
	return h.updatedAt
}

func (h *scriptElement) update() {
	h.updatedAt = time.Now()
}
