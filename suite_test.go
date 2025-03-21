package persistence

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/fulldump/biff"
)

type TestItem struct {
	*Id         `bson:",inline"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Subitems    []*SubItem `json:"subitems"`
	Counter     int        `json:"counter"`
}

var _ Identifier = TestItem{}

type SubItem struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
}

func SuitePersistencer(p Persistencer[TestItem], t *testing.T) {

	ctx := context.Background()

	t.Run("List empty", func(t *testing.T) {
		listResult, listErr := p.List(ctx)
		AssertNil(listErr)
		AssertEqual(len(listResult), 0)
	})

	item1 := &TestItem{
		Id:    NewId("1"),
		Title: "Title 1",
	}

	t.Run("Insert one", func(t *testing.T) {
		putErr := p.Put(ctx, item1)
		AssertNil(putErr)
	})

	t.Run("Retrieve one", func(t *testing.T) {
		getResult, getErr := p.Get(ctx, "1")
		AssertNil(getErr)
		AssertEqual(getResult.Id, item1.Id)
		AssertEqual(getResult.Id, item1.Id)
		item1 = getResult
	})

	t.Run("List one", func(t *testing.T) {
		listResult, listErr := p.List(ctx)

		AssertNil(listErr)
		AssertEqual(len(listResult), 1)
		AssertEqual(listResult[0].Id, item1.Id)
	})

	item1.Title = "Title 1 updated"

	t.Run("Update one", func(t *testing.T) {
		putErr := p.Put(ctx, item1)
		AssertNil(putErr)

		t.Run("Check list length = 1 and value is one", func(t *testing.T) {
			listResult, _ := p.List(ctx)
			AssertEqual(len(listResult), 1)
			AssertEqual(listResult[0].Id, item1.Id)
		})

	})

	item2 := &TestItem{
		Id:    NewId("2"),
		Title: "Title 2",
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
			AssertEqual(listResult[0].Id, item2.Id)
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

			p.Put(ctx, &TestItem{
				Id:    NewId(id),
				Title: id,
			})

			go func() {
				p.Delete(ctx, id)
				w.Done()
			}()
		}

		w.Wait()
	})

}

func SuiteOptimisticLocking(p Persistencer[TestItem], t *testing.T) {

	ctx := context.Background()

	t.Run("Concurrency - optimistic", func(t *testing.T) {

		err := p.Put(ctx, &TestItem{
			Id:      NewId("1"),
			Title:   "Title 1",
			Counter: 0,
		})
		AssertNil(err)

		itema, errGet := p.Get(ctx, "1")
		AssertNil(errGet)
		itema.Counter++

		itemb, errGet := p.Get(ctx, "1")
		AssertNil(errGet)
		itemb.Counter++

		erra := p.Put(ctx, itema)
		AssertNil(erra)

		errb := p.Put(ctx, itemb)
		AssertNotNil(errb)

		finalItem, err := p.Get(ctx, "1")
		AssertNil(err)
		fmt.Println(finalItem)
	})

	t.Run("Concurrency - optimistic2", func(t *testing.T) {
		w := &sync.WaitGroup{}

		err := p.Put(ctx, &TestItem{
			Id:      NewId("optimistic-2"),
			Title:   "Title 1",
			Counter: 0,
		})
		AssertNil(err)

		collisions := int32(0)
		workers := 50
		for i := 0; i < workers; i++ {
			w.Add(1)

			go func() {
				defer w.Done()

				for {
					item, _ := p.Get(ctx, "optimistic-2")
					item.Counter++

					errPut := p.Put(ctx, item)
					if errPut == ErrVersionGone {
						atomic.AddInt32(&collisions, 1)
						time.Sleep(time.Duration(rand.IntN(workers)) * time.Millisecond)
						continue
					}
					return
				}

			}()
		}

		w.Wait()

		fmt.Println("COLLISIONS:", collisions)

		finalItem, err := p.Get(ctx, "optimistic-2")
		AssertNil(err)
		AssertEqual(finalItem.Counter, workers)
	})

}
