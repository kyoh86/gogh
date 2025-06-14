package overlay_extract

import (
	"context"
	"fmt"
	"iter"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase implements the overlay extraction use case
type UseCase struct {
	gitService       git.GitService
	overlayService   overlay.OverlayService
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	referenceParser  repository.ReferenceParser
}

// NewUseCase creates a new overlay extraction use case
func NewUseCase(
	gitService git.GitService,
	overlayService overlay.OverlayService,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	referenceParser repository.ReferenceParser,
) *UseCase {
	return &UseCase{
		gitService:       gitService,
		overlayService:   overlayService,
		workspaceService: workspaceService,
		finderService:    finderService,
		referenceParser:  referenceParser,
	}
}

// Options for the extraction operation
type Options struct {
	Excluded     bool     // Whether to extract only excluded files
	FilePatterns []string // Pattern to match files for extraction; if empty, all excluded files are returned.
}

// ExtractResult represents a single untracked file that can be extracted
type ExtractResult struct {
	Reference    repository.Reference // Reference to the repository
	RelativePath string
	FilePath     string // Path of the untracked file
}

// Extract files in the repository and returns them.
// The caller is responsible for confirming and registering files as overlays
func (uc *UseCase) Execute(ctx context.Context, refs string, opts Options) iter.Seq2[*ExtractResult, error] {
	return func(yield func(*ExtractResult, error) bool) {
		ref, err := uc.referenceParser.Parse(refs)
		if err != nil {
			yield(nil, err)
			return
		}

		// Find the repository path
		repo, err := uc.finderService.FindByReference(ctx, uc.workspaceService, *ref)
		if err != nil {
			yield(nil, fmt.Errorf("failed to find repository: %w", err))
			return
		}

		files := uc.gitService.ListAllFiles
		if opts.Excluded {
			files = uc.gitService.ListExcludedFiles
		}

		// Read file contents
		repoPath := repo.FullPath()
		for file, err := range files(ctx, repoPath, opts.FilePatterns) {
			if err != nil {
				yield(nil, fmt.Errorf("failed to list untracked files: %w", err))
				return
			}
			rel, err := filepath.Rel(repoPath, file)
			if err != nil {
				yield(nil, fmt.Errorf("failed to get relative path for %s: %w", file, err))
				return
			}
			if !yield(&ExtractResult{
				Reference:    *ref,
				RelativePath: rel,
				FilePath:     file,
			}, nil) {
				return
			}
		}
	}
}
