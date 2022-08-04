package database

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *measurementDB

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// os.Exit skips defer calls, so we need to call another function
	code, err := run(m)
	if err != nil {
		log.Info().Err(err).Msg("errors in test")
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	// db, err := sql.Open("sqlite3", "file:../test/test.db?cache=shared")
	db, err := gorm.Open(sqlite.Open("../test/test.db"), &gorm.Config{})
	if err != nil {
		log.Info().Err(err).Msg("could not connect to database")
		return -1, err
	}

	// create the database and the table as a base for every test case
	// this is done here, so the test can run independently
	testDB = &measurementDB{db: db}
	err = testDB.SetupMeasurements()
	if err != nil {
		log.Info().Err(err).Msg("could not setup the database")
		return -1, err
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
