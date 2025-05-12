package commands

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/apex/log"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

//go:embed config_template.txt
var configTemplate string

func NewConfigCommand(svc *ServiceSet) *cobra.Command {
	return &cobra.Command{
		Use:     "config",
		Short:   "Show configurations",
		Aliases: []string{"conf", "setting", "context"},
		RunE: func(cmd *cobra.Command, _ []string) error {
			//TODO: support set-default-name subcommand
			//TODO: support set-flag subcommand
			//TODO: new command: migrate files to new format
			logger := log.FromContext(cmd.Context())
			t, err := template.New("gogh context").Parse(configTemplate)
			if err != nil {
				logger.WithError(err).Error("[Bug] Failed to parse template string")
				return nil
			}

			flags, err := encodeYAML(svc.flags)
			if err != nil {
				logger.Error("[Bug] Failed to load flags")
				return nil
			}
			var w strings.Builder
			if err := t.Execute(&w, map[string]any{
				"defaultNameFilePath": svc.defaultNameSource,
				"tokensFilePath":      svc.tokenSource,
				"flagsFilePath":       svc.flagsSource,
				"workspaceFilePath":   svc.workspaceSource,
				"roots":               svc.workspaceService.GetRoots(),
				"defaultHost":         svc.defaultNameService.GetDefaultHost(),
				"defaultNames":        svc.defaultNameService.GetMap(),
				"tokens":              svc.tokenService.Entries(),
				"flags":               flags,
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
