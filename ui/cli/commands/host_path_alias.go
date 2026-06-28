package commands

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/service"
	"github.com/kyoh86/gogh/v4/core/workspace"
	"github.com/spf13/cobra"
)

func NewHostPathAliasCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "host-path-alias",
		Aliases: []string{"host-path-aliases", "host-alias", "host-aliases"},
		Short:   "Manage host path aliases",
		RunE:    HostPathAliasListRunE(svc),
	}, nil
}

func NewHostPathAliasListCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "list",
		Short: "List host path aliases",
		Args:  cobra.NoArgs,
		RunE:  HostPathAliasListRunE(svc),
	}, nil
}

func HostPathAliasListRunE(svc *service.ServiceSet) func(*cobra.Command, []string) error {
	return func(_ *cobra.Command, _ []string) error {
		aliases := svc.WorkspaceService.GetHostPathAliases()
		if len(aliases) == 0 {
			return errors.New("no host path aliases found: you can set one by `gogh config host-path-alias set <host> <alias>`")
		}
		hosts := make([]string, 0, len(aliases))
		for host := range aliases {
			hosts = append(hosts, host)
		}
		slices.Sort(hosts)
		for _, host := range hosts {
			fmt.Printf("%s: %s\n", host, aliases[host])
		}
		return nil
	}
}

func NewHostPathAliasSetCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "set <host> <alias>",
		Short: "Set a host path alias",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			alias := args[1]
			if err := validateHostPathAlias(host, alias); err != nil {
				return err
			}
			aliases := svc.WorkspaceService.GetHostPathAliases()
			if aliases == nil {
				aliases = workspace.HostPathAliases{}
			}
			aliases[host] = alias
			svc.WorkspaceService.SetHostPathAliases(aliases)
			log.FromContext(cmd.Context()).Infof("Set host path alias: %s -> %s", host, alias)
			return nil
		},
	}, nil
}

func NewHostPathAliasRemoveCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:     "remove <host>",
		Aliases: []string{"rm"},
		Short:   "Remove a host path alias",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := args[0]
			aliases := svc.WorkspaceService.GetHostPathAliases()
			if _, ok := aliases[host]; !ok {
				return fmt.Errorf("host path alias for %q not found", host)
			}
			delete(aliases, host)
			svc.WorkspaceService.SetHostPathAliases(aliases)
			log.FromContext(cmd.Context()).Infof("Removed host path alias: %s", host)
			return nil
		},
	}, nil
}

func validateHostPathAlias(host, alias string) error {
	if host == "" {
		return errors.New("host must not be empty")
	}
	if alias == "" {
		return errors.New("alias must not be empty")
	}
	if host == alias {
		return errors.New("alias must be different from host")
	}
	if strings.ContainsAny(alias, `/\`) {
		return errors.New("alias must be a single path element")
	}
	return nil
}
