package hook

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/repository"
)

type Entry struct {
	Name          string
	RepoPattern   string
	TriggerEvent  Event
	OperationType OperationType
	OperationID   string
}

// Event defines the trigger of the hook, such as post-clone, post-fork, or post-create
type Event string

const (
	EventClone  Event = "post-clone"
	EventFork   Event = "post-fork"
	EventCreate Event = "post-create"
)

// OperationType defines the type of operation that the hook performs
type OperationType string

const (
	OperationTypeOverlay OperationType = "overlay"
	OperationTypeScript  OperationType = "script"
)

type Hook interface {
	ID() string
	UUID() uuid.UUID
	Name() string

	RepoPattern() string
	TriggerEvent() Event

	OperationType() OperationType
	OperationID() string

	Match(ref repository.Reference, event Event) (bool, error)
}

func NewHook(entry Entry) Hook {
	return hookElement{
		id:            uuid.Must(uuid.NewRandom()),
		name:          entry.Name,
		repoPattern:   entry.RepoPattern,
		triggerEvent:  entry.TriggerEvent,
		operationType: entry.OperationType,
		operationID:   entry.OperationID,
	}
}

type hookElement struct {
	id   uuid.UUID
	name string

	repoPattern  string
	triggerEvent Event

	operationType OperationType
	operationID   string
}

func (h hookElement) ID() string {
	return h.id.String()
}

func (h hookElement) UUID() uuid.UUID {
	return h.id
}

func (h hookElement) Name() string {
	return h.name
}

func (h hookElement) RepoPattern() string {
	return h.repoPattern
}

func (h hookElement) TriggerEvent() Event {
	return h.triggerEvent
}

func (h hookElement) OperationType() OperationType {
	return h.operationType
}

func (h hookElement) OperationID() string {
	return h.operationID
}

func (h hookElement) Match(ref repository.Reference, event Event) (bool, error) {
	if h.triggerEvent != event {
		return false, nil
	}
	if h.repoPattern == "" {
		return true, nil
	}
	return doublestar.Match(h.repoPattern, ref.String())
}
