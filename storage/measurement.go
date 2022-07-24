package storage

import (
	"context"
	"database/sql"
	"fmt"
)

// Measurement is the transport object for the values
type Measurement struct {
	ID        string  `json:"id"`
	Timestamp string  `json:"timestamp"`
	Sensor    string  `json:"sensor"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
}

// 'new' exported MeasurementStore interface that can be mocked
type MeasurementStore interface {
	SetupMeasurements(context.Context, Measurement) error
	InsertMeasurement(context.Context, Measurement) error
	// ....
}

// unexported SQL implementation
type measurementStore struct {
	db *sql.DB
}

// exported 'constructor'
func NewMeasurementStore(db *sql.DB) *measurementStore {
	return &measurementStore{
		db: db,
	}
}

func (s *measurementStore) SetupMeasurements(ctx context.Context, m Measurement) (err error) {
	return nil
}

func (s *measurementStore) InsertMeasurement(ctx context.Context, m Measurement) (entity Measurement, err error) {
	const stmt = "INSERT INTO measurements (ID, Timestamp, Sensor, Value, Unit) VALUES (?, ?, ?, ?, ?)"

	result, e := s.db.ExecContext(ctx, stmt, m.ID, m.Timestamp, m.Sensor, m.Value, m.Unit)
	if err != nil {
		e = fmt.Errorf("could not insert row: %w", err)
	}

	if _, e = result.RowsAffected(); err != nil {
		e = fmt.Errorf("could not get affected rows: %w", err)
	}

	if e != nil {
		return m, err
	}

	return entity, nil
}
