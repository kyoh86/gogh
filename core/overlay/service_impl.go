package overlay

import (
	"context"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/kyoh86/gogh/v4/core/set"
)

// serviceImpl is the default implementation of OverlayService.
// It manages overlays in memory and delegates content storage to ContentStore.
type serviceImpl struct {
	mu       sync.RWMutex
	overlays *set.Set[overlayElement]
	content  ContentStore
	dirty    bool
}

// NewOverlayService creates a new OverlayService with the given ContentStore.
func NewOverlayService(content ContentStore) OverlayService {
	return &serviceImpl{
		overlays: set.NewSet[overlayElement](),
		content:  content,
		dirty:    false,
	}
}

func (s *serviceImpl) List() iter.Seq2[Overlay, error] {
	return func(yield func(Overlay, error) bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		for ov := range s.overlays.Iter() {
			if !yield(ov, nil) {
				return
			}
		}
	}
}

// Add registers a new overlay and stores its content.
func (s *serviceImpl) Add(ctx context.Context, entry Entry) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	overlay := NewOverlay(entry).(overlayElement)
	if err := s.content.Save(ctx, overlay.ID(), entry.Content); err != nil {
		return "", err
	}
	if err := s.overlays.Add(overlay); err != nil {
		return "", fmt.Errorf("add overlay: %w", err)
	}
	s.dirty = true
	return overlay.ID(), nil
}

// Update the overlay content of an existing overlay (by ID).
func (s *serviceImpl) Update(ctx context.Context, idlike string, entry Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	overlay, err := s.overlays.GetBy(idlike)
	if err != nil {
		return fmt.Errorf("overlay not found: %w", err)
	}
	dirty := false
	if entry.Content != nil {
		if err := s.content.Save(ctx, overlay.ID(), entry.Content); err != nil {
			return err
		}
		dirty = true
	}
	if entry.Name != "" {
		overlay.name = entry.Name
		dirty = true
	}
	if entry.RelativePath != "" {
		overlay.relativePath = entry.RelativePath
		dirty = true
	}
	if dirty {
		s.overlays.Set(overlay)
		s.dirty = true
	}
	return nil
}

// Get retrieves a script by its ID.
func (s *serviceImpl) Get(ctx context.Context, idlike string) (Overlay, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.overlays.GetBy(idlike)
}

// Remove removes a overlay and its overlay content by ID.
func (s *serviceImpl) Remove(ctx context.Context, idlike string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	overlay, err := s.overlays.GetBy(idlike)
	if err != nil {
		return err
	}
	if err := s.content.Remove(ctx, overlay.ID()); err != nil {
		return err
	}
	if err := s.overlays.Remove(overlay); err != nil {
		return fmt.Errorf("remove overlay: %w", err)
	}
	s.dirty = true
	return nil
}

// Open retrieves the content of an overlay by its ID.
func (s *serviceImpl) Open(ctx context.Context, idlike string) (io.ReadCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	overlay, err := s.overlays.GetBy(idlike)
	if err != nil {
		return nil, err
	}
	return s.content.Open(ctx, overlay.ID())
}

// Load replaces the list of overlays (used for loading from persistent storage).
func (s *serviceImpl) Load(seq iter.Seq2[Overlay, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	overlays := set.NewSet[overlayElement]()
	for h, err := range seq {
		if err != nil {
			return err
		}
		if err := overlays.Add(overlayElement{
			id:           h.UUID(),
			name:         h.Name(),
			relativePath: h.RelativePath(),
		}); err != nil {
			return fmt.Errorf("add overlay: %w", err)
		}
	}
	s.overlays = overlays
	s.dirty = true
	return nil
}

func (s *serviceImpl) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dirty
}

func (s *serviceImpl) MarkSaved() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = false
}

var _ OverlayService = (*serviceImpl)(nil)
