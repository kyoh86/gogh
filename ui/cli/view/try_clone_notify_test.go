package view

import (
	"bytes"
	"context"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/kyoh86/gogh/v4/app/clone/try"
	"github.com/stretchr/testify/assert"
)

func TestTryCloneNotify(t *testing.T) {
	t.Run("StatusEmpty", func(t *testing.T) {
		// Set up buffer to capture log output
		var buf bytes.Buffer
		logger := &log.Logger{
			Handler: text.New(&buf),
			Level:   log.InfoLevel,
		}
		ctx := log.NewContext(context.Background(), logger)

		// Mock notify function to verify call
		var called bool
		mockNotify := func(status try.Status) error {
			called = true
			assert.Equal(t, try.StatusEmpty, status)
			return nil
		}

		// Execute the function under test
		notify := TryCloneNotify(ctx, mockNotify)
		err := notify(try.StatusEmpty)

		// Verify
		assert.NoError(t, err)
		assert.True(t, called, "original notify function should be called")
		assert.Contains(t, buf.String(), "Created empty repository")
	})

	t.Run("StatusRetry", func(t *testing.T) {
		// Set up buffer to capture log output
		var buf bytes.Buffer
		logger := &log.Logger{
			Handler: text.New(&buf),
			Level:   log.InfoLevel,
		}
		ctx := log.NewContext(context.Background(), logger)

		// Mock notify function to verify call
		var called bool
		mockNotify := func(status try.Status) error {
			called = true
			assert.Equal(t, try.StatusRetry, status)
			return nil
		}

		// Execute the function under test
		notify := TryCloneNotify(ctx, mockNotify)
		err := notify(try.StatusRetry)

		// Verify
		assert.NoError(t, err)
		assert.True(t, called, "original notify function should be called")
		assert.Contains(t, buf.String(), "Waiting the remote repository is ready")
	})

	t.Run("NilNotify", func(t *testing.T) {
		ctx := context.Background()

		// Test with nil notify function
		notify := TryCloneNotify(ctx, nil)
		err := notify(try.StatusEmpty)

		// Verify no error occurs
		assert.NoError(t, err)
	})
}
