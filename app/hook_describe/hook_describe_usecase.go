package hook_describe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/core/hook"
)

// Hook represents the hook type
type Hook = hook.Hook

// UseCaseJSON represents the use case for showing hooks in JSON format
type UseCaseJSON struct {
	enc *json.Encoder
}

// NewUseCaseJSON creates a new use case for showing hooks in JSON format
func NewUseCaseJSON(writer io.Writer) *UseCaseJSON {
	return &UseCaseJSON{enc: json.NewEncoder(writer)}
}

// Execute executes the use case to show a hook in JSON format
func (uc *UseCaseJSON) Execute(ctx context.Context, s Hook) error {
	return uc.enc.Encode(s) // TODO: Encoding hook
}

// UseCaseOneLine represents the use case for showing hooks in a single line format
type UseCaseOneLine struct {
	writer io.Writer
}

// NewUseCaseOneline creates a new use case for showing overlays in text format
func NewUseCaseOneLine(writer io.Writer) *UseCaseOneLine {
	return &UseCaseOneLine{writer: writer}
}

// Execute executes the use case to show a hook in a single line format
func (uc *UseCaseOneLine) Execute(ctx context.Context, s hook.Hook) error {
	_, err := fmt.Fprintf(
		uc.writer,
		"* [%s] %s for %s @ %s: %s(%s)\n",
		s.ID()[:8],
		s.Name(),
		s.RepoPattern(),
		s.TriggerEvent(),
		s.OperationType(),
		s.OperationID(),
	)
	return err
}
