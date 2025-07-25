package describe_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/describe"
	"github.com/kyoh86/gogh/v4/core/hook"
)

func TestJSONUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	hookUUID := uuid.New()
	hookID := hookUUID.String()
	hookName := "test-hook"
	repoPattern := "github.com/owner/*"
	triggerEvent := string(hook.EventPostClone)
	operationType := string(hook.OperationTypeOverlay)
	operationID := uuid.New()

	h := hook.ConcreteHook(hookUUID, hookName, repoPattern, triggerEvent, operationType, operationID)

	var buf bytes.Buffer
	uc := testtarget.NewJSONUsecase(&buf)

	err := uc.Execute(ctx, h)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["id"] != hookID {
		t.Errorf("Expected id %s, got %v", hookID, result["id"])
	}
	if result["name"] != hookName {
		t.Errorf("Expected name %s, got %v", hookName, result["name"])
	}
	if result["repo_pattern"] != repoPattern {
		t.Errorf("Expected repo_pattern %s, got %v", repoPattern, result["repo_pattern"])
	}
	if result["trigger_event"] != triggerEvent {
		t.Errorf("Expected trigger_event %s, got %v", triggerEvent, result["trigger_event"])
	}
	if result["operation_type"] != operationType {
		t.Errorf("Expected operation_type %s, got %v", operationType, result["operation_type"])
	}
	if result["operation_id"] != operationID.String() {
		t.Errorf("Expected operation_id %s, got %v", operationID.String(), result["operation_id"])
	}
}

func TestOnelineUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name            string
		hookName        string
		repoPattern     string
		triggerEvent    string
		operationType   string
		operationID     uuid.UUID
		expectedPattern string
	}{
		{
			name:            "With repo pattern",
			hookName:        "test-hook",
			repoPattern:     "github.com/owner/*",
			triggerEvent:    string(hook.EventPostClone),
			operationType:   string(hook.OperationTypeOverlay),
			operationID:     uuid.New(),
			expectedPattern: "github.com/owner/*",
		},
		{
			name:            "Without repo pattern",
			hookName:        "global-hook",
			repoPattern:     "",
			triggerEvent:    string(hook.EventPostFork),
			operationType:   string(hook.OperationTypeScript),
			operationID:     uuid.New(),
			expectedPattern: "*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hookUUID := uuid.New()
			hookID := hookUUID.String()

			h := hook.ConcreteHook(hookUUID, tc.hookName, tc.repoPattern, tc.triggerEvent, tc.operationType, tc.operationID)

			var buf bytes.Buffer
			uc := testtarget.NewOnelineUsecase(&buf)

			err := uc.Execute(ctx, h)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			output := buf.String()
			expectedParts := []string{
				"[" + hookID[:8] + "]",
				tc.hookName,
				"for repos(" + tc.expectedPattern + ")",
				"@" + tc.triggerEvent,
				tc.operationType + "(" + tc.operationID.String()[:8] + ")",
			}

			for _, part := range expectedParts {
				if !strings.Contains(output, part) {
					t.Errorf("Expected output to contain '%s', but it doesn't: %s", part, output)
				}
			}
		})
	}
}

func TestOnelineUsecase_Execute_AllEventTypes(t *testing.T) {
	ctx := context.Background()

	events := []hook.Event{
		hook.EventPostClone,
		hook.EventPostFork,
		hook.EventPostCreate,
	}

	for _, event := range events {
		t.Run(string(event), func(t *testing.T) {
			hookUUID := uuid.New()
			h := hook.ConcreteHook(
				hookUUID,
				"test-hook",
				"github.com/test/*",
				string(event),
				string(hook.OperationTypeOverlay),
				uuid.New(),
			)

			var buf bytes.Buffer
			uc := testtarget.NewOnelineUsecase(&buf)

			err := uc.Execute(ctx, h)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, "@"+string(event)) {
				t.Errorf("Expected output to contain '@%s', but it doesn't: %s", event, output)
			}
		})
	}
}

func TestOnelineUsecase_Execute_AllOperationTypes(t *testing.T) {
	ctx := context.Background()

	operations := []struct {
		opType hook.OperationType
		opID   uuid.UUID
	}{
		{hook.OperationTypeOverlay, uuid.New()},
		{hook.OperationTypeScript, uuid.New()},
	}

	for _, op := range operations {
		t.Run(string(op.opType), func(t *testing.T) {
			hookUUID := uuid.New()
			h := hook.ConcreteHook(
				hookUUID,
				"test-hook",
				"github.com/test/*",
				string(hook.EventPostClone),
				string(op.opType),
				op.opID,
			)

			var buf bytes.Buffer
			uc := testtarget.NewOnelineUsecase(&buf)

			err := uc.Execute(ctx, h)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			output := buf.String()
			expected := string(op.opType) + "(" + op.opID.String()[:8] + ")"
			if !strings.Contains(output, expected) {
				t.Errorf("Expected output to contain '%s', but it doesn't: %s", expected, output)
			}
		})
	}
}
