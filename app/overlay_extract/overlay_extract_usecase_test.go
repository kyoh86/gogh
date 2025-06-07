package overlay_extract_test

import (
	"context"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"testing"

	testtarget "github.com/kyoh86/gogh/v4/app/overlay_extract"
	"github.com/kyoh86/gogh/v4/core/git_mock"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/repository_mock"
	"github.com/kyoh86/gogh/v4/core/workspace_mock"
	"go.uber.org/mock/gomock"
)

func TestExecute(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		setupMocks    func(*git_mock.MockGitService, *overlay_mock.MockOverlayService, *workspace_mock.MockWorkspaceService, *workspace_mock.MockFinderService, *repository_mock.MockReferenceParser, string)
		refString     string
		options       testtarget.Options
		expectedCount int
		expectedError bool
		setupFiles    func(repoPath string) []string
	}{
		{
			name: "Success: When untracked files exist",
			setupMocks: func(git *git_mock.MockGitService, overlay *overlay_mock.MockOverlayService, ws *workspace_mock.MockWorkspaceService, finder *workspace_mock.MockFinderService, parser *repository_mock.MockReferenceParser, repoPath string) {
				// Setup reference
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				parser.EXPECT().
					Parse("github.com/kyoh86/gogh").
					Return(&ref, nil)

				// Setup repository location
				loc := repository.NewLocation(repoPath, "github.com", "kyoh86", "gogh")
				finder.EXPECT().
					FindByReference(gomock.Any(), ws, gomock.Eq(ref)).
					Return(loc, nil)

				// Setup untracked files with the actual file paths
				git.EXPECT().
					ListExcludedFiles(gomock.Any(), repoPath, gomock.Any()).
					Return(maps.All(map[string]error{
						filepath.Join(repoPath, "file1.txt"): nil,
						filepath.Join(repoPath, "file2.txt"): nil,
					}))
			},
			refString:     "github.com/kyoh86/gogh",
			options:       testtarget.Options{Excluded: true},
			expectedCount: 2,
			expectedError: false,
			setupFiles: func(repoPath string) []string {
				// Create test files
				file1Path := filepath.Join(repoPath, "file1.txt")
				file2Path := filepath.Join(repoPath, "file2.txt")

				err := os.WriteFile(file1Path, []byte("content of file 1"), 0644)
				if err != nil {
					panic(err)
				}

				err = os.WriteFile(file2Path, []byte("content of file 2"), 0644)
				if err != nil {
					panic(err)
				}

				return []string{file1Path, file2Path}
			},
		},
		{
			name: "Success: When no untracked files exist",
			setupMocks: func(git *git_mock.MockGitService, overlay *overlay_mock.MockOverlayService, ws *workspace_mock.MockWorkspaceService, finder *workspace_mock.MockFinderService, parser *repository_mock.MockReferenceParser, repoPath string) {
				// Setup reference
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				parser.EXPECT().
					Parse("github.com/kyoh86/gogh").
					Return(&ref, nil)

				// Setup repository location
				loc := repository.NewLocation(repoPath, "github.com", "kyoh86", "gogh")
				finder.EXPECT().
					FindByReference(gomock.Any(), ws, gomock.Eq(ref)).
					Return(loc, nil)

				// Setup empty untracked files
				git.EXPECT().
					ListExcludedFiles(gomock.Any(), repoPath, gomock.Any()).
					Return(maps.All(map[string]error{}))
			},
			refString:     "github.com/kyoh86/gogh",
			options:       testtarget.Options{Excluded: true},
			expectedCount: 0,
			expectedError: false,
			setupFiles: func(repoPath string) []string {
				// No files to create
				return []string{}
			},
		},
		{
			name: "Error: When reference parsing fails",
			setupMocks: func(git *git_mock.MockGitService, overlay *overlay_mock.MockOverlayService, ws *workspace_mock.MockWorkspaceService, finder *workspace_mock.MockFinderService, parser *repository_mock.MockReferenceParser, repoPath string) {
				parser.EXPECT().
					Parse("invalid/ref").
					Return(nil, errors.New("invalid reference format"))
			},
			refString:     "invalid/ref",
			options:       testtarget.Options{Excluded: true},
			expectedCount: 0,
			expectedError: true,
			setupFiles: func(repoPath string) []string {
				return []string{}
			},
		},
		{
			name: "Error: When repository finder fails",
			setupMocks: func(git *git_mock.MockGitService, overlay *overlay_mock.MockOverlayService, ws *workspace_mock.MockWorkspaceService, finder *workspace_mock.MockFinderService, parser *repository_mock.MockReferenceParser, repoPath string) {
				// Setup reference
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				parser.EXPECT().
					Parse("github.com/kyoh86/gogh").
					Return(&ref, nil)

				// Setup finder error
				finder.EXPECT().
					FindByReference(gomock.Any(), ws, gomock.Eq(ref)).
					Return(nil, errors.New("repository not found"))
			},
			refString:     "github.com/kyoh86/gogh",
			options:       testtarget.Options{Excluded: true},
			expectedCount: 0,
			expectedError: true,
			setupFiles: func(repoPath string) []string {
				return []string{}
			},
		},
		{
			name: "Error: When git service fails",
			setupMocks: func(git *git_mock.MockGitService, overlay *overlay_mock.MockOverlayService, ws *workspace_mock.MockWorkspaceService, finder *workspace_mock.MockFinderService, parser *repository_mock.MockReferenceParser, repoPath string) {
				// Setup reference
				ref := repository.NewReference("github.com", "kyoh86", "gogh")
				parser.EXPECT().
					Parse("github.com/kyoh86/gogh").
					Return(&ref, nil)

				// Setup repository location
				loc := repository.NewLocation(repoPath, "github.com", "kyoh86", "gogh")
				finder.EXPECT().
					FindByReference(gomock.Any(), ws, gomock.Eq(ref)).
					Return(loc, nil)

				// Setup git error
				git.EXPECT().
					ListExcludedFiles(gomock.Any(), repoPath, gomock.Any()).
					Return(maps.All(map[string]error{"": errors.New("git command failed")}))
			},
			refString:     "github.com/kyoh86/gogh",
			options:       testtarget.Options{Excluded: true},
			expectedCount: 0,
			expectedError: true,
			setupFiles: func(repoPath string) []string {
				return []string{}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup gomock controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Create a temporary directory for the test
			tempRepoDir := t.TempDir()

			// Setup test files
			filePaths := tt.setupFiles(tempRepoDir)

			// Create mocks
			mockGit := git_mock.NewMockGitService(ctrl)
			mockOverlay := overlay_mock.NewMockOverlayService(ctrl)
			mockWs := workspace_mock.NewMockWorkspaceService(ctrl)
			mockFinder := workspace_mock.NewMockFinderService(ctrl)
			mockParser := repository_mock.NewMockReferenceParser(ctrl)

			// Setup mocks with the temporary directory path
			tt.setupMocks(mockGit, mockOverlay, mockWs, mockFinder, mockParser, tempRepoDir)

			// Create the target UseCase
			useCase := testtarget.NewUseCase(mockGit, mockOverlay, mockWs, mockFinder, mockParser)

			// Execute the UseCase
			ctx := context.Background()
			result := useCase.Execute(ctx, tt.refString, tt.options)

			// Verify the results
			count := 0
			var err error
			for entry, iterErr := range result {
				if iterErr != nil {
					err = iterErr
					break
				}
				if entry != nil {
					count++
					// Check that the entry has the expected fields
					if entry.RelativePath == "" {
						t.Error("Expected RelativePath to be non-empty")
					}
					if entry.FilePath == "" {
						t.Error("Expected FilePath to be non-empty")
					}
					if entry.Reference.String() == "" {
						t.Error("Expected Reference to be non-empty")
					}

					// Verify the file content if it's a success case
					if len(filePaths) > 0 {
						// Read some content to verify it's valid
						buf, readErr := os.ReadFile(entry.FilePath)
						if readErr != nil && readErr.Error() != "EOF" {
							t.Errorf("Failed to read content: %v", readErr)
						}
						if len(buf) == 0 && len(filePaths) > 0 {
							t.Error("Expected to read some content but got nothing")
						}
					}
				}
			}

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.expectedError && count != tt.expectedCount {
				t.Errorf("Expected %d entries, got %d", tt.expectedCount, count)
			}
		})
	}
}
