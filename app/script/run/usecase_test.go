package run_test

import (
	"context"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/script/run"
)

func TestNewUsecase(t *testing.T) {
	uc := testtarget.NewUsecase()
	if uc == nil {
		t.Fatal("expected non-nil Usecase")
	}
}

func TestGlobals_ToLuaTable(t *testing.T) {
	t.Skip("Skipping Lua-related test to avoid CGO dependencies")
}

func TestScript(t *testing.T) {
	// Test that Script struct can be instantiated
	script := testtarget.Script{
		Code: "print('test')",
		Globals: testtarget.Globals{
			"test": "value",
			"num":  42,
		},
	}

	if script.Code != "print('test')" {
		t.Errorf("expected Code to be 'print('test')', got %q", script.Code)
	}

	if len(script.Globals) != 2 {
		t.Errorf("expected 2 globals, got %d", len(script.Globals))
	}
}

func TestGlobalsTypes(t *testing.T) {
	// Test different types in Globals
	globals := testtarget.Globals{
		"string": "test",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"map": map[string]any{
			"nested": "value",
		},
		"slice": []any{"a", "b", "c"},
	}

	// Verify types
	if v, ok := globals["string"].(string); !ok || v != "test" {
		t.Error("string type not preserved correctly")
	}

	if v, ok := globals["int"].(int); !ok || v != 42 {
		t.Error("int type not preserved correctly")
	}

	if v, ok := globals["float"].(float64); !ok || v != 3.14 {
		t.Error("float64 type not preserved correctly")
	}

	if v, ok := globals["bool"].(bool); !ok || v != true {
		t.Error("bool type not preserved correctly")
	}

	if v, ok := globals["map"].(map[string]any); !ok {
		t.Error("map type not preserved correctly")
	} else if nested, ok := v["nested"].(string); !ok || nested != "value" {
		t.Error("nested map value not preserved correctly")
	}

	if v, ok := globals["slice"].([]any); !ok || len(v) != 3 {
		t.Error("slice type not preserved correctly")
	}
}

func TestUsecase_Execute(t *testing.T) {
	// Skip actual Lua execution tests as they require CGO and Lua runtime
	t.Skip("Skipping Lua execution test to avoid CGO dependencies")

	// If we were to test this, we would:
	// 1. Create a Usecase
	// 2. Create a Script with simple Lua code
	// 3. Call Execute and verify no errors
	// 4. Test error cases like invalid Lua syntax

	ctx := context.Background()
	uc := testtarget.NewUsecase()

	// Example of what would be tested:
	script := testtarget.Script{
		Code: "print('Hello from Lua')",
		Globals: testtarget.Globals{
			"repo": map[string]any{
				"name": "test-repo",
			},
		},
	}

	// This would execute the Lua script
	_ = uc
	_ = ctx
	_ = script
}

func TestGlobalsMapOperations(t *testing.T) {
	// Test map operations on Globals
	globals := make(testtarget.Globals)

	// Test adding values
	globals["key1"] = "value1"
	globals["key2"] = 42

	// Test retrieving values
	if v, ok := globals["key1"]; !ok || v != "value1" {
		t.Error("failed to retrieve string value")
	}

	if v, ok := globals["key2"]; !ok || v != 42 {
		t.Error("failed to retrieve int value")
	}

	// Test updating values
	globals["key1"] = "updated"
	if v, ok := globals["key1"]; !ok || v != "updated" {
		t.Error("failed to update value")
	}

	// Test deleting values
	delete(globals, "key2")
	if _, ok := globals["key2"]; ok {
		t.Error("failed to delete value")
	}

	// Test length
	if len(globals) != 1 {
		t.Errorf("expected length 1, got %d", len(globals))
	}
}

func TestEmptyGlobals(t *testing.T) {
	// Test empty Globals
	globals := testtarget.Globals{}

	if len(globals) != 0 {
		t.Errorf("expected empty globals, got %d items", len(globals))
	}

	// Test nil Globals
	var nilGlobals testtarget.Globals
	if nilGlobals != nil {
		t.Error("expected nil globals to be nil")
	}
}

func TestGlobalsWithComplexStructures(t *testing.T) {
	// Test complex nested structures
	globals := testtarget.Globals{
		"repo": map[string]any{
			"name":  "gogh",
			"owner": "kyoh86",
			"host":  "github.com",
			"metadata": map[string]any{
				"stars":  100,
				"forks":  20,
				"topics": []string{"go", "cli", "github"},
			},
		},
		"hook": map[string]any{
			"id":   "hook-123",
			"name": "post-clone",
			"config": map[string]any{
				"enabled": true,
				"actions": []string{"overlay", "script"},
			},
		},
	}

	// Verify structure
	repo, ok := globals["repo"].(map[string]any)
	if !ok {
		t.Fatal("repo is not a map")
	}

	if name, ok := repo["name"].(string); !ok || name != "gogh" {
		t.Error("repo name not correct")
	}

	metadata, ok := repo["metadata"].(map[string]any)
	if !ok {
		t.Fatal("metadata is not a map")
	}

	if stars, ok := metadata["stars"].(int); !ok || stars != 100 {
		t.Error("stars count not correct")
	}

	topics, ok := metadata["topics"].([]string)
	if !ok || len(topics) != 3 {
		t.Error("topics not correct")
	}
}
