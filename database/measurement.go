// Package database implements the DB access for the CRUD operations.
package database

import (
	"MeasurementWeb/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MeasurementDo is the data object for the database.
type MeasurementDo struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;autoIncrement"`
	Timestamp time.Time
	Sensor    string
	Value     float64
	Unit      string
}

// MeasurementDB provides the CRUD storage functionality.
type MeasurementDB interface {
	// SetupMeasurements creates the database and the tables.
	SetupMeasurements() (err error)

	// CreateMeasurement adds a new measurement to the table. The ID must be empty and
	// will be defined within this function.
	CreateMeasurement(m MeasurementDo) (entity MeasurementDo, err error)

	// ReadMeasurement reads an measurement. When no record is found a empty struct
	// will be returned.
	ReadMeasurement(id uint) (entity MeasurementDo, err error)

	// TODO comment
	UpdateMeasurement(m MeasurementDo) (err error)

	// TODO comment
	DeleteMeasurement(id uint) (err error)
}

// NewMeasurementDB returns the DB object.
func NewMeasurementDB(config *utils.Conf) (db *measurementDB, err error) {

	database, err := gorm.Open(sqlite.Open(config.DataSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}

	db = &measurementDB{db: database}

	return db, nil
}

type measurementDB struct {
	db *gorm.DB
}

func (s *measurementDB) SetupMeasurements() (err error) {
	err = s.db.AutoMigrate(&MeasurementDo{})

	if err != nil {
		err = fmt.Errorf("setup: %w", err)
	}

	return err
}

func (s *measurementDB) CreateMeasurement(m MeasurementDo) (entity MeasurementDo, err error) {
	if m.ID != 0 {
		err = errors.New("the ID must be empty")
		return MeasurementDo{}, err
	}

	tx := s.db.Create(&m)

	err = tx.Error
	if err != nil {
		err = fmt.Errorf("create: %w", err)
	}

	return m, err
}

func (s *measurementDB) ReadMeasurement(id uint) (entity MeasurementDo, err error) {
	entity = MeasurementDo{}
	tx := s.db.First(&entity, "id = ?", id)

	err = tx.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return MeasurementDo{}, nil
		}

		err = fmt.Errorf("read: %w", err)
	}

	return entity, err
}

func (s *measurementDB) UpdateMeasurement(m MeasurementDo) (err error) {
	entity := MeasurementDo{}
	tx := s.db.First(&entity, "id = ?", m.ID)
	tx.Model(entity).Updates(m)

	if tx.RowsAffected == 0 {
		return errors.New("update: no rows affected")
	}

	err = tx.Error
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (s *measurementDB) DeleteMeasurement(id uint) (err error) {
	tx := s.db.Delete(&MeasurementDo{}, id)

	err = tx.Error
	if err != nil {
		err = fmt.Errorf("delete: %w", err)
	}

	return err
}
