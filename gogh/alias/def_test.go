package alias

import (
	"testing"
)

func TestDef(t *testing.T) {
	var d Def

	d.Set("alias1", "path1")
	d.Set("alias2", "path2")
	d.Set("alias3", "path3")
	d.Set("dup-alias1", "path1")
	d.Set("dup-alias2", "path2")
	d.Set("dup-alias3", "path3")
	if act := d.Lookup("alias1"); act != "path1" { // nolint
		t.Error("expect path1 is related for alias1, but not")
	}
	if act := d.Lookup("alias2"); act != "path2" {
		t.Error("expect path2 is related for alias2, but not")
	}
	if act := d.Lookup("alias3"); act != "path3" {
		t.Error("expect path3 is related for alias3, but not")
	}
	if act := d.Lookup("dup-alias1"); act != "path1" {
		t.Error("expect path1 is related for dup-alias1, but not")
	}
	if act := d.Lookup("dup-alias2"); act != "path2" {
		t.Error("expect path2 is related for dup-alias2, but not")
	}
	if act := d.Lookup("dup-alias3"); act != "path3" {
		t.Error("expect path3 is related for dup-alias3, but not")
	}
	assertSet(t, newSet("alias1", "dup-alias1"), newSet(d.Reverse("path1")...))
	assertSet(t, newSet("alias2", "dup-alias2"), newSet(d.Reverse("path2")...))
	assertSet(t, newSet("alias3", "dup-alias3"), newSet(d.Reverse("path3")...))

	d.Set("dup-alias3", "path1") // replace
	if act := d.Lookup("dup-alias3"); act != "path1" {
		t.Error("expect path1 is related for dup-alias3, but not")
	}

	assertSet(t, newSet("alias1", "dup-alias1", "dup-alias3"), newSet(d.Reverse("path1")...))
	assertSet(t, newSet("alias3"), newSet(d.Reverse("path3")...))

	d.Set("alias3", "path1") // replace
	if act := d.Lookup("alias3"); act != "path1" {
		t.Error("expect path1 is related for alias3, but not")
	}

	assertSet(t, newSet("alias1", "alias3", "dup-alias1", "dup-alias3"), newSet(d.Reverse("path1")...))
	assertSet(t, newSet(), newSet(d.Reverse("path3")...))
}
