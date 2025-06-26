package describe_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/overlay/describe"
	"github.com/kyoh86/gogh/v4/core/overlay"
	"github.com/kyoh86/gogh/v4/core/overlay_mock"
	"go.uber.org/mock/gomock"
)

func TestJSONUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	var buf bytes.Buffer
	uc := testtarget.NewJSONUsecase(&buf)

	err := uc.Execute(ctx, o)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != overlayID {
		t.Errorf("Expected id %s, got %v", overlayID, result["id"])
	}
	if result["name"] != overlayName {
		t.Errorf("Expected name %s, got %v", overlayName, result["name"])
	}
	if result["relative_path"] != relativePath {
		t.Errorf("Expected relative_path %s, got %v", relativePath, result["relative_path"])
	}
}

func TestOnelineUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	var buf bytes.Buffer
	uc := testtarget.NewOnelineUsecase(&buf)

	err := uc.Execute(ctx, o)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedPrefix := "[" + overlayID[:8] + "] " + overlayName + " for " + relativePath
	if !strings.HasPrefix(output, expectedPrefix) {
		t.Errorf("Expected output to start with %s, got %s", expectedPrefix, output)
	}
}

func TestJSONWithContentUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"
	overlayContent := "This is the overlay content"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
	mockOverlayService.EXPECT().Open(ctx, overlayID).Return(
		io.NopCloser(strings.NewReader(overlayContent)), nil,
	)

	var buf bytes.Buffer
	uc := testtarget.NewJSONWithContentUsecase(mockOverlayService, &buf)

	err := uc.Execute(ctx, o)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != overlayID {
		t.Errorf("Expected id %s, got %v", overlayID, result["id"])
	}
	if result["name"] != overlayName {
		t.Errorf("Expected name %s, got %v", overlayName, result["name"])
	}
	if result["relative_path"] != relativePath {
		t.Errorf("Expected relative_path %s, got %v", relativePath, result["relative_path"])
	}
	if result["content"] != overlayContent {
		t.Errorf("Expected content %s, got %v", overlayContent, result["content"])
	}
}

func TestJSONWithContentUsecase_Execute_OpenError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
	mockOverlayService.EXPECT().Open(ctx, overlayID).Return(
		nil, errors.New("open error"),
	)

	var buf bytes.Buffer
	uc := testtarget.NewJSONWithContentUsecase(mockOverlayService, &buf)

	err := uc.Execute(ctx, o)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "open overlay content") {
		t.Errorf("Expected error to contain 'open overlay content', got %v", err)
	}
}

func TestDetailUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"
	overlayContent := "This is the overlay content"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
	mockOverlayService.EXPECT().Open(ctx, overlayID).Return(
		io.NopCloser(strings.NewReader(overlayContent)), nil,
	)

	var buf bytes.Buffer
	uc := testtarget.NewDetailUsecase(mockOverlayService, &buf)

	err := uc.Execute(ctx, o)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedLines := []string{
		"ID: " + overlayID,
		"Name: " + overlayName,
		"Relative path: " + relativePath,
		"Content<<<",
		overlayContent,
		">>>Content",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nFull output:\n%s", expected, output)
		}
	}
}

func TestDetailUsecase_Execute_OpenError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
	mockOverlayService.EXPECT().Open(ctx, overlayID).Return(
		nil, errors.New("open error"),
	)

	var buf bytes.Buffer
	uc := testtarget.NewDetailUsecase(mockOverlayService, &buf)

	err := uc.Execute(ctx, o)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "open overlay content") {
		t.Errorf("Expected error to contain 'open overlay content', got %v", err)
	}
}

func TestDetailUsecase_Execute_ReadError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	overlayUUID := uuid.New()
	overlayID := overlayUUID.String()
	overlayName := "test-overlay"
	relativePath := "path/to/file.txt"

	o := overlay.ConcreteOverlay(overlayUUID, overlayName, relativePath)

	// Create a reader that will fail on Read
	failReader := &failingReader{err: errors.New("read error")}

	mockOverlayService := overlay_mock.NewMockOverlayService(ctrl)
	mockOverlayService.EXPECT().Open(ctx, overlayID).Return(
		io.NopCloser(failReader), nil,
	)

	var buf bytes.Buffer
	uc := testtarget.NewDetailUsecase(mockOverlayService, &buf)

	err := uc.Execute(ctx, o)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read overlay content") {
		t.Errorf("Expected error to contain 'read overlay content', got %v", err)
	}
}

// failingReader is a helper type that always fails on Read
type failingReader struct {
	err error
}

func (fr *failingReader) Read(p []byte) (n int, err error) {
	return 0, fr.err
}
