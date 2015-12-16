package cl

import (
	"fmt"
	"os"

	"github.com/bgentry/speakeasy"
)

var (
	stderr = os.Stderr
	stdin  = os.Stdin
)

// Ask the user to enter. prompt is a string to display before the user's input.
func Ask(prompt string) (string, error) {
	var input string
	if _, err := fmt.Fprintf(stderr, "%s: ", prompt); err != nil {
		return "", err
	}
	if _, err := fmt.Fscanln(stdin, &input); err != nil {
		return "", err
	}
	return input, nil
}

// Secret the user to enter a password with input hidden. prompt is a string to
// display before the user's input. Returns the provided password, or an error
// if the command failed.
func Secret(prompt string) (string, error) {
	return speakeasy.FAsk(stderr, prompt+": ")
}
