package persistence

import "testing"

func TestInMemory(t *testing.T) {

	p := NewInMemory[TestItem]()

	SuitePersistencer(p, t)
	SuiteOptimisticLocking(p, t)
}
