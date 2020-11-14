package cli

import "github.com/morikuni/aec"

func Question(text string) string {
	return aec.GreenF.With(aec.Bold).Apply("?") + " " + aec.Bold.Apply(text)
}
