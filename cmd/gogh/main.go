package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/infra/config"
	"github.com/kyoh86/gogh/v3/infra/logger"
	"github.com/kyoh86/gogh/v3/ui/cli"
)

var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func loadConfigOrExit[T any](name string, loader func() (T, error)) T {
	v, err := loader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load %s: %s\n", name, err)
		os.Exit(1)
	}
	return v
}

func main() {
	ctx := logger.NewLogger()

	conf := loadConfigOrExit("config", config.LoadConfig)
	tokens := loadConfigOrExit("tokens", config.LoadTokens)
	defaults := loadConfigOrExit("flags", config.LoadFlags)

	cmd := cli.NewApp(conf, tokens, defaults)
	cmd.Version = fmt.Sprintf("%s-%s (%s)", version, commit, date)

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.FromContext(ctx).Error(err.Error())
		os.Exit(1)
	}
}
