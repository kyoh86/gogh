package extra_remove

import (
	"context"
	"fmt"

	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/repository"
)

// UseCase represents the extra remove use case
type UseCase struct {
	extraService    extra.ExtraService
	referenceParser repository.ReferenceParser
}

// NewUseCase creates a new extra remove use case
func NewUseCase(
	extraService extra.ExtraService,
	referenceParser repository.ReferenceParser,
) *UseCase {
	return &UseCase{
		extraService:    extraService,
		referenceParser: referenceParser,
	}
}

// Options contains options for the extra remove operation
type Options struct {
	ID         string
	Name       string // For named extras
	Repository string // For auto extras
}

// Execute performs the extra remove operation
func (uc *UseCase) Execute(ctx context.Context, opts Options) error {
	switch {
	case opts.ID != "":
		// Remove by ID
		if err := uc.extraService.Remove(ctx, opts.ID); err != nil {
			return fmt.Errorf("removing extra by ID: %w", err)
		}
		fmt.Printf("Removed extra with ID %s\n", opts.ID)
	case opts.Name != "":
		// Remove named extra
		if err := uc.extraService.RemoveNamedExtra(ctx, opts.Name); err != nil {
			return fmt.Errorf("removing named extra: %w", err)
		}
		fmt.Printf("Removed named extra %q\n", opts.Name)
	case opts.Repository != "":
		// Remove auto extra
		ref, err := uc.referenceParser.Parse(opts.Repository)
		if err != nil {
			return fmt.Errorf("invalid repository reference: %w", err)
		}
		if err := uc.extraService.RemoveAutoExtra(ctx, *ref); err != nil {
			return fmt.Errorf("removing auto extra: %w", err)
		}
		fmt.Printf("Removed auto extra for repository %s\n", ref.String())
	default:
		return fmt.Errorf("one of --id, --name, or --repository must be specified")
	}

	return nil
}
