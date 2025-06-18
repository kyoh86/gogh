package script

import (
	"context"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/kyoh86/gogh/v4/core/set"
)

// serviceImpl is the concrete implementation of HookService.
type serviceImpl struct {
	mu      sync.RWMutex
	scripts *set.Set[scriptElement]
	content ScriptSourceStore
	dirty   bool
}

// NewScriptService creates a new HookService with the given content store.
func NewScriptService(content ScriptSourceStore) ScriptService {
	return &serviceImpl{
		scripts: set.NewSet[scriptElement](),
		content: content,
		dirty:   false,
	}
}

// List returns an iterator for all registered scripts.
func (s *serviceImpl) List() iter.Seq2[Script, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(Script, error) bool) {
		for script := range s.scripts.Iter() {
			if !yield(script, nil) {
				break
			}
		}
	}
}

// Add registers a new script and stores its script content.
func (s *serviceImpl) Add(ctx context.Context, entry Entry) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	script := NewScript(entry).(scriptElement)
	if err := s.content.Save(ctx, script.ID(), entry.Content); err != nil {
		return "", err
	}
	s.scripts.Add(script)
	s.dirty = true
	return script.ID(), nil
}

// Update the script content of an existing script (by ID).
func (s *serviceImpl) Update(ctx context.Context, idlike string, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	script, err := s.scripts.GetBy(idlike)
	if err != nil {
		return fmt.Errorf("script not found: %w", err)
	}
	dirty := false
	if entry.Content != nil {
		if err := s.content.Save(ctx, script.ID(), entry.Content); err != nil {
			return err
		}
		dirty = true
	}
	if entry.Name != "" {
		script.name = entry.Name
		dirty = true
	}
	if dirty {
		script.update()
		s.scripts.Set(script)
		s.dirty = true
	}
	return nil
}

// Get retrieves a script by its ID.
func (s *serviceImpl) Get(ctx context.Context, idlike string) (Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.scripts.GetBy(idlike)
}

// Remove removes a script and its script content by ID.
func (s *serviceImpl) Remove(ctx context.Context, idlike string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	script, err := s.scripts.GetBy(idlike)
	if err != nil {
		return err
	}
	if err := s.content.Remove(ctx, script.ID()); err != nil {
		return err
	}
	if err := s.scripts.Remove(script); err != nil {
		return fmt.Errorf("remove script: %w", err)
	}
	s.dirty = true
	return nil
}

// Open opens the script content for a given script by its ID.
func (s *serviceImpl) Open(ctx context.Context, idlike string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	script, err := s.scripts.GetBy(idlike)
	if err != nil {
		return nil, err
	}
	return s.content.Open(ctx, script.ID())
}

// Load replaces the list of scripts (used for loading from persistent storage).
func (s *serviceImpl) Load(seq iter.Seq2[Script, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	scripts := set.NewSet[scriptElement]()
	for h, err := range seq {
		if err != nil {
			return err
		}
		scripts.Add(scriptElement{
			id:        h.UUID(),
			name:      h.Name(),
			createdAt: h.CreatedAt(),
			updatedAt: h.UpdatedAt(),
		})
	}
	s.scripts = scripts
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

var _ ScriptService = (*serviceImpl)(nil)
