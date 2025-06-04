package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v4/app/overlay_list"
	"github.com/kyoh86/gogh/v4/app/overlay_show"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/typ"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func NewOverlayShowCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	var f struct {
		repoPattern  string
		forInit      bool
		relativePath string
	}
	checkFlags := func(ctx context.Context, _ []string) ([]overlay_list.OverlayEntry, error) {
		list, err := typ.CollectWithError(typ.FilterE(overlay_list.NewUseCase(svc.OverlayStore).Execute(ctx), func(entry *overlay_list.OverlayEntry) (bool, error) {
			if f.repoPattern != "" && f.repoPattern != entry.RepoPattern {
				return false, nil
			}
			if f.forInit && !entry.ForInit {
				return false, nil
			}
			if f.relativePath != "" && f.relativePath != entry.RelativePath {
				return false, nil
			}
			return true, nil
		}))
		if err != nil {
			return nil, fmt.Errorf("listing overlays: %w", err)
		}
		switch len(list) {
		case 0:
			return nil, fmt.Errorf("no overlays found matching the criteria")
		case 1:
			return []overlay_list.OverlayEntry{*list[0]}, nil
		}
		var opts []huh.Option[overlay_list.OverlayEntry]
		for _, entry := range list {
			opts = append(opts, huh.Option[overlay_list.OverlayEntry]{
				Key:   entry.String(),
				Value: *entry,
			})
		}
		var selected []overlay_list.OverlayEntry
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[overlay_list.OverlayEntry]().
				Title("Overlays to show").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return nil, err
		}
		return selected, nil
	}

	var width = 60
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show overlays",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			entries, err := checkFlags(ctx, args)
			if err != nil {
				return err
			}
			if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
				width = w
			}
			overlayShowUseCase := overlay_show.NewUseCase(svc.OverlayStore)
			for i, entry := range entries {
				if i > 0 {
					fmt.Println()
				}
				name := entry.String()
				fmt.Printf("%s %s\n", name, strings.Repeat("-", width-len(name)-1))
				if err := overlayShowUseCase.Execute(ctx, entry.RepoPattern, entry.ForInit, entry.RelativePath); err != nil {
					return fmt.Errorf("showing overlay %s: %w", entry.RelativePath, err)
				}
			}
			return nil
		},
	}

	return cmd, nil
}
