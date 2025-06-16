package overlay

import (
	"io"

	"github.com/google/uuid"
)

type Entry struct {
	Name         string
	RelativePath string
	Content      io.Reader
}

// Overlay represents the metadata for an overlay entry.
type Overlay interface {
	ID() string
	UUID() uuid.UUID
	Name() string
	RelativePath() string
}

func NewOverlay(entry Entry) Overlay {
	return overlayElement{
		id:           uuid.Must(uuid.NewRandom()),
		name:         entry.Name,
		relativePath: entry.RelativePath,
	}
}

type overlayElement struct {
	id           uuid.UUID
	name         string
	relativePath string
}

func (o overlayElement) ID() string {
	return o.id.String()
}

func (o overlayElement) UUID() uuid.UUID {
	return o.id
}

func (o overlayElement) Name() string {
	return o.name
}

func (o overlayElement) RelativePath() string {
	return o.relativePath
}
