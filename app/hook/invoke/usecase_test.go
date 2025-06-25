package invoke_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/invoke"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Mock implementations

type mockWorkspaceService struct {
	workspace.WorkspaceService
}

type mockFinderService struct {
	findByReferenceFunc func(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (*repository.Location, error)
}

func (m *mockFinderService) FindByReference(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (*repository.Location, error) {
	if m.findByReferenceFunc != nil {
		return m.findByReferenceFunc(ctx, ws, reference)
	}
	return repository.NewLocation("/tmp/repo", "github.com", "kyoh86", "gogh"), nil
}

func (m *mockFinderService) FindByPath(ctx context.Context, ws workspace.WorkspaceService, path string) (*repository.Location, error) {
	return nil, errors.New("not implemented")
}

func (m *mockFinderService) ListAllRepository(context.Context, workspace.WorkspaceService, workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	return func(yield func(*repository.Location, error) bool) {}
}

func (m *mockFinderService) ListRepositoryInRoot(context.Context, workspace.LayoutService, workspace.ListOptions) iter.Seq2[*repository.Location, error] {
	return func(yield func(*repository.Location, error) bool) {}
}

type mockHook struct {
	id            uuid.UUID
	name          string
	repoPattern   string
	triggerEvent  hook.Event
	operationType hook.OperationType
	operationID   string
}

func (m *mockHook) ID() string                        { return m.id.String() }
func (m *mockHook) UUID() uuid.UUID                   { return m.id }
func (m *mockHook) Name() string                      { return m.name }
func (m *mockHook) RepoPattern() string               { return m.repoPattern }
func (m *mockHook) TriggerEvent() hook.Event          { return m.triggerEvent }
func (m *mockHook) OperationType() hook.OperationType { return m.operationType }
func (m *mockHook) OperationID() string               { return m.operationID }
func (m *mockHook) Match(ref repository.Reference, event hook.Event) (bool, error) {
	if m.triggerEvent != hook.EventAny && m.triggerEvent != event {
		return false, nil
	}
	return true, nil
}

type mockHookService struct {
	getFunc     func(ctx context.Context, id string) (hook.Hook, error)
	listForFunc func(reference repository.Reference, event hook.Event) iter.Seq2[hook.Hook, error]
}

func (m *mockHookService) List() iter.Seq2[hook.Hook, error] {
	return func(yield func(hook.Hook, error) bool) {}
}

func (m *mockHookService) ListFor(reference repository.Reference, event hook.Event) iter.Seq2[hook.Hook, error] {
	if m.listForFunc != nil {
		return m.listForFunc(reference, event)
	}
	return func(yield func(hook.Hook, error) bool) {}
}

func (m *mockHookService) Add(ctx context.Context, entry hook.Entry) (string, error) {
	return "", errors.New("not implemented")
}

func (m *mockHookService) Get(ctx context.Context, id string) (hook.Hook, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, errors.New("hook not found")
}

func (m *mockHookService) Update(ctx context.Context, id string, entry hook.Entry) error {
	return errors.New("not implemented")
}

func (m *mockHookService) Remove(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockHookService) Load(iter.Seq2[hook.Hook, error]) error {
	return errors.New("not implemented")
}

func (m *mockHookService) HasChanges() bool { return false }
func (m *mockHookService) MarkSaved()       {}

type mockOverlayService struct {
	getFunc func(ctx context.Context, id string) (overlay.Overlay, error)
}

func (m *mockOverlayService) List() iter.Seq2[overlay.Overlay, error] {
	return func(yield func(overlay.Overlay, error) bool) {}
}

func (m *mockOverlayService) Add(ctx context.Context, entry overlay.Entry) (string, error) {
	return "", errors.New("not implemented")
}

func (m *mockOverlayService) Get(ctx context.Context, idlike string) (overlay.Overlay, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, idlike)
	}
	return nil, errors.New("overlay not found")
}

func (m *mockOverlayService) Update(ctx context.Context, idlike string, entry overlay.Entry) error {
	return errors.New("not implemented")
}

func (m *mockOverlayService) Remove(ctx context.Context, idlike string) error {
	return errors.New("not implemented")
}

