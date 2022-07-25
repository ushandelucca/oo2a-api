package service

import (
	"MeasurementWeb/storage"
	"fmt"
)

type MeasurementService interface {
	UpsertMeasurement(m storage.Measurement) (entity storage.Measurement, err error)
	SelectMeasurements(query string) (m []storage.Measurement, err error)
	DeleteMeasurement(id string) (err error)
}

// MeasurementService relies on the BookStore
type measurementService struct {
	store storage.MeasurementStore
}

// exported 'constructor'
func NewMeasurementService() *measurementService {
	return &measurementService{
		store: storage.NewMeasurementStore(nil),
	}
}

func (s *measurementService) UpsertMeasurement(m storage.Measurement) error {
	// do validation/business rule validation here
	// .. more stuff
	// finally, insert into the DB

	var entity storage.Measurement
	var err error

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
