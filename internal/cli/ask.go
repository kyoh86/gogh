package cli

import (
	"fmt"

	"github.com/kyoh86/ask"
)

func AskPassword(msg string) (string, error) {
	var password string
	if err := ask.Hidden(true).Message(Question(msg)).StringVar(&password).Do(); err != nil {
		return "", fmt.Errorf("asking password: %w", err)
	}
	return password, nil
}
