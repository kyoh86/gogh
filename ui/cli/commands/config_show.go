package commands

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/apex/log"
	"github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/app/service"
	"github.com/spf13/cobra"
)

//go:embed config_template.txt
var configTemplate string

func NewConfigShowCommand(_ context.Context, svc *service.ServiceSet) (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "show",
		Short: "Show configurations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger := log.FromContext(cmd.Context())
			t, err := template.New("gogh context").Parse(configTemplate)
			if err != nil {
				return fmt.Errorf("[Bug] invalid template string: %w", err)
			}

			flags, err := encodeYAML(svc.Flags)
			if err != nil {
				logger.Error("[Bug] failed to load flags")
				return nil
			}
			var w strings.Builder
			defaultNameSource, err := svc.DefaultNameStore.Source()
			if err != nil {
				return err
			}
			tokenSource, err := svc.TokenStore.Source()
			if err != nil {
				return err
			}
			flagsSource, err := svc.FlagsStore.Source()
			if err != nil {
				return err
			}
			workspaceSource, err := svc.WorkspaceStore.Source()
			if err != nil {
				return err
			}
			if err := t.Execute(&w, map[string]any{
				"defaultNameSource": defaultNameSource,
				"tokenSource":       tokenSource,
				"flagsSource":       flagsSource,
				"workspaceSource":   workspaceSource,
				"roots":             svc.WorkspaceService.GetRoots(),
				"defaultHost":       svc.DefaultNameService.GetDefaultHost(),
				"defaultNames":      svc.DefaultNameService.GetMap(),
				"tokens":            svc.TokenService.Entries(),
				"flags":             flags,
			}); err != nil {
				return fmt.Errorf("[Bug] failed to execute template: %w", err)
			}
			fmt.Println(w.String())
			return nil
		},
	}, nil
}

func encodeYAML(v any) (string, error) {
	var w strings.Builder
	if err := yaml.NewEncoder(&w).Encode(v); err != nil {
		return "", err
	}
	return regexp.MustCompile("(?m)^").ReplaceAllString(w.String(), "  "), nil
}
