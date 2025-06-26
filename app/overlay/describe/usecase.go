package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// Overlay represents the overlay type
type Overlay = overlay.Overlay

// JSONUsecase represents the use case for showing overlays in JSON format
type JSONUsecase struct {
	enc *json.Encoder
}

// NewJSONUsecase creates a new use case for showing overlays in JSON format
func NewJSONUsecase(writer io.Writer) *JSONUsecase {
	return &JSONUsecase{enc: json.NewEncoder(writer)}
}

// Execute executes the use case to show a overlay in JSON format
func (uc *JSONUsecase) Execute(ctx context.Context, s Overlay) error {
	return uc.enc.Encode(map[string]any{
		"id":            s.ID(),
		"name":          s.Name(),
		"relative_path": s.RelativePath(),
	})
}

// OnelineUsecase represents the use case for showing overlays in a single line format
type OnelineUsecase struct {
	writer io.Writer
}

// NewOnelineUsecase creates a new use case for showing overlays in text format
func NewOnelineUsecase(writer io.Writer) *OnelineUsecase {
	return &OnelineUsecase{writer: writer}
}

// Execute executes the use case to show a overlay in a single line format
func (uc *OnelineUsecase) Execute(ctx context.Context, s overlay.Overlay) error {
	_, err := fmt.Fprintf(uc.writer, "[%s] %s for %s\n", s.ID()[:8], s.Name(), s.RelativePath())
	return err
}

// JSONWithContentUsecase represents the use case for showing overlays in a single line format
type JSONWithContentUsecase struct {
	overlayService overlay.OverlayService
	enc            *json.Encoder
}

// NewJSONWithContentUsecase creates a new use case for showing overlays in text format
func NewJSONWithContentUsecase(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *JSONWithContentUsecase {
	return &JSONWithContentUsecase{
		overlayService: overlayService,
		enc:            json.NewEncoder(writer),
	}
}

// Execute executes the use case to show a overlay in a single line format
func (uc *JSONWithContentUsecase) Execute(ctx context.Context, s overlay.Overlay) error {
	src, err := uc.overlayService.Open(ctx, s.ID())
	if err != nil {
		return fmt.Errorf("open overlay content: %w", err)
	}
	defer src.Close()
	content, err := io.ReadAll(src)
	if err != nil {
		return fmt.Errorf("read overlay content: %w", err)
	}
	if err := uc.enc.Encode(map[string]any{
		"id":            s.ID(),
		"name":          s.Name(),
		"relative_path": s.RelativePath(),
		"content":       string(content),
	}); err != nil {
		return fmt.Errorf("encode overlay: %w", err)
	}
	return err
}

// DetailUsecase represents the use case for showing overlays in a single line format
type DetailUsecase struct {
	overlayService overlay.OverlayService
	writer         io.Writer
}

// NewDetailUsecase creates a new use case for showing overlays in text format
func NewDetailUsecase(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *DetailUsecase {
	return &DetailUsecase{
		overlayService: overlayService,
		writer:         writer,
	}
}

// Execute executes the use case to show a overlay in a single line format
func (uc *DetailUsecase) Execute(ctx context.Context, s overlay.Overlay) error {
	cnt, err := uc.overlayService.Open(ctx, s.ID())
	if err != nil {
		return fmt.Errorf("open overlay content: %w", err)
	}
	defer cnt.Close()
	fmt.Fprintf(uc.writer, "ID: %s\n", s.ID())
	fmt.Fprintf(uc.writer, "Name: %s\n", s.Name())
	fmt.Fprintf(uc.writer, "Relative path: %s\n", s.RelativePath())
	fmt.Fprintln(uc.writer, "Content<<<"+strings.Repeat("-", 20))
	if _, err := io.Copy(uc.writer, cnt); err != nil {
		return fmt.Errorf("read overlay content: %w", err)
	}
	fmt.Fprintln(uc.writer)
	fmt.Fprintln(uc.writer, ">>>Content", strings.Repeat("-", 20))
	return err
}
