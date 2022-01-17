package store

import (
	"fmt"
	"strings"
	"testing"
)

//// Tests Helper

// TestStore ...
func TestStore(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper() // Say - this is testing method
	config := NewConfig()
	config.DatabaseURL = databaseURL
	s := New(config)

	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {

		// check if we have new tables and remove all
		if len(tables) > 0 {
			if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}

		s.Close()
	}
}
