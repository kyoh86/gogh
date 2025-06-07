package typ

import (
	"errors"
	"maps"
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "filter even numbers",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(i int) bool { return i%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "filter positive numbers",
			input:     []int{-3, -2, -1, 0, 1, 2, 3},
			predicate: func(i int) bool { return i > 0 },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "filter all (true predicate)",
			input:     []int{1, 2, 3},
			predicate: func(i int) bool { return true },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "filter none (false predicate)",
			input:     []int{1, 2, 3},
			predicate: func(i int) bool { return false },
			expected:  []int{},
		},
		{
			name:      "filter empty input",
			input:     []int{},
			predicate: func(i int) bool { return true },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := slices.Values(tt.input)
			filtered := Filter(seq, tt.predicate)
			result := slices.Collect(filtered)

			if !slices.Equal(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFilter2(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]int
		predicate func(string, int) bool
		expected  map[string]int
	}{
		{
			name: "filter pairs where string length equals number",
			input: map[string]int{
				"a": 1, "bb": 2, "ccc": 3, "dddd": 4,
			},
			predicate: func(s string, i int) bool { return len(s) == i },
			expected: map[string]int{
				"a": 1, "bb": 2, "ccc": 3, "dddd": 4,
			},
		},
		{
			name: "filter pairs where string length greater than number",
			input: map[string]int{
				"a": 2, "bb": 2, "ccc": 2, "dddd": 2,
			},
			predicate: func(s string, i int) bool { return len(s) > i },
			expected: map[string]int{
				"ccc": 2, "dddd": 2,
			},
		},
		{
			name: "filter all (true predicate)",
			input: map[string]int{
				"a": 1, "b": 2, "c": 3,
			},
			predicate: func(s string, i int) bool { return true },
			expected: map[string]int{
				"a": 1, "b": 2, "c": 3,
			},
		},
		{
			name: "filter none (false predicate)",
			input: map[string]int{
				"a": 1, "b": 2, "c": 3,
			},
			predicate: func(s string, i int) bool { return false },
			expected:  map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := maps.All(tt.input)
			filtered := Filter2(seq, tt.predicate)
			result := maps.Collect(filtered)

			if !maps.Equal(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFilterE(t *testing.T) {
	errTest := errors.New("test error")

	tests := []struct {
		name      string
		input     map[int]error
		predicate func(int) (bool, error)
		expected  map[int]error
	}{
		{
			name: "filter even numbers",
			input: map[int]error{
				1: nil, 2: nil, 3: nil, 4: nil, 5: nil, 6: nil,
			},
			predicate: func(i int) (bool, error) { return i%2 == 0, nil },
			expected: map[int]error{
				2: nil, 4: nil, 6: nil,
			},
		},
		{
			name: "propagate input errors",
			input: map[int]error{
				1: nil, 2: errTest, 3: nil, 4: nil,
			},
			predicate: func(i int) (bool, error) { return true, nil },
			expected: map[int]error{
				1: nil, 2: errTest, 3: nil, 4: nil,
			},
		},
		{
			name: "return predicate errors",
			input: map[int]error{
				1: nil, 2: nil, 3: nil, 4: nil,
			},
			predicate: func(i int) (bool, error) {
				if i == 3 {
					return false, errTest
				}
				return true, nil
			},
			expected: map[int]error{
				1: nil, 2: nil, 3: errTest, 4: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := maps.All(tt.input)
			filtered := FilterE(seq, tt.predicate)
			result := maps.Collect(filtered)

			if !maps.Equal(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestWithNilError(t *testing.T) {
	input := []int{1, 2, 3}
	expected := map[int]error{
		1: nil, 2: nil, 3: nil,
	}

	seq := slices.Values(input)
	withErr := WithNilError(seq)
	result := maps.Collect(withErr)

	if !maps.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestWithError(t *testing.T) {
	input := []int{1, 2, 3}
	testErr := errors.New("test error")
	expected := map[int]error{
		1: testErr, 2: testErr, 3: testErr,
	}

	seq := slices.Values(input)
	withErr := WithError(seq, testErr)
	result := maps.Collect(withErr)

	if !maps.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCollectWithError(t *testing.T) {
	testErr := errors.New("test error")

	tests := []struct {
		name        string
		input       map[int]error
		expectedOut []int
		expectedErr error
	}{
		{
			name: "no errors",
			input: map[int]error{
				1: nil, 2: nil, 3: nil,
			},
			expectedOut: []int{1, 2, 3},
			expectedErr: nil,
		},
		{
			name: "with error",
			input: map[int]error{
				1: nil, 2: testErr, 3: nil, 4: nil,
			},
			expectedOut: nil,
			expectedErr: testErr,
		},
		{
			name:        "empty input",
			input:       map[int]error{},
			expectedOut: []int{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := maps.All(tt.input)
			result, err := CollectWithError(seq)

			// Sort result for stable comparison
			slices.Sort(result)

			if (tt.expectedErr == nil && err != nil) ||
				(tt.expectedErr != nil && err == nil) ||
				(tt.expectedErr != nil && err != nil && tt.expectedErr.Error() != err.Error()) {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}

			if !slices.Equal(result, tt.expectedOut) {
				t.Errorf("Expected %v, got %v", tt.expectedOut, result)
			}
		})
	}
}
