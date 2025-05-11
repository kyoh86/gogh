package commands

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/apex/log"
	"github.com/goccy/go-yaml"
	"github.com/kyoh86/gogh/v3/infra/config"
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
			logger := log.FromContext(cmd.Context())
			t, err := template.New("gogh context").Parse(configTemplate)
			if err != nil {
				logger.Error("[Bug] Failed to parse template string")
				return nil
			}
			tokensFilePath, err := config.TokensPathV0()
			if err != nil {
				return fmt.Errorf("failed to get tokens file path: %w", err)
			}
			flagsFilePath, err := config.FlagsPath()
			if err != nil {
				return fmt.Errorf("failed to get flags file path: %w", err)
			}

			defaultFlags, err := encodeYAML(svc.defaults)
			if err != nil {
				logger.Error("[Bug] Failed to build default flag map")
				return nil
			}
			var w strings.Builder
			if err := t.Execute(&w, map[string]any{
				"tokensFilePath":      tokensFilePath,
				"defaultFlagFilePath": flagsFilePath,
				"roots":               svc.workspaceService.GetRoots(),
				"tokens":              svc.tokenService.Entries(),
				"defaultFlags":        defaultFlags,
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
