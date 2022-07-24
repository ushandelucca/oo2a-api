package storage

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var testDB struct {
	db *sql.DB
}

func TestMain(m *testing.M) {
	// os.Exit skips defer calls
	// so we need to call another function
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	// pseudo-code, some implementation excluded:
	//
	// 1. create test.db if it does not exist
	// 2. run our DDL statements to create the required tables if they do not exist
	// 3. run our tests
	// 4. truncate the test db tables

	//info, err := os.Stat(filename)
	//if os.IsNotExist(err) {
	//	return false
	//}

	testDB.db, err = sql.Open("sqlite3", "file:../test.db?cache=shared")
	if err != nil {
		return -1, fmt.Errorf("could not connect to database: %w", err)
	}

	// truncates all test data after the tests are run
	defer func() {
		for _, t := range []string{"measurements"} {
			_, _ = testDB.db.Exec(fmt.Sprintf("DELETE FROM %s", t))
			// _, _ = testDB.db.Exec(fmt.Sprintf("DROP TABLE %s", t))
		}

		testDB.db.Close()
	}()

	return m.Run(), nil
}
