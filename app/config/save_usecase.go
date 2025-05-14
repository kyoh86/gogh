package config

import (
	"context"

	"github.com/kyoh86/gogh/v3/core/store"
)

type SaverService interface {
	Save(ctx context.Context, force bool) error
}

type Saver[T store.Content] struct {
	store   store.Store[T]
	content T
}

func NewSaver[T store.Content](store store.Store[T], content T) *Saver[T] {
	return &Saver[T]{
		store:   store,
		content: content,
	}
}

func (u *Saver[T]) Save(ctx context.Context, force bool) error {
	if !u.content.HasChanges() {
		return nil
	}
	if err := u.store.Save(ctx, u.content, force); err != nil {
		return err
	}
	u.content.MarkSaved()
	return nil
}
