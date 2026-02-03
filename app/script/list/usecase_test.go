package list_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/script/list"
	"github.com/kyoh86/gogh/v4/core/script"
	"github.com/kyoh86/gogh/v4/core/script_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Pass through scripts", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockScriptService := script_mock.NewMockScriptService(ctrl)
		scripts := []script.Script{
			script.ConcreteScript(
				uuid.New(),
				"test-script-1",
				time.Now().Add(-24*time.Hour),
				time.Now(),
			),
			script.ConcreteScript(
				uuid.New(),
				"test-script-2",
				time.Now().Add(-48*time.Hour),
				time.Now().Add(-12*time.Hour),
			),
		}
		mockScriptService.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
			for _, s := range scripts {
				if !yield(s, nil) {
					return
				}
			}
		})

		uc := testtarget.NewUsecase(mockScriptService)
		got := make([]script.Script, 0, len(scripts))
		for s, err := range uc.Execute(ctx) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got = append(got, s)
		}

		if len(got) != len(scripts) {
			t.Fatalf("got %d scripts, want %d", len(got), len(scripts))
		}
		for i := range got {
			if got[i].ID() != scripts[i].ID() {
				t.Fatalf("script[%d] id = %s, want %s", i, got[i].ID(), scripts[i].ID())
			}
		}
	})

	t.Run("Pass through errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockScriptService := script_mock.NewMockScriptService(ctrl)
		mockScriptService.EXPECT().List().Return(func(yield func(script.Script, error) bool) {
			yield(nil, errors.New("list error"))
		})

		uc := testtarget.NewUsecase(mockScriptService)
		for _, err := range uc.Execute(ctx) {
			if err == nil {
				t.Fatal("expected error")
			}
			return
		}
		t.Fatal("expected iterator result")
	})
}
