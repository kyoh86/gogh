package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

// nolint
var (
	version = "snapshot"
	commit  = "snapshot"
	date    = "snapshot"
)

func main() {
	app := kingpin.New("gogh", "GO GitHub project manager").Version(fmt.Sprintf("%s-%s (%s)", version, commit, date)).Author("kyoh86")
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
