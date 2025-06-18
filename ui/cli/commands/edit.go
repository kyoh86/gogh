package commands

import (
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

func edit(editor, fileName string) error {
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
