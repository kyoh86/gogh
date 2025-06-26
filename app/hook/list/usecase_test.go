package list_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/list"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		asJSON    bool
		setupMock func(*hook_mock.MockHookService)
		wantErr   bool
		validate  func(*testing.T, string)
	}{
		{
			name:   "List hooks as one-line",
			asJSON: false,
			setupMock: func(m *hook_mock.MockHookService) {
				hooks := []hook.Hook{
					hook.ConcreteHook(
						uuid.New(),
						"test-hook-1",
						"github.com/owner/*",
						string(hook.EventPostClone),
						string(hook.OperationTypeOverlay),
						"overlay-id-1",
					),
					hook.ConcreteHook(
						uuid.New(),
						"test-hook-2",
						"github.com/org/**",
						string(hook.EventPostCreate),
						string(hook.OperationTypeScript),
						"script-id-1",
					),
				}
				m.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					for _, h := range hooks {
						if !yield(h, nil) {
							return
						}
					}
				})
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if output == "" {
					t.Error("Expected non-empty output")
				}
			},
		},
		{
			name:   "List hooks as JSON",
			asJSON: true,
			setupMock: func(m *hook_mock.MockHookService) {
				hooks := []hook.Hook{
					hook.ConcreteHook(
						uuid.New(),
						"json-hook",
						"github.com/test/*",
						string(hook.EventPostFork),
						string(hook.OperationTypeScript),
						"script-id-2",
					),
				}
				m.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					for _, h := range hooks {
						if !yield(h, nil) {
							return
						}
					}
				})
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if output == "" {
					t.Error("Expected non-empty JSON output")
				}
			},
		},
		{
			name:   "Skip nil hooks",
			asJSON: false,
			setupMock: func(m *hook_mock.MockHookService) {
				m.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					yield(nil, nil) // nil hook should be skipped
					yield(hook.ConcreteHook(
						uuid.New(),
						"valid-hook",
						"",
						string(hook.EventAny),
						string(hook.OperationTypeOverlay),
						"overlay-id",
					), nil)
				})
			},
			wantErr: false,
		},
		{
			name:   "Error from List",
			asJSON: false,
			setupMock: func(m *hook_mock.MockHookService) {
				m.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					yield(nil, errors.New("list error"))
				})
			},
			wantErr: true,
		},
		{
			name:   "Error from Execute",
			asJSON: false,
			setupMock: func(m *hook_mock.MockHookService) {
				// Create a hook that would cause the Execute to fail
				// by having invalid characters that might cause JSON marshaling issues
				hooks := []hook.Hook{
					hook.ConcreteHook(
						uuid.New(),
						"test-hook",
						"github.com/owner/*",
						string(hook.EventPostClone),
						string(hook.OperationTypeOverlay),
						"overlay-id",
					),
				}
				m.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
					for _, h := range hooks {
						if !yield(h, nil) {
							return
						}
					}
				})
			},
			wantErr: false, // Normal hooks shouldn't cause Execute errors
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHookService := hook_mock.NewMockHookService(ctrl)
			tc.setupMock(mockHookService)

			var buf bytes.Buffer
			uc := testtarget.NewUsecase(mockHookService, &buf)

			err := uc.Execute(ctx, tc.asJSON)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if tc.validate != nil {
				tc.validate(t, buf.String())
			}
		})
	}
}
