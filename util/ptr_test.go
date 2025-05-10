package util_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v3/util"
)

func TestPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		val := 42
		ptr := testtarget.Ptr(val)
		if *ptr != val {
			t.Errorf("ptr() = %v, want %v", *ptr, val)
		}
	})

	t.Run("string", func(t *testing.T) {
		val := "hello"
		ptr := testtarget.Ptr(val)
		if *ptr != val {
			t.Errorf("ptr() = %v, want %v", *ptr, val)
		}
	})
}

func TestNilablePtr(t *testing.T) {
	t.Run("int non-zero", func(t *testing.T) {
		val := 42
		ptr := testtarget.NilablePtr(val)
		if ptr == nil || *ptr != val {
			t.Errorf("nilablePtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("int zero", func(t *testing.T) {
		val := 0
		ptr := testtarget.NilablePtr(val)
		if ptr != nil {
			t.Errorf("nilablePtr() = %v, want nil", ptr)
		}
	})

	t.Run("string non-empty", func(t *testing.T) {
		val := "hello"
		ptr := testtarget.NilablePtr(val)
		if ptr == nil || *ptr != val {
			t.Errorf("nilablePtr() = %v, want %v", ptr, val)
		}
	})

	t.Run("string empty", func(t *testing.T) {
		val := ""
		ptr := testtarget.NilablePtr(val)
		if ptr != nil {
			t.Errorf("nilablePtr() = %v, want nil", ptr)
		}
	})
}
