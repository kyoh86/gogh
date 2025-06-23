package set_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/set"
)

// TestItem implements an item with UUID() method for testing
type TestItem struct {
	id   uuid.UUID
	name string
}

func (t TestItem) UUID() uuid.UUID {
	return t.id
}

func TestNewSet(t *testing.T) {
	t.Run("Empty set", func(t *testing.T) {
		s := set.NewSet[TestItem]()
		if s.Len() != 0 {
			t.Errorf("expected empty set, got %d items", s.Len())
		}
	})

	t.Run("Set with initial items", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		item1 := TestItem{id: id1, name: "item1"}
		item2 := TestItem{id: id2, name: "item2"}

		s := set.NewSet(item1, item2)
		if s.Len() != 2 {
			t.Errorf("expected 2 items, got %d", s.Len())
		}
	})

	t.Run("Set with duplicate items", func(t *testing.T) {
		id := uuid.New()
		item1 := TestItem{id: id, name: "item1"}
		item2 := TestItem{id: id, name: "item2"}

		s := set.NewSet(item1, item2)
		if s.Len() != 1 {
			t.Errorf("expected 1 item (duplicates removed), got %d", s.Len())
		}
	})
}

func TestAdd(t *testing.T) {
	t.Run("Add new item", func(t *testing.T) {
		s := set.NewSet[TestItem]()
		item := TestItem{id: uuid.New(), name: "test"}

		err := s.Add(item)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Len() != 1 {
			t.Errorf("expected 1 item, got %d", s.Len())
		}
	})

	t.Run("Add duplicate item", func(t *testing.T) {
		id := uuid.New()
		item := TestItem{id: id, name: "test"}
		s := set.NewSet(item)

		err := s.Add(item)
		if err == nil {
			t.Fatal("expected error for duplicate, got nil")
		}
		if !errors.Is(err, set.ErrDuplicated) {
			t.Errorf("expected ErrDuplicated, got %v", err)
		}
	})
}

