//go:build cgo
// +build cgo

package sql

import (
	"database/sql"
	"log/slog"
	"testing"

	sqlite3 "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	logger := slog.New(slog.DiscardHandler)

	errCheck := func(err error) bool {
		sqlErr, ok := err.(sqlite3.Error)
		if !ok {
			return false
		}
		return sqlErr.ExtendedCode == sqlite3.ErrConstraintUnique
	}

	var sqliteMigrations []migration
	for _, m := range migrations {
		if m.flavor == nil || m.flavor == &flavorSQLite3 {
			sqliteMigrations = append(sqliteMigrations, m)
		}
	}

	c := &conn{db, &flavorSQLite3, logger, errCheck}
	for _, want := range []int{len(sqliteMigrations), 0} {
		got, err := c.migrate()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("expected %d migrations, got %d", want, got)
		}
	}
}
