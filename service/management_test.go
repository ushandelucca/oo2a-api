// Package service implements the management and the reporting use case.
package service

import (
	"MeasurementWeb/database"
	"testing"
)

func Test_managementService_SaveMeasurement(t *testing.T) {
	type fields struct {
		db database.MeasurementDB
	}
	type args struct {
		m MeasurementModel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &managementService{
				db: tt.fields.db,
			}
			if err := s.SaveMeasurement(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("managementService.SaveMeasurement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
