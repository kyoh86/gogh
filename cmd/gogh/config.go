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

var conf struct {
	Roots []expandedPath `yaml:"roots"`
}

func defaultRoot() string {
	return conf.Roots[0].expanded
}

func roots() []string {
	list := make([]string, 0, len(conf.Roots))
	for _, r := range conf.Roots {
		list = append(list, r.expanded)
	}
	return list
}

func setDefaultRoot(r string) error {
	rootList := make([]expandedPath, 0, len(conf.Roots))
	newDefault, err := parsePath(r)
	if err != nil {
		return err
	}
	rootList = append(rootList, newDefault)
	for _, root := range conf.Roots {
		if root.raw == r {
			continue
		}
		rootList = append(rootList, root)
	}
	conf.Roots = rootList
	return nil
}

func addRoots(rootList []string) error {
	for _, r := range rootList {
		newRoot, err := parsePath(r)
		if err != nil {
			return err
		}
		conf.Roots = append(conf.Roots, newRoot)
	}
	return nil
}

func removeRoot(r string) {
	rootList := make([]expandedPath, 0, len(conf.Roots))
	for _, root := range conf.Roots {
		if root.raw == r || root.expanded == r {
			continue
		}
		rootList = append(rootList, root)
	}
	conf.Roots = rootList
}

//go:embed config_template.txt
var configTemplate string

var configCommand = &cobra.Command{
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
		if err := t.Execute(&w, map[string]any{
			"configFilePath":      configFilePath,
			"tokensFilePath":      tokensFilePath,
			"defaultFlagFilePath": defaultFlagFilePath,
			"roots":               roots(),
			"tokens":              tokens.Entries(),
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
	facadeCommand.AddCommand(configCommand)
}
