package list_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	testtarget "github.com/kyoh86/gogh/v4/app/extra/list"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		asJSON    bool
		extraType string
		setupMock func(*extra_mock.MockExtraService)
		wantErr   bool
		validate  func(*testing.T, string)
	}{
		{
			name:      "List all extras as one-line",
			asJSON:    false,
			extraType: "all",
			setupMock: func(m *extra_mock.MockExtraService) {
				sourceRef := repository.NewReference("github.com", "owner", "repo")
				autoRef := repository.NewReference("github.com", "owner", "auto-repo")
				extras := []*extra.Extra{
					extra.NewNamedExtra(
						"test-id-1",
						"test-extra-1",
						sourceRef,
						[]extra.Item{{OverlayID: "overlay-1", HookID: "hook-1"}},
						time.Now(),
					),
					extra.NewAutoExtra(
						"test-id-2",
						autoRef,
						sourceRef,
						[]extra.Item{{OverlayID: "overlay-2", HookID: "hook-2"}},
						time.Now(),
					),
				}
				m.EXPECT().List(ctx).Return(func(yield func(*extra.Extra, error) bool) {
					for _, e := range extras {
						if !yield(e, nil) {
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
			name:      "List auto extras as JSON",
			asJSON:    true,
			extraType: "auto",
			setupMock: func(m *extra_mock.MockExtraService) {
				sourceRef := repository.NewReference("github.com", "owner", "repo")
				autoRef := repository.NewReference("github.com", "owner", "auto-repo")
				extras := []*extra.Extra{
					extra.NewAutoExtra(
						"auto-id-1",
						autoRef,
						sourceRef,
						[]extra.Item{{OverlayID: "overlay-auto", HookID: "hook-auto"}},
						time.Now(),
					),
				}
				m.EXPECT().ListByType(ctx, extra.TypeAuto).Return(func(yield func(*extra.Extra, error) bool) {
					for _, e := range extras {
						if !yield(e, nil) {
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
			name:      "List named extras",
			asJSON:    false,
			extraType: "named",
			setupMock: func(m *extra_mock.MockExtraService) {
				sourceRef := repository.NewReference("github.com", "owner", "repo")
				extras := []*extra.Extra{
					extra.NewNamedExtra(
						"named-id-1",
						"my-template",
						sourceRef,
						[]extra.Item{{OverlayID: "overlay-named", HookID: "hook-named"}},
						time.Now(),
					),
				}
				m.EXPECT().ListByType(ctx, extra.TypeNamed).Return(func(yield func(*extra.Extra, error) bool) {
					for _, e := range extras {
						if !yield(e, nil) {
							return
						}
					}
				})
			},
			wantErr: false,
		},
		{
			name:      "Skip nil extras",
			asJSON:    false,
			extraType: "all",
			setupMock: func(m *extra_mock.MockExtraService) {
				sourceRef := repository.NewReference("github.com", "owner", "repo")
				m.EXPECT().List(ctx).Return(func(yield func(*extra.Extra, error) bool) {
					yield(nil, nil) // nil extra should be skipped
					yield(extra.NewNamedExtra(
						"valid-id",
						"valid-extra",
						sourceRef,
						[]extra.Item{{OverlayID: "overlay-valid", HookID: "hook-valid"}},
						time.Now(),
					), nil)
				})
			},
			wantErr: false,
		},
		{
			name:      "Error from List",
			asJSON:    false,
			extraType: "all",
			setupMock: func(m *extra_mock.MockExtraService) {
				m.EXPECT().List(ctx).Return(func(yield func(*extra.Extra, error) bool) {
					yield(nil, errors.New("list error"))
				})
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockExtraService := extra_mock.NewMockExtraService(ctrl)
			tc.setupMock(mockExtraService)

			var buf bytes.Buffer
			uc := testtarget.NewUsecase(mockExtraService, &buf)

			err := uc.Execute(ctx, tc.asJSON, tc.extraType)
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
