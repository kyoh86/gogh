package hook_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestNewHookService(t *testing.T) {
	service := hook.NewHookService()
	if service == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestHookService_Add(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	entry := hook.Entry{
		Name:          "test-hook",
		RepoPattern:   "github.com/kyoh86/*",
		TriggerEvent:  hook.EventPostClone,
		OperationType: hook.OperationTypeOverlay,
		OperationID:   "overlay-123",
	}

	// Test successful addition
	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Verify the hook was added
	h, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get hook: %v", err)
	}
	if h.ID() != id {
		t.Errorf("expected ID %s, got %s", id, h.ID())
	}
	if h.Name() != entry.Name {
		t.Errorf("expected name %s, got %s", entry.Name, h.Name())
	}
	if h.RepoPattern() != entry.RepoPattern {
		t.Errorf("expected pattern %s, got %s", entry.RepoPattern, h.RepoPattern())
	}
	if h.TriggerEvent() != entry.TriggerEvent {
		t.Errorf("expected event %s, got %s", entry.TriggerEvent, h.TriggerEvent())
	}
	if h.OperationType() != entry.OperationType {
		t.Errorf("expected operation type %s, got %s", entry.OperationType, h.OperationType())
	}
	if h.OperationID() != entry.OperationID {
		t.Errorf("expected operation ID %s, got %s", entry.OperationID, h.OperationID())
	}
}

func TestHookService_Update(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	// Add initial hook
	entry := hook.Entry{
		Name:          "test-hook",
		RepoPattern:   "github.com/kyoh86/*",
		TriggerEvent:  hook.EventPostClone,
		OperationType: hook.OperationTypeOverlay,
		OperationID:   "overlay-123",
	}

	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add hook: %v", err)
	}

	// Test updating fields
	updateEntry := hook.Entry{
		Name:          "updated-hook",
		RepoPattern:   "github.com/updated/*",
		TriggerEvent:  hook.EventPostFork,
		OperationType: hook.OperationTypeScript,
		OperationID:   "script-456",
	}

	err = service.Update(ctx, id, updateEntry)
	if err != nil {
		t.Fatalf("failed to update hook: %v", err)
	}

	// Verify the updates
	h, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get updated hook: %v", err)
	}
	if h.Name() != updateEntry.Name {
		t.Errorf("expected name %s, got %s", updateEntry.Name, h.Name())
	}
	if h.RepoPattern() != updateEntry.RepoPattern {
		t.Errorf("expected pattern %s, got %s", updateEntry.RepoPattern, h.RepoPattern())
	}
	if h.TriggerEvent() != updateEntry.TriggerEvent {
		t.Errorf("expected event %s, got %s", updateEntry.TriggerEvent, h.TriggerEvent())
	}
	if h.OperationType() != updateEntry.OperationType {
		t.Errorf("expected operation type %s, got %s", updateEntry.OperationType, h.OperationType())
	}
	if h.OperationID() != updateEntry.OperationID {
		t.Errorf("expected operation ID %s, got %s", updateEntry.OperationID, h.OperationID())
	}

	// Test partial update
	partialUpdate := hook.Entry{
		Name: "partially-updated",
	}
	err = service.Update(ctx, id, partialUpdate)
	if err != nil {
		t.Fatalf("failed to partially update hook: %v", err)
	}

	h, err = service.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get partially updated hook: %v", err)
	}
	if h.Name() != "partially-updated" {
		t.Errorf("expected name %s, got %s", "partially-updated", h.Name())
	}
	// Other fields should remain unchanged
	if h.RepoPattern() != updateEntry.RepoPattern {
		t.Errorf("pattern should not have changed, got %s", h.RepoPattern())
	}

	// Test updating non-existent hook
	err = service.Update(ctx, "non-existent", updateEntry)
	if err == nil {
		t.Error("expected error when updating non-existent hook")
	}
}

