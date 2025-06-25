package extra_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestNewExtraService(t *testing.T) {
	service := extra.NewExtraService()
	if service == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestExtraService_AddAutoExtra(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	repo := repository.NewReference("github.com", "kyoh86", "gogh")
	source := repository.NewReference("github.com", "kyoh86", "gogh-extras")
	items := []extra.Item{
		{OverlayID: "overlay-1", HookID: "hook-1"},
		{OverlayID: "overlay-2", HookID: "hook-2"},
	}

	// Test successful addition
	id, err := service.AddAutoExtra(ctx, repo, source, items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Test duplicate addition
	_, err = service.AddAutoExtra(ctx, repo, source, items)
	if err != extra.ErrExtraAlreadyExists {
		t.Errorf("expected ErrExtraAlreadyExists, got %v", err)
	}

	// Verify the extra was added
	e, err := service.GetAutoExtra(ctx, repo)
	if err != nil {
		t.Fatalf("failed to get auto extra: %v", err)
	}
	if e.ID() != id {
		t.Errorf("expected ID %s, got %s", id, e.ID())
	}
	if e.Type() != extra.TypeAuto {
		t.Errorf("expected type %s, got %s", extra.TypeAuto, e.Type())
	}
	if e.Source().String() != source.String() {
		t.Errorf("expected source %s, got %s", source.String(), e.Source().String())
	}
	if len(e.Items()) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(e.Items()))
	}
}

func TestExtraService_AddNamedExtra(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	name := "my-template"
	source := repository.NewReference("github.com", "kyoh86", "gogh-templates")
	items := []extra.Item{
		{OverlayID: "overlay-3", HookID: "hook-3"},
	}

	// Test successful addition
	id, err := service.AddNamedExtra(ctx, name, source, items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty ID")
	}

	// Test duplicate addition
	_, err = service.AddNamedExtra(ctx, name, source, items)
	if err != extra.ErrExtraAlreadyExists {
		t.Errorf("expected ErrExtraAlreadyExists, got %v", err)
	}

	// Verify the extra was added
	e, err := service.GetNamedExtra(ctx, name)
	if err != nil {
		t.Fatalf("failed to get named extra: %v", err)
	}
	if e.ID() != id {
		t.Errorf("expected ID %s, got %s", id, e.ID())
	}
	if e.Type() != extra.TypeNamed {
		t.Errorf("expected type %s, got %s", extra.TypeNamed, e.Type())
	}
	if e.Name() != name {
		t.Errorf("expected name %s, got %s", name, e.Name())
	}
	if e.Repository() != nil {
		t.Error("expected nil repository for named extra")
	}
}

func TestExtraService_Get(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	// Add an auto extra
	repo := repository.NewReference("github.com", "kyoh86", "test")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	id, err := service.AddAutoExtra(ctx, repo, source, items)
	if err != nil {
		t.Fatalf("failed to add auto extra: %v", err)
	}

	// Test getting by ID
	e, err := service.Get(ctx, id)
	if err != nil {
		t.Fatalf("failed to get extra by ID: %v", err)
	}
	if e.ID() != id {
		t.Errorf("expected ID %s, got %s", id, e.ID())
	}

	// Test getting non-existent ID
	_, err = service.Get(ctx, "non-existent-id")
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound, got %v", err)
	}
}

func TestExtraService_RemoveAutoExtra(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	repo := repository.NewReference("github.com", "kyoh86", "test")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	id, err := service.AddAutoExtra(ctx, repo, source, items)
	if err != nil {
		t.Fatalf("failed to add auto extra: %v", err)
	}

	// Test removing existing extra
	err = service.RemoveAutoExtra(ctx, repo)
	if err != nil {
		t.Fatalf("failed to remove auto extra: %v", err)
	}

	// Verify it was removed
	_, err = service.GetAutoExtra(ctx, repo)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound after removal, got %v", err)
	}

	// Verify it was also removed from byID map
	_, err = service.Get(ctx, id)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound from Get after removal, got %v", err)
	}

	// Test removing non-existent extra
	err = service.RemoveAutoExtra(ctx, repo)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound, got %v", err)
	}
}

func TestExtraService_RemoveNamedExtra(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	name := "template"
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	id, err := service.AddNamedExtra(ctx, name, source, items)
	if err != nil {
		t.Fatalf("failed to add named extra: %v", err)
	}

	// Test removing existing extra
	err = service.RemoveNamedExtra(ctx, name)
	if err != nil {
		t.Fatalf("failed to remove named extra: %v", err)
	}

	// Verify it was removed
	_, err = service.GetNamedExtra(ctx, name)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound after removal, got %v", err)
	}

	// Verify it was also removed from byID map
	_, err = service.Get(ctx, id)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound from Get after removal, got %v", err)
	}

	// Test removing non-existent extra
	err = service.RemoveNamedExtra(ctx, name)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound, got %v", err)
	}
}

