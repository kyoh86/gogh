package hook

import (
	"context"
	"errors"
	"iter"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// hookServiceImpl is the concrete implementation of HookService.
type hookServiceImpl struct {
	mu    sync.RWMutex
	hooks []*Hook
	dirty bool
}

// NewHookService creates a new HookService with the given content store.
func NewHookService() HookService {
	return &hookServiceImpl{
		hooks: []*Hook{},
		dirty: false,
	}
}

// List returns an iterator for all registered hooks.
func (s *hookServiceImpl) List() iter.Seq2[*Hook, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(*Hook, error) bool) {
		for _, h := range s.hooks {
			el := *h
			if !yield(&el, nil) {
				break
			}
		}
	}
}

// Add registers a new hook and stores its script content.
func (s *hookServiceImpl) Add(
	ctx context.Context,
	name string,
	repoPattern string,
	triggerEvent Event,
	operationType OperationType,
	operationID string,
) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hook := &Hook{
		ID:            uuid.NewString(),
		Name:          name,
		RepoPattern:   repoPattern,
		TriggerEvent:  triggerEvent,
		OperationType: operationType,
		OperationID:   operationID,
	}
	s.hooks = append(s.hooks, hook)
	s.dirty = true
	return hook.ID, nil
}

// Update updates the script content of an existing hook (by ID).
func (s *hookServiceImpl) Update(
	ctx context.Context,
	id string,
	name string,
	repoPattern string,
	triggerEvent Event,
	operationType OperationType,
	operationID string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, hook, err := s.find(id)
	if err != nil {
		return err
	}
	if name != "" {
		hook.Name = name
		s.dirty = true
	}
	if repoPattern != "" {
		hook.RepoPattern = repoPattern
		s.dirty = true
	}
	if triggerEvent != "" {
		hook.TriggerEvent = triggerEvent
		s.dirty = true
	}
	if operationType != "" {
		hook.OperationType = operationType
		s.dirty = true
	}
	if operationID != "" {
		hook.OperationID = operationID
		s.dirty = true
	}
	return nil
}

// find searches for a hook by its ID and returns its index and the hook itself.
// If an exact match is found, it returns that hook. If no exact match exists but a single
// hook ID matches the prefix, it returns that hook. Returns an error if the ID is empty,
// multiple hooks match the prefix, or no matching hook is found.
func (s *hookServiceImpl) find(id string) (int, *Hook, error) {
	// NOTE: this function is internal and should not be locked externally.
	if len(id) == 0 {
		return -1, nil, errors.New("hook ID cannot be empty")
	}
	matched := -1
	for i, h := range s.hooks {
		if h.ID == id {
			return i, h, nil
		}
		if strings.HasPrefix(h.ID, id) {
			if matched >= 0 {
				return -1, nil, errors.New("multiple hooks found with similar ID")
			}
			matched = i
		}
	}
	if matched >= 0 {
		return matched, s.hooks[matched], nil
	}
	return -1, nil, errors.New("hook not found")
}

// Get retrieves a hook by its ID.
func (s *hookServiceImpl) Get(ctx context.Context, id string) (*Hook, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, h, err := s.find(id)
	return h, err
}

// Remove removes a hook and its script content by ID.
func (s *hookServiceImpl) Remove(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, _, err := s.find(id)
	if err != nil {
		return err
	}
	s.hooks = slices.Delete(s.hooks, i, i+1)
	s.dirty = true
	return nil
}

// Load replaces the list of hooks (used for loading from persistent storage).
func (s *hookServiceImpl) Load(seq iter.Seq2[*Hook, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var hooks []*Hook
	for h, err := range seq {
		if err != nil {
			return err
		}
		hooks = append(hooks, h)
	}
	s.hooks = hooks
	s.dirty = true
	return nil
}

// HasChanges returns true if there are unsaved changes.
func (s *hookServiceImpl) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dirty
}

// MarkSaved marks the current state as saved (no unsaved changes).
func (s *hookServiceImpl) MarkSaved() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = false
}

var _ HookService = (*hookServiceImpl)(nil)
