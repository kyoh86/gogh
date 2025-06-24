package apply_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay/apply"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

type readCloserMock struct {
	*bytes.Reader
	closed bool
}

func (r *readCloserMock) Close() error {
	r.closed = true
	return nil
}

func TestUseCase_Execute(t *testing.T) {
	// Setup temp directory for file operations
	tempDir, err := os.MkdirTemp("", "overlay_apply_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repoPath := filepath.Join(tempDir, "github.com", "kyoh86", "example")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("failed to create test repo dir: %v", err)
	}

	tests := []struct {
		name         string
		refs         string
		id           string
		relativePath string
		mockSetup    func(*gomock.Controller) (
			*workspace_mock.MockWorkspaceService,
			*workspace_mock.MockFinderService,
			*repository_mock.MockReferenceParser,
			*overlay_mock.MockOverlayService,
			io.ReadCloser,
		)
		wantErr bool
	}{
		{
			name:         "Success: Apply overlay to repository",
			refs:         "kyoh86/example",
			id:           "overlay1",
			relativePath: "config/settings.yaml",
			mockSetup: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
				*overlay_mock.MockOverlayService,
				io.ReadCloser,
			) {
				workspaceSvc := workspace_mock.NewMockWorkspaceService(ctrl)
				finderSvc := workspace_mock.NewMockFinderService(ctrl)
				refParser := repository_mock.NewMockReferenceParser(ctrl)
				overlaySvc := overlay_mock.NewMockOverlayService(ctrl)

				ref := repository.NewReference("github.com", "kyoh86", "example")
				refParser.EXPECT().ParseWithAlias("kyoh86/example").Return(&repository.ReferenceWithAlias{Reference: ref}, nil)

				repo := repository.NewLocation(
					repoPath,
					"github.com",
					"kyoh86",
					"example",
				)
				finderSvc.EXPECT().FindByReference(gomock.Any(), workspaceSvc, ref).Return(repo, nil)

				content := &readCloserMock{
					Reader: bytes.NewReader([]byte("overlay content")),
				}
				ov := overlay.NewOverlay(overlay.Entry{
					Name:         "overlay-name 1",
					RelativePath: "overlay/path/1",
				})
				overlaySvc.EXPECT().Get(gomock.Any(), "overlay1").Return(ov, nil)
				overlaySvc.EXPECT().Open(gomock.Any(), "overlay1").Return(content, nil)

				return workspaceSvc, finderSvc, refParser, overlaySvc, content
			},
			wantErr: false,
		},
		{
			name: "Error: Failed to parse reference",
			refs: "invalid/ref/format",
			mockSetup: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
				*overlay_mock.MockOverlayService,
				io.ReadCloser,
			) {
				workspaceSvc := workspace_mock.NewMockWorkspaceService(ctrl)
				finderSvc := workspace_mock.NewMockFinderService(ctrl)
				refParser := repository_mock.NewMockReferenceParser(ctrl)
				overlaySvc := overlay_mock.NewMockOverlayService(ctrl)

				refParser.EXPECT().ParseWithAlias("invalid/ref/format").Return(nil, errors.New("invalid reference format"))

				return workspaceSvc, finderSvc, refParser, overlaySvc, nil
			},
			wantErr: true,
		},
		{
			name: "Error: Repository not found",
			refs: "nonexistent/repo",
			mockSetup: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
				*overlay_mock.MockOverlayService,
				io.ReadCloser,
			) {
				workspaceSvc := workspace_mock.NewMockWorkspaceService(ctrl)
				finderSvc := workspace_mock.NewMockFinderService(ctrl)
				refParser := repository_mock.NewMockReferenceParser(ctrl)
				overlaySvc := overlay_mock.NewMockOverlayService(ctrl)

				ref := repository.NewReference("github.com", "nonexistent", "repo")
				refParser.EXPECT().ParseWithAlias("nonexistent/repo").Return(&repository.ReferenceWithAlias{Reference: ref}, nil)
				finderSvc.EXPECT().FindByReference(gomock.Any(), workspaceSvc, ref).Return(nil, nil)

				return workspaceSvc, finderSvc, refParser, overlaySvc, nil
			},
			wantErr: true,
		},
		{
			name:         "Error: Failed to open overlay content",
			refs:         "kyoh86/example",
			id:           "overlay4",
			relativePath: "config.yaml",
			mockSetup: func(ctrl *gomock.Controller) (
				*workspace_mock.MockWorkspaceService,
				*workspace_mock.MockFinderService,
				*repository_mock.MockReferenceParser,
				*overlay_mock.MockOverlayService,
				io.ReadCloser,
			) {
				workspaceSvc := workspace_mock.NewMockWorkspaceService(ctrl)
				finderSvc := workspace_mock.NewMockFinderService(ctrl)
				refParser := repository_mock.NewMockReferenceParser(ctrl)
				overlaySvc := overlay_mock.NewMockOverlayService(ctrl)

				ref := repository.NewReference("github.com", "kyoh86", "example")
				refParser.EXPECT().ParseWithAlias("kyoh86/example").Return(&repository.ReferenceWithAlias{Reference: ref}, nil)

				repo := repository.NewLocation(
					repoPath,
					"github.com",
					"kyoh86",
					"example",
				)
				finderSvc.EXPECT().FindByReference(gomock.Any(), workspaceSvc, ref).Return(repo, nil)

				ov := overlay.NewOverlay(overlay.Entry{
					Name:         "overlay-name 1",
					RelativePath: "overlay/path/1",
				})
				overlaySvc.EXPECT().Get(gomock.Any(), "overlay4").Return(ov, nil)
				overlaySvc.EXPECT().Open(gomock.Any(), "overlay4").Return(nil, errors.New("overlay not found"))

				return workspaceSvc, finderSvc, refParser, overlaySvc, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			workspaceSvc, finderSvc, refParser, overlaySvc, content := tt.mockSetup(ctrl)

			uc := testtarget.NewUseCase(workspaceSvc, finderSvc, refParser, overlaySvc)
			err := uc.Execute(context.Background(), tt.refs, tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("UseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify that ReadCloser was properly closed (if there was no error)
			if err == nil && content != nil {
				if mock, ok := content.(*readCloserMock); ok && !mock.closed {
					t.Error("ReadCloser was not closed")
				}
			}
		})
	}
}
