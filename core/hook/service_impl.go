package hook

import (
	"context"
	"errors"
	"io"
	"iter"
	"slices"
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

// ListHooks returns an iterator for all registered hooks.
func (s *hookServiceImpl) ListHooks() iter.Seq2[*Hook, error] {
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

// AddHook registers a new hook and stores its script content.
func (s *hookServiceImpl) AddHook(ctx context.Context, h Hook, content io.Reader) error {
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

// UpdateHook updates the script content of an existing hook (by ID).
func (s *hookServiceImpl) UpdateHook(ctx context.Context, h Hook, content io.Reader) error {
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

// GetHookByID retrieves a hook by its ID.
func (s *hookServiceImpl) GetHookByID(ctx context.Context, id string) (*Hook, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, h := range s.hooks {
		if h.ID == id {
			return h, nil
		}
	}
	return nil, errors.New("hook not found")
}

// RemoveHook removes a hook and its script content by ID.
func (s *hookServiceImpl) RemoveHook(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, h := range s.hooks {
		if h.ID == id {
			if h.ScriptPath != "" {
				_ = s.content.RemoveScript(ctx, h.ScriptPath)
			}
			s.hooks = slices.Delete(s.hooks, i, i+1)
			s.dirty = true
			return nil
		}
	}
	return errors.New("hook not found")
}

// OpenHookScript opens the script content for a given hook.
func (s *hookServiceImpl) OpenHookScript(ctx context.Context, h Hook) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if h.ScriptPath == "" {
		return nil, errors.New("hook has no script path")
	}
	return s.content.OpenScript(ctx, h.ScriptPath)
}

// SetHooks replaces the list of hooks (used for loading from persistent storage).
func (s *hookServiceImpl) SetHooks(seq iter.Seq2[*Hook, error]) error {
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
