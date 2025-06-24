package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kyoh86/gogh/v4/core/script"
)

// Script represents the script type
type Script = script.Script

// UseCaseJSON represents the use case for showing scripts in JSON format
type UseCaseJSON struct {
	enc *json.Encoder
}

// NewUseCaseJSON creates a new use case for showing scripts in JSON format
func NewUseCaseJSON(writer io.Writer) *UseCaseJSON {
	return &UseCaseJSON{enc: json.NewEncoder(writer)}
}

// Execute executes the use case to show a script in JSON format
func (uc *UseCaseJSON) Execute(ctx context.Context, s Script) error {
	return uc.enc.Encode(map[string]any{
		"id":         s.ID(),
		"name":       s.Name(),
		"created_at": s.CreatedAt(),
		"updated_at": s.UpdatedAt(),
	})
}

// UseCaseOneLine represents the use case for showing scripts in a single line format
type UseCaseOneLine struct {
	writer io.Writer
}

// NewUseCaseOneline creates a new use case for showing overlays in text format
func NewUseCaseOneLine(writer io.Writer) *UseCaseOneLine {
	return &UseCaseOneLine{writer: writer}
}

// Execute executes the use case to show a script in a single line format
func (uc *UseCaseOneLine) Execute(ctx context.Context, s script.Script) error {
	_, err := fmt.Fprintf(uc.writer, "[%s] %s @ %s\n", s.ID()[:8], s.Name(), s.UpdatedAt().Format("2006-01-02 15:04:05"))
	return err
}

// UseCaseJSONWithSource represents the use case for showing scripts in a single line format
type UseCaseJSONWithSource struct {
	scriptService script.ScriptService
	enc           *json.Encoder
}

// NewUseCaseJSONWithSource creates a new use case for showing overlays in text format
func NewUseCaseJSONWithSource(
	scriptService script.ScriptService,
	writer io.Writer,
) *UseCaseJSONWithSource {
	return &UseCaseJSONWithSource{
		scriptService: scriptService,
		enc:           json.NewEncoder(writer),
	}
}

// Execute executes the use case to show a script in a single line format
func (uc *UseCaseJSONWithSource) Execute(ctx context.Context, s script.Script) error {
	src, err := uc.scriptService.Open(ctx, s.ID())
	if err != nil {
		return fmt.Errorf("open script source: %w", err)
	}
	defer src.Close()
	source, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read script source: %w", err)
	}
	if err := uc.enc.Encode(map[string]any{
		"id":         s.ID(),
		"name":       s.Name(),
		"created_at": s.CreatedAt(),
		"updated_at": s.UpdatedAt(),
		"source":     string(source),
	}); err != nil {
		return fmt.Errorf("encode script: %w", err)
	}
	return err
}

// UseCaseDetail represents the use case for showing scripts in a single line format
type UseCaseDetail struct {
	scriptService script.ScriptService
	writer        io.Writer
}

// NewUseCaseDetail creates a new use case for showing overlays in text format
func NewUseCaseDetail(
	scriptService script.ScriptService,
	writer io.Writer,
) *UseCaseDetail {
	return &UseCaseDetail{
		scriptService: scriptService,
		writer:        writer,
	}
}

// Execute executes the use case to show a script in a single line format
func (uc *UseCaseDetail) Execute(ctx context.Context, s script.Script) error {
	src, err := uc.scriptService.Open(ctx, s.ID())
	if err != nil {
		return fmt.Errorf("open script source: %w", err)
	}
	defer src.Close()
	fmt.Fprintf(uc.writer, "ID: %s\n", s.ID())
	fmt.Fprintf(uc.writer, "Name: %s\n", s.Name())
	fmt.Fprintf(uc.writer, "Created at: %s\n", s.CreatedAt().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(uc.writer, "Updated at: %s\n", s.UpdatedAt().Format("2006-01-02 15:04:05"))
	fmt.Fprintln(uc.writer, "Source<<<"+strings.Repeat("-", 20))
	if _, err := io.Copy(uc.writer, src); err != nil {
		return fmt.Errorf("read script source: %w", err)
	}
	fmt.Fprintln(uc.writer)
	fmt.Fprintln(uc.writer, ">>>Source", strings.Repeat("-", 20))
	return err
}
