package script_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/script"
)

func TestNewScript(t *testing.T) {
	content := bytes.NewReader([]byte("print('hello world')"))
	entry := script.Entry{
		Name:    "test-script",
		Content: content,
	}

	beforeCreate := time.Now()
	s := script.NewScript(entry)
	afterCreate := time.Now()

	// Test Name
	if s.Name() != entry.Name {
		t.Errorf("Name() = %v, want %v", s.Name(), entry.Name)
	}

	// Test ID and UUID
	if s.ID() == "" {
		t.Error("ID() should not be empty")
	}
	if s.UUID() == uuid.Nil {
		t.Error("UUID() should not be nil")
	}
	// Verify ID is the string representation of UUID
	if s.ID() != s.UUID().String() {
		t.Errorf("ID() = %v, want %v", s.ID(), s.UUID().String())
	}

	// Test timestamps
	if s.CreatedAt().Before(beforeCreate) || s.CreatedAt().After(afterCreate) {
		t.Errorf("CreatedAt() = %v, want between %v and %v", s.CreatedAt(), beforeCreate, afterCreate)
	}
	if s.UpdatedAt().Before(beforeCreate) || s.UpdatedAt().After(afterCreate) {
		t.Errorf("UpdatedAt() = %v, want between %v and %v", s.UpdatedAt(), beforeCreate, afterCreate)
	}
	// CreatedAt and UpdatedAt should be the same for new scripts
	if !s.CreatedAt().Equal(s.UpdatedAt()) {
		t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want equal", s.CreatedAt(), s.UpdatedAt())
	}
}

func TestConcreteScript(t *testing.T) {
	id := uuid.New()
	name := "concrete-script"
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	s := script.ConcreteScript(id, name, createdAt, updatedAt)

	// Test ID and UUID
	if s.ID() != id.String() {
		t.Errorf("ID() = %v, want %v", s.ID(), id.String())
	}
	if s.UUID() != id {
		t.Errorf("UUID() = %v, want %v", s.UUID(), id)
	}

	// Test Name
	if s.Name() != name {
		t.Errorf("Name() = %v, want %v", s.Name(), name)
	}

	// Test timestamps
	if !s.CreatedAt().Equal(createdAt) {
		t.Errorf("CreatedAt() = %v, want %v", s.CreatedAt(), createdAt)
	}
	if !s.UpdatedAt().Equal(updatedAt) {
		t.Errorf("UpdatedAt() = %v, want %v", s.UpdatedAt(), updatedAt)
	}
}

func TestMultipleScriptsHaveUniqueIDs(t *testing.T) {
	entry := script.Entry{
		Name:    "test-script",
		Content: bytes.NewReader([]byte("test")),
	}

	scripts := make([]script.Script, 10)
	ids := make(map[string]bool)

	for i := range scripts {
		scripts[i] = script.NewScript(entry)
		id := scripts[i].ID()

		if ids[id] {
			t.Errorf("Duplicate ID found: %v", id)
		}
		ids[id] = true
	}
}

func TestScriptTimestamps(t *testing.T) {
	// Test that timestamps are properly set and maintained
	pastTime := time.Now().Add(-48 * time.Hour)
	recentTime := time.Now().Add(-1 * time.Hour)

	s := script.ConcreteScript(
		uuid.New(),
		"timestamp-test",
		pastTime,
		recentTime,
	)

	if !s.CreatedAt().Equal(pastTime) {
		t.Errorf("CreatedAt() = %v, want %v", s.CreatedAt(), pastTime)
	}
	if !s.UpdatedAt().Equal(recentTime) {
		t.Errorf("UpdatedAt() = %v, want %v", s.UpdatedAt(), recentTime)
	}

	// Verify CreatedAt is before UpdatedAt
	if !s.CreatedAt().Before(s.UpdatedAt()) {
		t.Error("CreatedAt should be before UpdatedAt for an updated script")
	}
}

func TestScriptWithEmptyName(t *testing.T) {
	// Test that scripts can have empty names
	entry := script.Entry{
		Name:    "",
		Content: bytes.NewReader([]byte("anonymous script")),
	}

	s := script.NewScript(entry)

	if s.Name() != "" {
		t.Errorf("Name() = %v, want empty string", s.Name())
	}

	// Other properties should still work
	if s.ID() == "" {
		t.Error("ID() should not be empty even with empty name")
	}
}

func TestScriptInterfaceImplementation(t *testing.T) {
	// Verify that scriptElement implements Script interface
	_ = script.NewScript(script.Entry{Name: "test"})
	_ = script.ConcreteScript(uuid.New(), "test", time.Now(), time.Now())
}
