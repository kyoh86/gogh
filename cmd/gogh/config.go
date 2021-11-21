package main

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

var config struct {
	Roots []expandedPath `yaml:"roots"`
}

func defaultRoot() string {
	return config.Roots[0].expanded
}

func roots() []string {
	list := make([]string, 0, len(config.Roots))
	for _, r := range config.Roots {
		list = append(list, r.expanded)
	}
	return list
}

func setDefaultRoot(r string) error {
	rootList := make([]expandedPath, 0, len(config.Roots))
	newDefault, err := parsePath(r)
	if err != nil {
		return err
	}
	rootList = append(rootList, newDefault)
	for _, root := range config.Roots {
		if root.raw == r {
			continue
		}
		rootList = append(rootList, root)
	}
	config.Roots = rootList
	return nil
}

func addRoots(rootList []string) error {
	for _, r := range rootList {
		newRoot, err := parsePath(r)
		if err != nil {
			return err
		}
		config.Roots = append(config.Roots, newRoot)
	}
	return nil
}

func removeRoot(r string) {
	rootList := make([]expandedPath, 0, len(config.Roots))
	for _, root := range config.Roots {
		if root.raw == r || root.expanded == r {
			continue
		}
		rootList = append(rootList, root)
	}
	config.Roots = rootList
}

//go:embed config_template.txt
var configTemplate string

var configCommand = &cobra.Command{
	Use:     "config",
	Short:   "Manage config",
	Aliases: []string{"conf", "setting", "context"},
	RunE: func(cmd *cobra.Command, _ []string) error {
		logger := log.FromContext(cmd.Context())
		t, err := template.New("gogh context").Parse(configTemplate)
		if err != nil {
			logger.Error("[Bug] Failed to parse template string")
			return nil
		}
		var serverIdentifiers []string
		{
			list, err := servers.List()
			if err != nil {
				return fmt.Errorf("listup servers: %w", err)
			}
			for _, s := range list {
				serverIdentifiers = append(serverIdentifiers, s.String())
			}
		}
		var defaultFlags string
		{
			var w strings.Builder
			if err := yaml.NewEncoder(&w).Encode(defaultFlag); err != nil {
				logger.Error("[Bug] Failed to build default flag map")
				return nil
			}
			defaultFlags = regexp.MustCompile("(?m)^").ReplaceAllString(w.String(), "  ")
		}
		var w strings.Builder
		if err := t.Execute(&w, map[string]interface{}{
			"configFilePath":      configFilePath,
			"serversFilePath":     serversFilePath,
			"defaultFlagFilePath": defaultFlagFilePath,
			"roots":               roots(),
			"servers":             serverIdentifiers,
			"defaultFlags":        defaultFlags,
		}); err != nil {
			log.FromContext(cmd.Context()).Error("[Bug] Failed to execute template string")
			return nil
		}
		fmt.Println(w.String())
		return nil
	},
}

func init() {
	setup()
	facadeCommand.AddCommand(configCommand)
}
