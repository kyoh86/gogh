// +build !debug

package env

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var (
	// ConfigurationsFile is the path for configurations file
	ConfigurationsFile = os.ExpandEnv("${HOME}/.gogh")
)

const (
	// AppDescription of GOGH
	AppDescription = "GO GitHub client"
	// LogLevel to output
	LogLevel = logrus.InfoLevel
)
