package gogh

import "context"

type Remote interface {
	List(ctx context.Context, params *RemoteListParam) ([]Project, error)
	Create(ctx context.Context, description Description) (Project, error)
	Remove(ctx context.Context, description Description) error
}

type RemoteListParam struct {
	Query string
}
