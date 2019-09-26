package command_test

import (
	"github.com/kyoh86/gogh/command"
	"github.com/kyoh86/gogh/config"
)

func ExampleConfigGetAll() {
	if err := command.ConfigGetAll(&config.Config{
		GitHub: config.GitHubConfig{
			Token: "tokenx1",
			Host:  "hostx1",
			User:  "kyoh86",
		},
		Log: config.LogConfig{
			Level:        "trace",
			Date:         config.TrueOption,
			Time:         config.FalseOption,
			MicroSeconds: config.TrueOption,
			LongFile:     config.TrueOption,
			ShortFile:    config.TrueOption,
			UTC:          config.TrueOption,
		},
		VRoot: []string{"/foo", "/bar"},
	}); err != nil {
		panic(err)
	}
	// Unordered output:
	// root = /foo:/bar
	// github.host = hostx1
	// github.user = kyoh86
	// github.token = tokenx1
	// log.level = trace
	// log.date = yes
	// log.time = no
	// log.microseconds = yes
	// log.longfile = yes
	// log.shortfile = yes
	// log.utc = yes
}