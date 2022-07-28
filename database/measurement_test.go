package database

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupMeasurements(t *testing.T) {
	// expecting the existing table
	_, err := testDB.db.Exec("SELECT * FROM measurements LIMIT 1")
	require.NoError(t, err)
	err = testDB.SetupMeasurements()
	require.NoError(t, err)

	// drop table
	_, _ = testDB.db.Exec("DROP TABLE measurements")
	_, err = testDB.db.Exec("SELECT * FROM measurements LIMIT 1")
	require.Error(t, err)

	// create table
	err = testDB.SetupMeasurements()
	require.NoError(t, err)

	// check table
	_, err = testDB.db.Exec("SELECT * FROM measurements LIMIT 1")
	require.NoError(t, err)
}

var createMeasurementTestCases = []struct {
	name            string
	arg             MeasurementDo
	wantMeasurement MeasurementDo
	wantErr         bool
}{
	{"case 1", MeasurementDo{"", "t1", "s1", 1.3, "%"}, MeasurementDo{"id", "t1", "s1", 1.3, "%"}, false},
	{"error when id", MeasurementDo{"id1", "t1", "s1", 1.3, "%"}, MeasurementDo{}, true},
}

func TestCreateMeasurement(t *testing.T) {
	for _, tt := range createMeasurementTestCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testDB.CreateMeasurement(tt.arg)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.True(t, (tt.wantMeasurement.ID != got.ID) != tt.wantErr)
			assert.Equal(t, tt.wantMeasurement.Timestamp, got.Timestamp)
			assert.Equal(t, tt.wantMeasurement.Sensor, got.Sensor)
			assert.Equal(t, tt.wantMeasurement.Value, got.Value)
			assert.Equal(t, tt.wantMeasurement.Unit, got.Unit)
		})
	}
}

var readDeleteTestCases = []struct {
	name    string
	arg     MeasurementDo
	wantErr bool
}{
	{"case 1", MeasurementDo{"", "t1", "s1", 1.3, "%"}, false},
	{"case 2", MeasurementDo{"", "t2", "s1", 2.1, "%"}, false},
}

func TestReadMeasurement(t *testing.T) {
	for _, tt := range readDeleteTestCases {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := testDB.CreateMeasurement(tt.arg)
			require.NoError(t, err)
			require.NotZero(t, entity.ID)

			got, err := testDB.ReadMeasurement(entity.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, entity) {
				t.Errorf("ReadMeasurement() = %v, want %v", got, entity)
			}

		})
	}
}

func TestDeleteMeasurement(t *testing.T) {
	for _, tt := range readDeleteTestCases {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := testDB.CreateMeasurement(tt.arg)
			require.NoError(t, err)
			require.NotZero(t, entity.ID)

			err = testDB.DeleteMeasurement(entity.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := testDB.ReadMeasurement(entity.ID)
			require.NoError(t, err)

			if !reflect.DeepEqual(got, emptyMeasurement) {
				t.Errorf("DeleteMeasurement() = %v, want %v", got, emptyMeasurement)
			}

		})
	}
}

var updateMeasurementTestCases = []struct {
	name            string
	arg             MeasurementDo
	wantMeasurement MeasurementDo
	wantErr         bool
}{
	{"case 1", MeasurementDo{"", "t1", "s1", 1, "%"}, MeasurementDo{"", "t1", "s1", 1.2, "%"}, false},
	{"case 2", MeasurementDo{"", "t1", "s1", 2, "%"}, MeasurementDo{"", "t1.1", "s1", 2.2, "%"}, false},
}

func TestUpdateMeasurement(t *testing.T) {
	for _, tt := range updateMeasurementTestCases {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := testDB.CreateMeasurement(tt.arg)
			require.NoError(t, err)
			require.NotZero(t, entity.ID)

			tt.wantMeasurement.ID = entity.ID
			_, err = testDB.UpdateMeasurement(tt.wantMeasurement)
			require.NoError(t, err)

			got, err := testDB.ReadMeasurement(entity.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.wantMeasurement) {
				t.Errorf("UpdateMeasurement() = %v, want %v", got, tt.wantMeasurement)
			}
		})
	}
}