func TestExtraService_Remove(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	// Add both auto and named extras
	repo := repository.NewReference("github.com", "kyoh86", "test")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	autoID, err := service.AddAutoExtra(ctx, repo, source, items)
	if err != nil {
		t.Fatalf("failed to add auto extra: %v", err)
	}

	namedID, err := service.AddNamedExtra(ctx, "template", source, items)
	if err != nil {
		t.Fatalf("failed to add named extra: %v", err)
	}

	// Test removing auto extra by ID
	err = service.Remove(ctx, autoID)
	if err != nil {
		t.Fatalf("failed to remove auto extra by ID: %v", err)
	}

	// Verify auto extra was removed from all maps
	_, err = service.Get(ctx, autoID)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound from Get, got %v", err)
	}
	_, err = service.GetAutoExtra(ctx, repo)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound from GetAutoExtra, got %v", err)
	}

	// Test removing named extra by ID
	err = service.Remove(ctx, namedID)
	if err != nil {
		t.Fatalf("failed to remove named extra by ID: %v", err)
	}

	// Verify named extra was removed from all maps
	_, err = service.Get(ctx, namedID)
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound from Get, got %v", err)
	}
	_, err = service.GetNamedExtra(ctx, "template")
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound from GetNamedExtra, got %v", err)
	}

	// Test removing non-existent ID
	err = service.Remove(ctx, "non-existent")
	if err != extra.ErrExtraNotFound {
		t.Errorf("expected ErrExtraNotFound, got %v", err)
	}
}

func TestExtraService_List(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	// Add multiple extras
	repo1 := repository.NewReference("github.com", "kyoh86", "test1")
	repo2 := repository.NewReference("github.com", "kyoh86", "test2")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	id1, _ := service.AddAutoExtra(ctx, repo1, source, items)
	id2, _ := service.AddAutoExtra(ctx, repo2, source, items)
	id3, _ := service.AddNamedExtra(ctx, "template1", source, items)
	id4, _ := service.AddNamedExtra(ctx, "template2", source, items)

	// List all extras
	var count int
	ids := make(map[string]bool)
	for e, err := range service.List(ctx) {
		if err != nil {
			t.Fatalf("unexpected error in List: %v", err)
		}
		ids[e.ID()] = true
		count++
	}

	if count != 4 {
		t.Errorf("expected 4 extras, got %d", count)
	}

	// Verify all IDs are present
	for _, id := range []string{id1, id2, id3, id4} {
		if !ids[id] {
			t.Errorf("expected ID %s in list", id)
		}
	}
}

func TestExtraService_ListByType(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	// Add multiple extras
	repo1 := repository.NewReference("github.com", "kyoh86", "test1")
	repo2 := repository.NewReference("github.com", "kyoh86", "test2")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	_, _ = service.AddAutoExtra(ctx, repo1, source, items)
	_, _ = service.AddAutoExtra(ctx, repo2, source, items)
	_, _ = service.AddNamedExtra(ctx, "template1", source, items)
	_, _ = service.AddNamedExtra(ctx, "template2", source, items)
	_, _ = service.AddNamedExtra(ctx, "template3", source, items)

	// List auto extras
	var autoCount int
	for e, err := range service.ListByType(ctx, extra.TypeAuto) {
		if err != nil {
			t.Fatalf("unexpected error in ListByType(auto): %v", err)
		}
		if e.Type() != extra.TypeAuto {
			t.Errorf("expected type %s, got %s", extra.TypeAuto, e.Type())
		}
		autoCount++
	}
	if autoCount != 2 {
		t.Errorf("expected 2 auto extras, got %d", autoCount)
	}

	// List named extras
	var namedCount int
	for e, err := range service.ListByType(ctx, extra.TypeNamed) {
		if err != nil {
			t.Fatalf("unexpected error in ListByType(named): %v", err)
		}
		if e.Type() != extra.TypeNamed {
			t.Errorf("expected type %s, got %s", extra.TypeNamed, e.Type())
		}
		namedCount++
	}
	if namedCount != 3 {
		t.Errorf("expected 3 named extras, got %d", namedCount)
	}

	// Test invalid type
	var errorCount int
	for _, err := range service.ListByType(ctx, extra.Type("invalid")) {
		if err == nil {
			t.Error("expected error for invalid type")
		}
		errorCount++
		break // Should only yield one error
	}
	if errorCount != 1 {
		t.Errorf("expected 1 error, got %d", errorCount)
	}
}

