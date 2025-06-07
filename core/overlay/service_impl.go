package overlay

import (
	"context"
	"io"
	"iter"
	"slices"
)

// serviceImpl is the default implementation of OverlayService.
// It manages overlays in memory and delegates content storage to ContentStore.
type serviceImpl struct {
	overlays     []Overlay
	contentStore ContentStore
	changed      bool
}

// NewOverlayService creates a new OverlayService with the given ContentStore.
func NewOverlayService(contentStore ContentStore) OverlayService {
	return &serviceImpl{
		overlays:     []Overlay{},
		contentStore: contentStore,
	}
}

func (s *serviceImpl) ListOverlays() iter.Seq2[*Overlay, error] {
	return func(yield func(*Overlay, error) bool) {
		for _, ov := range s.overlays {
			ov := ov
			if !yield(&ov, nil) {
				return
			}
		}
	}
}

func (s *serviceImpl) AddOverlay(ctx context.Context, ov Overlay, content io.Reader) error {
	location, err := s.contentStore.SaveContent(ctx, ov, content)
	if err != nil {
		return err
	}
	ov.ContentLocation = location
	s.overlays = append(s.overlays, ov)
	s.changed = true
	return nil
}

func (s *serviceImpl) RemoveOverlay(ctx context.Context, ov Overlay) error {
	idx := -1
	for i, v := range s.overlays {
		if v.RepoPattern == ov.RepoPattern &&
			v.ForInit == ov.ForInit &&
			v.RelativePath == ov.RelativePath {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil // Overlay not found
	}
	_ = s.contentStore.RemoveContent(ctx, s.overlays[idx].ContentLocation)
	s.overlays = slices.Delete(s.overlays, idx, idx+1)
	s.changed = true
	return nil
}

func (s *serviceImpl) OpenOverlayContent(ctx context.Context, ov Overlay) (io.ReadCloser, error) {
	for _, v := range s.overlays {
		if v.RepoPattern == ov.RepoPattern &&
			v.ForInit == ov.ForInit &&
			v.RelativePath == ov.RelativePath {
			return s.contentStore.OpenContent(ctx, v.ContentLocation)
		}
	}
	return nil, io.EOF
}

func (s *serviceImpl) HasChanges() bool {
	return s.changed
}

func (s *serviceImpl) MarkSaved() {
	s.changed = false
}

func (s *serviceImpl) SetOverlays(overlays iter.Seq2[*Overlay, error]) error {
	s.overlays = []Overlay{}
	for ov, err := range overlays {
		if err != nil {
			return err
		}
		s.overlays = append(s.overlays, *ov)
	}
	s.changed = false
	return nil
}

var _ OverlayService = (*serviceImpl)(nil)