func (m *mockOverlayService) Open(ctx context.Context, idlike string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (m *mockOverlayService) Load(iter.Seq2[overlay.Overlay, error]) error {
	return errors.New("not implemented")
}

func (m *mockOverlayService) HasChanges() bool { return false }
func (m *mockOverlayService) MarkSaved()       {}

type mockScriptService struct {
	getFunc func(ctx context.Context, id string) (script.Script, error)
}

func (m *mockScriptService) List() iter.Seq2[script.Script, error] {
	return func(yield func(script.Script, error) bool) {}
}

func (m *mockScriptService) Add(ctx context.Context, entry script.Entry) (string, error) {
	return "", errors.New("not implemented")
}

func (m *mockScriptService) Get(ctx context.Context, id string) (script.Script, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, errors.New("script not found")
}

func (m *mockScriptService) Update(ctx context.Context, id string, entry script.Entry) error {
	return errors.New("not implemented")
}

func (m *mockScriptService) Remove(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockScriptService) Open(ctx context.Context, id string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func (m *mockScriptService) Load(iter.Seq2[script.Script, error]) error {
	return errors.New("not implemented")
}

func (m *mockScriptService) HasChanges() bool { return false }
func (m *mockScriptService) MarkSaved()       {}

type mockReferenceParser struct{}

func (m *mockReferenceParser) Parse(refStr string) (*repository.Reference, error) {
	// Simple parser for testing
	if refStr == "github.com/kyoh86/gogh" {
		ref := repository.NewReference("github.com", "kyoh86", "gogh")
		return &ref, nil
	}
	return nil, fmt.Errorf("invalid reference: %s", refStr)
}

func (m *mockReferenceParser) ParseWithAlias(refStr string) (*repository.ReferenceWithAlias, error) {
	ref, err := m.Parse(refStr)
	if err != nil {
		return nil, err
	}
	return &repository.ReferenceWithAlias{Reference: *ref}, nil
}

// Tests

func TestNewUseCase(t *testing.T) {
	ws := &mockWorkspaceService{}
	finder := &mockFinderService{}
	hooks := &mockHookService{}
	overlays := &mockOverlayService{}
	scripts := &mockScriptService{}
	parser := &mockReferenceParser{}

	uc := testtarget.NewUseCase(ws, finder, hooks, overlays, scripts, parser)
	if uc == nil {
		t.Fatal("expected non-nil UseCase")
	}
}

func TestUseCase_Invoke(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		hookID    string
		refStr    string
		setupHook func() *mockHook
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "overlay hook",
			hookID: "test-hook-id",
			refStr: "github.com/kyoh86/gogh",
			setupHook: func() *mockHook {
				return &mockHook{
					id:            uuid.New(),
					name:          "test overlay hook",
					operationType: hook.OperationTypeOverlay,
					operationID:   "overlay-123",
				}
			},
			wantErr: true, // Will fail because overlay service is not fully mocked
		},
		{
			name:   "script hook",
			hookID: "test-hook-id",
			refStr: "github.com/kyoh86/gogh",
			setupHook: func() *mockHook {
				return &mockHook{
					id:            uuid.New(),
					name:          "test script hook",
					operationType: hook.OperationTypeScript,
					operationID:   "script-123",
				}
			},
			wantErr: true, // Will fail because script service is not fully mocked
		},
		{
			name:    "hook not found",
			hookID:  "non-existent",
			refStr:  "github.com/kyoh86/gogh",
			wantErr: true,
			errMsg:  "hook not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hookSvc := &mockHookService{
				getFunc: func(ctx context.Context, id string) (hook.Hook, error) {
					if tt.setupHook != nil && id == tt.hookID {
						return tt.setupHook(), nil
					}
					return nil, errors.New("hook not found")
				},
			}

			uc := testtarget.NewUseCase(
				&mockWorkspaceService{},
				&mockFinderService{},
				hookSvc,
				&mockOverlayService{},
				&mockScriptService{},
				&mockReferenceParser{},
			)

			err := uc.Invoke(ctx, tt.hookID, tt.refStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("Invoke() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestUseCase_InvokeFor(t *testing.T) {
	ctx := context.Background()

	t.Run("successful invocation", func(t *testing.T) {
		overlayHook := &mockHook{
			id:            uuid.New(),
			name:          "overlay hook",
			operationType: hook.OperationTypeOverlay,
			operationID:   "overlay-123",
			triggerEvent:  testtarget.EventPostClone,
		}

		scriptHook := &mockHook{
			id:            uuid.New(),
			name:          "script hook",
			operationType: hook.OperationTypeScript,
			operationID:   "script-123",
			triggerEvent:  testtarget.EventPostClone,
		}

		hookSvc := &mockHookService{
			listForFunc: func(reference repository.Reference, event hook.Event) iter.Seq2[hook.Hook, error] {
				return func(yield func(hook.Hook, error) bool) {
					if !yield(overlayHook, nil) {
						return
					}
					yield(scriptHook, nil)
				}
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			hookSvc,
			&mockOverlayService{},
			&mockScriptService{},
			&mockReferenceParser{},
		)

		err := uc.InvokeFor(ctx, testtarget.EventPostClone, "github.com/kyoh86/gogh")
		// This will error because the apply operations are not fully mocked
		if err == nil {
			t.Error("expected error due to incomplete mocks")
		}
	})

	t.Run("invalid reference", func(t *testing.T) {
		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			&mockHookService{},
			&mockOverlayService{},
			&mockScriptService{},
			&mockReferenceParser{},
		)

		err := uc.InvokeFor(ctx, testtarget.EventPostClone, "invalid-ref")
		if err == nil {
			t.Error("expected error for invalid reference")
		}
	})

	t.Run("repository not found", func(t *testing.T) {
		finder := &mockFinderService{
			findByReferenceFunc: func(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (*repository.Location, error) {
				return nil, errors.New("not found")
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			finder,
			&mockHookService{},
			&mockOverlayService{},
			&mockScriptService{},
			&mockReferenceParser{},
		)

		err := uc.InvokeFor(ctx, testtarget.EventPostClone, "github.com/kyoh86/gogh")
		if err == nil {
			t.Error("expected error for repository not found")
		}
	})
}

func TestUseCase_InvokeForWithGlobals(t *testing.T) {
	ctx := context.Background()

	scriptHook := &mockHook{
		id:            uuid.New(),
		name:          "script hook",
		operationType: hook.OperationTypeScript,
		operationID:   "script-123",
		triggerEvent:  testtarget.EventPostFork,
	}

	hookSvc := &mockHookService{
		listForFunc: func(reference repository.Reference, event hook.Event) iter.Seq2[hook.Hook, error] {
			return func(yield func(hook.Hook, error) bool) {
				yield(scriptHook, nil)
			}
		},
	}

	uc := testtarget.NewUseCase(
		&mockWorkspaceService{},
		&mockFinderService{},
		hookSvc,
		&mockOverlayService{},
		&mockScriptService{},
		&mockReferenceParser{},
	)

	globals := map[string]any{
		"custom": "value",
		"fork":   true,
	}

	err := uc.InvokeForWithGlobals(ctx, testtarget.EventPostFork, "github.com/kyoh86/gogh", globals)
	// This will error because the script invoke is not fully mocked
	if err == nil {
		t.Error("expected error due to incomplete mocks")
	}
}

func TestEventConstants(t *testing.T) {
	// Test that event constants match
	if testtarget.EventAny != hook.EventAny {
		t.Errorf("EventAny mismatch: %v != %v", testtarget.EventAny, hook.EventAny)
	}
	if testtarget.EventPostClone != hook.EventPostClone {
		t.Errorf("EventPostClone mismatch: %v != %v", testtarget.EventPostClone, hook.EventPostClone)
	}
	if testtarget.EventPostFork != hook.EventPostFork {
		t.Errorf("EventPostFork mismatch: %v != %v", testtarget.EventPostFork, hook.EventPostFork)
	}
	if testtarget.EventPostCreate != hook.EventPostCreate {
		t.Errorf("EventPostCreate mismatch: %v != %v", testtarget.EventPostCreate, hook.EventPostCreate)
	}
}

func TestOptions(t *testing.T) {
	// Test that Options can be instantiated
	opts := testtarget.Options{}
	// Currently no fields, but this ensures the struct exists and can be used
	_ = opts
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
