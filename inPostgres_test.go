package persistence

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/fulldump/biff"
)

func TestInPostgres(t *testing.T) {

	host := ""
	for _, h := range []string{"localhost", "postgres"} {

		db, err := sql.Open("postgres", "host="+h+" port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable")
		if err != nil {
			continue
		}

		err = db.Ping() // check if db exists
		if err != nil {
			continue
		}

		host = h
		break
	}

	if host == "" {
		t.Skipf("Postgres not available")
	}

	dbname := "test" + strconv.FormatInt(time.Now().UnixNano(), 10)

	p, err := NewInPostgres[Item]("host=" + host + " port=5432 user=postgres password=mysecretpassword dbname=" + dbname + " sslmode=disable")
	biff.AssertNil(err)

	SuitePersistencer(p, t)
	SuiteOptimisticLocking(p, t)
}
