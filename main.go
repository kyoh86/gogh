package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/Sirupsen/logrus"
	"github.com/doloopwhile/logrusltsv"
	"github.com/kyoh86/gogh/env"
	"github.com/kyoh86/gogh/gh"
	"github.com/kyoh86/gogh/gh/cf"
	"github.com/kyoh86/gogh/gh/pr"
)

func main() {
	logrus.SetLevel(env.LogLevel)
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrusltsv.Formatter{})

	app := kingpin.New(env.AppName, env.AppDescription).Author(env.Author)
	commands := map[string]gh.Command{}

	prCmd := app.Command("pull-request", "Access for pull-requests").Alias("pr").Alias("pulls")
	prListCmd := prCmd.Command("list", "List up pull-requests").Alias("ls")
	prCreateCmd := prCmd.Command("create", "Create a pull-request").Alias("new").Alias("make").Alias("n")

	commands[prListCmd.FullCommand()] = pr.ListCommand(prListCmd)
	commands[prCreateCmd.FullCommand()] = pr.CreateCommand(prCreateCmd)

	cfCmd := app.Command("config", "Configure the gogh").Alias("conf").Alias("configure")
	cfSetCmd := cfCmd.Command("set", "Set a configuration")

	commands[cfSetCmd.FullCommand()] = cf.SetCommand(cfSetCmd)

	full := kingpin.MustParse(app.Parse(os.Args[1:]))
	if err := commands[full](); err != nil {
		logrus.Fatal(err)
	}
}
