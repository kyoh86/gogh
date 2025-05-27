package typ_test

import (
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/typ"
)

func TestRemap(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		setupMap    map[string]int
		inputKey    string
		initialVal  int
		expectedVal int
		expectError bool
	}{
		{
			name: "remap with existing key",
			setupMap: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			},
			inputKey:    "two",
			initialVal:  0,
			expectedVal: 2,
			expectError: false,
		},
		{
			name: "remap with non-existent key",
			setupMap: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			},
			inputKey:    "four",
			initialVal:  0,
			expectedVal: 0, // Should remain unchanged
			expectError: true,
		},
		{
			name: "remap with empty key",
			setupMap: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			},
			inputKey:    "", // Zero value for string
			initialVal:  0,
			expectedVal: 0, // Should remain unchanged
			expectError: false,
		},
		{
			name: "remap with empty key that exists in map",
			setupMap: map[string]int{
				"":      42,
				"one":   1,
				"two":   2,
				"three": 3,
			},
			inputKey:    "", // Zero value for string but exists in map
			initialVal:  0,
			expectedVal: 0, // Should remain unchanged because zero value is skipped
			expectError: false,
		},
		{
			name: "overwrite existing value",
			setupMap: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			},
			inputKey:    "one",
			initialVal:  99,
			expectedVal: 1, // Should be overwritten with map value
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup value to be remapped
			value := tc.initialVal

			// Call Remap
			err := testtarget.Remap(&value, tc.setupMap, tc.inputKey)

			// Check error result
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}

			// Check value
			if value != tc.expectedVal {
				t.Errorf("Expected value to be %v, but got %v", tc.expectedVal, value)
			}
		})
	}
}

// TestRemapWithDifferentTypes tests Remap with different generic testtargete combinations
func TestRemapWithDifferentTypes(t *testing.T) {
	// Test with int keys
	t.Run("int keys", func(t *testing.T) {
		intMap := map[int]string{
			1: "one",
			2: "two",
			3: "three",
		}
		var value string
		err := testtarget.Remap(&value, intMap, 2)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if value != "two" {
			t.Errorf("Expected value to be 'two', but got '%s'", value)
		}

		// Test with zero value key (0)
		value = "should not change"
		err = testtarget.Remap(&value, intMap, 0)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if value != "should not change" {
			t.Errorf("Expected value to remain 'should not change', but got '%s'", value)
		}
	})

	// Test with struct keys
	t.Run("struct keys", func(t *testing.T) {
		type Key struct {
			ID   int
			Name string
		}

		structMap := map[Key]float64{
			{ID: 1, Name: "first"}:  1.1,
			{ID: 2, Name: "second"}: 2.2,
		}

		var value float64
		err := testtarget.Remap(&value, structMap, Key{ID: 2, Name: "second"})
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if value != 2.2 {
			t.Errorf("Expected value to be 2.2, but got %f", value)
		}

		// Test with non-existent key
		value = 0.0
		err = testtarget.Remap(&value, structMap, Key{ID: 3, Name: "third"})
		if err == nil {
			t.Errorf("Expected error but got nil")
		}
		if value != 0.0 {
			t.Errorf("Expected value to remain 0.0, but got %f", value)
		}
	})
}
