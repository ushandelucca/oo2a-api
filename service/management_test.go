// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"testing"
)

type args struct {
	m  MeasurementModel
	id string
}

var saveTests = []struct {
	name       string
	args       args
	mockDo     database.MeasurementDo
	mockEntity database.MeasurementDo
	mockError  error
	wantErr    bool
}{
	{"case 1", args{MeasurementModel{}, ""}, database.MeasurementDo{}, database.MeasurementDo{}, nil, true},
}

func TestSaveMeasurement(t *testing.T) {
	for _, tt := range saveTests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB.On("CreateMeasurement", tt.mockDo).Return(tt.mockEntity, tt.mockError)

			if err := testService.SaveMeasurement(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("managementService.SaveMeasurement() error = %v, wantErr %v", err, tt.wantErr)
			}

			// mockDB.AssertExpectations(t)
		})
	}
}