func TestHookService_Get(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	// Add a hook
	entry := hook.Entry{
		Name:          "test-hook",
		RepoPattern:   "github.com/kyoh86/*",
		TriggerEvent:  hook.EventPostClone,
		OperationType: hook.OperationTypeOverlay,
		OperationID:   "overlay-123",
	}

	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add hook: %v", err)
	}

	// Test getting by full ID
	h, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get hook by ID: %v", err)
	}
	if h.ID() != id {
		t.Errorf("expected ID %s, got %s", id, h.ID())
	}

	// Test getting by partial ID (prefix)
	if len(id) > 8 {
		partialID := id[:8]
		h, err = service.Get(ctx, partialID)
		if err != nil {
			t.Fatalf("failed to get hook by partial ID: %v", err)
		}
		if h.ID() != id {
			t.Errorf("expected ID %s, got %s", id, h.ID())
		}
	}

	// Test getting non-existent hook
	_, err = service.Get(ctx, "non-existent")
	if err == nil {
		t.Error("expected error when getting non-existent hook")
	}
}

func TestHookService_Remove(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	// Add a hook
	entry := hook.Entry{
		Name:          "test-hook",
		RepoPattern:   "github.com/kyoh86/*",
		TriggerEvent:  hook.EventPostClone,
		OperationType: hook.OperationTypeOverlay,
		OperationID:   "overlay-123",
	}

	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add hook: %v", err)
	}

	// Test removing existing hook
	err = service.Remove(ctx, id)
	if err != nil {
		t.Fatalf("failed to remove hook: %v", err)
	}

	// Verify it was removed
	_, err = service.Get(ctx, id)
	if err == nil {
		t.Error("expected error when getting removed hook")
	}

	// Test removing non-existent hook
	err = service.Remove(ctx, "non-existent")
	if err == nil {
		t.Error("expected error when removing non-existent hook")
	}
}

func TestHookService_List(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	// Add multiple hooks
	entries := []hook.Entry{
		{
			Name:          "hook1",
			RepoPattern:   "github.com/kyoh86/*",
			TriggerEvent:  hook.EventPostClone,
			OperationType: hook.OperationTypeOverlay,
			OperationID:   "overlay-1",
		},
		{
			Name:          "hook2",
			RepoPattern:   "github.com/golang/*",
			TriggerEvent:  hook.EventPostFork,
			OperationType: hook.OperationTypeScript,
			OperationID:   "script-1",
		},
		{
			Name:          "hook3",
			RepoPattern:   "*",
			TriggerEvent:  hook.EventPostCreate,
			OperationType: hook.OperationTypeOverlay,
			OperationID:   "overlay-2",
		},
	}

	var ids []string
	for _, entry := range entries {
		id, err := service.Add(ctx, entry)
		if err != nil {
			t.Fatalf("failed to add hook: %v", err)
		}
		ids = append(ids, id)
	}

	// List all hooks
	var count int
	foundIDs := make(map[string]bool)
	for h, err := range service.List() {
		if err != nil {
			t.Fatalf("unexpected error in List: %v", err)
		}
		foundIDs[h.ID()] = true
		count++
	}

	if count != len(entries) {
		t.Errorf("expected %d hooks, got %d", len(entries), count)
	}

	// Verify all IDs are present
	for _, id := range ids {
		if !foundIDs[id] {
			t.Errorf("expected ID %s in list", id)
		}
	}
}

