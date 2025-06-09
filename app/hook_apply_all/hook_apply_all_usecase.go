package hook_apply_all

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/hook_apply"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// UseCase represents the hook apply all use case
type UseCase struct {
	hookService      hook.HookService
	referenceParser  repository.ReferenceParser
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
}

// NewUseCase creates a new hook apply all use case
func NewUseCase(
	hookService hook.HookService,
	referenceParser repository.ReferenceParser,
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
) *UseCase {
	return &UseCase{
		hookService:      hookService,
		referenceParser:  referenceParser,
		workspaceService: workspaceService,
		finderService:    finderService,
	}
}

// Execute applies hooks to the given reference and use case and event.
func (uc *UseCase) Execute(ctx context.Context, refWithAlias string, useCase hook.UseCase, event hook.Event) error {
	hookApplyUseCase := hook_apply.NewUseCase(
		uc.hookService,
		uc.referenceParser,
		uc.workspaceService,
		uc.finderService,
	)
	ref, err := uc.referenceParser.ParseWithAlias(refWithAlias)
	if err != nil {
		return fmt.Errorf("parsing reference '%s': %w", refWithAlias, err)
	}
	for h, err := range uc.hookService.ListHooks() {
		if err != nil {
			return fmt.Errorf("list hooks error: %w", err)
		}
		match, err := h.Match(ref.Local(), useCase, event)
		if err != nil {
			return fmt.Errorf("hook %s match error: %w", h.Name, err)
		}
		if !match {
			continue
		}
		if err := hookApplyUseCase.Execute(ctx, h.ID, refWithAlias, map[string]any{
			"use_case": useCase,
			"event":    event,
		}); err != nil {
			return fmt.Errorf("applying hook '%s': %w", h.Name, err)
		}
	}
	return nil
}
