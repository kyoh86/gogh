package command

import (
	"log"
	"os"
)

func InitLog() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
}
