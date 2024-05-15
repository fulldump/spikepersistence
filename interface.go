package persistence

import (
	"context"
)

type ItemWithId struct {
	Id   string `json:"id"`
	Item `bson:",inline"`
}

type Persistencer interface {
	List(ctx context.Context) ([]*ItemWithId, error)
	Put(ctx context.Context, item *ItemWithId) error
	Get(ctx context.Context, id string) (*ItemWithId, error)
	Delete(ctx context.Context, id string) error
}
