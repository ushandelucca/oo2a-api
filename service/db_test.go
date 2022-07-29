package service

import (
	"MeasurementWeb/database"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

type dbMock struct {
	mock.Mock

	// other fields go here as normal
	WantErr error
}

var testService *managementService

func TestMain(m *testing.M) {
	// before tests
	testService = &managementService{db: &dbMock{}}
	code := m.Run()

	// after tests

	os.Exit(code)
}

func (o *dbMock) SetupMeasurements() (err error) {
	o.Called()
	return err
}

func (o *dbMock) CreateMeasurement(m database.MeasurementDo) (entity database.MeasurementDo, err error) {
	o.Called()
	return entity, err
}

func (o *dbMock) ReadMeasurement(id string) (entity database.MeasurementDo, err error) {
	o.Called()
	return entity, err
}

func (o *dbMock) UpdateMeasurement(m database.MeasurementDo) (entity database.MeasurementDo, err error) {
	o.Called()
	return entity, err
}

func (o *dbMock) DeleteMeasurement(id string) (err error) {
	o.Called()
	return err
}
