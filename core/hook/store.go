package hook

import "github.com/kyoh86/gogh/v4/core/store"

// HookStore persists hooks
type HookStore interface {
	store.Store[HookService]
}
