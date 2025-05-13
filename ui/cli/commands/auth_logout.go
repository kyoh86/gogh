package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
	"github.com/kyoh86/gogh/v3/app/auth_list"
	"github.com/kyoh86/gogh/v3/app/auth_logout"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

func NewAuthLogoutCommand(_ context.Context, svc *service.ServiceSet) *cobra.Command {
	listUseCase := auth_list.NewUseCase(svc.TokenService)
	logoutUseCase := auth_logout.NewUseCase(svc.TokenService)

	checkFlags := func(cmd *cobra.Command, args []string) ([]string, error) {
		if len(args) > 0 {
			return args, nil
		}
		entries, err := listUseCase.Execute(cmd.Context())
		if err != nil {
			return nil, fmt.Errorf("failed to list tokens: %w", err)
		}
		if len(entries) == 0 {
			return nil, errors.New("no valid token found")
		}
		opts := make([]huh.Option[string], 0, len(entries))
		for _, c := range entries {
			name := fmt.Sprintf("%s/%s", c.Host, c.Owner)
			opts = append(opts, huh.Option[string]{Key: name, Value: name})
		}

		var selected []string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Hosts to logout from").
				Options(opts...).
				Value(&selected),
		)).Run(); err != nil {
			return nil, err
		}
		return selected, nil
	}

	preprocessFlags := func(cmd *cobra.Command, args []string) [][2]string {
		rets := make([][2]string, 0, len(args))
		for _, target := range args {
			words := strings.SplitN(target, "/", 2)
			if len(words) != 2 {
				log.FromContext(cmd.Context()).WithField("target", target).Error("invalid target (must be host/owner)")
				continue
			}
			rets = append(rets, [2]string{words[0], words[1]})
		}
		return rets
	}

	return &cobra.Command{
		Use:     "logout",
		Aliases: []string{"signout", "remove"},
		Short:   "Logout from the host and owner",
		RunE: func(cmd *cobra.Command, args []string) error {
			indices, err := checkFlags(cmd, args)
			if err != nil {
				return err
			}

			owners := preprocessFlags(cmd, indices)

			for _, target := range owners {
				log.FromContext(cmd.Context()).WithField("target", target).Info("logout from")
				if err := logoutUseCase.Execute(target[0], target[1]); err != nil {
					log.FromContext(cmd.Context()).WithField("target", target).Errorf("failed to delete token: %s", err)
					continue
				}
			}
			return nil
		},
	}
}
