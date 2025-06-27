package hook

import (
	"context"
	"fmt"
	"iter"
	"sync"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/set"
)

// serviceImpl is the concrete implementation of HookService.
type serviceImpl struct {
	mu    sync.RWMutex
	hooks *set.Set[hookElement]
	dirty bool
}

// NewHookService creates a new HookService with the given content store.
func NewHookService() HookService {
	return &serviceImpl{
		hooks: set.NewSet[hookElement](),
		dirty: false,
	}
}

// List returns an iterator for all registered hooks.
func (s *serviceImpl) List() iter.Seq2[Hook, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(Hook, error) bool) {
		for h := range s.hooks.Iter() {
			if !yield(h, nil) {
				break
			}
		}
	}
}

// ListFor returns an iterator for hooks that match the given repository reference.
func (s *serviceImpl) ListFor(reference repository.Reference, event Event) iter.Seq2[Hook, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(Hook, error) bool) {
		for h := range s.hooks.Iter() {
			match, err := h.Match(reference, event)
			if err != nil {
				yield(h, err)
				return
			}
			if !match {
				continue
			}
			if !yield(h, nil) {
				break
			}
		}
	}
}

// Add registers a new hook and stores its script content.
func (s *serviceImpl) Add(ctx context.Context, entry Entry) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hook := NewHook(entry).(hookElement)
	if err := s.hooks.Add(hook); err != nil {
		return "", fmt.Errorf("add hook: %w", err)
	}
	s.dirty = true
	return hook.ID(), nil
}

// Update updates the script content of an existing hook (by ID).
func (s *serviceImpl) Update(ctx context.Context, idlike string, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hook, err := s.hooks.GetBy(idlike)
	if err != nil {
		return fmt.Errorf("hook not found: %w", err)
	}
	dirty := false
	if entry.Name != "" {
		hook.name = entry.Name
		dirty = true
	}
	if entry.RepoPattern != "" {
		hook.repoPattern = entry.RepoPattern
		dirty = true
	}
	if entry.TriggerEvent != "" {
		hook.triggerEvent = entry.TriggerEvent
		dirty = true
	}
	if entry.OperationType != "" {
		hook.operationType = entry.OperationType
		dirty = true
	}
	var zero uuid.UUID
	if entry.OperationID != zero {
		hook.operationID = entry.OperationID
		dirty = true
	}
	if dirty {
		s.hooks.Set(hook)
		s.dirty = true
	}
	return nil
}

// Get retrieves a script by its ID.
func (s *serviceImpl) Get(ctx context.Context, idlike string) (Hook, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hooks.GetBy(idlike)
}

// Remove removes a hook and its hook content by ID.
func (s *serviceImpl) Remove(ctx context.Context, idlike string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hook, err := s.hooks.GetBy(idlike)
	if err != nil {
		return err
	}
	if err := s.hooks.Remove(hook); err != nil {
		return fmt.Errorf("remove hook: %w", err)
	}
	s.dirty = true
	return nil
}

// Load replaces the list of hooks (used for loading from persistent storage).
func (s *serviceImpl) Load(seq iter.Seq2[Hook, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hooks := set.NewSet[hookElement]()
	for h, err := range seq {
		if err != nil {
			return err
		}
		if err := hooks.Add(hookElement{
			id:            h.UUID(),
			name:          h.Name(),
			repoPattern:   h.RepoPattern(),
			triggerEvent:  h.TriggerEvent(),
			operationType: h.OperationType(),
			operationID:   h.OperationUUID(),
		}); err != nil {
			return fmt.Errorf("load hook: %w", err)
		}
	}
	s.hooks = hooks
	s.dirty = true
	return nil
}

// HasChanges returns true if there are unsaved changes.
func (s *serviceImpl) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dirty
}

// MarkSaved marks the current state as saved (no unsaved changes).
func (s *serviceImpl) MarkSaved() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = false
}

var _ HookService = (*serviceImpl)(nil)
