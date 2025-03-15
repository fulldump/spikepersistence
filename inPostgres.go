package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	_ "github.com/lib/pq"
)

type InPostgres[T any] struct {
	connection string
	db         *sql.DB
}

func NewInPostgres[T any](connection string) (*InPostgres[T], error) {

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err // can not reach postgres, retry?
	}

	err = db.Ping() // check if db exists
	if err != nil {

		// Try to connect and create database
		fields := parseConnection(connection)
		dbname := fields["dbname"]
		fields["dbname"] = "postgres"
		connectionPostgres := connectionToString(fields)

		dbPostgres, err := sql.Open("postgres", connectionPostgres)
		if err != nil {
			return nil, err // could not connect as postgres
		}

		_, err = dbPostgres.Exec("create database " + dbname)
		if err != nil {
			return nil, err // could not create database
		}

		// Connect again with previous connection string
		db, err = sql.Open("postgres", connection)
		if err != nil {
			return nil, err // can not reach postgres, retry?
		}

		err = db.Ping() // check if db exists
		if err != nil {
			return nil, err // could not connecto to new database
		}
	}

	// ensure table `items`
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS "items" (
		    id       VARCHAR(36) PRIMARY KEY,
		    record   JSONB,
		    version  bigint
		);
	`)
	if err != nil {
		return nil, err // could not create database
	}

	return &InPostgres[T]{
		db:         db,
		connection: connection,
	}, nil
}

func connectionToString(fields map[string]string) string {
	pairs := []string{}

	for k, v := range fields {
		pairs = append(pairs, k+"="+v)
	}

	return strings.Join(pairs, " ")
}

func parseConnection(connection string) map[string]string {

	result := map[string]string{}

	for _, pair := range strings.Split(connection, " ") {
		parts := strings.SplitN(pair, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := ""
		if len(parts) > 1 {
			value = strings.TrimSpace(parts[1])
		}
		result[key] = value
	}

	return result
}

func (f *InPostgres[T]) List(ctx context.Context) ([]*ItemWithId[T], error) {

	rows, err := f.db.QueryContext(ctx, `SELECT id, record, version FROM "items";`)
	if err != nil {
		return nil, err
	}

	result := []*ItemWithId[T]{}
	for rows.Next() {
		id := []byte{}
		record := []byte{}
		version := int64(0)
		err := rows.Scan(&id, &record, &version)
		if err != nil {
			return nil, err
		}

		item := &ItemWithId[T]{
			Id: string(id),
			// Item:    T{},
			Version: version,
		}
		err = json.Unmarshal(record, &item.Item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (f *InPostgres[T]) Put(ctx context.Context, item *ItemWithId[T]) error {

	itemJson, err := json.Marshal(item.Item)
	if err != nil {
		return err
	}

	result, err := f.db.ExecContext(ctx, `
		INSERT INTO "items" (id, record, version) VALUES ($1, $2::jsonb, $4)
		ON CONFLICT (ID)
		DO UPDATE SET record = $2, version = $4 WHERE items.version = $3
	`, item.Id, string(itemJson), item.Version, item.Version+1)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrVersionGone
	}

	item.Version++

	return nil
}

func (f *InPostgres[T]) Get(ctx context.Context, id string) (*ItemWithId[T], error) {

	row := f.db.QueryRowContext(ctx, `
		SELECT  record, version FROM "items" WHERE id = $1;
	`, id)

	record := []byte{}
	version := int64(0)
	err := row.Scan(&record, &version)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	item := &ItemWithId[T]{
		Id:      id,
		Version: version,
	}
	err = json.Unmarshal(record, &item.Item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (f *InPostgres[T]) Delete(ctx context.Context, id string) error {

	_, err := f.db.ExecContext(ctx, `
		DELETE FROM "items" 
		WHERE id = $1;
	`, id)
	if err != nil {
		return err
	}

	return nil
}
