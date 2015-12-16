// +build debug test

package env

import "github.com/Sirupsen/logrus"

const (
	// ConfigurationsFile is the path for configurations file
	ConfigurationsFile = "./.gogh"
	// AppDescription of GOGH
	AppDescription = "GO GitHub client (debug)"
	// LogLevel to output
	LogLevel = logrus.DebugLevel
)
