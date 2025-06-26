package delete_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kyoh86/gogh/v4/app/delete"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	testCases := []struct {
		name          string
		refs          string
		options       delete.Options
		setupMocks    func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser)
		expectErr     bool
		expectErrText string
	}{
		{
			name: "delete local only",
			refs: "github.com/user/repo",
			options: delete.Options{
				Local:  true,
				Remote: false,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				ref := repository.NewReference("github.com", "user", "repo")
				match := repository.NewLocation("/path/to/repo", "github.com", "user", "repo")

				mrp.EXPECT().
					Parse("github.com/user/repo").
					Return(&ref, nil)

				mfs.EXPECT().
					FindByReference(gomock.Any(), mws, ref).
					Return(match, nil)

				// Not expecting a call to DeleteRepository since Remote is false
			},
			expectErr: false,
		},
		{
			name: "delete remote only",
			refs: "github.com/user/repo",
			options: delete.Options{
				Local:  false,
				Remote: true,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				ref := repository.NewReference("github.com", "user", "repo")

				mrp.EXPECT().
					Parse("github.com/user/repo").
					Return(&ref, nil)

				// Not expecting a call to FindByReference since Local is false

				mhs.EXPECT().
					DeleteRepository(gomock.Any(), ref).
					Return(nil)
			},
			expectErr: false,
		},
		{
			name: "delete both local and remote",
			refs: "github.com/user/repo",
			options: delete.Options{
				Local:  true,
				Remote: true,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				ref := repository.NewReference("github.com", "user", "repo")
				match := repository.NewLocation("/path/to/repo", "github.com", "user", "repo")

				mrp.EXPECT().
					Parse("github.com/user/repo").
					Return(&ref, nil)

				mfs.EXPECT().
					FindByReference(gomock.Any(), mws, ref).
					Return(match, nil)

				mhs.EXPECT().
					DeleteRepository(gomock.Any(), ref).
					Return(nil)
			},
			expectErr: false,
		},
		{
			name: "local repo not found",
			refs: "github.com/user/repo",
			options: delete.Options{
				Local:  true,
				Remote: false,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				ref := repository.NewReference("github.com", "user", "repo")

				mrp.EXPECT().
					Parse("github.com/user/repo").
					Return(&ref, nil)

				mfs.EXPECT().
					FindByReference(gomock.Any(), mws, ref).
					Return(nil, nil) // No match found but no error
			},
			expectErr: false, // Not finding a local repo isn't an error
		},
		{
			name: "invalid reference",
			refs: "invalid-reference",
			options: delete.Options{
				Local:  true,
				Remote: true,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				mrp.EXPECT().
					Parse("invalid-reference").
					Return(nil, errors.New("invalid repository reference"))
			},
			expectErr:     true,
			expectErrText: "invalid repository reference",
		},
		{
			name: "error finding local repository",
			refs: "github.com/user/repo",
			options: delete.Options{
				Local:  true,
				Remote: false,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				ref := repository.NewReference("github.com", "user", "repo")

				mrp.EXPECT().
					Parse("github.com/user/repo").
					Return(&ref, nil)

				mfs.EXPECT().
					FindByReference(gomock.Any(), mws, ref).
					Return(nil, errors.New("finder error"))
			},
			expectErr:     true,
			expectErrText: "deleting local: finding local repository: finder error",
		},
		{
			name: "error deleting remote repository",
			refs: "github.com/user/repo",
			options: delete.Options{
				Local:  false,
				Remote: true,
			},
			setupMocks: func(mws *workspace_mock.MockWorkspaceService, mfs *workspace_mock.MockFinderService, mhs *hosting_mock.MockHostingService, mrp *repository_mock.MockReferenceParser) {
				ref := repository.NewReference("github.com", "user", "repo")

				mrp.EXPECT().
					Parse("github.com/user/repo").
					Return(&ref, nil)

				mhs.EXPECT().
					DeleteRepository(gomock.Any(), ref).
					Return(errors.New("remote deletion error"))
			},
			expectErr:     true,
			expectErrText: "deleting remote: remote deletion error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWorkspaceService := workspace_mock.NewMockWorkspaceService(ctrl)
			mockFinderService := workspace_mock.NewMockFinderService(ctrl)
			mockHostingService := hosting_mock.NewMockHostingService(ctrl)
			mockReferenceParser := repository_mock.NewMockReferenceParser(ctrl)

			tc.setupMocks(mockWorkspaceService, mockFinderService, mockHostingService, mockReferenceParser)

			usecase := delete.NewUsecase(
				mockWorkspaceService,
				mockFinderService,
				mockHostingService,
				mockReferenceParser,
			)

			err := usecase.Execute(context.Background(), tc.refs, tc.options)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.expectErrText != "" && err.Error() != tc.expectErrText {
					t.Fatalf("Expected error %q but got %q", tc.expectErrText, err.Error())
				}
			} else if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}
		})
	}
}
