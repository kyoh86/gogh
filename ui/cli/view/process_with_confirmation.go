package view

import (
	"context"
	"errors"
	"iter"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
)

var ErrQuit = errors.New("quit the process")

// ProcessWithConfirmation processes each entry in the sequence with a confirmation prompt.
func ProcessWithConfirmation[T any](ctx context.Context, seq iter.Seq2[T, error], title func(T) string, process func(entry T) error) error {
	logger := log.FromContext(ctx)
	var all bool
	for entry, err := range seq {
		if err != nil {
			return err
		}
		if all {
			if err := process(entry); err != nil {
				return err
			}
			continue
		}
		var selected string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				CharLimit(1).
				Inline(true).
				Title(title(entry) + " ").
				Validate(func(s string) error {
					if s == "y" || s == "n" || s == "q" || s == "a" {
						return nil
					}
					return errors.New("invalid selection")
				}).
				Prompt("(y/n/q/a): ").
				Value(&selected),
		)).Run(); err != nil {
			return err
		}
		switch selected {
		case "a":
			all = true
			fallthrough
		case "y":
			if err := process(entry); err != nil {
				return err
			}
		case "n":
			logger.Info("Skipped")
		case "q":
			logger.Info("Quit")
			return ErrQuit
		}
	}
	return nil
}
