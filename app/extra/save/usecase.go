package save

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kyoh86/gogh/v4/app/hook/add"
	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/git"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Usecase for saving auto-apply extra
type Usecase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	gitService       git.GitService
	overlayService   overlay.OverlayService
	scriptService    script.ScriptService
	hookService      hook.HookService
	extraService     extra.ExtraService
	referenceParser  repository.ReferenceParser
}

// NewUsecase creates a new extra save use case
func NewUsecase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	gitService git.GitService,
	overlayService overlay.OverlayService,
	scriptService script.ScriptService,
	hookService hook.HookService,
	extraService extra.ExtraService,
	referenceParser repository.ReferenceParser,
) *Usecase {
	return &Usecase{
		workspaceService: workspaceService,
		finderService:    finderService,
		gitService:       gitService,
		overlayService:   overlayService,
		scriptService:    scriptService,
		hookService:      hookService,
		extraService:     extraService,
		referenceParser:  referenceParser,
	}
}

// Execute saves excluded files as auto-apply extra
func (uc *Usecase) Execute(ctx context.Context, repoStr string) error {
	result, err := uc.GetExcludedFiles(ctx, repoStr)
	if err != nil {
		return err
	}

	if len(result.Files) == 0 {
		return fmt.Errorf("no excluded files found in %s", repoStr)
	}

	return uc.SaveFiles(ctx, repoStr, result.Files)
}

// ExcludedFilesResult contains the result of GetExcludedFiles
type ExcludedFilesResult struct {
	RepositoryPath string
	Files          []string
}

// GetExcludedFiles returns list of excluded files for the repository
func (uc *Usecase) GetExcludedFiles(ctx context.Context, repoStr string) (*ExcludedFilesResult, error) {
	// Parse repository reference
	ref, err := uc.referenceParser.Parse(repoStr)
	if err != nil {
		return nil, fmt.Errorf("parsing repository reference: %w", err)
	}

	// Find repository location
	location, err := uc.finderService.FindByReference(ctx, uc.workspaceService, *ref)
	if err != nil {
		return nil, fmt.Errorf("finding repository: %w", err)
	}
	if location == nil {
		return nil, fmt.Errorf("repository not found: %s", repoStr)
	}

	// Get excluded files
	var excludedFiles []string
	for file, err := range uc.gitService.ListExcludedFiles(ctx, location.FullPath(), nil) {
		if err != nil {
			return nil, fmt.Errorf("listing excluded files: %w", err)
		}
		excludedFiles = append(excludedFiles, file)
	}

	return &ExcludedFilesResult{
		RepositoryPath: location.FullPath(),
		Files:          excludedFiles,
	}, nil
}

// SaveFiles saves specified files as auto-apply extra
func (uc *Usecase) SaveFiles(ctx context.Context, repoStr string, files []string) error {
	// Parse repository reference
	ref, err := uc.referenceParser.Parse(repoStr)
	if err != nil {
		return fmt.Errorf("parsing repository reference: %w", err)
	}

	// Find repository location
	location, err := uc.finderService.FindByReference(ctx, uc.workspaceService, *ref)
	if err != nil {
		return fmt.Errorf("finding repository: %w", err)
	}
	if location == nil {
		return fmt.Errorf("repository not found: %s", repoStr)
	}

	// Check if auto extra already exists
	existing, err := uc.extraService.GetAutoExtra(ctx, *ref)
	if err == nil && existing != nil {
		return fmt.Errorf("auto extra already exists for %s, use 'extras clear' first", repoStr)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files provided")
	}

	// Create overlay and hook for each file
	var items []extra.Item
	for _, file := range files {
		// ListExcludedFiles returns absolute paths, so use file directly
		fullPath := file

		// Make the path relative to the repository for the overlay
		relPath, err := filepath.Rel(location.FullPath(), file)
		if err != nil {
			return fmt.Errorf("making path relative: %w", err)
		}

		// Read file content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", file, err)
		}

		// Create overlay
		overlayID, err := uc.overlayService.Add(ctx, overlay.Entry{
			Name:         fmt.Sprintf("Auto extra: %s", relPath),
			RelativePath: relPath,
			Content:      bytes.NewReader(content),
		})
		if err != nil {
			// Rollback created overlays
			for _, item := range items {
				_ = uc.overlayService.Remove(ctx, item.OverlayID)
				_ = uc.hookService.Remove(ctx, item.HookID)
			}
			return fmt.Errorf("creating overlay for %s: %w", file, err)
		}

		// Create post-clone hook for this overlay
		hookAddUC := add.NewUsecase(uc.hookService, uc.overlayService, uc.scriptService)
		hookID, err := hookAddUC.Execute(ctx, add.Options{
			Name:          fmt.Sprintf("Auto extra for %s: %s", ref.String(), relPath),
			RepoPattern:   ref.String(),
			TriggerEvent:  string(hook.EventPostClone),
			OperationType: string(hook.OperationTypeOverlay),
			OperationID:   overlayID,
		})
		if err != nil {
			// Rollback
			_ = uc.overlayService.Remove(ctx, overlayID)
			for _, item := range items {
				_ = uc.overlayService.Remove(ctx, item.OverlayID)
				_ = uc.hookService.Remove(ctx, item.HookID)
			}
			return fmt.Errorf("creating hook for %s: %w", file, err)
		}

		items = append(items, extra.Item{
			OverlayID: overlayID,
			HookID:    hookID,
		})
	}

	// Save auto extra
	_, err = uc.extraService.AddAutoExtra(ctx, *ref, *ref, items)
	if err != nil {
		// Rollback
		for _, item := range items {
			_ = uc.overlayService.Remove(ctx, item.OverlayID)
			_ = uc.hookService.Remove(ctx, item.HookID)
		}
		return fmt.Errorf("saving auto extra: %w", err)
	}

	return nil
}
