package alias_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/gogh/alias"
)

func TestLookup(t *testing.T) {
	var l testtarget.Lookup

	l.Set("key1", "val1")
	assertSet(t, set("key1"), set(l.Keys()...))

	l.Set("key2", "val2")
	assertSet(t, set("key1", "key2"), set(l.Keys()...))

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
	assertSet(t, set("key1", "key2", "key3"), set(l.Keys()...))

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
	assertSet(t, set("key1", "key2"), set(l.Keys()...))
}

func TestReverse(t *testing.T) {
	var r testtarget.Reverse

	r.Set("key1", "val1-1")
	assertSet(t, set("val1-1"), r.Get("key1"))

	r.Set("key1", "val1-2")
	assertSet(t, set("val1-1", "val1-2"), r.Get("key1"))

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
	assertSet(t, set("val2-1"), r.Get("key2"))
}
