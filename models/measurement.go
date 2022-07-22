package models

import (
	"database/sql"
	"strconv"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./values.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// Measurement is the transport object for the values
type Measurement struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Sensor    string  `json:"sensor"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

func GetMeasurements(count int) ([]Measurement, error) {

	rows, err := DB.Query("SELECT ID, Timestamp, Sensor, Value, Unit from measurements LIMIT " + strconv.Itoa(count))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	measurements := make([]Measurement, 0)

	for rows.Next() {
		oneMeasurement := Measurement{}
		err = rows.Scan(&oneMeasurement.ID, &oneMeasurement.Timestamp, &oneMeasurement.Sensor, &oneMeasurement.Value, &oneMeasurement.Unit)

		if err != nil {
			return nil, err
		}

		measurements = append(measurements, oneMeasurement)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return measurements, err
}

func GetMeasurementById(id string) (Measurement, error) {

	stmt, err := DB.Prepare("SELECT ID, Timestamp, Sensor, Value, Unit from measurements WHERE ID = ?")

	if err != nil {
		return Measurement{}, err
	}

	measurement := Measurement{}

	sqlErr := stmt.QueryRow(id).Scan(&measurement.ID, &measurement.Timestamp, &measurement.Sensor, &measurement.Value, &measurement.Unit)

	if sqlErr != nil {
		if sqlErr == sql.ErrNoRows {
			return Measurement{}, nil
		}
		return Measurement{}, sqlErr
	}
	return measurement, nil
}

func AddMeasurement(newMeasurement Measurement) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("INSERT INTO measurements (ID, Timestamp, Sensor, Value, Unit) VALUES (?, ?, ?, ?, ?)")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	id := uuid.New()

	_, err = stmt.Exec(id, newMeasurement.Timestamp, newMeasurement.Sensor, newMeasurement.Value, newMeasurement.Unit)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func UpdateMeasurement(ourMeasurement Measurement, id int) (bool, error) {

	tx, err := DB.Begin()
	if err != nil {
		return false, err
	}

	stmt, err := tx.Prepare("UPDATE people SET first_name = ?, last_name = ?, email = ?, ip_address = ? WHERE Id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(ourMeasurement.Timestamp, ourMeasurement.Sensor, ourMeasurement.Value, ourMeasurement.Unit, id)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}

func DeleteMeasurement(personId int) (bool, error) {

	tx, err := DB.Begin()

	if err != nil {
		return false, err
	}

	stmt, err := DB.Prepare("DELETE from people where id = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(personId)

	if err != nil {
		return false, err
	}

	tx.Commit()

	return true, nil
}
