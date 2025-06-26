package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kyoh86/gogh/v4/core/extra"
)

// Extra represents the extra type
type Extra = extra.Extra

// JSONUsecase represents the use case for showing extras in JSON format
type JSONUsecase struct {
	enc *json.Encoder
}

// NewJSONUsecase creates a new use case for showing extras in JSON format
func NewJSONUsecase(writer io.Writer) *JSONUsecase {
	return &JSONUsecase{enc: json.NewEncoder(writer)}
}

// Execute executes the use case to show an extra in JSON format
func (uc *JSONUsecase) Execute(ctx context.Context, e Extra) error {
	data := map[string]any{
		"id":         e.ID(),
		"type":       e.Type(),
		"source":     e.Source().String(),
		"created_at": e.CreatedAt(),
		"items":      e.Items(),
	}

	switch e.Type() {
	case extra.TypeAuto:
		if repo := e.Repository(); repo != nil {
			data["repository"] = repo.String()
		}
	case extra.TypeNamed:
		data["name"] = e.Name()
	}

	return uc.enc.Encode(data)
}

// OnelineUsecase represents the use case for showing extras in a single line format
type OnelineUsecase struct {
	writer io.Writer
}

// NewOnelineUsecase creates a new use case for showing extras in text format
func NewOnelineUsecase(writer io.Writer) *OnelineUsecase {
	return &OnelineUsecase{writer: writer}
}

// Execute executes the use case to show an extra in a single line format
func (uc *OnelineUsecase) Execute(ctx context.Context, e extra.Extra) error {
	var nameOrRepo string
	switch e.Type() {
	case extra.TypeAuto:
		if repo := e.Repository(); repo != nil {
			nameOrRepo = repo.String()
		}
	case extra.TypeNamed:
		nameOrRepo = e.Name()
	}

	_, err := fmt.Fprintf(uc.writer, "[%s] %s %s (%d items)\n",
		e.ID()[:8],
		e.Type(),
		nameOrRepo,
		len(e.Items()),
	)
	return err
}

// DetailUsecase represents the use case for showing extras in detail format
type DetailUsecase struct {
	writer io.Writer
}

// NewDetailUsecase creates a new use case for showing extras in detail format
func NewDetailUsecase(writer io.Writer) *DetailUsecase {
	return &DetailUsecase{writer: writer}
}

// Execute executes the use case to show an extra in detail format
func (uc *DetailUsecase) Execute(ctx context.Context, e extra.Extra) error {
	fmt.Fprintf(uc.writer, "ID: %s\n", e.ID())
	fmt.Fprintf(uc.writer, "Type: %s\n", e.Type())

	switch e.Type() {
	case extra.TypeAuto:
		if repo := e.Repository(); repo != nil {
			fmt.Fprintf(uc.writer, "Repository: %s\n", repo.String())
		}
	case extra.TypeNamed:
		fmt.Fprintf(uc.writer, "Name: %s\n", e.Name())
	}

	fmt.Fprintf(uc.writer, "Source: %s\n", e.Source().String())
	fmt.Fprintf(uc.writer, "Created: %s\n", e.CreatedAt().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(uc.writer, "Items (%d):\n", len(e.Items()))

	for i, item := range e.Items() {
		fmt.Fprintf(uc.writer, "  %d. Overlay: %s\n", i+1, item.OverlayID)
		if item.HookID != "" {
			fmt.Fprintf(uc.writer, "     Hook: %s\n", item.HookID)
		}
	}
	fmt.Fprintln(uc.writer)

	return nil
}
