package persistence

import (
	"context"
	"sync"
)

type InMemory struct {
	Items map[string]*ItemWithId
	mutex sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		Items: map[string]*ItemWithId{},
	}
}

func (f *InMemory) List(ctx context.Context) ([]*ItemWithId, error) {

	f.mutex.RLock()
	defer f.mutex.RUnlock()

	result := make([]*ItemWithId, len(f.Items))

	i := -1
	for _, f := range f.Items {
		i++
		result[i] = f
	}

	return result, nil
}

func (f *InMemory) Put(ctx context.Context, item *ItemWithId) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.Items[item.Id] = item
	return nil
}

func (f *InMemory) Get(ctx context.Context, id string) (*ItemWithId, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.Items[id], nil
}

func (f *InMemory) Delete(ctx context.Context, id string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	delete(f.Items, id)
	return nil
}
