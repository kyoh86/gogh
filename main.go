package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Sirupsen/logrus"
	"github.com/doloopwhile/logrusltsv"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gh"
	"github.com/kyoh86/gogh/gh/pr"
)

func main() {
	logrus.SetLevel(env.LogLevel)
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrusltsv.Formatter{})

	app := kingpin.New(env.AppName, env.AppDescription).Author(env.Author)
	prCmd := app.Command("pull-request", "Access for pull-requests").Alias("pr").Alias("pulls")
	prlsCmd := prCmd.Command("list", "List up pull-requests").Alias("ls")

	commands := map[string]gh.Command{}
	commands[prlsCmd.FullCommand()] = pr.List(prlsCmd)

	full := kingpin.MustParse(app.Parse(os.Args[1:]))
	if err := commands[full](); err != nil {
		logrus.Fatal(err)
	}
}
