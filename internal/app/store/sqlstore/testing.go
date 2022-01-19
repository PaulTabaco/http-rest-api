package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

//// Tests Helper

// TestDB ...
func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper() // Say - this is testing method

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) { // callback to clean new tables
		if len(tables) > 0 {
			db.Exec("TRUNCATE  %s CASCADE", strings.Join(tables, ", "))
		}

		if len(tables) > 0 {
			fmt.Println("Tables not emtied !!!!!!!")
		}

		db.Close()
	}

}