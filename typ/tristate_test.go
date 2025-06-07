package typ_test

import (
	"fmt"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/typ"
)

// TestTristateValues ensures the constants have the expected values
func TestTristateValues(t *testing.T) {
	tests := []struct {
		value    testtarget.Tristate
		expected int
	}{
		{testtarget.TristateZero, 0},
		{testtarget.TristateTrue, 1},
		{testtarget.TristateFalse, 2},
	}

	for _, test := range tests {
		if int(test.value) != test.expected {
			t.Errorf("Expected %s to have value %d, got %d", tristateToString(test.value), test.expected, int(test.value))
		}
	}
}

// Helper function to convert Tristate to string for error messages
func tristateToString(t testtarget.Tristate) string {
	switch t {
	case testtarget.TristateZero:
		return "TristateZero"
	case testtarget.TristateTrue:
		return "TristateTrue"
	case testtarget.TristateFalse:
		return "TristateFalse"
	default:
		return fmt.Sprintf("Unknown Tristate(%d)", t)
	}
}

// TestAsBoolPtr tests the AsBoolPtr method
func TestAsBoolPtr(t *testing.T) {
	tests := []struct {
		value       testtarget.Tristate
		expectedPtr *bool
		expectErr   bool
	}{
		{testtarget.TristateZero, nil, false},
		{testtarget.TristateTrue, testtarget.Ptr(true), false},
		{testtarget.TristateFalse, testtarget.Ptr(false), false},
		{testtarget.Tristate(99), nil, true}, // Invalid value should return error
	}

	for _, test := range tests {
		result, err := test.value.AsBoolPtr()

		// Check error expectation
		if test.expectErr && err == nil {
			t.Errorf("Expected error for %s, got nil", tristateToString(test.value))
			continue
		}
		if !test.expectErr && err != nil {
			t.Errorf("Unexpected error for %s: %v", tristateToString(test.value), err)
			continue
		}

		// If we expect error, no need to check the result
		if test.expectErr {
			continue
		}

		// Check result value
		if test.expectedPtr == nil && result != nil {
			t.Errorf("Expected nil for %s, got %v", tristateToString(test.value), *result)
			continue
		}
		if test.expectedPtr != nil && result == nil {
			t.Errorf("Expected %v for %s, got nil", *test.expectedPtr, tristateToString(test.value))
			continue
		}
		if test.expectedPtr != nil && result != nil && *test.expectedPtr != *result {
			t.Errorf("Expected %v for %s, got %v", *test.expectedPtr, tristateToString(test.value), *result)
		}
	}
}
