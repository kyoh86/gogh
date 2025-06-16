package script

import (
	"context"
	"errors"
	"io"
	"iter"
	"slices"
	"strings"
	"sync"
)

// scriptServiceImpl is the concrete implementation of HookService.
type scriptServiceImpl struct {
	mu      sync.RWMutex
	scripts []*Script
	content ScriptStore
	dirty   bool
}

// NewScriptService creates a new HookService with the given content store.
func NewScriptService(content ScriptStore) ScriptService {
	return &scriptServiceImpl{
		scripts: []*Script{},
		content: content,
		dirty:   false,
	}
}

// List returns an iterator for all registered scripts.
func (s *scriptServiceImpl) List() iter.Seq2[*Script, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(*Script, error) bool) {
		for _, script := range s.scripts {
			el := *script
			if !yield(&el, nil) {
				break
			}
		}
	}
}

// Add registers a new script and stores its script content.
func (s *scriptServiceImpl) Add(ctx context.Context, name string, content io.Reader) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	script := &Script{
		Name: name,
	}
	script.init()
	if err := s.content.Save(ctx, script.ID, content); err != nil {
		return "", err
	}
	s.scripts = append(s.scripts, script)
	s.dirty = true
	return script.ID, nil
}

// Update the script content of an existing script (by ID).
func (s *scriptServiceImpl) Update(ctx context.Context, id, name string, content io.Reader) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, script, err := s.find(id)
	if err != nil {
		return err
	}
	dirty := false
	if content != nil {
		if err := s.content.Save(ctx, id, content); err != nil {
			return err
		}
		dirty = true
	}
	if name != "" {
		script.Name = name
		dirty = true
	}
	if dirty {
		script.update()
		s.scripts[i] = script
		s.dirty = true
	}
	return nil
}

// find searches for a script by its ID and returns its index and the script itself.
// If an exact match is found, it returns that script. If no exact match exists but a single
// script ID matches the prefix, it returns that script. Returns an error if the ID is empty,
// multiple scripts match the prefix, or no matching script is found.
func (s *scriptServiceImpl) find(id string) (int, *Script, error) {
	// NOTE: this function is internal and should not be locked externally.
	if len(id) == 0 {
		return -1, nil, errors.New("script ID cannot be empty")
	}
	matched := -1
	for i, h := range s.scripts {
		if h.ID == id {
			return i, h, nil
		}
		if strings.HasPrefix(h.ID, id) {
			if matched >= 0 {
				return -1, nil, errors.New("multiple scripts found with similar ID")
			}
			matched = i
		}
	}
	if matched >= 0 {
		return matched, s.scripts[matched], nil
	}
	return -1, nil, errors.New("script not found")
}

// Get retrieves a script by its ID.
func (s *scriptServiceImpl) Get(ctx context.Context, id string) (*Script, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, h, err := s.find(id)
	return h, err
}

// Remove removes a script and its script content by ID.
func (s *scriptServiceImpl) Remove(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, h, err := s.find(id)
	if err != nil {
		return err
	}
	if err := s.content.Remove(ctx, h.ID); err != nil {
		return err
	}
	s.scripts = slices.Delete(s.scripts, i, i+1)
	s.dirty = true
	return nil
}

// Open opens the script content for a given script.
func (s *scriptServiceImpl) Open(ctx context.Context, id string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, script, err := s.find(id)
	if err != nil {
		return nil, err
	}
	return s.content.Open(ctx, script.ID)
}

// Load replaces the list of scripts (used for loading from persistent storage).
func (s *scriptServiceImpl) Load(seq iter.Seq2[*Script, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var scripts []*Script
	for h, err := range seq {
		if err != nil {
			return err
		}
		scripts = append(scripts, h)
	}
	s.scripts = scripts
	s.dirty = true
	return nil
}

// HasChanges returns true if there are unsaved changes.
func (s *scriptServiceImpl) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dirty
}

// MarkSaved marks the current state as saved (no unsaved changes).
func (s *scriptServiceImpl) MarkSaved() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = false
}

var _ ScriptService = (*scriptServiceImpl)(nil)
