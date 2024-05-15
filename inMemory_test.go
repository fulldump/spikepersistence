package persistence

import "testing"

func TestInMemory(t *testing.T) {

	p := NewInMemory()

	SuitePersistencer(p, t)
}
