package invoke

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/app/overlay/apply"
	scriptinvoke "github.com/kyoh86/gogh/v4/app/script/invoke"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/workspace"
)

// Event defines the trigger of the hook, such as post-clone, post-fork, or post-create
type Event = hook.Event

const (
	EventAny        Event = hook.EventAny
	EventPostClone  Event = hook.EventPostClone
	EventPostFork   Event = hook.EventPostFork
	EventPostCreate Event = hook.EventPostCreate
)

type Options struct{}

type Usecase struct {
	workspaceService workspace.WorkspaceService
	finderService    workspace.FinderService
	hookService      hook.HookService
	overlayService   overlay.OverlayService
	scriptService    script.ScriptService
	referenceParser  repository.ReferenceParser
}

func NewUsecase(
	workspaceService workspace.WorkspaceService,
	finderService workspace.FinderService,
	hookService hook.HookService,
	overlayService overlay.OverlayService,
	scriptService script.ScriptService,
	referenceParser repository.ReferenceParser,
) *Usecase {
	return &Usecase{
		workspaceService: workspaceService,
		finderService:    finderService,
		hookService:      hookService,
		overlayService:   overlayService,
		scriptService:    scriptService,
		referenceParser:  referenceParser,
	}
}

// Invoke executes the hooks for the given hookID and refStr.
func (uc *Usecase) Invoke(ctx context.Context, hookID, refStr string) error {
	h, err := uc.hookService.Get(ctx, hookID)
	if err != nil {
		return err
	}
	switch h.OperationType() {
	case hook.OperationTypeOverlay:
		overlayApplyUsecase := apply.NewUsecase(
			uc.workspaceService,
			uc.finderService,
			uc.referenceParser,
			uc.overlayService,
		)
		return overlayApplyUsecase.Execute(ctx, refStr, h.OperationID())
	case hook.OperationTypeScript:
		scriptApplyUsecase := scriptinvoke.NewUsecase(
			uc.workspaceService,
			uc.finderService,
			uc.scriptService,
			uc.referenceParser,
		)
		return scriptApplyUsecase.Execute(ctx, refStr, h.OperationID(), map[string]any{
			"hook": map[string]any{
				"id":            h.ID(),
				"name":          h.Name(),
				"repoPattern":   h.RepoPattern(),
				"triggerEvent":  h.TriggerEvent(),
				"operationType": h.OperationType(),
				"operationId":   h.OperationID(),
			},
		})
	}
	return fmt.Errorf("unsupported hook operation type: %q", h.OperationType())
}

// InvokeFor executes all hooks that match the repository and the event
func (uc *Usecase) InvokeFor(ctx context.Context, event Event, refStr string) error {
	return uc.InvokeForWithGlobals(ctx, event, refStr, nil)
}

// InvokeForWithGlobals executes all hooks that match the repository and the event with additional globals
func (uc *Usecase) InvokeForWithGlobals(ctx context.Context, event Event, refStr string, globals map[string]any) error {
	refWithAlias, err := uc.referenceParser.ParseWithAlias(refStr)
	if err != nil {
		return fmt.Errorf("parsing repository reference: %w", err)
	}
	match, err := uc.finderService.FindByReference(ctx, uc.workspaceService, refWithAlias.Local())
	if err != nil {
		return fmt.Errorf("find repository location: %w", err)
	}
	if match == nil {
		return fmt.Errorf("repository not found for reference: %s", refStr)
	}
	overlayApplyUsecase := apply.NewUsecase(
		uc.workspaceService,
		uc.finderService,
		uc.referenceParser,
		uc.overlayService,
	)
	scriptApplyUsecase := scriptinvoke.NewUsecase(
		uc.workspaceService,
		uc.finderService,
		uc.scriptService,
		uc.referenceParser,
	)
	for h, err := range uc.hookService.ListFor(refWithAlias.Local(), event) {
		if err != nil {
			return err
		}
		switch h.OperationType() {
		case hook.OperationTypeOverlay:
			if err := overlayApplyUsecase.Apply(ctx, match, h.OperationID()); err != nil {
				return fmt.Errorf("applying overlay for the hook %s: %w", h.ID(), err)
			}
		case hook.OperationTypeScript:
			g := make(map[string]any)
			for k, v := range globals {
				g[k] = v
			}
			g["hook"] = map[string]any{
				"id":            h.ID(),
				"name":          h.Name(),
				"repoPattern":   h.RepoPattern(),
				"triggerEvent":  string(h.TriggerEvent()),
				"operationType": string(h.OperationType()),
				"operationId":   h.OperationID(),
			}
			if err := scriptApplyUsecase.Invoke(ctx, match, h.OperationID(), g); err != nil {
				return fmt.Errorf("invoking script for the hook %s: %w", h.ID(), err)
			}
		}
	}
	return nil
}
