package overlay_describe

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

// UseCaseJSON represents the use case for showing overlays in JSON format
type UseCaseJSON struct {
	enc *json.Encoder
}

// NewUseCaseJSON creates a new use case for showing overlays in JSON format
func NewUseCaseJSON(writer io.Writer) *UseCaseJSON {
	return &UseCaseJSON{enc: json.NewEncoder(writer)}
}

// Execute executes the use case to show a overlay in JSON format
func (uc *UseCaseJSON) Execute(ctx context.Context, s Overlay) error {
	return uc.enc.Encode(map[string]any{
		"id":            s.ID(),
		"name":          s.Name(),
		"relative_path": s.RelativePath(),
	})
}

// UseCaseOneLine represents the use case for showing overlays in a single line format
type UseCaseOneLine struct {
	writer io.Writer
}

// NewUseCaseOneline creates a new use case for showing overlays in text format
func NewUseCaseOneLine(writer io.Writer) *UseCaseOneLine {
	return &UseCaseOneLine{writer: writer}
}

// Execute executes the use case to show a overlay in a single line format
func (uc *UseCaseOneLine) Execute(ctx context.Context, s overlay.Overlay) error {
	_, err := fmt.Fprintf(uc.writer, "* [%s] %s for %s\n", s.ID()[:8], s.Name(), s.RelativePath())
	return err
}

// UseCaseJSONWithContent represents the use case for showing overlays in a single line format
type UseCaseJSONWithContent struct {
	overlayService overlay.OverlayService
	enc            *json.Encoder
}

// NewUseCaseJSONWithContent creates a new use case for showing overlays in text format
func NewUseCaseJSONWithContent(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *UseCaseJSONWithContent {
	return &UseCaseJSONWithContent{
		overlayService: overlayService,
		enc:            json.NewEncoder(writer),
	}
}

// Execute executes the use case to show a overlay in a single line format
func (uc *UseCaseJSONWithContent) Execute(ctx context.Context, s overlay.Overlay) error {
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

// UseCaseDetail represents the use case for showing overlays in a single line format
type UseCaseDetail struct {
	overlayService overlay.OverlayService
	writer         io.Writer
}

// NewUseCaseDetail creates a new use case for showing overlays in text format
func NewUseCaseDetail(
	overlayService overlay.OverlayService,
	writer io.Writer,
) *UseCaseDetail {
	return &UseCaseDetail{
		overlayService: overlayService,
		writer:         writer,
	}
}

// Execute executes the use case to show a overlay in a single line format
func (uc *UseCaseDetail) Execute(ctx context.Context, s overlay.Overlay) error {
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
