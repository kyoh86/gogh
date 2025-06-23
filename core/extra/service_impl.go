package extra

import (
	"context"
	"errors"
	"iter"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// serviceImpl is the default implementation of ExtraService.
// It manages extras in memory.
type serviceImpl struct {
	mu         sync.RWMutex
	autoExtra  map[string]*Extra // key: repository string
	namedExtra map[string]*Extra // key: name
	byID       map[string]*Extra // key: id
	dirty      bool
}

// NewExtraService creates a new ExtraService.
func NewExtraService() ExtraService {
	return &serviceImpl{
		autoExtra:  make(map[string]*Extra),
		namedExtra: make(map[string]*Extra),
		byID:       make(map[string]*Extra),
		dirty:      false,
	}
}

// HasChanges returns true if there are unsaved changes
func (s *serviceImpl) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.dirty
}

// MarkSaved marks the service as saved
func (s *serviceImpl) MarkSaved() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = false
}

func (s *serviceImpl) AddAutoExtra(ctx context.Context, repo repository.Reference, source repository.Reference, items []Item) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := repo.String()
	if _, exists := s.autoExtra[key]; exists {
		return "", ErrExtraAlreadyExists
	}

	id := uuid.Must(uuid.NewRandom()).String()
	e := NewAutoExtra(id, repo, source, items, time.Now())

	s.autoExtra[key] = e
	s.byID[id] = e
	s.dirty = true

	return id, nil
}

func (s *serviceImpl) AddNamedExtra(ctx context.Context, name string, source repository.Reference, items []Item) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.namedExtra[name]; exists {
		return "", ErrExtraAlreadyExists
	}

	id := uuid.Must(uuid.NewRandom()).String()
	e := NewNamedExtra(id, name, source, items, time.Now())

	s.namedExtra[name] = e
	s.byID[id] = e
	s.dirty = true

	return id, nil
}

func (s *serviceImpl) GetAutoExtra(ctx context.Context, repo repository.Reference) (*Extra, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.autoExtra[repo.String()]
	if !exists {
		return nil, ErrExtraNotFound
	}
	return e, nil
}

func (s *serviceImpl) GetNamedExtra(ctx context.Context, name string) (*Extra, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.namedExtra[name]
	if !exists {
		return nil, ErrExtraNotFound
	}
	return e, nil
}

func (s *serviceImpl) Get(ctx context.Context, id string) (*Extra, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, exists := s.byID[id]
	if !exists {
		return nil, ErrExtraNotFound
	}
	return e, nil
}

func (s *serviceImpl) RemoveAutoExtra(ctx context.Context, repo repository.Reference) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := repo.String()
	e, exists := s.autoExtra[key]
	if !exists {
		return ErrExtraNotFound
	}

	delete(s.autoExtra, key)
	delete(s.byID, e.ID())
	s.dirty = true
	return nil
}

func (s *serviceImpl) RemoveNamedExtra(ctx context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.namedExtra[name]
	if !exists {
		return ErrExtraNotFound
	}

	delete(s.namedExtra, name)
	delete(s.byID, e.ID())
	s.dirty = true
	return nil
}

func (s *serviceImpl) Remove(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, exists := s.byID[id]
	if !exists {
		return ErrExtraNotFound
	}

	switch e.Type() {
	case TypeAuto:
		if repo := e.Repository(); repo != nil {
			delete(s.autoExtra, repo.String())
		}
	case TypeNamed:
		delete(s.namedExtra, e.Name())
	}

	delete(s.byID, id)
	s.dirty = true
	return nil
}

func (s *serviceImpl) List(ctx context.Context) iter.Seq2[*Extra, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return func(yield func(*Extra, error) bool) {
		for _, e := range s.byID {
			if !yield(e, nil) {
				return
			}
		}
	}
}

func (s *serviceImpl) ListByType(ctx context.Context, extraType Type) iter.Seq2[*Extra, error] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return func(yield func(*Extra, error) bool) {
		switch extraType {
		case TypeAuto:
			for _, e := range s.autoExtra {
				if !yield(e, nil) {
					return
				}
			}
		case TypeNamed:
			for _, e := range s.namedExtra {
				if !yield(e, nil) {
					return
				}
			}
		default:
			yield(nil, errors.New("unknown extra type"))
		}
	}
}

// Load loads extras from an iterator
func (s *serviceImpl) Load(extras iter.Seq2[*Extra, error]) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear existing data
	s.autoExtra = make(map[string]*Extra)
	s.namedExtra = make(map[string]*Extra)
	s.byID = make(map[string]*Extra)

	for e, err := range extras {
		if err != nil {
			return err
		}

		s.byID[e.ID()] = e

		switch e.Type() {
		case TypeAuto:
			if repo := e.Repository(); repo != nil {
				s.autoExtra[repo.String()] = e
			}
		case TypeNamed:
			s.namedExtra[e.Name()] = e
		}
	}

	s.dirty = false
	return nil
}
