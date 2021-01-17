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
	var m testtarget.Set

	m.Set("key1")
	m.Set("key2")
	m.Set("key2") // dup
	m.Set("key3")
	if !m.Has("key1") {
		t.Error("expect has key1, but not")
	}
	if !m.Has("key2") {
		t.Error("expect has key2, but not")
	}
	if !m.Has("key3") {
		t.Error("expect has key3, but not")
	}
	assertSet(t, set("key1", "key2", "key3"), m)

	m.Del("key3")
	if !m.Has("key1") {
		t.Error("expect has key1, but not")
	}
	if !m.Has("key2") {
		t.Error("expect has key2, but not")
	}
	if m.Has("key3") {
		t.Error("expect does NOT have key3, but not")
	}
	assertSet(t, set("key1", "key2"), m)
}
