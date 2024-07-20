package task

import "context"

type Repository interface {
	ListActive(ctx context.Context) ([]Entity, error)
	ListDone(ctx context.Context) ([]Entity, error)
	Create(ctx context.Context, task Entity) error
	Done(ctx context.Context, id string) error
	Update(ctx context.Context, id string, task Entity) error
	Delete(ctx context.Context, id string) error
}
