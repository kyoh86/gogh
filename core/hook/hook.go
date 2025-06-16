package hook

import (
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// Event defines the timing of hook execution
type Event string

const (
	EventAny          Event = ""
	EventAfterClone   Event = "after-clone"
	EventAfterOverlay Event = "after-overlay"
	EventNever        Event = "never"
)

// UseCase defines the type of use case for the hook
type UseCase string

const (
	UseCaseAny    UseCase = ""
	UseCaseClone  UseCase = "clone"
	UseCaseFork   UseCase = "fork"
	UseCaseCreate UseCase = "create"
	UseCaseNever  UseCase = "never"
)

// Target defines the target for the hook, including repository pattern, use case, and event
type Target struct {
	RepoPattern string  `json:"repoPattern"` // Repository pattern (glob)
	UseCase     UseCase `json:"useCase"`
	Event       Event   `json:"event"`
}

func (t Target) eventString() string {
	if t.UseCase == UseCaseAny {
		return "*:" + string(t.Event)
	}
	return string(t.UseCase) + ":" + string(t.Event)
}

func (t Target) String() string {
	if t.RepoPattern == "" {
		return t.eventString() + ":*"
	}
	return t.eventString() + ":" + t.RepoPattern
}

func (t Target) MatchEvent(useCase UseCase, event Event) bool {
	if t.UseCase == UseCaseAny || useCase == UseCaseAny {
		return t.Event == event
	}
	return t.UseCase == useCase && t.Event == event
}

func (t Target) Match(ref repository.Reference, useCase UseCase, event Event) (bool, error) {
	if !t.MatchEvent(useCase, event) {
		return false, nil
	}
	if t.RepoPattern == "" {
		return true, nil
	}
	return doublestar.Match(t.RepoPattern, ref.String())
}

type Hook struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Target Target `json:"target"`

	ScriptPath string `json:"scriptPath"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h *Hook) CreatedNow() {
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
		h.UpdatedAt = h.CreatedAt
	}
}

func (h *Hook) UpdatedNow() {
	h.UpdatedAt = time.Now()
}

func (h Hook) Match(ref repository.Reference, useCase UseCase, event Event) (bool, error) {
	return h.Target.Match(ref, useCase, event)
}
