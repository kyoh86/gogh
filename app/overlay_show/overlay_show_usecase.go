package overlay_show

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/kyoh86/gogh/v4/core/overlay"
)

// JSONUseCase represents the showing use case
type JSONUseCase struct {
	enc *json.Encoder
}

// NewUseCaseJSON creates a new use case for showing overlays in JSON format
func NewUseCaseJSON(writer io.Writer) *JSONUseCase {
	return &JSONUseCase{enc: json.NewEncoder(writer)}
}

// Execute executes the use case to show an overlay in JSON format
func (uc *JSONUseCase) Execute(ctx context.Context, ov *overlay.Overlay) error {
	return uc.enc.Encode(ov)
}

// TextUseCase represents the showing use case
type TextUseCase struct {
	writer io.Writer
	width  int
}

// NewUseCaseText creates a new use case for showing overlays in text format
func NewUseCaseText(
	writer io.Writer,
	width int,
) *TextUseCase {
	return &TextUseCase{
		writer: writer,
		width:  width,
	}
}

// Execute executes the use case to show an overlay in text format
func (uc *TextUseCase) Execute(ctx context.Context, ov *overlay.Overlay) error {
	fmt.Fprintf(uc.writer, "Repository pattern: %s\n", ov.RepoPattern)
	if ov.ForInit {
		fmt.Fprintln(uc.writer, "For Init: Yes")
	}
	fmt.Fprintf(uc.writer, "Overlay path: %s\n", ov.RelativePath)
	return nil
}

// ContentUseCase represents the showing use case
type ContentUseCase struct {
	overlayService overlay.OverlayService
	writer         io.Writer
	width          int
}

// NewUseCaseContent creates a new use case for showing overlay content
func NewUseCaseContent(
	overlayService overlay.OverlayService,
	writer io.Writer,
	width int,
) *ContentUseCase {
	return &ContentUseCase{
		overlayService: overlayService,
		writer:         writer,
		width:          width,
	}
}

// Execute executes the use case to show overlay content
func (uc *ContentUseCase) Execute(ctx context.Context, ov *overlay.Overlay) error {
	fmt.Fprintf(uc.writer, "Repository pattern: %s\n", ov.RepoPattern)
	if ov.ForInit {
		fmt.Fprintln(uc.writer, "For Init: Yes")
	}
	fmt.Fprintf(uc.writer, "Overlay path: %s\n", ov.RelativePath)
	// Open the overlay content
	content, err := uc.overlayService.OpenOverlayContent(ctx, *ov)
	if err != nil {
		return fmt.Errorf("opening overlay for repo-pattern '%s': %w", ov.RepoPattern, err)
	}
	defer content.Close()
	if _, err := io.Copy(uc.writer, content); err != nil {
		return fmt.Errorf("copying overlay content to stdout: %w", err)
	}
	fmt.Fprintln(uc.writer)
	fmt.Fprintf(uc.writer, "%s\n", strings.Repeat("-", uc.width))
	return nil
}
