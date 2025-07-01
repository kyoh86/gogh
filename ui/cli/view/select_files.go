package view

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
)

// FileSelection represents the result of file selection
type FileSelection struct {
	Selected []string
	Skipped  []string
}

// SelectFiles shows an interactive file selection dialog
func SelectFiles(ctx context.Context, repoPath string, files []string) (*FileSelection, error) {
	logger := log.FromContext(ctx)

	// Show file list with checkboxes
	var selected []string
	var options []huh.Option[string]

	for _, file := range files {
		relPath, err := filepath.Rel(repoPath, file)
		if err != nil {
			relPath = file
		}
		options = append(options, huh.Option[string]{
			Key:   relPath,
			Value: file,
		})
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(fmt.Sprintf("Select files to save as auto-apply extra (%d files found)", len(files))).
				Description("Use space to select/deselect, enter to confirm").
				Options(options...).
				Value(&selected),
		),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	// Build skipped list
	selectedMap := make(map[string]bool)
	for _, f := range selected {
		selectedMap[f] = true
	}

	var skipped []string
	for _, f := range files {
		if !selectedMap[f] {
			skipped = append(skipped, f)
		}
	}

	// Show summary
	if len(selected) > 0 {
		logger.Infof("Selected %d files:", len(selected))
		for _, f := range selected {
			relPath, _ := filepath.Rel(repoPath, f)
			logger.Infof("  - %s", relPath)
		}
	}

	if len(skipped) > 0 {
		logger.Infof("Skipped %d files", len(skipped))
	}

	return &FileSelection{
		Selected: selected,
		Skipped:  skipped,
	}, nil
}

// ConfirmFilesIterative shows files one by one for confirmation
func ConfirmFilesIterative(ctx context.Context, repoPath string, files []string) (*FileSelection, error) {
	logger := log.FromContext(ctx)

	var selected []string
	var skipped []string
	var all bool

	for i, file := range files {
		relPath, err := filepath.Rel(repoPath, file)
		if err != nil {
			relPath = file
		}

		if all {
			selected = append(selected, file)
			continue
		}

		var choice string
		prompt := fmt.Sprintf("[%d/%d] %s", i+1, len(files), relPath)

		if err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				CharLimit(1).
				Inline(true).
				Title(prompt + " ").
				Validate(func(s string) error {
					s = strings.ToLower(s)
					if s == "y" || s == "n" || s == "q" || s == "a" {
						return nil
					}
					return fmt.Errorf("invalid selection, please enter y/n/q/a")
				}).
				Prompt("(y/n/q/a): ").
				Value(&choice),
		)).Run(); err != nil {
			return nil, err
		}

		switch strings.ToLower(choice) {
		case "a":
			all = true
			selected = append(selected, file)
			// Add remaining files
			selected = append(selected, files[i+1:]...)
			logger.Infof("Selected all remaining files")
		case "y":
			selected = append(selected, file)
		case "n":
			skipped = append(skipped, file)
		case "q":
			logger.Info("Quit file selection")
			// Add remaining files to skipped
			skipped = append(skipped, files[i+1:]...)
			return &FileSelection{
				Selected: selected,
				Skipped:  skipped,
			}, ErrQuit
		}
	}

	return &FileSelection{
		Selected: selected,
		Skipped:  skipped,
	}, nil
}
