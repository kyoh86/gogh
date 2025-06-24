package apply_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/extra/apply"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/extra_mock"
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
		opts      testtarget.Options
		setupMock func(*gomock.Controller) (
			*extra_mock.MockExtraService,
			*overlay_mock.MockOverlayService,
			*workspace_mock.MockWorkspaceService,
			*workspace_mock.MockFinderService,
			*repository_mock.MockReferenceParser,
		)
		wantErr bool
	}{
		{
			name: "Empty name error",
			opts: testtarget.Options{
				Name:       "",
				TargetRepo: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
			) {
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				return es, os, ws, fs, rp
			},
			wantErr: true,
		},
		{
			name: "Named extra not found",
			opts: testtarget.Options{
				Name:       "nonexistent",
				TargetRepo: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
			) {
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				es.EXPECT().GetNamedExtra(ctx, "nonexistent").Return(nil, errors.New("not found"))

				return es, os, ws, fs, rp
			},
			wantErr: true,
		},
		{
			name: "Invalid repository reference",
			opts: testtarget.Options{
				Name:       "my-extra",
				TargetRepo: "invalid-ref",
			},
			setupMock: func(ctrl *gomock.Controller) (
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
			) {
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "source", "repo")
				namedExtra := extra.NewNamedExtra(
					uuid.New().String(),
					"my-extra",
					sourceRef,
					[]extra.Item{},
					time.Now(),
				)

				es.EXPECT().GetNamedExtra(ctx, "my-extra").Return(namedExtra, nil)
				rp.EXPECT().Parse("invalid-ref").Return(nil, errors.New("invalid reference"))

				return es, os, ws, fs, rp
			},
			wantErr: true,
		},
		{
			name: "Repository not found",
			opts: testtarget.Options{
				Name:       "my-extra",
				TargetRepo: "github.com/owner/notfound",
			},
			setupMock: func(ctrl *gomock.Controller) (
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
			) {
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "owner", "notfound")
				namedExtra := extra.NewNamedExtra(
					uuid.New().String(),
					"my-extra",
					sourceRef,
					[]extra.Item{},
					time.Now(),
				)

				es.EXPECT().GetNamedExtra(ctx, "my-extra").Return(namedExtra, nil)
				rp.EXPECT().Parse("github.com/owner/notfound").Return(&targetRef, nil)
				fs.EXPECT().FindByReference(ctx, ws, targetRef).Return(nil, errors.New("not found"))

				return es, os, ws, fs, rp
			},
			wantErr: true,
		},
		{
			name: "Current directory not in repository",
			opts: testtarget.Options{
				Name:       "my-extra",
				TargetRepo: "",
			},
			setupMock: func(ctrl *gomock.Controller) (
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
			) {
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "source", "repo")
				namedExtra := extra.NewNamedExtra(
					uuid.New().String(),
					"my-extra",
					sourceRef,
					[]extra.Item{},
					time.Now(),
				)

				es.EXPECT().GetNamedExtra(ctx, "my-extra").Return(namedExtra, nil)
				fs.EXPECT().FindByPath(ctx, ws, ".").Return(nil, errors.New("not in repository"))

				return es, os, ws, fs, rp
			},
			wantErr: true,
		},
		{
			name: "Overlay not found",
			opts: testtarget.Options{
				Name:       "my-extra",
				TargetRepo: "github.com/owner/repo",
			},
			setupMock: func(ctrl *gomock.Controller) (
				*extra_mock.MockExtraService,
				*overlay_mock.MockOverlayService,
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
			) {
				es := extra_mock.NewMockExtraService(ctrl)
				os := overlay_mock.NewMockOverlayService(ctrl)
				ws := workspace_mock.NewMockWorkspaceService(ctrl)
				fs := workspace_mock.NewMockFinderService(ctrl)
				rp := repository_mock.NewMockReferenceParser(ctrl)

				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := repository.NewReference("github.com", "owner", "repo")
				location := repository.NewLocation(
					"/root/github.com/owner/repo",
					"github.com",
					"owner",
					"repo",
				)

				overlayID := uuid.New().String()
				items := []extra.Item{{OverlayID: overlayID, HookID: ""}}

				namedExtra := extra.NewNamedExtra(
					uuid.New().String(),
					"my-extra",
					sourceRef,
					items,
					time.Now(),
				)

				es.EXPECT().GetNamedExtra(ctx, "my-extra").Return(namedExtra, nil)
				rp.EXPECT().Parse("github.com/owner/repo").Return(&targetRef, nil)
				fs.EXPECT().FindByReference(ctx, ws, targetRef).Return(location, nil)
				os.EXPECT().Get(ctx, overlayID).Return(nil, errors.New("overlay not found"))

				return es, os, ws, fs, rp
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			es, os, ws, fs, rp := tc.setupMock(ctrl)
			uc := testtarget.NewUseCase(es, os, ws, fs, rp)

			err := uc.Execute(ctx, tc.opts)
			if (err != nil) != tc.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestUseCase_Execute_FileSystemOperations(t *testing.T) {
	// Note: This test is skipped because it requires actual file system operations
	// which are not mockable in the current implementation.
	// The UseCase uses os.MkdirAll, os.OpenFile, and io.Copy directly.
	t.Skip("File system operations are not mockable in current implementation")
}
