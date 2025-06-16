package script

import (
	"time"

	"github.com/google/uuid"
)

type Script struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (h *Script) init() {
	if h.ID == "" {
		h.ID = uuid.NewString()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
		h.UpdatedAt = h.CreatedAt
	}
}

func (h *Script) update() {
	h.UpdatedAt = time.Now()
}
