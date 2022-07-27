// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"fmt"
)

type ManagementService interface {
	InitMeasurement() (err error)
	UpsertMeasurement(m database.Measurement) (entity database.Measurement, err error)
	DeleteMeasurement(id string) (err error)
}

// exported 'constructor'
func NewManagementService() *managementService {
	// ---------------------config---vvv
	s := database.NewMeasurementStore(nil)

	return &managementService{store: s}
}

// managementService relies on the MeasurementStore -> mockable for the tests without DB
type managementService struct {
	store database.MeasurementStore
}

func (s *managementService) InitMeasurement() (err error) {
	err = s.store.SetupMeasurements()
	return err
}

func (s *managementService) UpsertMeasurement(m database.Measurement) (err error) {
	// do validation/business rule validation here
	// finally, insert into the DB

	var entity database.Measurement

	if m.ID == "" {
		entity, err = s.store.CreateMeasurement(m)
	} else {
		entity, err = s.store.UpdateMeasurement(m)
	}

	if err != nil {
		return err
	}

	fmt.Sprintln("result: %w", entity)

	return nil
}
