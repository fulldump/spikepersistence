package persistence

import (
	"testing"
)

func TestInInception(t *testing.T) {

	p := NewInInceptionDB[Item](&ConfigInceptionDB{
		Base:       "http://localhost:1212/v1",
		Collection: "test",
	})

	SuitePersistencer(p, t)
	// SuiteOptimisticLocking(p, t) // working on this!
}
