package persistence

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestInInception(t *testing.T) {

	collection := "testing-" + uuid.NewString()
	var p *InInceptionDB[TestItem]

	for _, base := range []string{"http://inceptiondb:1212/v1", "http://localhost:1212/v1"} {
		p = NewInInceptionDB[TestItem](&ConfigInceptionDB{
			Base:       base,
			Collection: collection,
		})
		_, err := p.List(context.Background())
		if err == nil {
			break
		}
	}

	SuitePersistencer(p, t)
	SuiteOptimisticLocking(p, t) // working on this!
}
