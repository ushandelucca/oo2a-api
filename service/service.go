// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"MeasurementWeb/utils"
	"fmt"
	"log"
	"time"

	"github.com/dranikpg/dto-mapper"
)

type Unit int

const (
	Percent Unit = iota
	lph
)

// Measurement is the transport object for the measurement values.
type MeasurementModel struct {
	Timestamp time.Time
	Sensor    string
	Value     float64
	Unit      Unit
}

// MeasurementDto is the transport object for the measurement values.
type Oo2aModel struct {
	ValuesBegin time.Time
	ValuesEnd   time.Time

	LevelCurrent float64
	LevelValues  []float64
	LevelUnit    Unit

	PrecipitationCurrent float64
	PrecipitationValues  []float64
	PrecipitationUnit    Unit
}

type ManagementService interface {
	InitMeasurement() (err error)
	SaveMeasurement(m database.MeasurementDo) (entity database.MeasurementDo, err error)
}

// exported 'constructor'
func NewManagementService(config *utils.Conf) *managementService {
	db, err := database.NewMeasurementDB(config)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &managementService{db: db}
}

// managementService relies on the MeasurementDB which is mockable for the tests without a DB
type managementService struct {
	db database.MeasurementDB
}

func (s *managementService) InitMeasurement() (err error) {
	err = s.db.SetupMeasurements()
	return err
}

func (s *managementService) SaveMeasurement(m MeasurementModel) (err error) {
	// do validation/business rule validation here
	// finally, insert into the DB

	do := database.MeasurementDo{}
	// TODO fix the mapping e.g. Timstamp to string
	err = dto.Map(&do, m)
	if err != nil {
		return err
	}

	if do.ID == "" {
		_, err = s.db.CreateMeasurement(do)
	} else {
		_, err = s.db.UpdateMeasurement(do)
	}

	if err != nil {
		return err
	}

	fmt.Sprintln("result: %w", do)

	return nil
}
