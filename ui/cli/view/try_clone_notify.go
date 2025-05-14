package view

import (
	"context"

	"github.com/apex/log"
	"github.com/kyoh86/gogh/v3/app/service"
)

func TryCloneNotify(
	ctx context.Context,
	notify service.TryCloneNotify,
) service.TryCloneNotify {
	return func(n service.TryCloneStatus) error {
		switch n {
		case service.TryCloneStatusEmpty:
			log.FromContext(ctx).Info("created empty repository")
		case service.TryCloneStatusRetry:
			log.FromContext(ctx).Info("waiting the remote repository is ready")
		}
		if notify != nil {
			return notify(n)
		}
		return nil
	}
}