func TestExtraService_Load(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	// Create extras to load
	repo := repository.NewReference("github.com", "kyoh86", "test")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	autoExtra := extra.NewAutoExtra("auto-1", repo, source, items, time.Now())
	namedExtra := extra.NewNamedExtra("named-1", "template", source, items, time.Now())

	// Load extras
	extras := func(yield func(*extra.Extra, error) bool) {
		yield(autoExtra, nil)
		yield(namedExtra, nil)
	}

	err := service.Load(extras)
	if err != nil {
		t.Fatalf("failed to load extras: %v", err)
	}

	// Verify auto extra was loaded
	e, err := service.GetAutoExtra(ctx, repo)
	if err != nil {
		t.Fatalf("failed to get auto extra after load: %v", err)
	}
	if e.ID() != "auto-1" {
		t.Errorf("expected ID auto-1, got %s", e.ID())
	}

	// Verify named extra was loaded
	e, err = service.GetNamedExtra(ctx, "template")
	if err != nil {
		t.Fatalf("failed to get named extra after load: %v", err)
	}
	if e.ID() != "named-1" {
		t.Errorf("expected ID named-1, got %s", e.ID())
	}

	// Test loading with error
	extrasWithError := func(yield func(*extra.Extra, error) bool) {
		yield(nil, errors.New("load error"))
	}

	err = service.Load(extrasWithError)
	if err == nil {
		t.Error("expected error when loading with error")
	}
}

func TestExtraService_HasChanges(t *testing.T) {
	ctx := context.Background()
	service := extra.NewExtraService()

	// Initially no changes
	if service.HasChanges() {
		t.Error("expected no changes initially")
	}

	// Add extra
	repo := repository.NewReference("github.com", "kyoh86", "test")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{{OverlayID: "o1", HookID: "h1"}}

	_, err := service.AddAutoExtra(ctx, repo, source, items)
	if err != nil {
		t.Fatalf("failed to add extra: %v", err)
	}

	// Should have changes after add
	if !service.HasChanges() {
		t.Error("expected changes after add")
	}

	// Mark saved
	service.MarkSaved()
	if service.HasChanges() {
		t.Error("expected no changes after MarkSaved")
	}

	// Remove extra
	err = service.RemoveAutoExtra(ctx, repo)
	if err != nil {
		t.Fatalf("failed to remove extra: %v", err)
	}

	// Should have changes after remove
	if !service.HasChanges() {
		t.Error("expected changes after remove")
	}
}

func TestExtra_Methods(t *testing.T) {
	repo := repository.NewReference("github.com", "kyoh86", "test")
	source := repository.NewReference("github.com", "kyoh86", "source")
	items := []extra.Item{
		{OverlayID: "o1", HookID: "h1"},
		{OverlayID: "o2", HookID: "h2"},
	}
	createdAt := time.Now()

	// Test auto extra
	autoExtra := extra.NewAutoExtra("auto-id", repo, source, items, createdAt)

	if autoExtra.ID() != "auto-id" {
		t.Errorf("expected ID auto-id, got %s", autoExtra.ID())
	}
	if autoExtra.Type() != extra.TypeAuto {
		t.Errorf("expected type %s, got %s", extra.TypeAuto, autoExtra.Type())
	}
	if autoExtra.Name() != "" {
		t.Errorf("expected empty name for auto extra, got %s", autoExtra.Name())
	}
	if autoExtra.Repository() == nil || autoExtra.Repository().String() != repo.String() {
		t.Error("unexpected repository value for auto extra")
	}
	if len(autoExtra.Items()) != len(items) {
		t.Errorf("expected %d items, got %d", len(items), len(autoExtra.Items()))
	}
	if autoExtra.Source().String() != source.String() {
		t.Errorf("expected source %s, got %s", source.String(), autoExtra.Source().String())
	}
	if !autoExtra.CreatedAt().Equal(createdAt) {
		t.Errorf("expected createdAt %v, got %v", createdAt, autoExtra.CreatedAt())
	}

	// Test named extra
	namedExtra := extra.NewNamedExtra("named-id", "template", source, items, createdAt)

	if namedExtra.ID() != "named-id" {
		t.Errorf("expected ID named-id, got %s", namedExtra.ID())
	}
	if namedExtra.Type() != extra.TypeNamed {
		t.Errorf("expected type %s, got %s", extra.TypeNamed, namedExtra.Type())
	}
	if namedExtra.Name() != "template" {
		t.Errorf("expected name template, got %s", namedExtra.Name())
	}
	if namedExtra.Repository() != nil {
		t.Error("expected nil repository for named extra")
	}

	// Test that Items() returns a copy
	itemsCopy := namedExtra.Items()
	itemsCopy[0].OverlayID = "modified"
	if namedExtra.Items()[0].OverlayID == "modified" {
		t.Error("Items() should return a copy, not the original slice")
	}
}

func TestExtraErrors(t *testing.T) {
	if extra.ErrExtraNotFound.Error() != "extra not found" {
		t.Errorf("unexpected error message: %v", extra.ErrExtraNotFound)
	}
	if extra.ErrExtraAlreadyExists.Error() != "extra already exists" {
		t.Errorf("unexpected error message: %v", extra.ErrExtraAlreadyExists)
	}
}
