package alias_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/gogh/alias"
)

func set(items ...string) testtarget.Set {
	return testtarget.NewSet(items...)
}

func assertSet(t *testing.T, expect, actual testtarget.Set) {
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
	var s testtarget.Set

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
	assertSet(t, set("key1", "key2", "key3"), s)

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
	assertSet(t, set("key1", "key2"), s)
}
