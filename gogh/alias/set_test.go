package alias

import (
	"testing"
)

func assertSet(t *testing.T, expect, actual set) {
	t.Helper()
	for key := range expect {
		if !actual.Has(key) {
			t.Errorf("expect has %q but not", key)
		}
	}
	for key := range actual {
		if !expect.Has(key) {
			t.Errorf("unexpected %q is exist", key)
		}
	}
}

func TestSet(t *testing.T) {
	var s set

	t.Run("EmptySetNotPanic", func(t *testing.T) {
		t.Run("Del", func(t *testing.T) {
			t.Parallel()
			var empty set
			empty.Del("key0")
		})
		t.Run("List", func(t *testing.T) {
			t.Parallel()
			var empty set
			list := empty.List()
			if len(list) > 0 {
				t.Errorf("empty set returns not empty list: %q", list)
			}
		})
		t.Run("Has", func(t *testing.T) {
			t.Parallel()
			var empty set
			if empty.Has("key0") {
				t.Errorf("empty set has key0")
			}
		})
		t.Run("Set", func(t *testing.T) {
			t.Parallel()
			var empty set
			empty.Set("key0")
		})
	})

	s.Set("key1")
	s.Set("key2")
	s.Set("key2") // dup
	s.Set("key3")
	if !s.Has("key1") {
		t.Error("expect has key1, but not")
	}
	if !s.Has("key2") {
		t.Error("expect has key2, but not")
	}
	if !s.Has("key3") {
		t.Error("expect has key3, but not")
	}
	assertSet(t, newSet("key1", "key2", "key3"), s)

	s.Del("key3")
	if !s.Has("key1") {
		t.Error("expect has key1, but not")
	}
	if !s.Has("key2") {
		t.Error("expect has key2, but not")
	}
	if s.Has("key3") {
		t.Error("expect does NOT have key3, but not")
	}
	assertSet(t, newSet("key1", "key2"), s)
}