func TestHookService_ListFor(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	// Add hooks with different patterns and events
	entries := []hook.Entry{
		{
			Name:          "specific-repo",
			RepoPattern:   "github.com/kyoh86/gogh",
			TriggerEvent:  hook.EventPostClone,
			OperationType: hook.OperationTypeOverlay,
			OperationID:   "overlay-1",
		},
		{
			Name:          "wildcard-owner",
			RepoPattern:   "github.com/kyoh86/*",
			TriggerEvent:  hook.EventPostClone,
			OperationType: hook.OperationTypeScript,
			OperationID:   "script-1",
		},
		{
			Name:          "any-repo",
			RepoPattern:   "",
			TriggerEvent:  hook.EventPostFork,
			OperationType: hook.OperationTypeOverlay,
			OperationID:   "overlay-2",
		},
		{
			Name:          "different-event",
			RepoPattern:   "github.com/kyoh86/gogh",
			TriggerEvent:  hook.EventPostCreate,
			OperationType: hook.OperationTypeScript,
			OperationID:   "script-2",
		},
		{
			Name:          "any-event",
			RepoPattern:   "github.com/kyoh86/gogh",
			TriggerEvent:  hook.EventAny,
			OperationType: hook.OperationTypeOverlay,
			OperationID:   "overlay-3",
		},
	}

	for _, entry := range entries {
		_, err := service.Add(ctx, entry)
		if err != nil {
			t.Fatalf("failed to add hook: %v", err)
		}
	}

	// Test ListFor with specific repo and event
	ref := repository.NewReference("github.com", "kyoh86", "gogh")
	var matchingHooks []string
	for h, err := range service.ListFor(ref, hook.EventPostClone) {
		if err != nil {
			t.Fatalf("unexpected error in ListFor: %v", err)
		}
		matchingHooks = append(matchingHooks, h.Name())
	}

	// Should match: specific-repo, wildcard-owner, any-event
	expectedMatches := []string{"specific-repo", "wildcard-owner", "any-event"}
	if len(matchingHooks) != len(expectedMatches) {
		t.Errorf("expected %d matching hooks, got %d", len(expectedMatches), len(matchingHooks))
	}

	for _, expected := range expectedMatches {
		found := false
		for _, actual := range matchingHooks {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected hook %s to match", expected)
		}
	}

	// Test with different event
	matchingHooks = nil
	for h, err := range service.ListFor(ref, hook.EventPostFork) {
		if err != nil {
			t.Fatalf("unexpected error in ListFor: %v", err)
		}
		matchingHooks = append(matchingHooks, h.Name())
	}

	// Should match: any-repo, any-event
	expectedMatches = []string{"any-repo", "any-event"}
	if len(matchingHooks) != len(expectedMatches) {
		t.Errorf("expected %d matching hooks for fork event, got %d", len(expectedMatches), len(matchingHooks))
	}
}

func TestHookService_Load(t *testing.T) {
	service := hook.NewHookService()

	// Create hooks to load
	hooks := []hook.Hook{
		hook.ConcreteHook(
			uuid.New(),
			"loaded-hook-1",
			"github.com/kyoh86/*",
			string(hook.EventPostClone),
			string(hook.OperationTypeOverlay),
			"overlay-1",
		),
		hook.ConcreteHook(
			uuid.New(),
			"loaded-hook-2",
			"github.com/golang/*",
			string(hook.EventPostFork),
			string(hook.OperationTypeScript),
			"script-1",
		),
	}

	// Load hooks
	loadSeq := func(yield func(hook.Hook, error) bool) {
		for _, h := range hooks {
			if !yield(h, nil) {
				return
			}
		}
	}

	err := service.Load(loadSeq)
	if err != nil {
		t.Fatalf("failed to load hooks: %v", err)
	}

	// Verify hooks were loaded
	var count int
	for h, err := range service.List() {
		if err != nil {
			t.Fatalf("unexpected error listing loaded hooks: %v", err)
		}
		count++

		// Find matching hook
		found := false
		for _, original := range hooks {
			if h.ID() == original.ID() {
				found = true
				if h.Name() != original.Name() {
					t.Errorf("expected name %s, got %s", original.Name(), h.Name())
				}
				break
			}
		}
		if !found {
			t.Errorf("unexpected hook with ID %s", h.ID())
		}
	}

	if count != len(hooks) {
		t.Errorf("expected %d hooks after load, got %d", len(hooks), count)
	}

	// Test loading with error
	errorSeq := func(yield func(hook.Hook, error) bool) {
		yield(nil, errors.New("load error"))
	}

	err = service.Load(errorSeq)
	if err == nil {
		t.Error("expected error when loading with error")
	}
}

func TestHookService_HasChanges(t *testing.T) {
	ctx := context.Background()
	service := hook.NewHookService()

	// Initially no changes
	if service.HasChanges() {
		t.Error("expected no changes initially")
	}

	// Add hook
	entry := hook.Entry{
		Name:          "test-hook",
		RepoPattern:   "github.com/kyoh86/*",
		TriggerEvent:  hook.EventPostClone,
		OperationType: hook.OperationTypeOverlay,
		OperationID:   "overlay-123",
	}

	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add hook: %v", err)
	}

	// Should have changes after add
	if !service.HasChanges() {
		t.Error("expected changes after add")
	}

	// Mark saved
	service.MarkSaved()
	if service.HasChanges() {
		t.Error("expected no changes after MarkSaved")
	}

	// Update hook
	err = service.Update(ctx, id, hook.Entry{Name: "updated"})
	if err != nil {
		t.Fatalf("failed to update hook: %v", err)
	}

	// Should have changes after update
	if !service.HasChanges() {
		t.Error("expected changes after update")
	}

	// Mark saved again
	service.MarkSaved()

	// Remove hook
	err = service.Remove(ctx, id)
	if err != nil {
		t.Fatalf("failed to remove hook: %v", err)
	}

	// Should have changes after remove
	if !service.HasChanges() {
		t.Error("expected changes after remove")
	}
}

