package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/overlay_list"
	"github.com/kyoh86/gogh/v4/app/overlay_show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewOverlayDescribeCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		repoPattern  string
		forInit      bool
		relativePath string
	}
	checkFlags := func(ctx context.Context, _ []string) ([]overlay_list.Overlay, error) {
		overlays, err := typ.CollectWithError(typ.FilterE(overlay_list.NewUseCase(svc.OverlayService).Execute(ctx), func(ov *overlay_list.Overlay) (bool, error) {
			if f.repoPattern != "" && f.repoPattern != ov.RepoPattern {
				return false, nil
			}
			if f.forInit && !ov.ForInit {
				return false, nil
			}
			if f.relativePath != "" && f.relativePath != ov.RelativePath {
				return false, nil
			}
			return true, nil
		}))
		if err != nil {
			return nil, fmt.Errorf("listing overlays: %w", err)
		}
		switch len(overlays) {
		case 0:
			return nil, fmt.Errorf("no overlays found matching the criteria")
		case 1:
			return []overlay_list.Overlay{*overlays[0]}, nil
		}
		var opts []huh.Option[overlay_list.Overlay]
		for _, ov := range overlays {
			opts = append(opts, huh.Option[overlay_list.Overlay]{
				Key:   ov.String(),
				Value: *ov,
			})
		}
		var selected []overlay_list.Overlay
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[overlay_list.Overlay]().
				Title("Overlays to describe").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return nil, err
		}
		return selected, nil
	}

	width := 60
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe overlays",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			overlays, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				width = w
			}
			overlayShowUseCase := overlay_show.NewUseCaseContent(svc.OverlayService, cmd.OutOrStdout(), width)
			for i, ov := range overlays {
				if i > 0 {
					fmt.Println()
				}
				if err := overlayShowUseCase.Execute(ctx, &ov); err != nil {
					return fmt.Errorf("describing overlay %s: %w", ov.RelativePath, err)
				}
			}
			return nil
		},
	}

	return cmd, nil
}
