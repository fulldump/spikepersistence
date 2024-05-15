package persistence

import (
	"context"
	"fmt"
	"sync"
	"testing"

	. "github.com/fulldump/biff"
)

func SuitePersistencer(p Persistencer, t *testing.T) {

	ctx := context.Background()

	t.Run("List empty", func(t *testing.T) {
		listResult, listErr := p.List(ctx)
		AssertNil(listErr)
		AssertEqual(len(listResult), 0)
	})

	item1 := &ItemWithId{
		Id: "1",
		Item: Item{
			Title: "Title 1",
		},
	}

	t.Run("Insert one", func(t *testing.T) {
		putErr := p.Put(ctx, item1)
		AssertNil(putErr)
	})

	t.Run("Retrieve one", func(t *testing.T) {
		getResult, getErr := p.Get(ctx, "1")
		AssertNil(getErr)
		AssertEqual(getResult, item1)
	})

	t.Run("List one", func(t *testing.T) {
		listResult, listErr := p.List(ctx)

		AssertNil(listErr)
		AssertEqual(len(listResult), 1)
		AssertEqual(listResult[0], item1)
	})

	item1updated := &ItemWithId{
		Id: "1",
		Item: Item{
			Title: "Title 1 updated",
		},
	}

	t.Run("Update one", func(t *testing.T) {
		putErr := p.Put(ctx, item1updated)
		AssertNil(putErr)

		t.Run("Check list length = 1 and value is one", func(t *testing.T) {
			listResult, _ := p.List(ctx)
			AssertEqual(len(listResult), 1)
			AssertEqual(listResult[0], item1updated)
		})

	})

	item2 := &ItemWithId{
		Id: "2",
		Item: Item{
			Title: "Title 2",
		},
	}

	t.Run("Insert two", func(t *testing.T) {
		putErr := p.Put(ctx, item2)
		AssertNil(putErr)

		t.Run("Check list length = 2", func(t *testing.T) {
			listResult, listErr := p.List(ctx)

			AssertNil(listErr)
			AssertEqual(len(listResult), 2)
		})

	})

	t.Run("Delete one", func(t *testing.T) {
		err := p.Delete(ctx, "1")
		AssertNil(err)

		t.Run("Check list length = 1 and value is two", func(t *testing.T) {
			listResult, listErr := p.List(ctx)

			AssertNil(listErr)
			AssertEqual(len(listResult), 1)
			AssertEqual(listResult[0], item2)
		})

		t.Run("Check one does not longer exist", func(t *testing.T) {
			getResult, getErr := p.Get(ctx, "1")
			AssertNil(getErr)
			AssertNil(getResult)
		})
	})

	t.Run("Concurrency", func(t *testing.T) {
		w := &sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			w.Add(1)

			id := fmt.Sprintf("item-%d", i)

			p.Put(ctx, &ItemWithId{
				Id: id,
				Item: Item{
					Title: id,
				},
			})

			go func() {
				p.Delete(ctx, id)
				w.Done()
			}()
		}

		w.Wait()
	})
}