func TestHookService_Match(t *testing.T) {
	tests := []struct {
		name      string
		hook      hook.Hook
		ref       repository.Reference
		event     hook.Event
		wantMatch bool
		wantErr   bool
	}{
		{
			name: "exact pattern match with matching event",
			hook: hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   "github.com/kyoh86/gogh",
				TriggerEvent:  hook.EventPostClone,
				OperationType: hook.OperationTypeOverlay,
				OperationID:   "overlay-1",
			}),
			ref:       repository.NewReference("github.com", "kyoh86", "gogh"),
			event:     hook.EventPostClone,
			wantMatch: true,
		},
		{
			name: "wildcard pattern match",
			hook: hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   "github.com/kyoh86/*",
				TriggerEvent:  hook.EventPostClone,
				OperationType: hook.OperationTypeOverlay,
				OperationID:   "overlay-1",
			}),
			ref:       repository.NewReference("github.com", "kyoh86", "dotfiles"),
			event:     hook.EventPostClone,
			wantMatch: true,
		},
		{
			name: "empty pattern matches any repo",
			hook: hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   "",
				TriggerEvent:  hook.EventPostClone,
				OperationType: hook.OperationTypeOverlay,
				OperationID:   "overlay-1",
			}),
			ref:       repository.NewReference("gitlab.com", "user", "project"),
			event:     hook.EventPostClone,
			wantMatch: true,
		},
		{
			name: "event mismatch",
			hook: hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   "github.com/kyoh86/gogh",
				TriggerEvent:  hook.EventPostClone,
				OperationType: hook.OperationTypeOverlay,
				OperationID:   "overlay-1",
			}),
			ref:       repository.NewReference("github.com", "kyoh86", "gogh"),
			event:     hook.EventPostFork,
			wantMatch: false,
		},
		{
			name: "any event matches all",
			hook: hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   "github.com/kyoh86/gogh",
				TriggerEvent:  hook.EventAny,
				OperationType: hook.OperationTypeOverlay,
				OperationID:   "overlay-1",
			}),
			ref:       repository.NewReference("github.com", "kyoh86", "gogh"),
			event:     hook.EventPostCreate,
			wantMatch: true,
		},
		{
			name: "pattern mismatch",
			hook: hook.NewHook(hook.Entry{
				Name:          "test",
				RepoPattern:   "github.com/golang/*",
				TriggerEvent:  hook.EventPostClone,
				OperationType: hook.OperationTypeOverlay,
				OperationID:   "overlay-1",
			}),
			ref:       repository.NewReference("github.com", "kyoh86", "gogh"),
			event:     hook.EventPostClone,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := tt.hook.Match(tt.ref, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if match != tt.wantMatch {
				t.Errorf("Match() = %v, want %v", match, tt.wantMatch)
			}
		})
	}
}

func TestHookConstants(t *testing.T) {
	// Test Event constants
	if hook.EventAny != "" {
		t.Errorf("EventAny should be empty string, got %q", hook.EventAny)
	}
	if hook.EventPostClone != "post-clone" {
		t.Errorf("EventPostClone = %q, want %q", hook.EventPostClone, "post-clone")
	}
	if hook.EventPostFork != "post-fork" {
		t.Errorf("EventPostFork = %q, want %q", hook.EventPostFork, "post-fork")
	}
	if hook.EventPostCreate != "post-create" {
		t.Errorf("EventPostCreate = %q, want %q", hook.EventPostCreate, "post-create")
	}

	// Test OperationType constants
	if hook.OperationTypeOverlay != "overlay" {
		t.Errorf("OperationTypeOverlay = %q, want %q", hook.OperationTypeOverlay, "overlay")
	}
	if hook.OperationTypeScript != "script" {
		t.Errorf("OperationTypeScript = %q, want %q", hook.OperationTypeScript, "script")
	}
}
