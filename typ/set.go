package typ

import (
	"iter"
	"slices"
	"sync"
)

type Set[T comparable, U interface{ ID() T }] struct {
	mu    sync.RWMutex
	ids   []T
	items map[T]U
}

func NewSet[T comparable, U interface{ ID() T }](sources ...U) *Set[T, U] {
	if len(sources) == 0 {
		return &Set[T, U]{
			items: make(map[T]U),
		}
	}
	items := make(map[T]U, len(sources))
	ids := make([]T, 0, len(sources))
	for _, item := range sources {
		id := item.ID()
		if _, exists := items[id]; !exists {
			items[id] = item
			ids = append(ids, id)
		}
	}
	return &Set[T, U]{items: items, ids: ids}
}

func (s *Set[T, U]) Add(item U) (added bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := item.ID()
	if _, exists := s.items[id]; exists {
		return false
	}
	s.items[id] = item
	s.ids = append(s.ids, id)
	return true
}

func (s *Set[T, U]) Set(item U) (updated bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := item.ID()
	if _, exists := s.items[id]; !exists {
		return false
	}
	s.items[id] = item
	return true
}

func (s *Set[T, U]) Remove(item U) (removed bool) {
	return s.RemoveByID(item.ID())
}

func (s *Set[T, U]) RemoveByID(id T) (removed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.items[id]; !exists {
		return false
	}
	delete(s.items, id)
	idx := slices.Index(s.ids, id)
	s.ids = slices.Delete(s.ids, idx, idx+1)
	return true
}

func (s *Set[T, U]) Has(item U) bool {
	return s.HasID(item.ID())
}

func (s *Set[T, U]) HasID(id T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[id]
	return exists
}

func (s *Set[T, U]) At(i int) U {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i < 0 {
		i += len(s.ids)
	}
	id := s.ids[i]
	return s.items[id]
}

func (s *Set[T, U]) GetByID(id T) (U, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, exists := s.items[id]
	return item, exists
}

func (s *Set[T, U]) List() []U {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]U, 0, len(s.items))
	for _, id := range s.ids {
		result = append(result, s.items[id])
	}
	return result
}

func (s *Set[T, U]) Iter() iter.Seq[U] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return func(yield func(U) bool) {
		for _, id := range s.ids {
			if !yield(s.items[id]) {
				break
			}
		}
	}
}

func (s *Set[T, U]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]U)
	s.ids = nil
}

func (s *Set[T, U]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
