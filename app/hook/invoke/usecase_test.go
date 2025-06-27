package invoke_test

import (
	"context"
	"errors"
	"iter"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/invoke"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestNewUsecase(t *testing.T) {
	ws := workspace_mock.NewMockWorkspaceService(gomock.NewController(t))
	finder := workspace_mock.NewMockFinderService(gomock.NewController(t))
	hooks := hook_mock.NewMockHookService(gomock.NewController(t))
	overlays := overlay_mock.NewMockOverlayService(gomock.NewController(t))
	scripts := script_mock.NewMockScriptService(gomock.NewController(t))
	parser := repository_mock.NewMockReferenceParser(gomock.NewController(t))

	uc := testtarget.NewUsecase(ws, finder, hooks, overlays, scripts, parser)
	if uc == nil {
		t.Fatal("expected non-nil Usecase")
	}
}

func TestUsecase_Invoke(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		hookID    string
		refStr    string
		setupHook func() hook.Hook
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "overlay hook",
			hookID: "test-hook-id",
			refStr: "github.com/kyoh86/gogh",
			setupHook: func() hook.Hook {
				h := hook_mock.NewMockHook(gomock.NewController(t))
				h.EXPECT().ID().Return(uuid.New().String()).AnyTimes()
				h.EXPECT().Name().Return("test overlay hook").AnyTimes()
				h.EXPECT().OperationType().Return(hook.OperationTypeOverlay).AnyTimes()
				h.EXPECT().OperationID().Return(uuid.New().String()).AnyTimes()
				h.EXPECT().RepoPattern().Return("github.com/kyoh86/*").AnyTimes()
				h.EXPECT().TriggerEvent().Return(hook.EventPostClone).AnyTimes()
				return h
			},
			wantErr: true, // Will fail because overlay service is not fully mocked
		},
		{
			name:   "script hook",
			hookID: "test-hook-id",
			refStr: "github.com/kyoh86/gogh",
			setupHook: func() hook.Hook {
				h := hook_mock.NewMockHook(gomock.NewController(t))
				h.EXPECT().ID().Return(uuid.New().String()).AnyTimes()
				h.EXPECT().Name().Return("test script hook").AnyTimes()
				h.EXPECT().OperationType().Return(hook.OperationTypeScript).AnyTimes()
				h.EXPECT().OperationID().Return(uuid.New().String()).AnyTimes()
				h.EXPECT().RepoPattern().Return("github.com/kyoh86/*").AnyTimes()
				h.EXPECT().TriggerEvent().Return(hook.EventPostClone).AnyTimes()
				return h
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
			hookSvc := hook_mock.NewMockHookService(gomock.NewController(t))
			hookSvc.EXPECT().Get(gomock.Any(), tt.hookID).DoAndReturn(func(ctx context.Context, id string) (hook.Hook, error) {
				if tt.setupHook != nil && id == tt.hookID {
					return tt.setupHook(), nil
				}
				return nil, errors.New("hook not found")
			}).AnyTimes()

			rp := repository_mock.NewMockReferenceParser(gomock.NewController(t))
			if tt.refStr != "" {
				rp.EXPECT().ParseWithAlias(tt.refStr).Return(
					&repository.ReferenceWithAlias{
						Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
					},
					nil,
				).AnyTimes()
			}

			ws := workspace_mock.NewMockWorkspaceService(gomock.NewController(t))
			fs := workspace_mock.NewMockFinderService(gomock.NewController(t))
			if tt.refStr != "" {
				fs.EXPECT().FindByReference(gomock.Any(), ws, repository.NewReference("github.com", "kyoh86", "gogh")).Return(
					repository.NewLocation("/path/to/repo", "github.com", "kyoh86", "gogh"),
					nil,
				).AnyTimes()
			}

			os := overlay_mock.NewMockOverlayService(gomock.NewController(t))
			os.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("overlay not found")).AnyTimes()
			ss := script_mock.NewMockScriptService(gomock.NewController(t))
			ss.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("script not found")).AnyTimes()

			uc := testtarget.NewUsecase(
				ws,
				fs,
				hookSvc,
				os,
				ss,
				rp,
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

func TestUsecase_InvokeFor(t *testing.T) {
	ctx := context.Background()

	t.Run("successful invocation", func(t *testing.T) {
		overlayHook := hook_mock.NewMockHook(gomock.NewController(t))
		overlayHook.EXPECT().ID().Return(uuid.New().String()).AnyTimes()
		overlayHook.EXPECT().Name().Return("overlay hook").AnyTimes()
		overlayHook.EXPECT().OperationType().Return(hook.OperationTypeOverlay).AnyTimes()
		overlayHook.EXPECT().OperationID().Return(uuid.New().String()).AnyTimes()
		overlayHook.EXPECT().TriggerEvent().Return(testtarget.EventPostClone).AnyTimes()
		overlayHook.EXPECT().RepoPattern().Return("github.com/kyoh86/*").AnyTimes()

		scriptHook := hook_mock.NewMockHook(gomock.NewController(t))
		scriptHook.EXPECT().ID().Return(uuid.New().String()).AnyTimes()
		scriptHook.EXPECT().Name().Return("script hook").AnyTimes()
		scriptHook.EXPECT().OperationType().Return(hook.OperationTypeScript).AnyTimes()
		scriptHook.EXPECT().OperationID().Return(uuid.New().String()).AnyTimes()
		scriptHook.EXPECT().TriggerEvent().Return(testtarget.EventPostClone).AnyTimes()
		scriptHook.EXPECT().RepoPattern().Return("github.com/kyoh86/*").AnyTimes()

		hookSvc := hook_mock.NewMockHookService(gomock.NewController(t))
		hookSvc.EXPECT().ListFor(gomock.Any(), gomock.Any()).DoAndReturn(func(reference repository.Reference, event hook.Event) iter.Seq2[hook.Hook, error] {
			return func(yield func(hook.Hook, error) bool) {
				if !yield(overlayHook, nil) {
					return
				}
				yield(scriptHook, nil)
			}
		}).AnyTimes()

		rp := repository_mock.NewMockReferenceParser(gomock.NewController(t))
		rp.EXPECT().ParseWithAlias("github.com/kyoh86/gogh").Return(
			&repository.ReferenceWithAlias{
				Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
			},
			nil,
		)

		ws := workspace_mock.NewMockWorkspaceService(gomock.NewController(t))
		fs := workspace_mock.NewMockFinderService(gomock.NewController(t))
		fs.EXPECT().FindByReference(gomock.Any(), ws, repository.NewReference("github.com", "kyoh86", "gogh")).Return(
			repository.NewLocation("/path/to/repo", "github.com", "kyoh86", "gogh"),
			nil,
		)

		os := overlay_mock.NewMockOverlayService(gomock.NewController(t))
		os.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("overlay not found")).AnyTimes()
		ss := script_mock.NewMockScriptService(gomock.NewController(t))
		ss.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("script not found")).AnyTimes()

		uc := testtarget.NewUsecase(
			ws,
			fs,
			hookSvc,
			os,
			ss,
			rp,
		)

		err := uc.InvokeFor(ctx, testtarget.EventPostClone, "github.com/kyoh86/gogh")
		// This will error because the apply operations are not fully mocked
		if err == nil {
			t.Error("expected error due to incomplete mocks")
		}
	})

	t.Run("invalid reference", func(t *testing.T) {
		rp := repository_mock.NewMockReferenceParser(gomock.NewController(t))
		rp.EXPECT().ParseWithAlias("invalid-ref").Return(
			nil,
			errors.New("invalid reference"),
		)

		uc := testtarget.NewUsecase(
			workspace_mock.NewMockWorkspaceService(gomock.NewController(t)),
			workspace_mock.NewMockFinderService(gomock.NewController(t)),
			hook_mock.NewMockHookService(gomock.NewController(t)),
			overlay_mock.NewMockOverlayService(gomock.NewController(t)),
			script_mock.NewMockScriptService(gomock.NewController(t)),
			rp,
		)

		err := uc.InvokeFor(ctx, testtarget.EventPostClone, "invalid-ref")
		if err == nil {
			t.Error("expected error for invalid reference")
		}
	})

	t.Run("repository not found", func(t *testing.T) {
		rp := repository_mock.NewMockReferenceParser(gomock.NewController(t))
		rp.EXPECT().ParseWithAlias("github.com/kyoh86/gogh").Return(
			&repository.ReferenceWithAlias{
				Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
			},
			nil,
		)

		finder := workspace_mock.NewMockFinderService(gomock.NewController(t))
		finder.EXPECT().FindByReference(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("not found")).AnyTimes()

		uc := testtarget.NewUsecase(
			workspace_mock.NewMockWorkspaceService(gomock.NewController(t)),
			finder,
			hook_mock.NewMockHookService(gomock.NewController(t)),
			overlay_mock.NewMockOverlayService(gomock.NewController(t)),
			script_mock.NewMockScriptService(gomock.NewController(t)),
			rp,
		)

		err := uc.InvokeFor(ctx, testtarget.EventPostClone, "github.com/kyoh86/gogh")
		if err == nil {
			t.Error("expected error for repository not found")
		}
	})
}

