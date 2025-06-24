package extra_show_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/extra_show"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	createdAt := time.Now()

	testCases := []struct {
		name       string
		identifier string
		asJSON     bool
		setupMock  func(*gomock.Controller) *extra_mock.MockExtraService
		wantErr    bool
		validate   func(*testing.T, string)
	}{
		{
			name:       "Show auto extra by ID as detail",
			identifier: uuid.New().String(),
			asJSON:     false,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				ref := repository.NewReference("github.com", "owner", "repo")
				sourceRef := repository.NewReference("github.com", "source", "extra-repo")
				items := []extra.Item{
					{OverlayID: "overlay-1", HookID: ""},
					{OverlayID: "", HookID: "hook-1"},
				}
				e := extra.NewAutoExtra(
					uuid.New().String(),
					ref,
					sourceRef,
					items,
					createdAt,
				)
				es.EXPECT().Get(ctx, gomock.Any()).Return(e, nil)
				return es
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should contain extra details
				if !strings.Contains(output, "Type: auto") {
					t.Error("Expected output to contain 'Type: auto'")
				}
				if !strings.Contains(output, "github.com/owner/repo") {
					t.Error("Expected output to contain repository reference")
				}
				if !strings.Contains(output, "Items (") {
					t.Error("Expected output to contain 'Items ('")
				}
			},
		},
		{
			name:       "Show named extra by ID as JSON",
			identifier: uuid.New().String(),
			asJSON:     true,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				ref := repository.NewReference("github.com", "source", "repo")
				items := []extra.Item{
					{OverlayID: "overlay-1", HookID: ""},
				}
				e := extra.NewNamedExtra(
					uuid.New().String(),
					"my-extra",
					ref,
					items,
					createdAt,
				)
				es.EXPECT().Get(ctx, gomock.Any()).Return(e, nil)
				return es
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should be valid JSON
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check type
				if data["type"] != "named" {
					t.Errorf("Expected type 'named', got %v", data["type"])
				}
				// Check name exists
				if _, ok := data["name"]; !ok {
					t.Error("Expected 'name' field in JSON output")
				}
			},
		},
		{
			name:       "Show extra by name (fallback from ID)",
			identifier: "my-extra-name",
			asJSON:     false,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				ref := repository.NewReference("github.com", "source", "repo")
				items := []extra.Item{
					{OverlayID: "overlay-1", HookID: ""},
				}
				e := extra.NewNamedExtra(
					uuid.New().String(),
					"my-extra-name",
					ref,
					items,
					createdAt,
				)
				// First try as ID fails
				es.EXPECT().Get(ctx, "my-extra-name").Return(nil, errors.New("not found"))
				// Then try as name succeeds
				es.EXPECT().GetNamedExtra(ctx, "my-extra-name").Return(e, nil)
				return es
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "Type: named") {
					t.Error("Expected output to contain 'Type: named'")
				}
				if !strings.Contains(output, "Name: my-extra-name") {
					t.Error("Expected output to contain name")
				}
			},
		},
		{
			name:       "Extra not found by ID or name",
			identifier: "non-existent",
			asJSON:     false,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				// Both attempts fail
				es.EXPECT().Get(ctx, "non-existent").Return(nil, errors.New("not found"))
				es.EXPECT().GetNamedExtra(ctx, "non-existent").Return(nil, errors.New("not found"))
				return es
			},
			wantErr: true,
		},
		{
			name:       "Show auto extra with empty items as JSON",
			identifier: uuid.New().String(),
			asJSON:     true,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				ref := repository.NewReference("github.com", "owner", "repo")
				sourceRef := repository.NewReference("github.com", "source", "extra-repo")
				e := extra.NewAutoExtra(
					uuid.New().String(),
					ref,
					sourceRef,
					[]extra.Item{}, // Empty items
					createdAt,
				)
				es.EXPECT().Get(ctx, gomock.Any()).Return(e, nil)
				return es
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check items is empty array
				if items, ok := data["items"].([]interface{}); ok {
					if len(items) != 0 {
						t.Errorf("Expected empty items array, got %d items", len(items))
					}
				} else {
					t.Error("Expected 'items' field to be an array")
				}
			},
		},
		{
			name:       "Show extra with multiple items as detail",
			identifier: uuid.New().String(),
			asJSON:     false,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				ref := repository.NewReference("github.com", "owner", "repo")
				sourceRef := repository.NewReference("github.com", "source", "extra-repo")
				items := []extra.Item{
					{OverlayID: "overlay-1", HookID: ""},
					{OverlayID: "overlay-2", HookID: ""},
					{OverlayID: "", HookID: "hook-1"},
					{OverlayID: "", HookID: "hook-2"},
					{OverlayID: "overlay-3", HookID: "hook-3"},
				}
				e := extra.NewAutoExtra(
					uuid.New().String(),
					ref,
					sourceRef,
					items,
					createdAt,
				)
				es.EXPECT().Get(ctx, gomock.Any()).Return(e, nil)
				return es
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				// Should list all items
				if !strings.Contains(output, "overlay-1") {
					t.Error("Expected output to contain overlay-1")
				}
				if !strings.Contains(output, "hook-2") {
					t.Error("Expected output to contain hook-2")
				}
			},
		},
		{
			name:       "Show named extra from different source as JSON",
			identifier: "cross-repo-extra",
			asJSON:     true,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				ref := repository.NewReference("gitlab.com", "different", "source")
				items := []extra.Item{
					{OverlayID: "overlay-x", HookID: ""},
				}
				e := extra.NewNamedExtra(
					uuid.New().String(),
					"cross-repo-extra",
					ref,
					items,
					createdAt,
				)
				es.EXPECT().Get(ctx, "cross-repo-extra").Return(nil, errors.New("not found"))
				es.EXPECT().GetNamedExtra(ctx, "cross-repo-extra").Return(e, nil)
				return es
			},
			wantErr: false,
			validate: func(t *testing.T, output string) {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(output), &data); err != nil {
					t.Errorf("Expected valid JSON output, got error: %v", err)
				}
				// Check source repository (it's a string, not a map)
				if source, ok := data["source"].(string); ok {
					if !strings.Contains(source, "gitlab.com") {
						t.Errorf("Expected source to contain 'gitlab.com', got %v", source)
					}
				} else {
					t.Error("Expected 'source' field in JSON output")
				}
			},
		},
		{
			name:       "Invalid UUID as identifier",
			identifier: "not-a-uuid",
			asJSON:     false,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				// First try as ID fails (not a valid UUID)
				es.EXPECT().Get(ctx, "not-a-uuid").Return(nil, errors.New("invalid ID"))
				// Try as name also fails
				es.EXPECT().GetNamedExtra(ctx, "not-a-uuid").Return(nil, errors.New("not found"))
				return es
			},
			wantErr: true,
		},
		{
			name:       "Empty identifier",
			identifier: "",
			asJSON:     false,
			setupMock: func(ctrl *gomock.Controller) *extra_mock.MockExtraService {
				es := extra_mock.NewMockExtraService(ctrl)
				es.EXPECT().Get(ctx, "").Return(nil, errors.New("empty ID"))
				es.EXPECT().GetNamedExtra(ctx, "").Return(nil, errors.New("empty name"))
				return es
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var buf bytes.Buffer
			es := tc.setupMock(ctrl)
			uc := extra_show.NewUseCase(es, &buf)

			err := uc.Execute(ctx, tc.identifier, tc.asJSON)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && tc.validate != nil {
				tc.validate(t, buf.String())
			}
		})
	}
}

func TestUseCase_Execute_ServiceErrors(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var buf bytes.Buffer
	es := extra_mock.NewMockExtraService(ctrl)
	uc := extra_show.NewUseCase(es, &buf)

	// Test service returning unexpected error
	es.EXPECT().Get(ctx, "test-id").Return(nil, errors.New("database connection error"))
	es.EXPECT().GetNamedExtra(ctx, "test-id").Return(nil, errors.New("database connection error"))

	err := uc.Execute(ctx, "test-id", false)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "extra not found") {
		t.Errorf("Expected error to contain 'extra not found', got %v", err.Error())
	}
}
