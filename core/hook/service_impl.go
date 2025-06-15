package hook

import (
	"context"
	"errors"
	"io"
	"iter"
	"slices"
	"strings"
	"sync"
)

// hookServiceImpl is the concrete implementation of HookService.
type hookServiceImpl struct {
	mu      sync.RWMutex
	hooks   []*Hook
	content HookScriptStore
	dirty   bool
}

// NewHookService creates a new HookService with the given content store.
func NewHookService(content HookScriptStore) HookService {
	return &hookServiceImpl{
		hooks:   []*Hook{},
		content: content,
		dirty:   false,
	}
}

// List returns an iterator for all registered hooks.
func (s *hookServiceImpl) List() iter.Seq2[*Hook, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(*Hook, error) bool) {
		for _, h := range s.hooks {
			if !yield(h, nil) {
				break
			}
		}
	}
}

// Add registers a new hook and stores its script content.
func (s *hookServiceImpl) Add(ctx context.Context, h Hook, content io.Reader) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	scriptPath, err := s.content.SaveScript(ctx, h, content)
	if err != nil {
		return err
	}
	h.ScriptPath = scriptPath
	s.hooks = append(s.hooks, &h)
	s.dirty = true
	return nil
}

// Update updates the script content of an existing hook (by ID).
func (s *hookServiceImpl) Update(ctx context.Context, h Hook, content io.Reader) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, hook := range s.hooks {
		if hook.ID == h.ID {
			if hook.ScriptPath != "" {
				_ = s.content.RemoveScript(ctx, hook.ScriptPath)
			}
			scriptPath, err := s.content.SaveScript(ctx, h, content)
			if err != nil {
				return err
			}
			h.ScriptPath = scriptPath
			s.hooks[i] = &h
			s.dirty = true
			return nil
		}
	}
	return errors.New("hook not found")
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
	i, h, err := s.find(id)
	if err != nil {
		return err
	}
	if h.ScriptPath != "" {
		_ = s.content.RemoveScript(ctx, h.ScriptPath)
	}
	s.hooks = slices.Delete(s.hooks, i, i+1)
	s.dirty = true
	return nil
}

// Open opens the script content for a given hook.
func (s *hookServiceImpl) Open(ctx context.Context, id string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, hook, err := s.find(id)
	if err != nil {
		return nil, err
	}
	return s.content.OpenScript(ctx, hook.ScriptPath)
}

// Set replaces the list of hooks (used for loading from persistent storage).
func (s *hookServiceImpl) Set(seq iter.Seq2[*Hook, error]) error {
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
