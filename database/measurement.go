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

var emptyMeasurement = MeasurementDo{}

// MeasurementDB provides the CRUD storage functionality.
type MeasurementDB interface {
	// SetupMeasurements creates the database and the tables.
	SetupMeasurements() (err error)

	// CreateMeasurement adds a new measurement to the table. The ID must be empty and
	// will be defined within this function.
	CreateMeasurement(m MeasurementDo) (entity MeasurementDo, err error)

	// ReadMeasurement reads an measurement
	ReadMeasurement(id uint) (entity MeasurementDo, err error)
	UpdateMeasurement(m MeasurementDo) (entity MeasurementDo, err error)
	DeleteMeasurement(id uint) (err error)
}

// NewMeasurementDB returns the config object.
func NewMeasurementDB(config *utils.Conf) (db *measurementDB, err error) {

	database, err := gorm.Open(sqlite.Open(config.DataSourceName), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
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
		err = fmt.Errorf("setup: could not create table: %w", err)
	}

	return err
}

func (s *measurementDB) CreateMeasurement(m MeasurementDo) (entity MeasurementDo, err error) {
	if m.ID != 0 {
		err = errors.New("the ID must be empty")
		return emptyMeasurement, err
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
		err = fmt.Errorf("read: %w", err)
	}

	return entity, err
}

func (s *measurementDB) UpdateMeasurement(m MeasurementDo) (entity MeasurementDo, err error) {

	tx := s.db.First(&entity, "id = ?", m.ID)
	tx.Model(entity).Updates(m)

	err = tx.Error

	if err != nil {
		err = fmt.Errorf("update: %w", err)
	}

	return m, err
}

func (s *measurementDB) DeleteMeasurement(id uint) (err error) {
	entity := MeasurementDo{}
	tx := s.db.First(&entity, "id = ?", id)

	err = tx.Error

	if err != nil {
		err = fmt.Errorf("delete: %w", err)
	}

	return err
}
