package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kyoh86/gogh/v4/core/extra"
	"github.com/kyoh86/gogh/v4/core/repository"
	"github.com/kyoh86/gogh/v4/core/store"
)

type extraItemEntry struct {
	OverlayID string `toml:"overlay_id"`
	HookID    string `toml:"hook_id"`
}

type extraEntry struct {
	ID         string           `toml:"id"`
	Type       string           `toml:"type"`
	Name       string           `toml:"name,omitempty"`
	Repository string           `toml:"repository,omitempty"`
	Items      []extraItemEntry `toml:"items"`
	Source     string           `toml:"source"`
	CreatedAt  time.Time        `toml:"created_at"`
}

type extraData struct {
	Extra []extraEntry `toml:"extra"`
}

// ExtraStore stores extra configuration
type ExtraStore struct{}

// NewExtraStore creates a new extra store
func NewExtraStore() *ExtraStore {
	return &ExtraStore{}
}

// Source returns the path to the extra store file
func (s *ExtraStore) Source() (string, error) {
	path, err := AppContextPathFunc("GOGH_EXTRA_PATH", os.UserConfigDir, "extra.v4.toml")
	if err != nil {
		return "", fmt.Errorf("search extra path: %w", err)
	}
	return path, nil
}

// Load loads extra configuration
func (s *ExtraStore) Load(ctx context.Context, initial func() extra.ExtraService) (extra.ExtraService, error) {
	src, err := s.Source()
	if err != nil {
		return nil, fmt.Errorf("get extra store source: %w", err)
	}

	data, err := loadTOMLFile[extraData](src)
	if err != nil {
		if os.IsNotExist(err) {
			svc := initial()
			svc.MarkSaved()
			return svc, nil
		}
		return nil, fmt.Errorf("load extra store: %w", err)
	}

	svc := initial()
	parser := repository.NewReferenceParser("", "")

	// Convert loaded data to extras
	extras := make([]*extra.Extra, 0, len(data.Extra))
	for _, entry := range data.Extra {
		sourceRef, err := parser.Parse(entry.Source)
		if err != nil {
			return nil, fmt.Errorf("parsing source reference %q: %w", entry.Source, err)
		}

		// Convert items
		items := make([]extra.Item, len(entry.Items))
		for i, item := range entry.Items {
			items[i] = extra.Item{
				OverlayID: item.OverlayID,
				HookID:    item.HookID,
			}
		}

		var e *extra.Extra
		switch extra.Type(entry.Type) {
		case extra.TypeAuto:
			repoRef, err := parser.Parse(entry.Repository)
			if err != nil {
				return nil, fmt.Errorf("parsing repository reference %q: %w", entry.Repository, err)
			}
			e = extra.NewAutoExtra(entry.ID, *repoRef, *sourceRef, items, entry.CreatedAt)
		case extra.TypeNamed:
			e = extra.NewNamedExtra(entry.ID, entry.Name, *sourceRef, items, entry.CreatedAt)
		default:
			return nil, fmt.Errorf("unknown extra type: %s", entry.Type)
		}
		extras = append(extras, e)
	}

	// Load extras into service
	if err := svc.Load(func(yield func(*extra.Extra, error) bool) {
		for _, e := range extras {
			if !yield(e, nil) {
				return
			}
		}
	}); err != nil {
		return nil, fmt.Errorf("loading extras: %w", err)
	}

	svc.MarkSaved()
	return svc, nil
}

// Save saves extra configuration
func (s *ExtraStore) Save(ctx context.Context, svc extra.ExtraService, force bool) error {
	if !force && !svc.HasChanges() {
		return nil
	}

	data := extraData{
		Extra: []extraEntry{},
	}

	for e, err := range svc.List(ctx) {
		if err != nil {
			return err
		}

		// Convert items
		items := make([]extraItemEntry, 0, len(e.Items()))
		for _, item := range e.Items() {
			items = append(items, extraItemEntry{
				OverlayID: item.OverlayID,
				HookID:    item.HookID,
			})
		}

		entry := extraEntry{
			ID:        e.ID(),
			Type:      string(e.Type()),
			Items:     items,
			Source:    e.Source().String(),
			CreatedAt: e.CreatedAt(),
		}

		switch e.Type() {
		case extra.TypeAuto:
			if repo := e.Repository(); repo != nil {
				entry.Repository = repo.String()
			}
		case extra.TypeNamed:
			entry.Name = e.Name()
		}

		data.Extra = append(data.Extra, entry)
	}

	src, err := s.Source()
	if err != nil {
		return fmt.Errorf("get extra store source: %w", err)
	}

	if err := saveTOMLFile(src, data); err != nil {
		return fmt.Errorf("save extra store: %w", err)
	}

	svc.MarkSaved()
	return nil
}

var _ store.Store[extra.ExtraService] = (*ExtraStore)(nil)
