package alias

import (
	"testing"
)

func TestLookup(t *testing.T) {
	var l lookup

	t.Run("EmptyLookupNotPanic", func(t *testing.T) {
		t.Run("Del", func(t *testing.T) {
			t.Parallel()
			var empty lookup
			empty.Del("key0")
		})
		t.Run("Get", func(t *testing.T) {
			t.Parallel()
			var empty lookup
			list := empty.Get("key0")
			if len(list) > 0 {
				t.Errorf("empty lookup returns not empty: %q", list)
			}
		})
		t.Run("Has", func(t *testing.T) {
			t.Parallel()
			var empty lookup
			if empty.Has("key0") {
				t.Errorf("empty lookup has key0:val0")
			}
		})
		t.Run("Set", func(t *testing.T) {
			t.Parallel()
			var empty lookup
			empty.Set("key0", "val0")
		})
	})

	l.Set("key1", "val1")

	l.Set("key2", "val2")

	l.Set("key3", "val3")
	if !l.Has("key1") {
		t.Error("expect has key1, but not")
	}
	if act := l.Get("key1"); act != "val1" {
		t.Error("expect val1 is related for key1, but not")
	}
	if !l.Has("key2") {
		t.Error("expect has key2, but not")
	}
	if act := l.Get("key2"); act != "val2" {
		t.Error("expect val2 is related for key2, but not")
	}
	if !l.Has("key3") {
		t.Error("expect has key3, but not")
	}
	if act := l.Get("key3"); act != "val3" {
		t.Error("expect val3 is related for key3, but not")
	}

	l.Del("key3")
	if !l.Has("key1") {
		t.Error("expect has key1, but not")
	}
	if act := l.Get("key1"); act != "val1" {
		t.Error("expect val1 is related for key1, but not")
	}
	if !l.Has("key2") {
		t.Error("expect has key2, but not")
	}
	if act := l.Get("key2"); act != "val2" {
		t.Error("expect val2 is related for key2, but not")
	}
	if l.Has("key3") {
		t.Error("expect does NOT have key3, but not")
	}
	if act := l.Get("key3"); act != "" {
		t.Error("expect empty value is related for key3, but not")
	}
}

func TestReverse(t *testing.T) {
	var r reverse

	t.Run("EmptyReverseNotPanic", func(t *testing.T) {
		t.Run("Del", func(t *testing.T) {
			t.Parallel()
			var empty reverse
			empty.Del("key0", "val0")
		})
		t.Run("Get", func(t *testing.T) {
			t.Parallel()
			var empty reverse
			list := empty.Get("key0")
			if len(list) > 0 {
				t.Errorf("empty reverse returns not empty: %q", list)
			}
		})
		t.Run("Has", func(t *testing.T) {
			t.Parallel()
			var empty reverse
			if empty.Has("key0", "val0") {
				t.Errorf("empty reverse has key0:val0")
			}
		})
		t.Run("Set", func(t *testing.T) {
			t.Parallel()
			var empty reverse
			empty.Set("key0", "val0")
		})
	})

	r.Set("key1", "val1-1")
	assertSet(t, newSet("val1-1"), r.Get("key1"))

	r.Set("key1", "val1-2")
	assertSet(t, newSet("val1-1", "val1-2"), r.Get("key1"))

	r.Set("key2", "val2-1")

	r.Set("key2", "val2-2")
	if !r.Has("key1", "val1-1") {
		t.Error("expect has val1-1, but not")
	}
	if !r.Get("key1").Has("val1-1") {
		t.Error("expect val1-1 is related for val1-1, but not")
	}
	if !r.Has("key1", "val1-2") {
		t.Error("expect has val1-2, but not")
	}
	if !r.Get("key1").Has("val1-2") {
		t.Error("expect val1-2 is related for val1-2, but not")
	}
	if !r.Has("key2", "val2-1") {
		t.Error("expect has val2-1, but not")
	}
	if !r.Get("key2").Has("val2-1") {
		t.Error("expect val2-1 is related for val2-1, but not")
	}
	if !r.Has("key2", "val2-2") {
		t.Error("expect has val2-2, but not")
	}
	if !r.Get("key2").Has("val2-2") {
		t.Error("expect val2-2 is related for val2-2, but not")
	}

	r.Del("key2", "val2-2")
	if !r.Has("key1", "val1-1") {
		t.Error("expect has val1-1, but not")
	}
	if !r.Get("key1").Has("val1-1") {
		t.Error("expect val1-1 is related for val1-1, but not")
	}
	if !r.Has("key1", "val1-2") {
		t.Error("expect has val1-2, but not")
	}
	if !r.Get("key1").Has("val1-2") {
		t.Error("expect val1-2 is related for val1-2, but not")
	}
	if !r.Has("key2", "val2-1") {
		t.Error("expect has val2-1, but not")
	}
	if !r.Get("key2").Has("val2-1") {
		t.Error("expect val2-1 is related for val2-1, but not")
	}
	if r.Has("key2", "val2-2") {
		t.Error("expect does NOT have key2, but not")
	}
	if r.Get("key2").Has("val2-2") {
		t.Error("expect val2-2 is not related for key2, but not")
	}
	assertSet(t, newSet("val2-1"), r.Get("key2"))

	r.Del("key2", "val2-1")
	if !r.Has("key1", "val1-1") {
		t.Error("expect has val1-1, but not")
	}
	if !r.Get("key1").Has("val1-1") {
		t.Error("expect val1-1 is related for val1-1, but not")
	}
	if !r.Has("key1", "val1-2") {
		t.Error("expect has val1-2, but not")
	}
	if !r.Get("key1").Has("val1-2") {
		t.Error("expect val1-2 is related for val1-2, but not")
	}
	if r.Has("key2", "val2-1") {
		t.Error("expect has val2-1, but not")
	}
	if r.Get("key2").Has("val2-1") {
		t.Error("expect val2-2 is not related for key2, but not")
	}
	assertSet(t, newSet(), r.Get("key2"))
}