func TestSet(t *testing.T) {
	t.Run("Set new item", func(t *testing.T) {
		s := set.NewSet[TestItem]()
		item := TestItem{id: uuid.New(), name: "test"}

		s.Set(item)
		if s.Len() != 1 {
			t.Errorf("expected 1 item, got %d", s.Len())
		}
	})

	t.Run("Set existing item (update)", func(t *testing.T) {
		id := uuid.New()
		item1 := TestItem{id: id, name: "original"}
		s := set.NewSet(item1)

		item2 := TestItem{id: id, name: "updated"}
		s.Set(item2)

		if s.Len() != 1 {
			t.Errorf("expected 1 item, got %d", s.Len())
		}

		retrieved, err := s.Get(id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved.name != "updated" {
			t.Errorf("expected updated name, got %s", retrieved.name)
		}
	})
}

func TestAt(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	item1 := TestItem{id: id1, name: "item1"}
	item2 := TestItem{id: id2, name: "item2"}
	s := set.NewSet(item1, item2)

	t.Run("Positive index", func(t *testing.T) {
		item := s.At(0)
		if item.id != id1 {
			t.Errorf("expected first item, got %v", item)
		}
	})

	t.Run("Negative index", func(t *testing.T) {
		item := s.At(-1)
		if item.id != id2 {
			t.Errorf("expected last item, got %v", item)
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("Remove existing item", func(t *testing.T) {
		item := TestItem{id: uuid.New(), name: "test"}
		s := set.NewSet(item)

		err := s.Remove(item)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Len() != 0 {
			t.Errorf("expected empty set, got %d items", s.Len())
		}
	})

	t.Run("Remove non-existing item", func(t *testing.T) {
		s := set.NewSet[TestItem]()
		item := TestItem{id: uuid.New(), name: "test"}

		err := s.Remove(item)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, set.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestRemoveBy(t *testing.T) {
	id := uuid.New()
	item := TestItem{id: id, name: "test"}
	t.Run("Remove by full UUID", func(t *testing.T) {
		s := set.NewSet(item)
		err := s.RemoveBy(id.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Len() != 0 {
			t.Errorf("expected empty set, got %d items", s.Len())
		}
	})

	t.Run("Remove by UUID prefix", func(t *testing.T) {
		s := set.NewSet(item)
		prefix := id.String()[:8]
		err := s.RemoveBy(prefix)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if s.Len() != 0 {
			t.Errorf("expected empty set, got %d items", s.Len())
		}
	})

	t.Run("Remove by empty ID", func(t *testing.T) {
		s := set.NewSet(item)
		err := s.RemoveBy("")
		if err == nil {
			t.Fatal("expected error for empty ID, got nil")
		}
	})

	t.Run("Remove by ambiguous prefix", func(t *testing.T) {
		// Create two items with similar prefixes
		id1 := uuid.MustParse("12345678-1234-1234-1234-123456789012")
		id2 := uuid.MustParse("12345678-5678-5678-5678-567890123456")
		item1 := TestItem{id: id1, name: "item1"}
		item2 := TestItem{id: id2, name: "item2"}
		s := set.NewSet(item1, item2)

		err := s.RemoveBy("12345678")
		if err == nil {
			t.Fatal("expected error for ambiguous prefix, got nil")
		}
		if !errors.Is(err, set.ErrMultipleFound) {
			t.Errorf("expected ErrMultipleFound, got %v", err)
		}
	})
}

func TestHas(t *testing.T) {
	item := TestItem{id: uuid.New(), name: "test"}
	s := set.NewSet(item)

	t.Run("Has existing item", func(t *testing.T) {
		if !s.Has(item) {
			t.Error("expected Has to return true for existing item")
		}
	})

	t.Run("Has non-existing item", func(t *testing.T) {
		nonExisting := TestItem{id: uuid.New(), name: "other"}
		if s.Has(nonExisting) {
			t.Error("expected Has to return false for non-existing item")
		}
	})
}

func TestHasBy(t *testing.T) {
	id := uuid.New()
	item := TestItem{id: id, name: "test"}
	s := set.NewSet(item)

	t.Run("HasBy full UUID", func(t *testing.T) {
		has, err := s.HasBy(id.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !has {
			t.Error("expected HasBy to return true for existing item")
		}
	})

	t.Run("HasBy UUID prefix", func(t *testing.T) {
		prefix := id.String()[:8]
		has, err := s.HasBy(prefix)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !has {
			t.Error("expected HasBy to return true for existing item")
		}
	})

	t.Run("HasBy non-existing", func(t *testing.T) {
		has, err := s.HasBy("non-existing")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if has {
			t.Error("expected HasBy to return false for non-existing item")
		}
	})
}

func TestGet(t *testing.T) {
	id := uuid.New()
	item := TestItem{id: id, name: "test"}
	s := set.NewSet(item)

	t.Run("Get existing item", func(t *testing.T) {
		retrieved, err := s.Get(id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved.id != id {
			t.Errorf("expected item with ID %v, got %v", id, retrieved.id)
		}
	})

	t.Run("Get non-existing item", func(t *testing.T) {
		_, err := s.Get(uuid.New())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, set.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestGetBy(t *testing.T) {
	id := uuid.New()
	item := TestItem{id: id, name: "test"}
	s := set.NewSet(item)

	t.Run("GetBy full UUID", func(t *testing.T) {
		retrieved, err := s.GetBy(id.String())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved.id != id {
			t.Errorf("expected item with ID %v, got %v", id, retrieved.id)
		}
	})

	t.Run("GetBy UUID prefix", func(t *testing.T) {
		prefix := id.String()[:8]
		retrieved, err := s.GetBy(prefix)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if retrieved.id != id {
			t.Errorf("expected item with ID %v, got %v", id, retrieved.id)
		}
	})
}

func TestList(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	item1 := TestItem{id: id1, name: "item1"}
	item2 := TestItem{id: id2, name: "item2"}
	s := set.NewSet(item1, item2)

	list := s.List()
	if len(list) != 2 {
		t.Errorf("expected 2 items, got %d", len(list))
	}

	// Check that items are in order
	if list[0].id != id1 || list[1].id != id2 {
		t.Error("items are not in expected order")
	}
}

func TestIter(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()
	item1 := TestItem{id: id1, name: "item1"}
	item2 := TestItem{id: id2, name: "item2"}
	s := set.NewSet(item1, item2)

	count := 0
	for item := range s.Iter() {
		count++
		if item.id != id1 && item.id != id2 {
			t.Errorf("unexpected item: %v", item)
		}
	}

	if count != 2 {
		t.Errorf("expected 2 iterations, got %d", count)
	}
}

func TestClear(t *testing.T) {
	item1 := TestItem{id: uuid.New(), name: "item1"}
	item2 := TestItem{id: uuid.New(), name: "item2"}
	s := set.NewSet(item1, item2)

	s.Clear()
	if s.Len() != 0 {
		t.Errorf("expected empty set after clear, got %d items", s.Len())
	}
}
