package database

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	// db, err := sql.Open("sqlite3", "file:../test/test.db?cache=shared")
	db, err := gorm.Open(sqlite.Open("../test/test.db"), &gorm.Config{})
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
		for _, t := range []string{"measurement_dos"} {
			testDB.db.Exec(fmt.Sprintf("DELETE FROM %s", t))
			testDB.db.Exec(fmt.Sprintf("DROP TABLE %s", t))
		}

		// testDB.db.Close()
	}()

	return m.Run(), nil
}
