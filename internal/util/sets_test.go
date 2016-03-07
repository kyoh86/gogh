package util

import "testing"

func TestStringBoolMapKeys(t *testing.T) {
	{
		nilResult := StringBoolMapKeys(nil)
		if len(nilResult) != 0 {
			t.Errorf("Expected be 0, but %d", len(nilResult))
		}
	}
	{
		emptyResult := StringBoolMapKeys(map[string]bool{})
		if len(emptyResult) != 0 {
			t.Errorf("Expected be 0, but %d", len(emptyResult))
		}
	}
	{
		valueResult := StringBoolMapKeys(map[string]bool{
			"foo": true,
			"bar": false,
			"":    true,
		})
		if len(valueResult) != 3 {
			t.Errorf("Expected be 3, but %d", len(valueResult))
		}
	}
}
