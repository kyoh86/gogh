package overlay

import "github.com/kyoh86/gogh/v4/core/store"

// OverlayStore interface extends store.Store[OverlayService] for persistent overlays.
type OverlayStore interface {
	store.Store[OverlayService]
}
