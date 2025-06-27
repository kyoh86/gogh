package update

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/script"
)

type Options struct {
	Name          string
	RepoPattern   string
	TriggerEvent  string
	OperationType string
	OperationID   string
}

type Usecase struct {
	hookService    hook.HookService
	overlayService overlay.OverlayService
	scriptService  script.ScriptService
}

func NewUsecase(hookService hook.HookService, overlayService overlay.OverlayService, scriptService script.ScriptService) *Usecase {
	return &Usecase{
		hookService:    hookService,
		overlayService: overlayService,
		scriptService:  scriptService,
	}
}

func (uc *Usecase) Execute(ctx context.Context, id string, opts Options) error {
	// Resolve OperationID to full UUID if operation type and ID are specified
	var resolvedID uuid.UUID
	if opts.OperationType != "" && opts.OperationID != "" {
		var err error
		resolvedID, err = uc.resolveOperationID(ctx, opts.OperationType, opts.OperationID)
		if err != nil {
			return err
		}
	}

	h := hook.Entry{
		Name:          opts.Name,
		RepoPattern:   opts.RepoPattern,
		TriggerEvent:  hook.Event(opts.TriggerEvent),
		OperationType: hook.OperationType(opts.OperationType),
		OperationID:   resolvedID,
	}
	return uc.hookService.Update(ctx, id, h)
}

func (uc *Usecase) resolveOperationID(ctx context.Context, opType string, idlike string) (uuid.UUID, error) {
	var id uuid.UUID
	switch hook.OperationType(opType) {
	case hook.OperationTypeOverlay:
		overlay, err := uc.overlayService.Get(ctx, idlike)
		if err != nil {
			return id, fmt.Errorf("failed to resolve overlay ID: %w", err)
		}
		return overlay.UUID(), nil
	case hook.OperationTypeScript:
		script, err := uc.scriptService.Get(ctx, idlike)
		if err != nil {
			return id, fmt.Errorf("failed to resolve script ID: %w", err)
		}
		return script.UUID(), nil
	default:
		return id, fmt.Errorf("invalid operation type: %s", opType)
	}
}
