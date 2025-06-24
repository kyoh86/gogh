package extra_describe_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/extra_describe"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/repository"
)

func TestUseCaseJSON_Execute_AutoExtra(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	repo := repository.NewReference("github.com", "owner", "repo")
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{
		{OverlayID: "overlay1", HookID: "hook1"},
		{OverlayID: "overlay2", HookID: "hook2"},
	}

	e := extra.NewAutoExtra(extraID, repo, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseJSON(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != extraID {
		t.Errorf("Expected id %s, got %v", extraID, result["id"])
	}
	if result["type"] != string(extra.TypeAuto) {
		t.Errorf("Expected type %s, got %v", extra.TypeAuto, result["type"])
	}
	if result["repository"] != repo.String() {
		t.Errorf("Expected repository %s, got %v", repo.String(), result["repository"])
	}
	if result["source"] != source.String() {
		t.Errorf("Expected source %s, got %v", source.String(), result["source"])
	}

	resultItems, ok := result["items"].([]interface{})
	if !ok || len(resultItems) != 2 {
		t.Errorf("Expected 2 items, got %v", result["items"])
	}
}

func TestUseCaseJSON_Execute_NamedExtra(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	extraName := "test-extra"
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{
		{OverlayID: "overlay1", HookID: "hook1"},
	}

	e := extra.NewNamedExtra(extraID, extraName, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseJSON(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != extraID {
		t.Errorf("Expected id %s, got %v", extraID, result["id"])
	}
	if result["type"] != string(extra.TypeNamed) {
		t.Errorf("Expected type %s, got %v", extra.TypeNamed, result["type"])
	}
	if result["name"] != extraName {
		t.Errorf("Expected name %s, got %v", extraName, result["name"])
	}
	if _, hasRepo := result["repository"]; hasRepo {
		t.Error("Named extra should not have repository field")
	}
}

func TestUseCaseOneLine_Execute_AutoExtra(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	repo := repository.NewReference("github.com", "owner", "repo")
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{
		{OverlayID: "overlay1", HookID: "hook1"},
		{OverlayID: "overlay2", HookID: "hook2"},
	}

	e := extra.NewAutoExtra(extraID, repo, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseOneLine(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedParts := []string{
		"[" + extraID[:8] + "]",
		string(extra.TypeAuto),
		repo.String(),
		"(2 items)",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected output to contain '%s', but it doesn't: %s", part, output)
		}
	}
}

func TestUseCaseOneLine_Execute_NamedExtra(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	extraName := "test-extra"
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{
		{OverlayID: "overlay1", HookID: "hook1"},
	}

	e := extra.NewNamedExtra(extraID, extraName, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseOneLine(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedParts := []string{
		"[" + extraID[:8] + "]",
		string(extra.TypeNamed),
		extraName,
		"(1 items)",
	}

	for _, part := range expectedParts {
		if !strings.Contains(output, part) {
			t.Errorf("Expected output to contain '%s', but it doesn't: %s", part, output)
		}
	}
}

func TestUseCaseDetail_Execute_AutoExtra(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	repo := repository.NewReference("github.com", "owner", "repo")
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{
		{OverlayID: "overlay1", HookID: "hook1"},
		{OverlayID: "overlay2", HookID: ""},
	}

	e := extra.NewAutoExtra(extraID, repo, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseDetail(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedLines := []string{
		"ID: " + extraID,
		"Type: " + string(extra.TypeAuto),
		"Repository: " + repo.String(),
		"Source: " + source.String(),
		"Created: " + createdAt.Format("2006-01-02 15:04:05"),
		"Items (2):",
		"1. Overlay: overlay1",
		"   Hook: hook1",
		"2. Overlay: overlay2",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nFull output:\n%s", expected, output)
		}
	}

	// Check that hook2 is NOT present (since it's empty)
	if strings.Contains(output, "Hook: hook2") {
		t.Error("Expected empty hook ID to not be displayed")
	}
}

func TestUseCaseDetail_Execute_NamedExtra(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	extraName := "test-extra"
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{
		{OverlayID: "overlay1", HookID: "hook1"},
	}

	e := extra.NewNamedExtra(extraID, extraName, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseDetail(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedLines := []string{
		"ID: " + extraID,
		"Type: " + string(extra.TypeNamed),
		"Name: " + extraName,
		"Source: " + source.String(),
		"Created: " + createdAt.Format("2006-01-02 15:04:05"),
		"Items (1):",
		"1. Overlay: overlay1",
		"   Hook: hook1",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nFull output:\n%s", expected, output)
		}
	}

	// Check that repository is NOT present for named extra
	if strings.Contains(output, "Repository:") {
		t.Error("Named extra should not display repository")
	}
}

func TestUseCaseDetail_Execute_EmptyItems(t *testing.T) {
	ctx := context.Background()

	extraID := uuid.New().String()
	extraName := "empty-extra"
	source := repository.NewReference("github.com", "source", "repo")
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	items := []extra.Item{}

	e := extra.NewNamedExtra(extraID, extraName, source, items, createdAt)

	var buf bytes.Buffer
	uc := extra_describe.NewUseCaseDetail(&buf)

	err := uc.Execute(ctx, *e)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Items (0):") {
		t.Error("Expected 'Items (0):' for empty items")
	}
}
