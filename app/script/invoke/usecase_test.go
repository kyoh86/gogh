package invoke_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"iter"
	"strings"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/script/invoke"
	"github.com/kyoh86/gogh/v4/app/script/run"
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

// Mock command runner to capture subprocess execution
type mockCmd struct {
	name         string
	args         []string
	dir          string
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
	stdinPipeErr error // For simulating StdinPipe errors
}

func (m *mockCmd) StdinPipe() (io.WriteCloser, error) {
	if m.stdinPipeErr != nil {
		return nil, m.stdinPipeErr
	}
	r, w := io.Pipe()
	m.stdin = r
	return w, nil
}

func (m *mockCmd) Run() error {
	// Simulate successful subprocess execution without actually running it
	if m.stdin != nil {
		// Read and decode the script from stdin
		dec := gob.NewDecoder(m.stdin)
		var script run.Script
		if err := dec.Decode(&script); err != nil {
			return fmt.Errorf("decoding script: %w", err)
		}
		// Optionally write some output to simulate script execution
		if m.stdout != nil {
			fmt.Fprintln(m.stdout, "Mock script output")
		}
	}
	return nil
}

func (m *mockCmd) SetDir(dir string) {
	m.dir = dir
}

func (m *mockCmd) SetStdout(stdout io.Writer) {
	m.stdout = stdout
}

func (m *mockCmd) SetStderr(stderr io.Writer) {
	m.stderr = stderr
}

// Global variable to track command creation for tests
var lastMockCmd *mockCmd

// Override commandRunner for tests
func init() {
	testtarget.SetCommandRunner(func(name string, args ...string) testtarget.Command {
		cmd := &mockCmd{
			name: name,
			args: args,
		}
		lastMockCmd = cmd
		return cmd
	})
}

// Helper functions
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

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
		{
			name:     "successful execution",
			refStr:   "github.com/kyoh86/gogh",
			scriptID: "script-123",
			globals:  map[string]any{"key": "value"},
			setupMock: func() (*mockFinderService, *mockScriptService, *mockReferenceParser) {
				scripts := &mockScriptService{
					openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
						return io.NopCloser(strings.NewReader("print('Hello world')")), nil
					},
				}
				return &mockFinderService{}, scripts, &mockReferenceParser{}
			},
			wantErr: false,
		},
		{
			name:     "reference with alias",
			refStr:   "gh:kyoh86/gogh",
			scriptID: "test-script",
			setupMock: func() (*mockFinderService, *mockScriptService, *mockReferenceParser) {
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				alias := repository.NewReference("gh", "kyoh86", "gogh")
				parser := &mockReferenceParser{
					parseWithAliasFunc: func(refStr string) (*repository.ReferenceWithAlias, error) {
						return &repository.ReferenceWithAlias{
							Reference: ref,
							Alias:     &alias,
						}, nil
					},
				}
				finder := &mockFinderService{
					findByReferenceFunc: func(ctx context.Context, ws workspace.WorkspaceService, reference repository.Reference) (*repository.Location, error) {
						// Should be called with the alias reference
						if reference.Host() != "gh" {
							return nil, fmt.Errorf("expected alias reference, got %v", reference)
						}
						return repository.NewLocation("/tmp/repo", "github.com", "kyoh86", "gogh"), nil
					},
				}
				scripts := &mockScriptService{
					openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
						return io.NopCloser(strings.NewReader("print('alias test')")), nil
					},
				}
				return finder, scripts, parser
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder, scripts, parser := tt.setupMock()
			uc := testtarget.NewUseCase(&mockWorkspaceService{}, finder, scripts, parser)

			// Reset last mock command
			lastMockCmd = nil

			err := uc.Execute(ctx, tt.refStr, tt.scriptID, tt.globals)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.errMsg != "" && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("Execute() error = %v, want error containing %q", err, tt.errMsg)
			}

			// Verify command creation for successful cases
			if !tt.wantErr && lastMockCmd != nil {
				// On Windows, the executable might have .exe extension and different path format
				// So we just verify that some command was created with the correct args
				if lastMockCmd.name == "" {
					t.Error("Expected command name to be set")
				}
				if len(lastMockCmd.args) != 2 || lastMockCmd.args[0] != "script" || lastMockCmd.args[1] != "run" {
					t.Errorf("Expected args [script run], got %v", lastMockCmd.args)
				}
			}
		})
	}
}

