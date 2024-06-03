package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	_ "github.com/lib/pq"
)

type InPostgres struct {
	connection string
	db         *sql.DB
}

func NewInPostgres(connection string) (*InPostgres, error) {

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
		    record   JSONB
		);
	`)
	if err != nil {
		return nil, err // could not create database
	}

	return &InPostgres{
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

func (f *InPostgres) List(ctx context.Context) ([]*ItemWithId, error) {

	rows, err := f.db.QueryContext(ctx, `SELECT record FROM "items";`)
	if err != nil {
		return nil, err
	}

	result := []*ItemWithId{}
	for rows.Next() {
		record := []byte{}
		err := rows.Scan(&record)
		if err != nil {
			return nil, err
		}

		item := &ItemWithId{}
		err = json.Unmarshal(record, &item)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (f *InPostgres) Put(ctx context.Context, item *ItemWithId) error {

	itemJson, err := json.Marshal(item)
	if err != nil {
		return err
	}

	_, err = f.db.ExecContext(ctx, `
		INSERT INTO "items" (id, record) VALUES ($1, $2::jsonb)
		ON CONFLICT (ID)
		DO UPDATE SET record = $2
	`, item.Id, string(itemJson))
	if err != nil {
		return err
	}

	// todo: check `result.RowsAffected() == 1` ???

	return nil
}

func (f *InPostgres) Get(ctx context.Context, id string) (*ItemWithId, error) {

	row := f.db.QueryRowContext(ctx, `
		SELECT record FROM "items" WHERE id = $1;
	`, id)

	record := []byte{}
	err := row.Scan(&record)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	item := &ItemWithId{}
	err = json.Unmarshal(record, item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (f *InPostgres) Delete(ctx context.Context, id string) error {

	_, err := f.db.ExecContext(ctx, `
		DELETE FROM "items" 
		WHERE id = $1;
	`, id)
	if err != nil {
		return err
	}

	return nil
}
