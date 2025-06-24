package extra_create_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/extra_create"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name      string
		opts      extra_create.Options
		setupMock func(*gomock.Controller) (
			*workspace_mock.MockWorkspaceService,
			*workspace_mock.MockFinderService,
			*extra_mock.MockExtraService,
			*overlay_mock.MockOverlayService,
			*repository_mock.MockReferenceParser,
		)
		wantErr bool
	}{
		{
			name: "Successfully create named extra",
			opts: extra_create.Options{
				Name:         "my-extra",
				SourceRepo:   "github.com/owner/repo",
				OverlayNames: []string{"overlay1", "overlay2"},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "owner", "repo")

				// Parse source repository
				rp.EXPECT().Parse("github.com/owner/repo").Return(&sourceRef, nil)

				// Get overlays
				overlay1UUID := uuid.New()
				overlay2UUID := uuid.New()
				os.EXPECT().Get(ctx, "overlay1").Return(
					overlay.ConcreteOverlay(overlay1UUID, "overlay1", "file1.txt"), nil,
				)
				os.EXPECT().Get(ctx, "overlay2").Return(
					overlay.ConcreteOverlay(overlay2UUID, "overlay2", "file2.txt"), nil,
				)

				// Create named extra
				es.EXPECT().AddNamedExtra(ctx, "my-extra", sourceRef, gomock.Any()).DoAndReturn(
					func(ctx context.Context, name string, source repository.Reference, items []extra.Item) (string, error) {
						if len(items) != 2 {
							t.Errorf("Expected 2 items, got %d", len(items))
						}
						if items[0].OverlayID != overlay1UUID.String() {
							t.Errorf("Expected first overlay ID %s, got %s", overlay1UUID.String(), items[0].OverlayID)
						}
						if items[1].OverlayID != overlay2UUID.String() {
							t.Errorf("Expected second overlay ID %s, got %s", overlay2UUID.String(), items[1].OverlayID)
						}
						// Check that hook IDs are empty for named extras
						if items[0].HookID != "" || items[1].HookID != "" {
							t.Error("Named extras should not have hook IDs")
						}
						return uuid.New().String(), nil
					},
				)

				return ws, fs, es, os, rp
			},
			wantErr: false,
		},
		{
			name: "Empty name error",
			opts: extra_create.Options{
				Name:         "",
				SourceRepo:   "github.com/owner/repo",
				OverlayNames: []string{"overlay1"},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				return ws, fs, es, os, rp
			},
			wantErr: true,
		},
		{
			name: "Empty source repository error",
			opts: extra_create.Options{
				Name:         "my-extra",
				SourceRepo:   "",
				OverlayNames: []string{"overlay1"},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				return ws, fs, es, os, rp
			},
			wantErr: true,
		},
		{
			name: "Invalid repository reference",
			opts: extra_create.Options{
				Name:         "my-extra",
				SourceRepo:   "invalid-ref",
				OverlayNames: []string{"overlay1"},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				rp.EXPECT().Parse("invalid-ref").Return(nil, errors.New("invalid reference"))

				return ws, fs, es, os, rp
			},
			wantErr: true,
		},
		{
			name: "Overlay not found",
			opts: extra_create.Options{
				Name:         "my-extra",
				SourceRepo:   "github.com/owner/repo",
				OverlayNames: []string{"nonexistent"},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&sourceRef, nil)
				os.EXPECT().Get(ctx, "nonexistent").Return(nil, errors.New("not found"))

				return ws, fs, es, os, rp
			},
			wantErr: true,
		},
		{
			name: "AddNamedExtra fails",
			opts: extra_create.Options{
				Name:         "my-extra",
				SourceRepo:   "github.com/owner/repo",
				OverlayNames: []string{"overlay1"},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "owner", "repo")
				overlay1UUID := uuid.New()

				rp.EXPECT().Parse("github.com/owner/repo").Return(&sourceRef, nil)
				os.EXPECT().Get(ctx, "overlay1").Return(
					overlay.ConcreteOverlay(overlay1UUID, "overlay1", "file1.txt"), nil,
				)
				es.EXPECT().AddNamedExtra(ctx, "my-extra", sourceRef, gomock.Any()).Return(
					"", errors.New("already exists"),
				)

				return ws, fs, es, os, rp
			},
			wantErr: true,
		},
		{
			name: "Create with no overlays",
			opts: extra_create.Options{
				Name:         "empty-extra",
				SourceRepo:   "github.com/owner/repo",
				OverlayNames: []string{},
			},
			setupMock: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*repository_mock.MockReferenceParser,
			) {
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "owner", "repo")

				rp.EXPECT().Parse("github.com/owner/repo").Return(&sourceRef, nil)
				es.EXPECT().AddNamedExtra(ctx, "empty-extra", sourceRef, gomock.Any()).DoAndReturn(
					func(ctx context.Context, name string, source repository.Reference, items []extra.Item) (string, error) {
						if len(items) != 0 {
							t.Errorf("Expected 0 items, got %d", len(items))
						}
						return uuid.New().String(), nil
					},
				)

				return ws, fs, es, os, rp
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ws, fs, es, os, rp := tc.setupMock(ctrl)
			uc := extra_create.NewUseCase(ws, fs, es, os, rp)

			err := uc.Execute(ctx, tc.opts)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