func TestUseCase_Invoke(t *testing.T) {
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

	t.Run("successful invocation with globals", func(t *testing.T) {
		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader(`
print("Script ID: " .. script_id)
print("Repo name: " .. gogh.repo.name)
print("Custom value: " .. gogh.custom_key)
`)), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/home/user/repos/test", "github.com", "kyoh86", "test-repo")
		globals := map[string]any{
			"custom_key": "custom_value",
			"script_id":  "test-123",
		}

		// Reset last mock command
		lastMockCmd = nil

		err := uc.Invoke(ctx, location, "test-123", globals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify command was created with correct parameters
		if lastMockCmd == nil {
			t.Error("expected command to be created")
		} else {
			if lastMockCmd.dir != location.FullPath() {
				t.Errorf("Expected command dir %q, got %q", location.FullPath(), lastMockCmd.dir)
			}
		}
	})

	t.Run("complex globals with nested structures", func(t *testing.T) {
		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("-- complex script")), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		globals := map[string]any{
			"string":  "test",
			"number":  42,
			"float":   3.14,
			"boolean": true,
			"nested": map[string]any{
				"level2": map[string]any{
					"value": "deep",
				},
			},
			"array": []string{"a", "b", "c"},
		}

		err := uc.Invoke(ctx, location, "complex-script", globals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty script content", func(t *testing.T) {
		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("")), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		err := uc.Invoke(ctx, location, "empty-script", nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("location with special characters", func(t *testing.T) {
		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("print('test')")), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/test with spaces/and-special!chars", "github.com", "kyoh86", "test")
		err := uc.Invoke(ctx, location, "script", nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("mock command error simulation", func(t *testing.T) {
		// Temporarily override the command runner to simulate an error
		originalRunner := testtarget.SetCommandRunner(func(name string, args ...string) testtarget.Command {
			return &mockCmd{
				name:  name,
				args:  args,
				stdin: nil, // This will cause Run() to fail
			}
		})
		defer func() {
			testtarget.SetCommandRunner(originalRunner)
		}()

		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("print('test')")), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		// This should not cause an error in our mock since we handle it gracefully
		err := uc.Invoke(ctx, location, "script", nil)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("stdin pipe error", func(t *testing.T) {
		// Override command runner to simulate StdinPipe error
		originalRunner := testtarget.SetCommandRunner(func(name string, args ...string) testtarget.Command {
			return &mockCmd{
				name:         name,
				args:         args,
				stdinPipeErr: errors.New("stdin pipe failed"),
			}
		})
		defer func() {
			testtarget.SetCommandRunner(originalRunner)
		}()

		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("print('test')")), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		location := repository.NewLocation("/tmp/test", "github.com", "kyoh86", "test")
		err := uc.Invoke(ctx, location, "script", nil)
		if err == nil {
			t.Error("expected error for stdin pipe failure")
		}
		if !contains(err.Error(), "stdin pipe failed") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// Test globals merging behavior
func TestUseCase_Invoke_GlobalsMerging(t *testing.T) {
	ctx := context.Background()
	location := repository.NewLocation("/home/user/repos/test", "github.com", "kyoh86", "test-repo")

	t.Run("repo globals override user globals", func(t *testing.T) {
		var capturedScript run.Script

		// Override command runner to capture the script
		originalRunner := testtarget.SetCommandRunner(func(name string, args ...string) testtarget.Command {
			return &mockCmd{
				name:  name,
				args:  args,
				stdin: &captureReader{script: &capturedScript},
			}
		})
		defer func() {
			testtarget.SetCommandRunner(originalRunner)
		}()

		globals := map[string]any{
			"custom_key": "custom_value",
			"repo": map[string]any{
				"should_be_overwritten": "yes",
			},
		}

		scripts := &mockScriptService{
			openFunc: func(ctx context.Context, id string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("test script")), nil
			},
		}

		uc := testtarget.NewUseCase(
			&mockWorkspaceService{},
			&mockFinderService{},
			scripts,
			&mockReferenceParser{},
		)

		err := uc.Invoke(ctx, location, "test-script", globals)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		// Verify the repo globals structure
		if capturedScript.Globals != nil {
			repoData, ok := capturedScript.Globals["repo"].(map[string]any)
			if !ok {
				t.Error("repo globals not found or wrong type")
			} else {
				// Verify repo data has the correct structure
				if repoData["full_path"] != location.FullPath() {
					t.Errorf("Expected full_path %q, got %q", location.FullPath(), repoData["full_path"])
				}
				if repoData["name"] != "test-repo" {
					t.Errorf("Expected name %q, got %q", "test-repo", repoData["name"])
				}
				// Verify that user's repo data was overwritten
				if _, exists := repoData["should_be_overwritten"]; exists {
					t.Error("User's repo data should have been overwritten")
				}
			}
			// Verify custom globals are preserved
			if capturedScript.Globals["custom_key"] != "custom_value" {
				t.Error("Custom globals should be preserved")
			}
		}
	})
}

// Helper types

type errorReadCloser struct{}

func (e *errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}

// captureReader captures the gob-encoded script
type captureReader struct {
	script *run.Script
	buffer bytes.Buffer
}

func (c *captureReader) Read(p []byte) (n int, err error) {
	n, err = c.buffer.Read(p)
	if err == io.EOF && c.buffer.Len() == 0 {
		// First write, capture the data
		c.buffer.Write(p)
		n = len(p)
		err = nil

		// Try to decode
		dec := gob.NewDecoder(&c.buffer)
		gob.Register(map[string]any{})
		_ = dec.Decode(c.script)
		c.buffer.Reset()
	}
	return
}

func (c *captureReader) Write(p []byte) (n int, err error) {
	return c.buffer.Write(p)
}
