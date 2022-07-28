// Package database implements the DB access for the CRUD operations.
package database

import (
	"MeasurementWeb/utils"
	"database/sql"
	"errors"
	"fmt"

	"github.com/teris-io/shortid"
)

// MeasurementDo is the data object for the database.
type MeasurementDo struct {
	ID        string
	Timestamp string
	Sensor    string
	Value     float64
	Unit      string
}

var emptyMeasurement = MeasurementDo{}

// MeasurementDB provides the CRUD storage functionality.
type MeasurementDB interface {
	// SetupMeasurements creates the database and the tables.
	SetupMeasurements() (err error)

	// CreateMeasurement adds a new measurement to the table. The ID must be empty and
	// will be defined within this function.
	CreateMeasurement(m MeasurementDo) (entity MeasurementDo, err error)

	// ReadMeasurement reads an measurement
	ReadMeasurement(id string) (entity MeasurementDo, err error)
	UpdateMeasurement(m MeasurementDo) (entity MeasurementDo, err error)
	DeleteMeasurement(id string) (err error)
}

// TODO: comments
// exported 'constructor'
func NewMeasurementDB(config *utils.Conf) (db *measurementDB, err error) {
	database, err := sql.Open("sqlite3", config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	db = &measurementDB{db: database}

	return db, nil
}

// unexported SQL implementation
type measurementDB struct {
	db *sql.DB
}

func (s *measurementDB) SetupMeasurements() (err error) {
	const sql = "CREATE TABLE \"measurements\" ( \"ID\" TEXT UNIQUE, \"Timestamp\" TEXT, \"Sensor\" TEXT, \"Value\" NUMERIC, \"Unit\" TEXT, PRIMARY KEY(\"ID\") )"

	// use config

	// exit if table already exits
	_, e := s.db.Exec("SELECT * FROM measurements LIMIT 1")
	if e == nil {
		return nil
	}

	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("setup: could not begin transaction: %w", e)
		return err
	}

	stmt, e := tx.Prepare(sql)

	if e != nil {
		err = fmt.Errorf("setup: could not prepare transaction: %w", e)
		return err
	}

	defer stmt.Close()

	_, e = stmt.Exec()

	if e != nil {
		err = fmt.Errorf("setup: could not execute statement: %w", e)
		return err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("setup: could not commit transaction: %w", e)
		return err
	}

	return nil
}

func (s *measurementDB) CreateMeasurement(m MeasurementDo) (entity MeasurementDo, err error) {
	const sql = "INSERT INTO measurements (ID, Timestamp, Sensor, Value, Unit) VALUES (?, ?, ?, ?, ?)"

	if m.ID != "" {
		err = errors.New("the ID must be empty")
		return emptyMeasurement, err
	}

	var e error
	m.ID, e = shortid.Generate()

	if e != nil {
		err = fmt.Errorf("insert: could not generate id: %w", e)
		return emptyMeasurement, err
	}

	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("insert: could not begin transaction: %w", e)
		return emptyMeasurement, err
	}

	stmt, e := tx.Prepare(sql)

	if e != nil {
		err = fmt.Errorf("insert: could not prepare transaction: %w", e)
		return emptyMeasurement, err
	}

	defer stmt.Close()

	_, e = stmt.Exec(m.ID, m.Timestamp, m.Sensor, m.Value, m.Unit)

	if e != nil {
		err = fmt.Errorf("insert: could not execute statement: %w", e)
		return emptyMeasurement, err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("insert: could not commit transaction: %w", e)
		return emptyMeasurement, err
	}

	return m, nil
}

func (s *measurementDB) ReadMeasurement(id string) (entity MeasurementDo, err error) {
	stmt, err := s.db.Prepare("SELECT ID, Timestamp, Sensor, Value, Unit from measurements WHERE ID = ?")

	if err != nil {
		return emptyMeasurement, err
	}

	measurement := emptyMeasurement

	e := stmt.QueryRow(id).Scan(&measurement.ID, &measurement.Timestamp, &measurement.Sensor, &measurement.Value, &measurement.Unit)

	if e != nil {
		if e == sql.ErrNoRows {
			return emptyMeasurement, nil
		}
		err = fmt.Errorf("read: could not query: %w", e)
		return emptyMeasurement, err
	}

	defer stmt.Close()

	return measurement, nil
}

func (s *measurementDB) UpdateMeasurement(m MeasurementDo) (entity MeasurementDo, err error) {
	const sql = "UPDATE measurements SET Timestamp = ?, Sensor = ?, Value = ?, Unit = ? WHERE ID = ?"

	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("update: could not begin transaction: %w", e)
		return emptyMeasurement, err
	}

	stmt, e := tx.Prepare(sql)

	if e != nil {
		err = fmt.Errorf("update: could not prepare transaction: %w", e)
		return emptyMeasurement, err
	}

	defer stmt.Close()

	_, e = stmt.Exec(m.Timestamp, m.Sensor, m.Value, m.Unit, m.ID)

	if e != nil {
		err = fmt.Errorf("update: could not execute statement: %w", e)
		return emptyMeasurement, err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("update: could not commit transaction: %w", e)
		return emptyMeasurement, err
	}

	return m, nil
}

func (s *measurementDB) DeleteMeasurement(id string) (err error) {
	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("delete: could not begin transaction: %w", e)
		return err
	}

	stmt, e := tx.Prepare("DELETE from measurements where id = ?")

	if e != nil {
		err = fmt.Errorf("delete: could not prepare transaction: %w", e)
		return err
	}

	defer stmt.Close()

	_, e = stmt.Exec(id)

	if e != nil {
		err = fmt.Errorf("delete: could not execute statement: %w", e)
		return err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("delete: could not commit transaction: %w", e)
		return err
	}

	return nil
}
