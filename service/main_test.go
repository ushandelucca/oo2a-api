package service

import (
	"MeasurementWeb/database"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
)

type measurementMockDb struct {
	mock.Mock
}

var mockDB *measurementMockDb
var testService *managementService

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// prepare
	mockDB = &measurementMockDb{}
	testService = &managementService{db: mockDB}

	// run the tests
	code := m.Run()

	// after the tests

	os.Exit(code)
}

func (o *measurementMockDb) SetupMeasurements() (err error) {
	args := o.Called()
	return args.Error(0)
}

func (o *measurementMockDb) CreateMeasurement(m database.MeasurementDo) (entity database.MeasurementDo, err error) {
	args := o.Called(m)
	return args.Get(0).(database.MeasurementDo), args.Error(1)
}

func (o *measurementMockDb) ReadMeasurement(id uint) (entity database.MeasurementDo, err error) {
	args := o.Called(id)
	return args.Get(0).(database.MeasurementDo), args.Error(1)
}

func (o *measurementMockDb) UpdateMeasurement(m database.MeasurementDo) (err error) {
	args := o.Called(m)
	return args.Error(1)
}

func (o *measurementMockDb) DeleteMeasurement(id uint) (err error) {
	args := o.Called(id)
	return args.Error(0)
}
