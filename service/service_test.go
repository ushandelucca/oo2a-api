// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"testing"
)

var saveTests = []struct {
	name       string
	dto        MeasurementDto
	mockDo     database.MeasurementDo
	mockEntity database.MeasurementDo
	mockError  error
	wantErr    bool
}{
	{"case 1", MeasurementDto{}, database.MeasurementDo{}, database.MeasurementDo{}, nil, true},
	//	{"case 1", MeasurementDto{}, database.MeasurementDo{Unit: "Percent"}, database.MeasurementDo{}, nil, true},
	// {"case 2", args{MeasurementDto{Timestamp: time.Date(2022, 3, 2, 10, 44, 48, 21, time.Local), Sensor: "s1", Value: 2, Unit: Percent}, ""}, database.MeasurementDo{Unit: "Percent"}, database.MeasurementDo{}, nil, false},
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
