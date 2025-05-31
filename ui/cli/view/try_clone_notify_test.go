package view

import (
	"bytes"
	"context"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/kyoh86/gogh/v4/app/try_clone"
	"github.com/stretchr/testify/assert"
)

func TestTryCloneNotify(t *testing.T) {
	t.Run("StatusEmpty", func(t *testing.T) {
		// ログ出力をキャプチャするためのバッファを設定
		var buf bytes.Buffer
		logger := &log.Logger{
			Handler: text.New(&buf),
			Level:   log.InfoLevel,
		}
		ctx := log.NewContext(context.Background(), logger)

		// 呼び出し確認用のモックnotify関数
		var called bool
		mockNotify := func(status try_clone.Status) error {
			called = true
			assert.Equal(t, try_clone.StatusEmpty, status)
			return nil
		}

		// テスト対象の関数を実行
		notify := TryCloneNotify(ctx, mockNotify)
		err := notify(try_clone.StatusEmpty)

		// 検証
		assert.NoError(t, err)
		assert.True(t, called, "元のnotify関数が呼び出されること")
		assert.Contains(t, buf.String(), "created empty repository")
	})

	t.Run("StatusRetry", func(t *testing.T) {
		// ログ出力をキャプチャするためのバッファを設定
		var buf bytes.Buffer
		logger := &log.Logger{
			Handler: text.New(&buf),
			Level:   log.InfoLevel,
		}
		ctx := log.NewContext(context.Background(), logger)

		// 呼び出し確認用のモックnotify関数
		var called bool
		mockNotify := func(status try_clone.Status) error {
			called = true
			assert.Equal(t, try_clone.StatusRetry, status)
			return nil
		}

		// テスト対象の関数を実行
		notify := TryCloneNotify(ctx, mockNotify)
		err := notify(try_clone.StatusRetry)

		// 検証
		assert.NoError(t, err)
		assert.True(t, called, "元のnotify関数が呼び出されること")
		assert.Contains(t, buf.String(), "waiting the remote repository is ready")
	})

	t.Run("NilNotify", func(t *testing.T) {
		ctx := context.Background()

		// nilのnotify関数でテスト
		notify := TryCloneNotify(ctx, nil)
		err := notify(try_clone.StatusEmpty)

		// エラーが発生しないことを確認
		assert.NoError(t, err)
	})
}
