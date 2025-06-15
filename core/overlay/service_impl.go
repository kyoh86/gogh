package overlay

import (
	"context"
	"io"
	"iter"

	"github.com/kyoh86/gogh/v4/typ"
)

// serviceImpl is the default implementation of OverlayService.
// It manages overlays in memory and delegates content storage to ContentStore.
type serviceImpl struct {
	overlays     *typ.Set[string, Overlay]
	contentStore ContentStore
	changed      bool
}

// NewOverlayService creates a new OverlayService with the given ContentStore.
func NewOverlayService(contentStore ContentStore) OverlayService {
	return &serviceImpl{
		overlays:     typ.NewSet[string, Overlay](),
		contentStore: contentStore,
	}
}

func (s *serviceImpl) List() iter.Seq2[*Overlay, error] {
	return func(yield func(*Overlay, error) bool) {
		for ov := range s.overlays.Iter() {
			ov := ov
			if !yield(&ov, nil) {
				return
			}
		}
	}
}

func (s *serviceImpl) Add(ctx context.Context, ov Overlay, content io.Reader) error {
	location, err := s.contentStore.SaveContent(ctx, ov, content)
	if err != nil {
		return err
	}
	ov.ContentLocation = location
	s.overlays.Add(ov)
	s.changed = true
	return nil
}

func (s *serviceImpl) Remove(ctx context.Context, ov Overlay) error {
	id := ov.ID()
	found, exists := s.overlays.GetByID(id)
	if !exists {
		return nil
	}
	s.changed = true
	s.overlays.RemoveByID(id)
	return s.contentStore.RemoveContent(ctx, found.ContentLocation)
}

func (s *serviceImpl) Open(ctx context.Context, ov Overlay) (io.ReadCloser, error) {
	found, exists := s.overlays.GetByID(ov.ID())
	if exists {
		return s.contentStore.OpenContent(ctx, found.ContentLocation)
	}
	return nil, io.EOF
}

func (s *serviceImpl) HasChanges() bool {
	return s.changed
}

func (s *serviceImpl) MarkSaved() {
	s.changed = false
}

func (s *serviceImpl) Set(overlays []Overlay) error {
	newSet := typ.NewSet(overlays...)
	s.overlays = newSet
	s.changed = true
	return nil
}

var _ OverlayService = (*serviceImpl)(nil)
