package hook

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
)

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

type Hook struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	RepoPattern  string `json:"repoPattern"`  // Repository pattern (glob)
	TriggerEvent Event  `json:"triggerEvent"` // Event that triggers the hook

	OperationType OperationType `json:"operationType"`
	OperationID   string
}

func (h Hook) String() string {
	return string(h.TriggerEvent) + "@" + h.RepoPattern
}

func (h Hook) Match(ref repository.Reference, useCase Event) (bool, error) {
	if h.TriggerEvent != useCase {
		return false, nil
	}
	if h.RepoPattern == "" {
		return true, nil
	}
	return doublestar.Match(h.RepoPattern, ref.String())
}
