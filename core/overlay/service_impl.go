package overlay

import (
	"context"
	"errors"
	"io"
	"iter"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// serviceImpl is the default implementation of OverlayService.
// It manages overlays in memory and delegates content storage to ContentStore.
type serviceImpl struct {
	mu           sync.RWMutex
	overlays     []*Overlay
	contentStore ContentStore
	dirty        bool
}

// NewOverlayService creates a new OverlayService with the given ContentStore.
func NewOverlayService(contentStore ContentStore) OverlayService {
	return &serviceImpl{
		contentStore: contentStore,
	}
}

func (s *serviceImpl) List() iter.Seq2[*Overlay, error] {
	return func(yield func(*Overlay, error) bool) {
		for _, ov := range s.overlays {
			el := *ov
			if !yield(&el, nil) {
				return
			}
		}
	}
}

func (s *serviceImpl) Add(
	ctx context.Context,
	name string,
	relativePath string,
	content io.Reader,
) (string, error) {
	ov := &Overlay{
		ID:           uuid.NewString(),
		Name:         name,
		RelativePath: relativePath,
	}
	if err := s.contentStore.Save(ctx, ov.ID, content); err != nil {
		return "", err
	}
	s.overlays = append(s.overlays, ov)
	s.dirty = true
	return ov.ID, nil
}

// find searches for a overlay by its ID and returns its index and the overlay itself.
// If an exact match is found, it returns that overlay. If no exact match exists but a single
// overlay ID matches the prefix, it returns that overlay. Returns an error if the ID is empty,
// multiple overlays match the prefix, or no matching overlay is found.
func (s *serviceImpl) find(id string) (int, *Overlay, error) {
	// NOTE: this function is internal and should not be locked externally.
	if len(id) == 0 {
		return -1, nil, errors.New("overlay ID cannot be empty")
	}
	matched := -1
	for i, h := range s.overlays {
		if h.ID == id {
			return i, h, nil
		}
		if strings.HasPrefix(h.ID, id) {
			if matched >= 0 {
				return -1, nil, errors.New("multiple overlays found with similar ID")
			}
			matched = i
		}
	}
	if matched >= 0 {
		return matched, s.overlays[matched], nil
	}
	return -1, nil, errors.New("overlay not found")
}

func (s *serviceImpl) Remove(ctx context.Context, overlayID string) error {
	i, _, err := s.find(overlayID)
	if err != nil {
		return err
	}
	s.overlays = slices.Delete(s.overlays, i, i+1)
	s.dirty = true
	return nil
}

func (s *serviceImpl) Open(ctx context.Context, overlayID string) (io.ReadCloser, error) {
	return s.contentStore.Open(ctx, overlayID)
}

func (s *serviceImpl) HasChanges() bool {
	return s.dirty
}

func (s *serviceImpl) MarkSaved() {
	s.dirty = false
}

// Load replaces the list of overlays (used for loading from persistent storage).
func (s *serviceImpl) Load(seq iter.Seq2[*Overlay, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	var overlays []*Overlay
	for h, err := range seq {
		if err != nil {
			return err
		}
		overlays = append(overlays, h)
	}
	s.overlays = overlays
	s.dirty = true
	return nil
}

var _ OverlayService = (*serviceImpl)(nil)
