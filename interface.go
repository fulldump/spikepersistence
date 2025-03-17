package persistence

import (
	"context"
	"errors"
)

type Identifier interface {
	GetId() string
	GetVersion() int64
	SetVersion(version int64)
}

type Id struct {
	Id      string `json:"id" bson:"_id"`
	Version int64  `json:"version" bson:"version"`
}

func NewId(id string) *Id {
	return &Id{
		Id:      id,
		Version: 0,
	}
}

func (i *Id) GetId() string {
	return i.Id
}

func (i *Id) GetVersion() int64 {
	return i.Version
}

func (i *Id) SetVersion(version int64) {
	i.Version = version
}

var ErrVersionGone = errors.New("version gone")

type Persistencer[T Identifier] interface {
	List(ctx context.Context) ([]*T, error)
	Put(ctx context.Context, item *T) error
	Get(ctx context.Context, id string) (*T, error)
	Delete(ctx context.Context, id string) error
}
