package typ_test

import (
	"sync"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/typ"
)

// TestItem implements the ID() method required by Set
type TestItem struct {
	id    string
	value int
}

func (i TestItem) ID() string {
	return i.id
}

func TestNewSet(t *testing.T) {
	// Test empty constructor
	emptySet := testtarget.NewSet[string, TestItem]()
	if emptySet.Len() != 0 {
		t.Errorf("Expected empty set, got length %d", emptySet.Len())
	}

	// Test with initial items
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
		{id: "c", value: 3},
	}
	set := testtarget.NewSet(items...)
	if set.Len() != 3 {
		t.Errorf("Expected set with 3 items, got %d", set.Len())
	}

	// Test with duplicate items
	dupeItems := []TestItem{
		{id: "a", value: 1},
		{id: "a", value: 2}, // Duplicate ID
		{id: "b", value: 3},
	}
	dupeSet := testtarget.NewSet(dupeItems...)
	if dupeSet.Len() != 2 {
		t.Errorf("Expected set with 2 items after deduplication, got %d", dupeSet.Len())
	}
}

func TestSet_Add(t *testing.T) {
	set := testtarget.NewSet[string, TestItem]()

	// Test adding new item
	added := set.Add(TestItem{id: "a", value: 1})
	if !added || set.Len() != 1 {
		t.Errorf("Expected Add to return true and length 1, got %v and %d", added, set.Len())
	}

	// Test adding duplicate item
	added = set.Add(TestItem{id: "a", value: 2})
	if added || set.Len() != 1 {
		t.Errorf("Expected Add to return false and length to remain 1, got %v and %d", added, set.Len())
	}
}

func TestSet_Set(t *testing.T) {
	set := testtarget.NewSet(TestItem{id: "a", value: 1})

	// Test updating existing item
	updated := set.Set(TestItem{id: "a", value: 2})
	if !updated || set.Len() != 1 {
		t.Errorf("Expected Set to return true and length 1, got %v and %d", updated, set.Len())
	}

	item, exists := set.GetByID("a")
	if !exists || item.value != 2 {
		t.Errorf("Expected item with value 2, got %v", item)
	}

	// Test updating non-existing item
	updated = set.Set(TestItem{id: "b", value: 3})
	if updated || set.Len() != 1 {
		t.Errorf("Expected Set to return false for non-existing item, got %v and length %d", updated, set.Len())
	}
}

func TestSet_Remove(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
	}
	set := testtarget.NewSet(items...)

	// Test removing existing item
	removed := set.Remove(TestItem{id: "a", value: 99}) // value doesn't matter, only ID
	if !removed || set.Len() != 1 {
		t.Errorf("Expected Remove to return true and length 1, got %v and %d", removed, set.Len())
	}

	// Test removing non-existing item
	removed = set.Remove(TestItem{id: "c", value: 3})
	if removed || set.Len() != 1 {
		t.Errorf("Expected Remove to return false and length to remain 1, got %v and %d", removed, set.Len())
	}
}

func TestSet_RemoveByID(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
	}
	set := testtarget.NewSet(items...)

	// Test removing existing ID
	removed := set.RemoveByID("a")
	if !removed || set.Len() != 1 {
		t.Errorf("Expected RemoveByID to return true and length 1, got %v and %d", removed, set.Len())
	}

	// Test removing non-existing ID
	removed = set.RemoveByID("c")
	if removed || set.Len() != 1 {
		t.Errorf("Expected RemoveByID to return false and length to remain 1, got %v and %d", removed, set.Len())
	}
}

func TestSet_Has_HasID(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
	}
	set := testtarget.NewSet(items...)

	// Test Has with existing item
	if !set.Has(TestItem{id: "a", value: 99}) { // value doesn't matter
		t.Errorf("Expected Has to return true for existing ID")
	}

	// Test Has with non-existing item
	if set.Has(TestItem{id: "c", value: 3}) {
		t.Errorf("Expected Has to return false for non-existing ID")
	}

	// Test HasID with existing ID
	if !set.HasID("b") {
		t.Errorf("Expected HasID to return true for existing ID")
	}

	// Test HasID with non-existing ID
	if set.HasID("c") {
		t.Errorf("Expected HasID to return false for non-existing ID")
	}
}

func TestSet_At(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
		{id: "c", value: 3},
	}
	set := testtarget.NewSet(items...)

	// Test At with positive index
	item := set.At(1)
	if item.ID() != "b" || item.value != 2 {
		t.Errorf("Expected item with ID 'b' and value 2, got %v", item)
	}

	// Test At with negative index
	item = set.At(-1)
	if item.ID() != "c" || item.value != 3 {
		t.Errorf("Expected item with ID 'c' and value 3, got %v", item)
	}
}

func TestSet_GetByID(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
	}
	set := testtarget.NewSet(items...)

	// Test GetByID with existing ID
	item, exists := set.GetByID("b")
	if !exists || item.ID() != "b" || item.value != 2 {
		t.Errorf("Expected item with ID 'b' and value 2, got exists=%v, item=%v", exists, item)
	}

	// Test GetByID with non-existing ID
	item, exists = set.GetByID("c")
	if exists {
		t.Errorf("Expected GetByID to return false for non-existing ID, got %v", item)
	}
}

func TestSet_List(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
		{id: "c", value: 3},
	}
	set := testtarget.NewSet(items...)

	list := set.List()
	if len(list) != 3 {
		t.Errorf("Expected list with 3 items, got %d", len(list))
	}

	// Verify order is preserved
	expectedIDs := []string{"a", "b", "c"}
	for i, item := range list {
		if item.ID() != expectedIDs[i] {
			t.Errorf("Expected ID %s at position %d, got %s", expectedIDs[i], i, item.ID())
		}
	}
}

func TestSet_Iter(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
		{id: "c", value: 3},
	}
	set := testtarget.NewSet(items...)

	var result []TestItem
	for item := range set.Iter() {
		result = append(result, item)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 items from iterator, got %d", len(result))
	}
}

func TestSet_Clear_Len(t *testing.T) {
	items := []TestItem{
		{id: "a", value: 1},
		{id: "b", value: 2},
	}
	set := testtarget.NewSet(items...)

	if set.Len() != 2 {
		t.Errorf("Expected length 2, got %d", set.Len())
	}

	set.Clear()
	if set.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", set.Len())
	}

	if set.HasID("a") {
		t.Errorf("Expected set to be empty after Clear")
	}
}

func TestSet_Concurrency(t *testing.T) {
	set := testtarget.NewSet[string, TestItem]()
	var wg sync.WaitGroup

	// Test concurrent add operations
	for i := range 100 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			set.Add(TestItem{id: string(rune('a' + i%26)), value: i})
		}(i)
	}

	wg.Wait()

	// Maximum unique IDs should be 26 (a-z)
	if set.Len() > 26 {
		t.Errorf("Expected at most 26 items, got %d", set.Len())
	}
}
