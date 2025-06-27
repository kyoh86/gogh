package hook_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestNewHook(t *testing.T) {
	entry := hook.Entry{
		Name:          "test-hook",
		RepoPattern:   "github.com/owner/*",
		TriggerEvent:  hook.EventPostClone,
		OperationType: hook.OperationTypeOverlay,
		OperationID:   uuid.New(),
	}

	h := hook.NewHook(entry)

	if h.Name() != entry.Name {
		t.Errorf("Name() = %v, want %v", h.Name(), entry.Name)
	}
	if h.RepoPattern() != entry.RepoPattern {
		t.Errorf("RepoPattern() = %v, want %v", h.RepoPattern(), entry.RepoPattern)
	}
	if h.TriggerEvent() != entry.TriggerEvent {
		t.Errorf("TriggerEvent() = %v, want %v", h.TriggerEvent(), entry.TriggerEvent)
	}
	if h.OperationType() != entry.OperationType {
		t.Errorf("OperationType() = %v, want %v", h.OperationType(), entry.OperationType)
	}
	if h.OperationUUID() != entry.OperationID {
		t.Errorf("OperationID() = %v, want %v", h.OperationID(), entry.OperationID)
	}
	if h.ID() == "" {
		t.Error("ID() should not be empty")
	}
	if h.UUID() == uuid.Nil {
		t.Error("UUID() should not be nil")
	}
}

func TestConcreteHook(t *testing.T) {
	id := uuid.New()
	name := "concrete-hook"
	repoPattern := "github.com/test/**"
	triggerEvent := string(hook.EventPostFork)
	operationType := string(hook.OperationTypeScript)
	operationID := uuid.New()

	h := hook.ConcreteHook(id, name, repoPattern, triggerEvent, operationType, operationID)

	if h.ID() != id.String() {
		t.Errorf("ID() = %v, want %v", h.ID(), id.String())
	}
	if h.UUID() != id {
		t.Errorf("UUID() = %v, want %v", h.UUID(), id)
	}
	if h.Name() != name {
		t.Errorf("Name() = %v, want %v", h.Name(), name)
	}
	if h.RepoPattern() != repoPattern {
		t.Errorf("RepoPattern() = %v, want %v", h.RepoPattern(), repoPattern)
	}
	if h.TriggerEvent() != hook.EventPostFork {
		t.Errorf("TriggerEvent() = %v, want %v", h.TriggerEvent(), hook.EventPostFork)
	}
	if h.OperationType() != hook.OperationTypeScript {
		t.Errorf("OperationType() = %v, want %v", h.OperationType(), hook.OperationTypeScript)
	}
	if h.OperationUUID() != operationID {
		t.Errorf("OperationID() = %v, want %v", h.OperationID(), operationID)
	}
}

func TestHook_Match(t *testing.T) {
	testCases := []struct {
		name         string
		repoPattern  string
		triggerEvent hook.Event
		ref          repository.Reference
		event        hook.Event
		wantMatch    bool
		wantErr      bool
	}{
		{
			name:         "Exact match with event",
			repoPattern:  "github.com/owner/repo",
			triggerEvent: hook.EventPostClone,
			ref:          repository.NewReference("github.com", "owner", "repo"),
			event:        hook.EventPostClone,
			wantMatch:    true,
			wantErr:      false,
		},
		{
			name:         "Pattern match with wildcard",
			repoPattern:  "github.com/owner/*",
			triggerEvent: hook.EventPostCreate,
			ref:          repository.NewReference("github.com", "owner", "any-repo"),
			event:        hook.EventPostCreate,
			wantMatch:    true,
			wantErr:      false,
		},
		{
			name:         "Pattern match with double wildcard",
			repoPattern:  "github.com/**",
			triggerEvent: hook.EventPostFork,
			ref:          repository.NewReference("github.com", "any-owner", "any-repo"),
			event:        hook.EventPostFork,
			wantMatch:    true,
			wantErr:      false,
		},
		{
			name:         "Event mismatch",
			repoPattern:  "github.com/owner/repo",
			triggerEvent: hook.EventPostClone,
			ref:          repository.NewReference("github.com", "owner", "repo"),
			event:        hook.EventPostFork,
			wantMatch:    false,
			wantErr:      false,
		},
		{
			name:         "Pattern mismatch",
			repoPattern:  "github.com/owner/*",
			triggerEvent: hook.EventPostClone,
			ref:          repository.NewReference("github.com", "other", "repo"),
			event:        hook.EventPostClone,
			wantMatch:    false,
			wantErr:      false,
		},
		{
			name:         "Any event matches all",
			repoPattern:  "github.com/owner/repo",
			triggerEvent: hook.EventAny,
			ref:          repository.NewReference("github.com", "owner", "repo"),
			event:        hook.EventPostCreate,
			wantMatch:    true,
			wantErr:      false,
		},
		{
			name:         "Empty pattern matches all",
			repoPattern:  "",
			triggerEvent: hook.EventPostClone,
			ref:          repository.NewReference("github.com", "any", "repo"),
			event:        hook.EventPostClone,
			wantMatch:    true,
			wantErr:      false,
		},
		{
			name:         "Invalid pattern",
			repoPattern:  "[invalid",
			triggerEvent: hook.EventPostClone,
			ref:          repository.NewReference("github.com", "owner", "repo"),
			event:        hook.EventPostClone,
			wantMatch:    false,
			wantErr:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   tc.repoPattern,
				TriggerEvent:  tc.triggerEvent,
				OperationType: hook.OperationTypeOverlay,
			})

			match, err := h.Match(tc.ref, tc.event)
			if (err != nil) != tc.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if match != tc.wantMatch {
				t.Errorf("Match() = %v, want %v", match, tc.wantMatch)
			}
		})
	}
}

func TestEventConstants(t *testing.T) {
	// Test that event constants have expected values
	if hook.EventAny != "" {
		t.Errorf("EventAny = %v, want empty string", hook.EventAny)
	}
	if hook.EventPostClone != "post-clone" {
		t.Errorf("EventPostClone = %v, want 'post-clone'", hook.EventPostClone)
	}
	if hook.EventPostFork != "post-fork" {
		t.Errorf("EventPostFork = %v, want 'post-fork'", hook.EventPostFork)
	}
	if hook.EventPostCreate != "post-create" {
		t.Errorf("EventPostCreate = %v, want 'post-create'", hook.EventPostCreate)
	}
}

func TestOperationTypeConstants(t *testing.T) {
	// Test that operation type constants have expected values
	if hook.OperationTypeOverlay != "overlay" {
		t.Errorf("OperationTypeOverlay = %v, want 'overlay'", hook.OperationTypeOverlay)
	}
	if hook.OperationTypeScript != "script" {
		t.Errorf("OperationTypeScript = %v, want 'script'", hook.OperationTypeScript)
	}
}
