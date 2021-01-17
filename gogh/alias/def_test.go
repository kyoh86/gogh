package alias_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/gogh/alias"
)

func TestDef(t *testing.T) {
	var m testtarget.Def

	m.Set("alias1", "path1")
	m.Set("alias2", "path2")
	m.Set("alias3", "path3")
	m.Set("dup-alias1", "path1")
	m.Set("dup-alias2", "path2")
	m.Set("dup-alias3", "path3")
	if act := m.Lookup("alias1"); act != "path1" {
		t.Error("expect path1 is related for alias1, but not")
	}
	if act := m.Lookup("alias2"); act != "path2" {
		t.Error("expect path2 is related for alias2, but not")
	}
	if act := m.Lookup("alias3"); act != "path3" {
		t.Error("expect path3 is related for alias3, but not")
	}
	if act := m.Lookup("dup-alias1"); act != "path1" {
		t.Error("expect path1 is related for dup-alias1, but not")
	}
	if act := m.Lookup("dup-alias2"); act != "path2" {
		t.Error("expect path2 is related for dup-alias2, but not")
	}
	if act := m.Lookup("dup-alias3"); act != "path3" {
		t.Error("expect path3 is related for dup-alias3, but not")
	}
	assertSet(t, set("alias1", "dup-alias1"), m.Reverse("path1"))
	assertSet(t, set("alias2", "dup-alias2"), m.Reverse("path2"))
	assertSet(t, set("alias3", "dup-alias3"), m.Reverse("path3"))

	m.Set("dup-alias3", "path1") // replace
	if act := m.Lookup("dup-alias3"); act != "path1" {
		t.Error("expect path1 is related for dup-alias3, but not")
	}

	assertSet(t, set("alias1", "dup-alias1", "dup-alias3"), m.Reverse("path1"))
	assertSet(t, set("alias3"), m.Reverse("path3"))

	m.Set("alias3", "path1") // replace
	if act := m.Lookup("alias3"); act != "path1" {
		t.Error("expect path1 is related for alias3, but not")
	}

	assertSet(t, set("alias1", "alias3", "dup-alias1", "dup-alias3"), m.Reverse("path1"))
	assertSet(t, set(), m.Reverse("path3"))
}
