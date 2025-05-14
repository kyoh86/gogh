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

func NewConfigCommand(_ context.Context, svc *service.ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Show configurations",
		Aliases: []string{"conf", "setting", "context"},
		RunE: func(cmd *cobra.Command, _ []string) error {
			logger := log.FromContext(cmd.Context())
			t, err := template.New("gogh context").Parse(configTemplate)
			if err != nil {
				logger.WithError(err).Error("[Bug] Failed to parse template string")
				return nil
			}

			flags, err := encodeYAML(svc.Flags)
			if err != nil {
				logger.Error("[Bug] Failed to load flags")
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
				log.FromContext(cmd.Context()).Error("[Bug] Failed to execute template string")
				return nil
			}
			fmt.Println(w.String())
			return nil
		},
	}
}

func encodeYAML(v interface{}) (string, error) {
	var w strings.Builder
	if err := yaml.NewEncoder(&w).Encode(v); err != nil {
		return "", err
	}
	return regexp.MustCompile("(?m)^").ReplaceAllString(w.String(), "  "), nil
}
