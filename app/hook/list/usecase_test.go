package list_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	testtarget "github.com/kyoh86/gogh/v4/app/hook/list"
	"github.com/kyoh86/gogh/v4/core/hook"
	"github.com/kyoh86/gogh/v4/core/hook_mock"
	"go.uber.org/mock/gomock"
)

func TestUsecase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Pass through hooks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHookService := hook_mock.NewMockHookService(ctrl)
		hooks := []hook.Hook{
			hook.ConcreteHook(
				uuid.New(),
				"test-hook-1",
				"github.com/owner/*",
				string(hook.EventPostClone),
				string(hook.OperationTypeOverlay),
				uuid.New(),
			),
			hook.ConcreteHook(
				uuid.New(),
				"test-hook-2",
				"github.com/org/**",
				string(hook.EventPostCreate),
				string(hook.OperationTypeScript),
				uuid.New(),
			),
		}
		mockHookService.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
			for _, h := range hooks {
				if !yield(h, nil) {
					return
				}
			}
		})

		uc := testtarget.NewUsecase(mockHookService)
		got := make([]hook.Hook, 0, len(hooks))
		for h, err := range uc.Execute(ctx) {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got = append(got, h)
		}

		if len(got) != len(hooks) {
			t.Fatalf("got %d hooks, want %d", len(got), len(hooks))
		}
		for i := range got {
			if got[i].ID() != hooks[i].ID() {
				t.Fatalf("hook[%d] id = %s, want %s", i, got[i].ID(), hooks[i].ID())
			}
		}
	})

	t.Run("Pass through errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHookService := hook_mock.NewMockHookService(ctrl)
		mockHookService.EXPECT().List().Return(func(yield func(hook.Hook, error) bool) {
			yield(nil, errors.New("list error"))
		})

		uc := testtarget.NewUsecase(mockHookService)
		for _, err := range uc.Execute(ctx) {
			if err == nil {
				t.Fatal("expected error")
			}
			return
		}
		t.Fatal("expected iterator result")
	})
}
