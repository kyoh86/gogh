package view

import (
	"context"
	"iter"

	"github.com/apex/log"
	"github.com/charmbracelet/huh"
)

func ProcessWithConfirmation[T any](ctx context.Context, seq iter.Seq2[T, error], title func(T) string, process func(entry T) error) error {
	logger := log.FromContext(ctx)
	for entry, err := range seq {
		if err != nil {
			return err
		}
		var selected string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewSelect[string]().
				Title(title(entry)).
				Options(huh.Option[string]{
					Key:   "y",
					Value: "Yes",
				}, huh.Option[string]{
					Key:   "n",
					Value: "No",
				}, huh.Option[string]{
					Key:   "q",
					Value: "Quit",
				}).
				Value(&selected),
		)).Run(); err != nil {
			return err
		}
		switch selected {
		case "Yes", "y":
			if err := process(entry); err != nil {
				return err
			}
		case "No", "n":
			logger.Info("Skipped")
		case "Quit", "q":
			logger.Info("Quit")
			return nil
		}
	}
	return nil
}
