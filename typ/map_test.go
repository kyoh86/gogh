package typ_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/typ"
)

func TestMapSet(t *testing.T) {
	// Create a new empty map
	var m testtarget.Map[string, int]

	// Test setting a value on nil map (should initialize the map)
	m.Set("one", 1)

	// Verify the value was set
	if v, ok := m["one"]; !ok || v != 1 {
		t.Errorf("Expected map to contain key 'one' with value 1, got %v, exists: %v", v, ok)
	}

	// Test overwriting an existing value
	m.Set("one", 100)

	// Verify the value was updated
	if v, ok := m["one"]; !ok || v != 100 {
		t.Errorf("Expected map to contain key 'one' with updated value 100, got %v, exists: %v", v, ok)
	}

	// Test setting another value
	m.Set("two", 2)

	// Verify both values exist
	if v, ok := m["one"]; !ok || v != 100 {
		t.Errorf("Expected map to still contain key 'one' with value 100, got %v, exists: %v", v, ok)
	}
	if v, ok := m["two"]; !ok || v != 2 {
		t.Errorf("Expected map to contain key 'two' with value 2, got %v, exists: %v", v, ok)
	}
}

func TestMapDelete(t *testing.T) {
	// Test delete on nil map (should not panic)
	var nilMap testtarget.Map[string, int]
	nilMap.Delete("nonexistent") // Should not panic

	// Create a map with data
	m := testtarget.Map[string, int]{
		"one": 1,
		"two": 2,
	}

	// Test deleting an existing key
	m.Delete("one")

	// Verify the key was deleted
	if _, ok := m["one"]; ok {
		t.Errorf("Expected key 'one' to be deleted, but it still exists")
	}

	// Verify other keys remain
	if v, ok := m["two"]; !ok || v != 2 {
		t.Errorf("Expected key 'two' to remain with value 2, got %v, exists: %v", v, ok)
	}

	// Test deleting a non-existent key (should not panic)
	m.Delete("nonexistent") // Should not panic

	// Verify the map is unchanged
	if len(m) != 1 {
		t.Errorf("Expected map to have 1 element after deleting non-existent key, got %d", len(m))
	}
}

func TestMapHas(t *testing.T) {
	// Test Has on nil map
	var nilMap testtarget.Map[string, int]
	if nilMap.Has("any") {
		t.Error("Expected Has to return false for nil map, got true")
	}

	// Create a map with data
	m := testtarget.Map[string, int]{
		"one": 1,
		"two": 2,
	}

	// Test Has with existing keys
	if !m.Has("one") {
		t.Error("Expected Has to return true for existing key 'one', got false")
	}
	if !m.Has("two") {
		t.Error("Expected Has to return true for existing key 'two', got false")
	}

	// Test Has with non-existent key
	if m.Has("three") {
		t.Error("Expected Has to return false for non-existent key 'three', got true")
	}

	// Test with zero value
	m["zero"] = 0
	if !m.Has("zero") {
		t.Error("Expected Has to return true for key with zero value, got false")
	}
}

func TestMapTryGet(t *testing.T) {
	// Test TryGet on nil map
	var nilMap testtarget.Map[string, int]
	if v, ok := nilMap.TryGet("any"); ok || v != 0 {
		t.Errorf("Expected TryGet to return (0, false) for nil map, got (%v, %v)", v, ok)
	}

	// Create a map with data
	m := testtarget.Map[string, int]{
		"one":  1,
		"two":  2,
		"zero": 0,
	}

	// Test TryGet with existing keys
	if v, ok := m.TryGet("one"); !ok || v != 1 {
		t.Errorf("Expected TryGet to return (1, true) for key 'one', got (%v, %v)", v, ok)
	}
	if v, ok := m.TryGet("two"); !ok || v != 2 {
		t.Errorf("Expected TryGet to return (2, true) for key 'two', got (%v, %v)", v, ok)
	}

	// Test TryGet with non-existent key
	if v, ok := m.TryGet("three"); ok || v != 0 {
		t.Errorf("Expected TryGet to return (0, false) for non-existent key, got (%v, %v)", v, ok)
	}

	// Test with zero value
	if v, ok := m.TryGet("zero"); !ok || v != 0 {
		t.Errorf("Expected TryGet to return (0, true) for zero value key, got (%v, %v)", v, ok)
	}
}

func TestMapGetOrSet(t *testing.T) {
	// Test GetOrSet on nil map
	var nilMap testtarget.Map[string, int]
	value := nilMap.GetOrSet("new", 42)

	// Verify the value was set and returned
	if value != 42 {
		t.Errorf("Expected GetOrSet to return 42 for new key in nil map, got %v", value)
	}
	if v, ok := nilMap["new"]; !ok || v != 42 {
		t.Errorf("Expected GetOrSet to set value 42 for new key, got (%v, %v)", v, ok)
	}

	// Create a map with data
	m := testtarget.Map[string, int]{
		"one": 1,
	}

	// Test GetOrSet with existing key (should return existing value)
	value = m.GetOrSet("one", 100)
	if value != 1 {
		t.Errorf("Expected GetOrSet to return existing value 1, got %v", value)
	}
	if v := m["one"]; v != 1 {
		t.Errorf("Expected value to remain 1 after GetOrSet on existing key, got %v", v)
	}

	// Test GetOrSet with new key (should set and return new value)
	value = m.GetOrSet("two", 2)
	if value != 2 {
		t.Errorf("Expected GetOrSet to return new value 2, got %v", value)
	}
	if v := m["two"]; v != 2 {
		t.Errorf("Expected GetOrSet to set value 2 for new key, got %v", v)
	}
}
