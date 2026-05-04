package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cli/safeexec"
	"github.com/kballard/go-shellquote"
)

func lookPath(name string) ([]string, error) {
	exe, err := safeexec.LookPath(name)
	if err != nil {
		return nil, err
	}
	return []string{exe}, nil
}

func edit(fileName string) error {
	editor := os.Getenv("GOGH_EDITOR")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			return fmt.Errorf("GOGH_EDITOR and EDITOR environment variable is not set")
		}
	}

	words, err := shellquote.Split(editor)
	if err != nil {
		return err
	}
	words = append(words, fileName)
	editorExe, err := lookPath(words[0])
	if err != nil {
		return err
	}
	words = append(editorExe, words[1:]...)

	cmdEdit := exec.Command(words[0], words[1:]...)
	cmdEdit.Env = os.Environ()
	cmdEdit.Stdin = os.Stdin
	cmdEdit.Stdout = os.Stdout
	cmdEdit.Stderr = os.Stderr
	return cmdEdit.Run()
}
