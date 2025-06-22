package set

import (
	"errors"
	"iter"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Set[T interface{ UUID() uuid.UUID }] struct {
	mu    sync.RWMutex
	ids   []uuid.UUID
	items map[uuid.UUID]T
}

var (
	ErrNotFound      = errors.New("not found")
	ErrDuplicated    = errors.New("duplicated id")
	ErrMultipleFound = errors.New("multiple items found with similar ID")
)

func NewSet[T interface{ UUID() uuid.UUID }](sources ...T) *Set[T] {
	if len(sources) == 0 {
		return &Set[T]{
			items: make(map[uuid.UUID]T),
		}
	}
	items := make(map[uuid.UUID]T, len(sources))
	ids := make([]uuid.UUID, 0, len(sources))
	for _, item := range sources {
		id := item.UUID()
		if _, exists := items[id]; !exists {
			items[id] = item
			ids = append(ids, id)
		}
	}
	return &Set[T]{items: items, ids: ids}
}

// find searches for an element by its ID-like string and returns its UUID.
func (s *Set[T]) find(idlike string) (uuid.UUID, error) {
	// NOTE: this function is internal and should not be locked externally.
	if len(idlike) == 0 {
		return uuid.UUID{}, errors.New("ID cannot be empty")
	}
	id, err := uuid.Parse(idlike)
	if err == nil {
		_, exists := s.items[id]
		if exists {
			return id, nil
		}
		return uuid.UUID{}, ErrNotFound
	}
	var matched uuid.UUID
	var found bool
	for _, id := range s.ids {
		if strings.HasPrefix(id.String(), idlike) {
			if found {
				return uuid.UUID{}, ErrMultipleFound
			}
			matched = id
			found = true
		}
	}
	if found {
		return matched, nil
	}
	return uuid.UUID{}, ErrNotFound
}

func (s *Set[T]) Add(item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := item.UUID()
	if _, exists := s.items[id]; exists {
		return ErrDuplicated
	}
	s.items[id] = item
	s.ids = append(s.ids, id)
	return nil
}

func (s *Set[T]) Set(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := item.UUID()
	if _, exists := s.items[id]; !exists {
		s.items[id] = item
		s.ids = append(s.ids, id)
		return
	}
	s.items[id] = item
}

func (s *Set[T]) At(i int) T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i < 0 {
		i += len(s.ids)
	}
	id := s.ids[i]
	return s.items[id]
}

func (s *Set[T]) remove(item T) error {
	// NOTE: this function is internal and should not be locked externally.
	id := item.UUID()
	if _, exists := s.items[id]; !exists {
		return ErrNotFound
	}
	delete(s.items, id)
	idx := slices.Index(s.ids, id)
	s.ids = slices.Delete(s.ids, idx, idx+1)
	return nil
}

func (s *Set[T]) Remove(item T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.remove(item)
}

func (s *Set[T]) RemoveBy(idlike string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, err := s.find(idlike)
	if err != nil {
		return err
	}
	return s.remove(s.items[id])
}

func (s *Set[T]) has(id uuid.UUID) bool {
	// NOTE: this function is internal and should not be locked externally.
	_, exists := s.items[id]
	return exists
}

func (s *Set[T]) Has(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.has(item.UUID())
}

func (s *Set[T]) HasBy(idlike string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, err := s.find(idlike)
	if err != nil {
		return false, err
	}
	return s.has(id), nil
}

func (s *Set[T]) get(id uuid.UUID) (T, error) {
	item, exists := s.items[id]
	if !exists {
		return item, ErrNotFound
	}
	return item, nil
}

func (s *Set[T]) Get(id uuid.UUID) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.get(id)
}

func (s *Set[T]) GetBy(idlike string) (t T, retErr error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, err := s.find(idlike)
	if err != nil {
		retErr = err
		return
	}
	return s.get(id)
}

func (s *Set[T]) List() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]T, 0, len(s.items))
	for _, id := range s.ids {
		result = append(result, s.items[id])
	}
	return result
}

func (s *Set[T]) Iter() iter.Seq[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(T) bool) {
		for _, id := range s.ids {
			if !yield(s.items[id]) {
				break
			}
		}
	}
}

func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[uuid.UUID]T)
	s.ids = nil
}

func (s *Set[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
