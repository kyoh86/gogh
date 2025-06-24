package list_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/script/list"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name       string
		asJSON     bool
		withSource bool
		setupMock  func(*script_mock.MockScriptService)
		wantErr    bool
		validate   func(*testing.T, string)
	}{
		{
			name:       "List scripts as one-line",
			asJSON:     false,
			withSource: false,
			setupMock: func(m *script_mock.MockScriptService) {
				scripts := []script.Script{
					script.ConcreteScript(
						uuid.New(),
						"test-script-1",
						time.Now().Add(-24*time.Hour),
						time.Now(),
					),
					script.ConcreteScript(
						uuid.New(),
						"test-script-2",
						time.Now().Add(-48*time.Hour),
						time.Now().Add(-12*time.Hour),
					),
				}
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					for _, s := range scripts {
						if !yield(s, nil) {
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
			name:       "List scripts as JSON",
			asJSON:     true,
			withSource: false,
			setupMock: func(m *script_mock.MockScriptService) {
				scripts := []script.Script{
					script.ConcreteScript(
						uuid.New(),
						"json-script",
						time.Now(),
						time.Now(),
					),
				}
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					for _, s := range scripts {
						if !yield(s, nil) {
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
			name:       "List scripts with source as detail",
			asJSON:     false,
			withSource: true,
			setupMock: func(m *script_mock.MockScriptService) {
				scripts := []script.Script{
					script.ConcreteScript(
						uuid.New(),
						"detail-script",
						time.Now(),
						time.Now(),
					),
				}
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					for _, s := range scripts {
						if !yield(s, nil) {
							return
						}
					}
				})
				// Mock for detail view that needs to open script content
				m.EXPECT().Open(gomock.Any(), gomock.Any()).Return(
					io.NopCloser(bytes.NewReader([]byte("print('test script')"))), nil,
				).AnyTimes()
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if output == "" {
					t.Error("Expected non-empty detail output")
				}
			},
		},
		{
			name:       "List scripts as JSON with source",
			asJSON:     true,
			withSource: true,
			setupMock: func(m *script_mock.MockScriptService) {
				scripts := []script.Script{
					script.ConcreteScript(
						uuid.New(),
						"json-detail-script",
						time.Now(),
						time.Now(),
					),
				}
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					for _, s := range scripts {
						if !yield(s, nil) {
							return
						}
					}
				})
				// Mock for JSON with source view that needs to open script content
				m.EXPECT().Open(gomock.Any(), gomock.Any()).Return(
					io.NopCloser(bytes.NewReader([]byte("print('json test script')"))), nil,
				).AnyTimes()
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if output == "" {
					t.Error("Expected non-empty JSON with source output")
				}
			},
		},
		{
			name:       "Skip nil scripts",
			asJSON:     false,
			withSource: false,
			setupMock: func(m *script_mock.MockScriptService) {
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					yield(nil, nil) // nil script should be skipped
					yield(script.ConcreteScript(
						uuid.New(),
						"valid-script",
						time.Now(),
						time.Now(),
					), nil)
				})
			},
			wantErr: false,
		},
		{
			name:       "Error from List",
			asJSON:     false,
			withSource: false,
			setupMock: func(m *script_mock.MockScriptService) {
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					yield(nil, errors.New("list error"))
				})
			},
			wantErr: true,
		},
		{
			name:       "Error from Execute with source",
			asJSON:     false,
			withSource: true,
			setupMock: func(m *script_mock.MockScriptService) {
				scripts := []script.Script{
					script.ConcreteScript(
						uuid.New(),
						"error-script",
						time.Now(),
						time.Now(),
					),
				}
				m.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
					for _, s := range scripts {
						if !yield(s, nil) {
							return
						}
					}
				})
				// Mock Open to return error for detail view
				m.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil, errors.New("open script error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockScriptService := script_mock.NewMockScriptService(ctrl)
			tc.setupMock(mockScriptService)

			var buf bytes.Buffer
			uc := testtarget.NewUseCase(mockScriptService, &buf)

			err := uc.Execute(ctx, tc.asJSON, tc.withSource)
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
