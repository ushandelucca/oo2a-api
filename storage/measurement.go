// Package storage implements the db access for the CRUD operations.
package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/teris-io/shortid"
)

// Measurement is the transport object for the values
type Measurement struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Sensor    string  `json:"sensor"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

var emptyMeasurement = Measurement{}

// MeasurementStore provides the CRUD storage functionality.
type MeasurementStore interface {
	// SetupMeasurements creates the database and the tables.
	SetupMeasurements() (err error)

	// CreateMeasurement adds a new measurement to the table. The ID must be empty and
	// will be defined within this function.
	CreateMeasurement(m Measurement) (entity Measurement, err error)
	ReadMeasurement(id string) (entity Measurement, err error)
	UpdateMeasurement(m Measurement) (entity Measurement, err error)
	DeleteMeasurement(id string) (err error)
}

// exported 'constructor'
func NewMeasurementStore(db *sql.DB) *measurementStore {
	return &measurementStore{
		db: db,
	}
}

// unexported SQL implementation
type measurementStore struct {
	db *sql.DB
}

func (s *measurementStore) SetupMeasurements() (err error) {
	const sql = "CREATE TABLE \"measurements\" ( \"ID\" TEXT UNIQUE, \"Timestamp\" TEXT, \"Sensor\" TEXT, \"Value\" NUMERIC, \"Unit\" TEXT, PRIMARY KEY(\"ID\") )"

	// use config

	// exit if table already exits
	_, e := s.db.Exec("SELECT * FROM measurements LIMIT 1")
	if e == nil {
		return nil
	}

	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("could not begin transaction: %w", e)
		return err
	}

	stmt, e := tx.Prepare(sql)

	if e != nil {
		err = fmt.Errorf("could not prepare transaction: %w", e)
		return err
	}

	defer stmt.Close()

	_, e = stmt.Exec()

	if e != nil {
		err = fmt.Errorf("could not execute statement: %w", e)
		return err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("could not commit transaction: %w", e)
		return err
	}

	return nil
}

func (s *measurementStore) executeTx(sql string, m Measurement) (err error) {
	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("could not begin transaction: %w", e)
		return err
	}

	stmt, e := tx.Prepare(sql)

	if e != nil {
		err = fmt.Errorf("could not prepare transaction: %w", e)
		return err
	}

	defer stmt.Close()

	_, e = stmt.Exec(m.ID, m.Timestamp, m.Sensor, m.Value, m.Unit)

	if e != nil {
		err = fmt.Errorf("could not execute statement: %w", e)
		return err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("could not commit transaction: %w", e)
		return err
	}

	return nil
}

func (s *measurementStore) CreateMeasurement(m Measurement) (entity Measurement, err error) {
	const sql = "INSERT INTO measurements (ID, Timestamp, Sensor, Value, Unit) VALUES (?, ?, ?, ?, ?)"

	if m.ID != "" {
		err = errors.New("the ID must be empty")
		return emptyMeasurement, err
	}

	var e error
	m.ID, e = shortid.Generate()

	if e != nil {
		err = fmt.Errorf("could not generate id: %w", e)
		return emptyMeasurement, err
	}

	e = s.executeTx(sql, m)

	if e != nil {
		err = fmt.Errorf("insert: %w", e)
		return emptyMeasurement, err
	}

	return m, nil
}

func (s *measurementStore) ReadMeasurement(id string) (entity Measurement, err error) {
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
		err = fmt.Errorf("could not query: %w", e)
		return emptyMeasurement, err
	}

	return measurement, nil
}

func (s *measurementStore) UpdateMeasurement(m Measurement) (entity Measurement, err error) {
	const sql = "UPDATE measurements SET Timestamp = ?, Sensor = ?, Value = ?, Unit = ? WHERE ID = ?"

	e := s.executeTx(sql, m)

	if e != nil {
		err = fmt.Errorf("update: %w", e)
		return emptyMeasurement, err
	}

	return m, nil
}

func (s *measurementStore) DeleteMeasurement(id string) (err error) {
	tx, e := s.db.Begin()

	if e != nil {
		err = fmt.Errorf("could not begin transaction: %w", e)
		return err
	}

	stmt, e := tx.Prepare("DELETE from measurements where id = ?")

	if e != nil {
		err = fmt.Errorf("could not prepare transaction: %w", e)
		return err
	}

	defer stmt.Close()

	_, e = stmt.Exec(id)

	if e != nil {
		err = fmt.Errorf("could not execute statement: %w", e)
		return err
	}

	e = tx.Commit()

	if e != nil {
		err = fmt.Errorf("could not commit transaction: %w", e)
		return err
	}

	return nil
}
