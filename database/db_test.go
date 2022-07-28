package database

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var testDB *measurementDB

func TestMain(m *testing.M) {
	// os.Exit skips defer calls, so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	db, err := sql.Open("sqlite3", "file:../test/test.db?cache=shared")
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	// create the database and the table as a base for every test case
	// this is done here, so the test can run independently
	testDB = &measurementDB{db: db}
	err = testDB.SetupMeasurements()
	if err != nil {
		return -1, fmt.Errorf("could not setup the database: %w", err)
	}

	// truncates all test data after the tests are run
	defer func() {
		for _, t := range []string{"measurements"} {
			_, _ = testDB.db.Exec(fmt.Sprintf("DELETE FROM %s", t))
			_, _ = testDB.db.Exec(fmt.Sprintf("DROP TABLE %s", t))
		}

		testDB.db.Close()
	}()

	return m.Run(), nil
}
