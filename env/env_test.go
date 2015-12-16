package env

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestEnv(t *testing.T) {
	if ConfigurationsFile != "./.gogh" {
		t.Error("ConfigurationsFile is not for debug")
	}
	if AppDescription != "GO GitHub client (debug)" {
		t.Error("AppDescription is not for debug")
	}
	if LogLevel != logrus.DebugLevel {
		t.Error("LogLevel is not for debug")
	}
}
