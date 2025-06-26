package view

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/clone/try"
)

// TryCloneNotify is a wrapper for the TryCloneNotify function to log the status.
func TryCloneNotify(
	ctx context.Context,
	notify try.Notify,
) try.Notify {
	return func(n try.Status) error {
		switch n {
		case try.StatusEmpty:
			log.FromContext(ctx).Info("created empty repository")
		case try.StatusRetry:
			log.FromContext(ctx).Info("waiting the remote repository is ready")
		}
		if notify != nil {
			return notify(n)
		}
		return nil
	}
}
