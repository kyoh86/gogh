package overlay

// Overlay represents the metadata for an overlay entry.
type Overlay struct {
	ID           string
	Name         string `json:"name"`         // Name of the overlay
	RelativePath string `json:"relativePath"` // Relative path in the repository where the overlay file will be placed
}
