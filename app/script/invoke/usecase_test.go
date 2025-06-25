package invoke_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"os/exec"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/script/invoke"
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

type mockScriptService struct {
	openFunc func(ctx context.Context, id string) (io.ReadCloser, error)
}

func (m *mockScriptService) List() iter.Seq2[script.Script, error] {
	return func(yield func(script.Script, error) bool) {}
}

func (m *mockScriptService) Add(ctx context.Context, entry script.Entry) (string, error) {
	return "", errors.New("not implemented")
}

func (m *mockScriptService) Get(ctx context.Context, id string) (script.Script, error) {
	return nil, errors.New("not implemented")
}

func (m *mockScriptService) Update(ctx context.Context, id string, entry script.Entry) error {
	return errors.New("not implemented")
}

func (m *mockScriptService) Remove(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (m *mockScriptService) Open(ctx context.Context, id string) (io.ReadCloser, error) {
	if m.openFunc != nil {
		return m.openFunc(ctx, id)
	}
	return nil, errors.New("script not found")
}

func (m *mockScriptService) Load(iter.Seq2[script.Script, error]) error {
	return errors.New("not implemented")
}

func (m *mockScriptService) HasChanges() bool { return false }
func (m *mockScriptService) MarkSaved()       {}

type mockReferenceParser struct {
	parseWithAliasFunc func(refStr string) (*repository.ReferenceWithAlias, error)
}

func (m *mockReferenceParser) Parse(refStr string) (*repository.Reference, error) {
	// Simple parser for testing
	if refStr == "github.com/kyoh86/gogh" {
		ref := repository.NewReference("github.com", "kyoh86", "gogh")
		return &ref, nil
	}
	return nil, fmt.Errorf("invalid reference: %s", refStr)
}

func (m *mockReferenceParser) ParseWithAlias(refStr string) (*repository.ReferenceWithAlias, error) {
	if m.parseWithAliasFunc != nil {
		return m.parseWithAliasFunc(refStr)
	}
	ref, err := m.Parse(refStr)
	if err != nil {
		return nil, err
	}
	return &repository.ReferenceWithAlias{Reference: *ref}, nil
}

// Helper to create a ReadCloser from a string - removed as unused

// Tests

func TestNewUseCase(t *testing.T) {
	ws := &mockWorkspaceService{}
	finder := &mockFinderService{}
	scripts := &mockScriptService{}
	parser := &mockReferenceParser{}

	uc := testtarget.NewUseCase(ws, finder, scripts, parser)
	if uc == nil {
		t.Fatal("expected non-nil UseCase")
	}
}

func TestUseCase_Execute(t *testing.T) {
	// Skip this test if we're not in a real environment
	if _, err := exec.LookPath(os.Args[0]); err != nil {
		t.Skip("Test requires real executable")
	}

	ctx := context.Background()

	tests := []struct {
		name      string
		refStr    string
		scriptID  string
		globals   map[string]any
		setupMock func() (*mockFinderService, *mockScriptService, *mockReferenceParser)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "invalid reference",
			refStr:   "invalid-ref",
			scriptID: "script-123",
			setupMock: func() (*mockFinderService, *mockScriptService, *mockReferenceParser) {
				return &mockFinderService{}, &mockScriptService{}, &mockReferenceParser{}
			},
			wantErr: true,
			errMsg:  "parsing repository reference",
		},
		{
			name:     "repository not found",
			refStr:   "github.com/kyoh86/gogh",
			scriptID: "script-123",
			setupMock: func() (*mockFinderService, *mockScriptService, *mockReferenceParser) {
				finder := &mockFinderService{
					findByReferenceFunc: func(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (*repository.Location, error) {
						return nil, errors.New("not found")
					},
				}
				return finder, &mockScriptService{}, &mockReferenceParser{}
			},
			wantErr: true,
			errMsg:  "find repository location",
		},
		{
			name:     "script not found",
			refStr:   "github.com/kyoh86/gogh",
			scriptID: "script-123",
			setupMock: func() (*mockFinderService, *mockScriptService, *mockReferenceParser) {
				scripts := &mockScriptService{
					openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
						return nil, errors.New("script not found")
					},
				}
				return &mockFinderService{}, scripts, &mockReferenceParser{}
			},
			wantErr: true,
			errMsg:  "open script",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder, scripts, parser := tt.setupMock()
			uc := testtarget.NewUseCase(&mockWorkspaceService{}, finder, scripts, parser)

			err := uc.Execute(ctx, tt.refStr, tt.scriptID, tt.globals)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("Execute() error = %v, want error containing %q", err, tt.errMsg)
			}
		})
	}
}

func TestUseCase_Invoke(t *testing.T) {
	// Skip this test if we're not in a real environment
	if _, err := exec.LookPath(os.Args[0]); err != nil {
		t.Skip("Test requires real executable")
	}

	ctx := context.Background()

	t.Run("nil location", func(t *testing.T) {
		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			&mockScriptService{},
			&mockReferenceParser{},
		)

		err := uc.Invoke(ctx, nil, "script-123", nil)
		if err == nil {
			t.Error("expected error for nil location")
		}
		if !contains(err.Error(), "repository not found") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("script open error", func(t *testing.T) {
		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return nil, errors.New("open failed")
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/repo", "github.com", "kyoh86", "gogh")
		err := uc.Invoke(ctx, location, "script-123", nil)
		if err == nil {
			t.Error("expected error")
		}
		if !contains(err.Error(), "open script") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("script read error", func(t *testing.T) {
		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return &errorReadCloser{}, nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/repo", "github.com", "kyoh86", "gogh")
		err := uc.Invoke(ctx, location, "script-123", nil)
		if err == nil {
			t.Error("expected error")
		}
		if !contains(err.Error(), "read script") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("with globals", func(t *testing.T) {
		t.Skip("Skipping test that requires subprocess execution")
	})
}

// Helper types and functions

type errorReadCloser struct{}

func (e *errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}

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
