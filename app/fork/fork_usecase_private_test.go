package fork

import (
	"errors"
	"strings"
	"testing"

	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/hosting_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCase_parseRefs(t *testing.T) {
	testCases := []struct {
		name           string
		source         string
		target         string
		setupMocks     func(ctrl *gomock.Controller, uc *UseCase)
		expectedSource repository.Reference
		expectedTarget *repository.ReferenceWithAlias
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			name:   "valid source and target",
			source: "github.com/source/repo",
			target: "github.com/target/repo",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				uc.referenceParser = mockParser

				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := &repository.ReferenceWithAlias{
					Reference: repository.NewReference("github.com", "target", "repo"),
				}

				mockParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mockParser.EXPECT().
					ParseWithAlias("github.com/target/repo").
					Return(targetRef, nil)
			},
			expectedSource: repository.NewReference("github.com", "source", "repo"),
			expectedTarget: &repository.ReferenceWithAlias{
				Reference: repository.NewReference("github.com", "target", "repo"),
			},
			expectedErr: false,
		},
		{
			name:   "invalid source",
			source: "invalid-source",
			target: "github.com/target/repo",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				uc.referenceParser = mockParser

				mockParser.EXPECT().
					Parse("invalid-source").
					Return(nil, errors.New("invalid source"))
			},
			expectedErr:    true,
			expectedErrMsg: "invalid source",
		},
		{
			name:   "empty target uses default owner",
			source: "github.com/source/repo",
			target: "",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				mockDefaultNameService := repository_mock.NewMockDefaultNameService(ctrl)
				uc.referenceParser = mockParser
				uc.defaultNameService = mockDefaultNameService

				sourceRef := repository.NewReference("github.com", "source", "repo")
				mockParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)

				mockDefaultNameService.EXPECT().
					GetDefaultOwnerFor("github.com").
					Return("default-owner", nil)
			},
			expectedSource: repository.NewReference("github.com", "source", "repo"),
			expectedTarget: &repository.ReferenceWithAlias{
				Reference: repository.NewReference("github.com", "default-owner", "repo"),
			},
			expectedErr: false,
		},
		{
			name:   "invalid target",
			source: "github.com/source/repo",
			target: "invalid-target",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				uc.referenceParser = mockParser

				sourceRef := repository.NewReference("github.com", "source", "repo")
				mockParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)

				mockParser.EXPECT().
					ParseWithAlias("invalid-target").
					Return(nil, errors.New("invalid target"))
			},
			expectedErr:    true,
			expectedErrMsg: "invalid target",
		},
		{
			name:   "different hosts",
			source: "github.com/source/repo",
			target: "gitlab.com/target/repo",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				uc.referenceParser = mockParser

				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := &repository.ReferenceWithAlias{
					Reference: repository.NewReference("gitlab.com", "target", "repo"),
				}

				mockParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mockParser.EXPECT().
					ParseWithAlias("gitlab.com/target/repo").
					Return(targetRef, nil)
			},
			expectedErr:    true,
			expectedErrMsg: "the host of the forked repository must be the same as the original repository",
		},
		{
			name:   "empty owner",
			source: "github.com/source/repo",
			target: "github.com//repo",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				uc.referenceParser = mockParser

				sourceRef := repository.NewReference("github.com", "source", "repo")
				targetRef := &repository.ReferenceWithAlias{
					Reference: repository.NewReference("github.com", "", "repo"),
				}

				mockParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)
				mockParser.EXPECT().
					ParseWithAlias("github.com//repo").
					Return(targetRef, nil)
			},
			expectedErr:    true,
			expectedErrMsg: "the owner of the forked repository must be specified",
		},
		{
			name:   "error getting default owner",
			source: "github.com/source/repo",
			target: "",
			setupMocks: func(ctrl *gomock.Controller, uc *UseCase) {
				mockParser := repository_mock.NewMockReferenceParser(ctrl)
				mockDefaultNameService := repository_mock.NewMockDefaultNameService(ctrl)
				uc.referenceParser = mockParser
				uc.defaultNameService = mockDefaultNameService

				sourceRef := repository.NewReference("github.com", "source", "repo")
				mockParser.EXPECT().
					Parse("github.com/source/repo").
					Return(&sourceRef, nil)

				mockDefaultNameService.EXPECT().
					GetDefaultOwnerFor("github.com").
					Return("", errors.New("test"))
			},
			expectedErr:    true,
			expectedErrMsg: "getting default owner for",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := &UseCase{
				hostingService:     hosting_mock.NewMockHostingService(ctrl),
				workspaceService:   workspace_mock.NewMockWorkspaceService(ctrl),
				defaultNameService: repository_mock.NewMockDefaultNameService(ctrl),
				referenceParser:    repository_mock.NewMockReferenceParser(ctrl),
				gitService:         git_mock.NewMockGitService(ctrl),
			}

			tc.setupMocks(ctrl, uc)

			sourceRef, targetRef, err := uc.parseRefs(tc.source, tc.target)

			if tc.expectedErr {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}
				if tc.expectedErrMsg != "" && err.Error() != tc.expectedErrMsg && !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Fatalf("Expected error containing %q but got %q", tc.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}
				if sourceRef.Host() != tc.expectedSource.Host() ||
					sourceRef.Owner() != tc.expectedSource.Owner() ||
					sourceRef.Name() != tc.expectedSource.Name() {
					t.Errorf("Source reference mismatch: got %v, want %v", sourceRef, tc.expectedSource)
				}
				if targetRef.Reference.Host() != tc.expectedTarget.Reference.Host() ||
					targetRef.Reference.Owner() != tc.expectedTarget.Reference.Owner() ||
					targetRef.Reference.Name() != tc.expectedTarget.Reference.Name() {
					t.Errorf("Target reference mismatch: got %v, want %v", targetRef, tc.expectedTarget)
				}
			}
		})
	}
}
