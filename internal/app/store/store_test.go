package store_test

import (
	"os"
	"testing"
)

var ( // a global var in this pacage - store_test
	databaseURL string
)

func TestMain(m *testing.M) { // TestMaim predefined, runs one time
	databaseURL = os.Getenv("DATABASE_URL") // try look in enviropment
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=restapi_test sslmode=disable"
	}
	os.Exit(m.Run())
}
