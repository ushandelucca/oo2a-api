// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validationsTests = []struct {
	name      string
	dto       MeasurementDto
	numValErr int
}{
	{"test 1", MeasurementDto{}, 4},
	{"test 2", MeasurementDto{ID: 1}, 4},
	{"test 3", MeasurementDto{ID: 1, Timestamp: time.Now()}, 3},
	{"test 4", MeasurementDto{ID: 0, Timestamp: time.Now(), Sensor: "s1"}, 2},
	{"test 5", MeasurementDto{ID: 0, Timestamp: time.Now(), Sensor: "s1", Value: 0.01}, 1},
	{"test 6", MeasurementDto{Timestamp: time.Date(2022, 3, 2, 10, 44, 48, 21, time.Local), Sensor: "s1", Value: 2, Unit: Percent}, 0},
}

func TestSaveNewMeasurementValidations(t *testing.T) {
	for _, tt := range validationsTests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.numValErr == 0 {
				mockDB.On("CreateMeasurement", database.MeasurementDo{Timestamp: time.Date(2022, 3, 2, 10, 44, 48, 21, time.Local), Sensor: "s1", Value: 2, Unit: Percent.String()}).Return(database.MeasurementDo{}, nil)
			}

			err := testService.SaveMeasurement(tt.dto)

			var errs []string
			if err != nil {
				var valErrs validator.ValidationErrors
				if errors.As(err, &valErrs) {
					for _, e := range err.(validator.ValidationErrors) {
						errs = append(errs, e.Error())
					}
				}
			}

			assert.Equal(t, tt.numValErr, len(errs))
			mockDB.AssertExpectations(t)
		})
	}
}

var saveTests = []struct {
	name       string
	dto        MeasurementDto
	mockDo     database.MeasurementDo
	mockEntity database.MeasurementDo
	mockError  error
	wantErr    bool
}{
	{"test 1", MeasurementDto{}, database.MeasurementDo{}, database.MeasurementDo{}, nil, true},
	// {"test 2", args{MeasurementDto{Timestamp: time.Date(2022, 3, 2, 10, 44, 48, 21, time.Local), Sensor: "s1", Value: 2, Unit: Percent}, ""}, database.MeasurementDo{Unit: "Percent"}, database.MeasurementDo{}, nil, false},
}

func TestSaveNewMeasurement(t *testing.T) {
	for _, tt := range saveTests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				mockDB.On("CreateMeasurement", tt.mockDo).Return(tt.mockEntity, tt.mockError)
			}

			if err := testService.SaveMeasurement(tt.dto); (err != nil) != tt.wantErr {
				t.Errorf("managementService.SaveMeasurement() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
