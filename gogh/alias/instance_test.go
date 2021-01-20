package alias

import (
	"testing"
)

func TestInstance(t *testing.T) {
	Set("alias1", "path1")
	Set("alias2", "path2")
	Set("alias3", "path3")
	Set("dup-alias1", "path1")
	Set("dup-alias2", "path2")
	Set("dup-alias3", "path3")

	for i := 0; i < 2; i++ { // loop to test marshal / unmarshal
		if act := Lookup("alias1"); act != "path1" { // nolint
			t.Error("expect path1 is related for alias1, but not")
		}
		if act := Lookup("alias2"); act != "path2" {
			t.Error("expect path2 is related for alias2, but not")
		}
		if act := Lookup("alias3"); act != "path3" {
			t.Error("expect path3 is related for alias3, but not")
		}
		if act := Lookup("dup-alias1"); act != "path1" {
			t.Error("expect path1 is related for dup-alias1, but not")
		}
		if act := Lookup("dup-alias2"); act != "path2" {
			t.Error("expect path2 is related for dup-alias2, but not")
		}
		if act := Lookup("dup-alias3"); act != "path3" {
			t.Error("expect path3 is related for dup-alias3, but not")
		}
		assertSet(t, newSet("alias1", "dup-alias1"), newSet(Reverse("path1")...))
		assertSet(t, newSet("alias2", "dup-alias2"), newSet(Reverse("path2")...))
		assertSet(t, newSet("alias3", "dup-alias3"), newSet(Reverse("path3")...))

		{
			if err := SaveInstance("./testdata/testtemporary.yaml"); err != nil {
				t.Fatalf("failed to marshal Def object: %s", err.Error())
			}
			if err := LoadInstance("./testdata/testtemporary.yaml"); err != nil {
				t.Fatalf("failed to marshal Def object: %s", err.Error())
			}
		}
	}

	Set("dup-alias3", "path1") // replace
	if act := Lookup("dup-alias3"); act != "path1" {
		t.Error("expect path1 is related for dup-alias3, but not")
	}

	assertSet(t, newSet("alias1", "dup-alias1", "dup-alias3"), newSet(Reverse("path1")...))
	assertSet(t, newSet("alias3"), newSet(Reverse("path3")...))

	Set("alias3", "path1") // replace
	if act := Lookup("alias3"); act != "path1" {
		t.Error("expect path1 is related for alias3, but not")
	}

	assertSet(t, newSet("alias1", "alias3", "dup-alias1", "dup-alias3"), newSet(Reverse("path1")...))
	assertSet(t, newSet(), newSet(Reverse("path3")...))

	Del("alias1")
	if act := Lookup("alias1"); act != "" {
		t.Errorf("expect path1 is not related for any path, but %q found", act)
	}
	assertSet(t, newSet("alias3", "dup-alias1", "dup-alias3"), newSet(Reverse("path1")...))

	Del("alias1") //  Del is idempotent
	if act := Lookup("alias1"); act != "" {
		t.Errorf("expect path1 is not related for any path, but %q found", act)
	}
	assertSet(t, newSet("alias3", "dup-alias1", "dup-alias3"), newSet(Reverse("path1")...))

	if err := SaveInstance("./testdata/dummyfile/test.yaml"); err == nil {
		t.Error("expect error when instance is saved to the path under the file")
	}
	if err := LoadInstance("./testdata/not-exist.yaml"); err != nil {
		t.Errorf("failed to load instance from not-exist file: %s", err.Error())
	}
}
