package script_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/script"
)

// mockScriptSourceStore implements ScriptSourceStore for testing
type mockScriptSourceStore struct {
	contents  map[string][]byte
	saveErr   error
	openErr   error
	removeErr error
}

func newMockScriptSourceStore() *mockScriptSourceStore {
	return &mockScriptSourceStore{
		contents: make(map[string][]byte),
	}
}

func (m *mockScriptSourceStore) Save(ctx context.Context, scriptID string, content io.Reader) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	data, err := io.ReadAll(content)
	if err != nil {
		return err
	}
	m.contents[scriptID] = data
	return nil
}

func (m *mockScriptSourceStore) Open(ctx context.Context, scriptID string) (io.ReadCloser, error) {
	if m.openErr != nil {
		return nil, m.openErr
	}
	content, exists := m.contents[scriptID]
	if !exists {
		return nil, errors.New("script not found")
	}
	return io.NopCloser(bytes.NewReader(content)), nil
}

func (m *mockScriptSourceStore) Remove(ctx context.Context, scriptID string) error {
	if m.removeErr != nil {
		return m.removeErr
	}
	delete(m.contents, scriptID)
	return nil
}

func TestNewScriptService(t *testing.T) {
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)
	if service == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestScriptService_Add(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	t.Run("success", func(t *testing.T) {
		entry := script.Entry{
			Name:    "test-script",
			Content: strings.NewReader("print('hello world')"),
		}

		id, err := service.Add(ctx, entry)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id == "" {
			t.Error("expected non-empty ID")
		}

		// Verify content was saved
		content, exists := store.contents[id]
		if !exists {
			t.Error("expected content to be saved")
		}
		if string(content) != "print('hello world')" {
			t.Errorf("expected content 'print('hello world')', got %s", string(content))
		}

		// Verify script was added
		s, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("failed to get script: %v", err)
		}
		if s.ID() != id {
			t.Errorf("expected ID %s, got %s", id, s.ID())
		}
		if s.Name() != "test-script" {
			t.Errorf("expected name 'test-script', got %s", s.Name())
		}
	})

	t.Run("save error", func(t *testing.T) {
		storeWithErr := newMockScriptSourceStore()
		storeWithErr.saveErr = errors.New("save failed")
		serviceWithErr := script.NewScriptService(storeWithErr)

		entry := script.Entry{
			Name:    "fail-script",
			Content: strings.NewReader("content"),
		}

		_, err := serviceWithErr.Add(ctx, entry)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestScriptService_Update(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Add initial script
	entry := script.Entry{
		Name:    "test-script",
		Content: strings.NewReader("original content"),
	}
	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add script: %v", err)
	}

	// Get original timestamps
	originalScript, _ := service.Get(ctx, id)
	originalCreatedAt := originalScript.CreatedAt()
	originalUpdatedAt := originalScript.UpdatedAt()

	// Wait a bit to ensure timestamps are different
	time.Sleep(10 * time.Millisecond)

	t.Run("update content only", func(t *testing.T) {
		updateEntry := script.Entry{
			Content: strings.NewReader("updated content"),
		}

		err := service.Update(ctx, id, updateEntry)
		if err != nil {
			t.Fatalf("failed to update script: %v", err)
		}

		// Verify content was updated
		content := store.contents[id]
		if string(content) != "updated content" {
			t.Errorf("expected content 'updated content', got %s", string(content))
		}

		// Verify script metadata
		s, _ := service.Get(ctx, id)
		if s.Name() != "test-script" {
			t.Errorf("name should not have changed, got %s", s.Name())
		}
		if !s.CreatedAt().Equal(originalCreatedAt) {
			t.Error("CreatedAt should not change on update")
		}
		if s.UpdatedAt().Equal(originalUpdatedAt) {
			t.Error("UpdatedAt should have changed")
		}
	})

	t.Run("update name only", func(t *testing.T) {
		time.Sleep(10 * time.Millisecond)
		updateEntry := script.Entry{
			Name: "renamed-script",
		}

		err := service.Update(ctx, id, updateEntry)
		if err != nil {
			t.Fatalf("failed to update script: %v", err)
		}

		// Verify name was updated
		s, _ := service.Get(ctx, id)
		if s.Name() != "renamed-script" {
			t.Errorf("expected name 'renamed-script', got %s", s.Name())
		}
	})

	t.Run("update both", func(t *testing.T) {
		time.Sleep(10 * time.Millisecond)
		updateEntry := script.Entry{
			Name:    "final-script",
			Content: strings.NewReader("final content"),
		}

		err := service.Update(ctx, id, updateEntry)
		if err != nil {
			t.Fatalf("failed to update script: %v", err)
		}

		s, _ := service.Get(ctx, id)
		if s.Name() != "final-script" {
			t.Errorf("expected name 'final-script', got %s", s.Name())
		}

		content := store.contents[id]
		if string(content) != "final content" {
			t.Errorf("expected content 'final content', got %s", string(content))
		}
	})

	t.Run("update non-existent", func(t *testing.T) {
		err := service.Update(ctx, "non-existent", script.Entry{Name: "new"})
		if err == nil {
			t.Error("expected error when updating non-existent script")
		}
	})

	t.Run("save error", func(t *testing.T) {
		storeWithErr := newMockScriptSourceStore()
		serviceWithErr := script.NewScriptService(storeWithErr)

		// Add a script first
		entry := script.Entry{
			Name:    "test",
			Content: strings.NewReader("content"),
		}
		id, _ := serviceWithErr.Add(ctx, entry)

		// Now make Save fail
		storeWithErr.saveErr = errors.New("save failed")

		err := serviceWithErr.Update(ctx, id, script.Entry{Content: strings.NewReader("new")})
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestScriptService_Get(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Add a script
	entry := script.Entry{
		Name:    "test-script",
		Content: strings.NewReader("content"),
	}
	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add script: %v", err)
	}

	t.Run("get by full ID", func(t *testing.T) {
		s, err := service.Get(ctx, id)
		if err != nil {
			t.Fatalf("failed to get script: %v", err)
		}
		if s.ID() != id {
			t.Errorf("expected ID %s, got %s", id, s.ID())
		}
	})

	t.Run("get by partial ID", func(t *testing.T) {
		if len(id) > 8 {
			partialID := id[:8]
			s, err := service.Get(ctx, partialID)
			if err != nil {
				t.Fatalf("failed to get script by partial ID: %v", err)
			}
			if s.ID() != id {
				t.Errorf("expected ID %s, got %s", id, s.ID())
			}
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		_, err := service.Get(ctx, "non-existent")
		if err == nil {
			t.Error("expected error when getting non-existent script")
		}
	})
}

func TestScriptService_Remove(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Add a script
	entry := script.Entry{
		Name:    "test-script",
		Content: strings.NewReader("content"),
	}
	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add script: %v", err)
	}

	t.Run("remove existing", func(t *testing.T) {
		err := service.Remove(ctx, id)
		if err != nil {
			t.Fatalf("failed to remove script: %v", err)
		}

		// Verify script was removed
		_, err = service.Get(ctx, id)
		if err == nil {
			t.Error("expected error when getting removed script")
		}

		// Verify content was removed
		_, exists := store.contents[id]
		if exists {
			t.Error("expected content to be removed")
		}
	})

	t.Run("remove non-existent", func(t *testing.T) {
		err := service.Remove(ctx, "non-existent")
		if err == nil {
			t.Error("expected error when removing non-existent script")
		}
	})

	t.Run("remove error from store", func(t *testing.T) {
		storeWithErr := newMockScriptSourceStore()
		storeWithErr.removeErr = errors.New("remove failed")
		serviceWithErr := script.NewScriptService(storeWithErr)

		// Add a script
		entry := script.Entry{
			Name:    "test",
			Content: strings.NewReader("content"),
		}
		id, _ := serviceWithErr.Add(ctx, entry)

		err := serviceWithErr.Remove(ctx, id)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestScriptService_Open(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Add a script
	content := "print('hello world')"
	entry := script.Entry{
		Name:    "test-script",
		Content: strings.NewReader(content),
	}
	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add script: %v", err)
	}

	t.Run("open existing", func(t *testing.T) {
		reader, err := service.Open(ctx, id)
		if err != nil {
			t.Fatalf("failed to open script: %v", err)
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read content: %v", err)
		}
		if string(data) != content {
			t.Errorf("expected content %s, got %s", content, string(data))
		}
	})

	t.Run("open non-existent script", func(t *testing.T) {
		_, err := service.Open(ctx, "non-existent")
		if err == nil {
			t.Error("expected error when opening non-existent script")
		}
	})

	t.Run("open error from store", func(t *testing.T) {
		storeWithErr := newMockScriptSourceStore()
		storeWithErr.openErr = errors.New("open failed")
		serviceWithErr := script.NewScriptService(storeWithErr)

		// Add a script
		entry := script.Entry{
			Name:    "test",
			Content: strings.NewReader("content"),
		}
		id, _ := serviceWithErr.Add(ctx, entry)

		_, err := serviceWithErr.Open(ctx, id)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestScriptService_List(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Add multiple scripts
	scripts := []script.Entry{
		{Name: "script1", Content: strings.NewReader("content1")},
		{Name: "script2", Content: strings.NewReader("content2")},
		{Name: "script3", Content: strings.NewReader("content3")},
	}

	var ids []string
	for _, entry := range scripts {
		id, err := service.Add(ctx, entry)
		if err != nil {
			t.Fatalf("failed to add script: %v", err)
		}
		ids = append(ids, id)
	}

	// List all scripts
	var count int
	foundIDs := make(map[string]bool)
	for s, err := range service.List() {
		if err != nil {
			t.Fatalf("unexpected error in List: %v", err)
		}
		foundIDs[s.ID()] = true
		count++
	}

	if count != len(scripts) {
		t.Errorf("expected %d scripts, got %d", len(scripts), count)
	}

	// Verify all IDs are present
	for _, id := range ids {
		if !foundIDs[id] {
			t.Errorf("expected ID %s in list", id)
		}
	}
}

func TestScriptService_Load(t *testing.T) {
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Create scripts to load
	now := time.Now()
	scripts := []script.Script{
		script.ConcreteScript(uuid.New(), "loaded-script-1", now.Add(-24*time.Hour), now.Add(-1*time.Hour)),
		script.ConcreteScript(uuid.New(), "loaded-script-2", now.Add(-48*time.Hour), now.Add(-2*time.Hour)),
		script.ConcreteScript(uuid.New(), "loaded-script-3", now.Add(-72*time.Hour), now),
	}

	// Load scripts
	loadSeq := func(yield func(script.Script, error) bool) {
		for _, s := range scripts {
			if !yield(s, nil) {
				return
			}
		}
	}

	err := service.Load(loadSeq)
	if err != nil {
		t.Fatalf("failed to load scripts: %v", err)
	}

	// Verify scripts were loaded
	var count int
	for s, err := range service.List() {
		if err != nil {
			t.Fatalf("unexpected error listing loaded scripts: %v", err)
		}
		count++

		// Find matching script
		found := false
		for _, original := range scripts {
			if s.ID() == original.ID() {
				found = true
				if s.Name() != original.Name() {
					t.Errorf("expected name %s, got %s", original.Name(), s.Name())
				}
				if !s.CreatedAt().Equal(original.CreatedAt()) {
					t.Errorf("expected CreatedAt %v, got %v", original.CreatedAt(), s.CreatedAt())
				}
				if !s.UpdatedAt().Equal(original.UpdatedAt()) {
					t.Errorf("expected UpdatedAt %v, got %v", original.UpdatedAt(), s.UpdatedAt())
				}
				break
			}
		}
		if !found {
			t.Errorf("unexpected script with ID %s", s.ID())
		}
	}

	if count != len(scripts) {
		t.Errorf("expected %d scripts after load, got %d", len(scripts), count)
	}

	// Test loading with error
	errorSeq := func(yield func(script.Script, error) bool) {
		yield(nil, errors.New("load error"))
	}

	err = service.Load(errorSeq)
	if err == nil {
		t.Error("expected error when loading with error")
	}
}

func TestScriptService_HasChanges(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Initially no changes
	if service.HasChanges() {
		t.Error("expected no changes initially")
	}

	// Add script
	entry := script.Entry{
		Name:    "test-script",
		Content: strings.NewReader("content"),
	}
	id, err := service.Add(ctx, entry)
	if err != nil {
		t.Fatalf("failed to add script: %v", err)
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

	// Update script
	err = service.Update(ctx, id, script.Entry{Name: "updated"})
	if err != nil {
		t.Fatalf("failed to update script: %v", err)
	}

	// Should have changes after update
	if !service.HasChanges() {
		t.Error("expected changes after update")
	}

	// Mark saved again
	service.MarkSaved()

	// Remove script
	err = service.Remove(ctx, id)
	if err != nil {
		t.Fatalf("failed to remove script: %v", err)
	}

	// Should have changes after remove
	if !service.HasChanges() {
		t.Error("expected changes after remove")
	}
}

func TestScriptService_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	store := newMockScriptSourceStore()
	service := script.NewScriptService(store)

	// Add initial scripts
	var ids []string
	for i := 0; i < 5; i++ {
		entry := script.Entry{
			Name:    "script-" + string(rune('0'+i)),
			Content: strings.NewReader("content"),
		}
		id, err := service.Add(ctx, entry)
		if err != nil {
			t.Fatalf("failed to add script: %v", err)
		}
		ids = append(ids, id)
	}

	// Run concurrent operations
	done := make(chan bool)

	// Reader goroutine
	go func() {
		for i := 0; i < 10; i++ {
			for _, err := range service.List() {
				if err != nil {
					t.Errorf("list error: %v", err)
				}
			}
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Getter goroutine
	go func() {
		for i := 0; i < 10; i++ {
			for _, id := range ids {
				_, err := service.Get(ctx, id)
				if err != nil {
					// Script might be removed by another goroutine
					continue
				}
			}
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Updater goroutine
	go func() {
		for i := 0; i < 5; i++ {
			if i < len(ids) {
				_ = service.Update(ctx, ids[i], script.Entry{Name: "updated-" + string(rune('0'+i))})
			}
			time.Sleep(2 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}
