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

func main() {
	ctx := logger.NewLogger()

	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %s\n", err)
		os.Exit(1)
	}
	tokens, err := config.LoadTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load tokens: %s\n", err)
		os.Exit(1)
	}
	defaults, err := config.LoadFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load flags: %s\n", err)
		os.Exit(1)
	}

	cmd := cli.NewApp(conf, tokens, defaults)
	cmd.Version = fmt.Sprintf("%s-%s (%s)", version, commit, date)

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.FromContext(ctx).Error(err.Error())
		os.Exit(1)
	}
}
