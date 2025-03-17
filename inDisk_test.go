package persistence

import (
	"testing"

	"github.com/fulldump/biff"
)

func TestInDisk(t *testing.T) {

	p, err := NewInDisk[TestItem](t.TempDir())

	biff.AssertNil(err)

	SuitePersistencer(p, t)
	SuiteOptimisticLocking(p, t)
}