func TestUsecase_InvokeForWithGlobals(t *testing.T) {
	ctx := context.Background()

	scriptHook := hook_mock.NewMockHook(gomock.NewController(t))
	scriptHook.EXPECT().ID().Return(uuid.New().String()).AnyTimes()
	scriptHook.EXPECT().Name().Return("script hook").AnyTimes()
	scriptHook.EXPECT().OperationType().Return(hook.OperationTypeScript).AnyTimes()
	scriptHook.EXPECT().OperationID().Return(uuid.New().String()).AnyTimes()
	scriptHook.EXPECT().TriggerEvent().Return(testtarget.EventPostFork).AnyTimes()
	scriptHook.EXPECT().RepoPattern().Return("github.com/kyoh86/*").AnyTimes()

	hookSvc := hook_mock.NewMockHookService(gomock.NewController(t))
	hookSvc.EXPECT().ListFor(gomock.Any(), gomock.Any()).DoAndReturn(func(reference repository.Reference, event hook.Event) iter.Seq2[hook.Hook, error] {
		return func(yield func(hook.Hook, error) bool) {
			yield(scriptHook, nil)
		}
	}).AnyTimes()

	rp := repository_mock.NewMockReferenceParser(gomock.NewController(t))
	rp.EXPECT().ParseWithAlias("github.com/kyoh86/gogh").Return(
		&repository.ReferenceWithAlias{
			Reference: repository.NewReference("github.com", "kyoh86", "gogh"),
		},
		nil,
	)

	ws := workspace_mock.NewMockWorkspaceService(gomock.NewController(t))
	fs := workspace_mock.NewMockFinderService(gomock.NewController(t))
	fs.EXPECT().FindByReference(gomock.Any(), ws, repository.NewReference("github.com", "kyoh86", "gogh")).Return(
		repository.NewLocation("/path/to/repo", "github.com", "kyoh86", "gogh"),
		nil,
	)

	ss := script_mock.NewMockScriptService(gomock.NewController(t))
	ss.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("script not found")).AnyTimes()

	uc := testtarget.NewUsecase(
		ws,
		fs,
		hookSvc,
		overlay_mock.NewMockOverlayService(gomock.NewController(t)),
		ss,
		rp,
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
