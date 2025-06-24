package extra_test

import (
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestNewAutoExtra(t *testing.T) {
	id := "auto-extra-123"
	repo := repository.NewReference("github.com", "owner", "repo")
	source := repository.NewReference("github.com", "source", "template")
	items := []extra.Item{
		{OverlayID: "overlay-1", HookID: "hook-1"},
		{OverlayID: "overlay-2", HookID: "hook-2"},
	}
	createdAt := time.Now()

	e := extra.NewAutoExtra(id, repo, source, items, createdAt)

	// Test ID
	if e.ID() != id {
		t.Errorf("ID() = %v, want %v", e.ID(), id)
	}

	// Test Type
	if e.Type() != extra.TypeAuto {
		t.Errorf("Type() = %v, want %v", e.Type(), extra.TypeAuto)
	}

	// Test Name (should be empty for auto extra)
	if e.Name() != "" {
		t.Errorf("Name() = %v, want empty string", e.Name())
	}

	// Test Repository
	gotRepo := e.Repository()
	if gotRepo == nil {
		t.Fatal("Repository() should not be nil for auto extra")
	}
	if gotRepo.String() != repo.String() {
		t.Errorf("Repository() = %v, want %v", gotRepo.String(), repo.String())
	}

	// Test Items
	gotItems := e.Items()
	if len(gotItems) != len(items) {
		t.Errorf("Items() length = %v, want %v", len(gotItems), len(items))
	}
	for i, item := range gotItems {
		if item.OverlayID != items[i].OverlayID || item.HookID != items[i].HookID {
			t.Errorf("Items()[%d] = %v, want %v", i, item, items[i])
		}
	}

	// Test Source
	if e.Source().String() != source.String() {
		t.Errorf("Source() = %v, want %v", e.Source().String(), source.String())
	}

	// Test CreatedAt
	if !e.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt() = %v, want %v", e.CreatedAt(), createdAt)
	}
}

func TestNewNamedExtra(t *testing.T) {
	id := "named-extra-456"
	name := "my-template"
	source := repository.NewReference("github.com", "source", "template")
	items := []extra.Item{
		{OverlayID: "overlay-3", HookID: "hook-3"},
	}
	createdAt := time.Now()

	e := extra.NewNamedExtra(id, name, source, items, createdAt)

	// Test ID
	if e.ID() != id {
		t.Errorf("ID() = %v, want %v", e.ID(), id)
	}

	// Test Type
	if e.Type() != extra.TypeNamed {
		t.Errorf("Type() = %v, want %v", e.Type(), extra.TypeNamed)
	}

	// Test Name
	if e.Name() != name {
		t.Errorf("Name() = %v, want %v", e.Name(), name)
	}

	// Test Repository (should be nil for named extra)
	if e.Repository() != nil {
		t.Errorf("Repository() = %v, want nil", e.Repository())
	}

	// Test Items
	gotItems := e.Items()
	if len(gotItems) != len(items) {
		t.Errorf("Items() length = %v, want %v", len(gotItems), len(items))
	}

	// Test Source
	if e.Source().String() != source.String() {
		t.Errorf("Source() = %v, want %v", e.Source().String(), source.String())
	}

	// Test CreatedAt
	if !e.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt() = %v, want %v", e.CreatedAt(), createdAt)
	}
}

func TestExtraItemsCopy(t *testing.T) {
	// Test that Items() returns a copy, not the original slice
	items := []extra.Item{
		{OverlayID: "overlay-1", HookID: "hook-1"},
		{OverlayID: "overlay-2", HookID: "hook-2"},
	}

	e := extra.NewNamedExtra(
		"test-id",
		"test-name",
		repository.NewReference("github.com", "owner", "repo"),
		items,
		time.Now(),
	)

	gotItems := e.Items()

	// Modify the returned slice
	gotItems[0].OverlayID = "modified"

	// Original should not be affected
	gotItemsAgain := e.Items()
	if gotItemsAgain[0].OverlayID == "modified" {
		t.Error("Items() should return a copy, not the original slice")
	}
}

func TestExtraRepositoryCopy(t *testing.T) {
	// Test that Repository() returns a copy, not the original reference
	repo := repository.NewReference("github.com", "owner", "repo")

	e := extra.NewAutoExtra(
		"test-id",
		repo,
		repository.NewReference("github.com", "source", "template"),
		[]extra.Item{{OverlayID: "o1", HookID: "h1"}},
		time.Now(),
	)

	gotRepo1 := e.Repository()
	gotRepo2 := e.Repository()

	// Should be equal but not the same pointer
	if gotRepo1.String() != gotRepo2.String() {
		t.Error("Repository() should return equal references")
	}
	if gotRepo1 == gotRepo2 {
		t.Error("Repository() should return a copy, not the same pointer")
	}
}

func TestExtraWithEmptyItems(t *testing.T) {
	// Test that extras can have empty items
	e := extra.NewNamedExtra(
		"empty-items",
		"template",
		repository.NewReference("github.com", "owner", "repo"),
		[]extra.Item{},
		time.Now(),
	)

	items := e.Items()
	if len(items) != 0 {
		t.Errorf("Items() length = %v, want 0", len(items))
	}
}

func TestExtraTypeConstants(t *testing.T) {
	// Test that type constants have expected values
	if extra.TypeAuto != "auto" {
		t.Errorf("TypeAuto = %v, want 'auto'", extra.TypeAuto)
	}
	if extra.TypeNamed != "named" {
		t.Errorf("TypeNamed = %v, want 'named'", extra.TypeNamed)
	}
}
