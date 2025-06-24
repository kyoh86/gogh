package script_describe_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyoh86/gogh/v4/app/script_describe"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUseCaseJSON_Execute(t *testing.T) {
	ctx := context.Background()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseJSON(&buf)

	err := uc.Execute(ctx, s)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != scriptID {
		t.Errorf("Expected id %s, got %v", scriptID, result["id"])
	}
	if result["name"] != scriptName {
		t.Errorf("Expected name %s, got %v", scriptName, result["name"])
	}
}

func TestUseCaseOneLine_Execute(t *testing.T) {
	ctx := context.Background()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseOneLine(&buf)

	err := uc.Execute(ctx, s)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedPrefix := "[" + scriptID[:8] + "] " + scriptName + " @ " + updatedAt.Format("2006-01-02 15:04:05")
	if !strings.HasPrefix(output, expectedPrefix) {
		t.Errorf("Expected output to start with %s, got %s", expectedPrefix, output)
	}
}

func TestUseCaseJSONWithSource_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)
	scriptSource := "print('Hello, World!')"

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	mockScriptService := script_mock.NewMockScriptService(ctrl)
	mockScriptService.EXPECT().Open(ctx, scriptID).Return(
		io.NopCloser(strings.NewReader(scriptSource)), nil,
	)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseJSONWithSource(mockScriptService, &buf)

	err := uc.Execute(ctx, s)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != scriptID {
		t.Errorf("Expected id %s, got %v", scriptID, result["id"])
	}
	if result["name"] != scriptName {
		t.Errorf("Expected name %s, got %v", scriptName, result["name"])
	}
	if result["source"] != scriptSource {
		t.Errorf("Expected source %s, got %v", scriptSource, result["source"])
	}
}

func TestUseCaseJSONWithSource_Execute_OpenError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	mockScriptService := script_mock.NewMockScriptService(ctrl)
	mockScriptService.EXPECT().Open(ctx, scriptID).Return(
		nil, errors.New("open error"),
	)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseJSONWithSource(mockScriptService, &buf)

	err := uc.Execute(ctx, s)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "open script source") {
		t.Errorf("Expected error to contain 'open script source', got %v", err)
	}
}

func TestUseCaseDetail_Execute(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)
	scriptSource := "print('Hello, World!')"

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	mockScriptService := script_mock.NewMockScriptService(ctrl)
	mockScriptService.EXPECT().Open(ctx, scriptID).Return(
		io.NopCloser(strings.NewReader(scriptSource)), nil,
	)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseDetail(mockScriptService, &buf)

	err := uc.Execute(ctx, s)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	output := buf.String()
	expectedLines := []string{
		"ID: " + scriptID,
		"Name: " + scriptName,
		"Created at: " + createdAt.Format("2006-01-02 15:04:05"),
		"Updated at: " + updatedAt.Format("2006-01-02 15:04:05"),
		"Source<<<",
		scriptSource,
		">>>Source",
	}

	for _, expected := range expectedLines {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', but it doesn't.\nFull output:\n%s", expected, output)
		}
	}
}

func TestUseCaseDetail_Execute_OpenError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	mockScriptService := script_mock.NewMockScriptService(ctrl)
	mockScriptService.EXPECT().Open(ctx, scriptID).Return(
		nil, errors.New("open error"),
	)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseDetail(mockScriptService, &buf)

	err := uc.Execute(ctx, s)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "open script source") {
		t.Errorf("Expected error to contain 'open script source', got %v", err)
	}
}

func TestUseCaseDetail_Execute_ReadError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scriptUUID := uuid.New()
	scriptID := scriptUUID.String()
	scriptName := "test-script"
	createdAt := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2023, 1, 2, 15, 30, 0, 0, time.UTC)

	s := script.ConcreteScript(scriptUUID, scriptName, createdAt, updatedAt)

	// Create a reader that will fail on Read
	failReader := &failingReader{err: errors.New("read error")}

	mockScriptService := script_mock.NewMockScriptService(ctrl)
	mockScriptService.EXPECT().Open(ctx, scriptID).Return(
		io.NopCloser(failReader), nil,
	)

	var buf bytes.Buffer
	uc := script_describe.NewUseCaseDetail(mockScriptService, &buf)

	err := uc.Execute(ctx, s)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read script source") {
		t.Errorf("Expected error to contain 'read script source', got %v", err)
	}
}

// failingReader is a helper type that always fails on Read
type failingReader struct {
	err error
}

func (fr *failingReader) Read(p []byte) (n int, err error) {
	return 0, fr.err
}
