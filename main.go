package main

import (
	"context"
	"os"

	"github.com/kyoh86/gogh/internal/cmd"
	"github.com/kyoh86/gogh/internal/cmd/pr"
	"github.com/kyoh86/gogh/internal/repo"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("gogh", "GitHub CLI Client").Author("kyoh86")
	commands := map[string]cmd.Command{}

	var directory string
	wd, err := os.Getwd()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get working directory")
	}
	app.Flag("directory", "Run as if git was started in <path> instead of the current working directory.").Short('C').Default(wd).StringVar(&directory)

	r, err := repo.Open(directory)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to open working git repository")
	}

	var debug bool
	app.Flag("debug", "Print debug logs").Hidden().BoolVar(&debug)
	logrus.SetLevel(logrus.InfoLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.SetOutput(os.Stderr)

	prCmd := app.Command("pull-request", "Access for pull-requests").Alias("pr").Alias("pulls")
	prListCmd := prCmd.Command("list", "List up pull-requests").Alias("ls")
	prCreateCmd := prCmd.Command("create", "Create a pull-request").Alias("new").Default()

	commands[prListCmd.FullCommand()] = pr.ListCommand(prListCmd, r)
	commands[prCreateCmd.FullCommand()] = pr.CreateCommand(prCreateCmd, r)

	full := kingpin.MustParse(app.Parse(os.Args[1:]))
	if err := commands[full](context.Background()); err != nil {
		logrus.Fatal(err)
	}
}
