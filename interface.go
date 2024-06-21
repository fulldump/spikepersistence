package persistence

import (
	"context"
	"errors"
)

type ItemWithId struct {
	Id      string `json:"id"`
	Item    `bson:",inline"`
	Version int64 `json:"-" bson:"version"`
}

var ErrVersionGone = errors.New("version gone")

type Persistencer interface {
	List(ctx context.Context) ([]*ItemWithId, error)
	Put(ctx context.Context, item *ItemWithId) error
	Get(ctx context.Context, id string) (*ItemWithId, error)
	Delete(ctx context.Context, id string) error
}
