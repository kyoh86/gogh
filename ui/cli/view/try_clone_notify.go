package view

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v4/app/try_clone"
)

// TryCloneNotify is a wrapper for the TryCloneNotify function to log the status.
func TryCloneNotify(
	ctx context.Context,
	notify try_clone.Notify,
) try_clone.Notify {
	return func(n try_clone.Status) error {
		switch n {
		case try_clone.StatusEmpty:
			log.FromContext(ctx).Info("created empty repository")
		case try_clone.StatusRetry:
			log.FromContext(ctx).Info("waiting the remote repository is ready")
		}
		if notify != nil {
			return notify(n)
		}
		return nil
	}
}
