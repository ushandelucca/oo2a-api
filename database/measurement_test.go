package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupMeasurements(t *testing.T) {
	// expecting the existing table
	tx := testDB.db.Exec("SELECT * FROM measurement_dos LIMIT 1")
	require.NoError(t, tx.Error)
	err := testDB.SetupMeasurements()
	require.NoError(t, err)

	// drop table
	testDB.db.Exec("DROP TABLE measurement_dos")
	tx = testDB.db.Exec("SELECT * FROM measurement_dos LIMIT 1")
	require.Error(t, tx.Error)

	// create table
	err = testDB.SetupMeasurements()
	require.NoError(t, err)

	// check table
	tx = testDB.db.Exec("SELECT * FROM measurement_dos LIMIT 1")
	require.NoError(t, tx.Error)
}

var createMeasurementTestCases = []struct {
	name            string
	arg             MeasurementDo
	wantMeasurement MeasurementDo
	wantErr         bool
}{
	{"case 1", MeasurementDo{Timestamp: time.Date(2022, 01, 1, 10, 44, 48, 120, time.Local), Sensor: "s1", Value: .3, Unit: "%"}, MeasurementDo{Timestamp: time.Date(2022, 01, 1, 10, 44, 48, 120, time.Local), Sensor: "s1", Value: .3, Unit: "%"}, false},
	{"error when id", MeasurementDo{ID: 1, Timestamp: time.Date(2022, 01, 2, 10, 44, 48, 120, time.Local), Sensor: "s1", Value: 1.3, Unit: "%"}, MeasurementDo{}, true},
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
	{"case 1", MeasurementDo{Timestamp: time.Date(2022, 02, 1, 10, 44, 48, 1, time.Local), Sensor: "s1", Value: 1.3, Unit: "%"}, false},
	{"case 2", MeasurementDo{Timestamp: time.Date(2022, 02, 2, 10, 44, 48, 2, time.Local), Sensor: "s1", Value: 2.1, Unit: "%"}, false},
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

			// if !reflect.DeepEqual(got, entity) {
			// 	t.Errorf("ReadMeasurement() = %v, want %v", got, entity)
			// }

			assert.True(t, (tt.arg.ID != got.ID) != tt.wantErr)
			assert.WithinDuration(t, tt.arg.Timestamp, got.Timestamp, 0)
			assert.Equal(t, tt.arg.Sensor, got.Sensor)
			assert.Equal(t, tt.arg.Value, got.Value)
			assert.Equal(t, tt.arg.Unit, got.Unit)
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

			assert.True(t, (tt.arg.ID != got.ID) != tt.wantErr)
			assert.WithinDuration(t, tt.arg.Timestamp, got.Timestamp, 0)
			assert.Equal(t, tt.arg.Sensor, got.Sensor)
			assert.Equal(t, tt.arg.Value, got.Value)
			assert.Equal(t, tt.arg.Unit, got.Unit)
		})
	}
}

var updateMeasurementTestCases = []struct {
	name            string
	arg             MeasurementDo
	wantMeasurement MeasurementDo
	wantErr         bool
}{
	{"case 1", MeasurementDo{Timestamp: time.Date(2022, 3, 1, 10, 44, 48, 20, time.Local), Sensor: "s1", Value: 1, Unit: "%"}, MeasurementDo{Timestamp: time.Date(2022, 3, 1, 10, 44, 48, 20, time.Local), Sensor: "s1", Value: 1.2, Unit: "%"}, false},
	{"case 2", MeasurementDo{Timestamp: time.Date(2022, 3, 2, 10, 44, 48, 21, time.Local), Sensor: "s1", Value: 2, Unit: "%"}, MeasurementDo{Timestamp: time.Date(2022, 3, 2, 10, 44, 48, 21, time.Local), Sensor: "s1", Value: 0.2, Unit: "%"}, false},
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

			assert.Equal(t, tt.wantMeasurement.ID, got.ID)
			assert.WithinDuration(t, tt.wantMeasurement.Timestamp, got.Timestamp, 0)
			assert.Equal(t, tt.wantMeasurement.Sensor, got.Sensor)
			assert.Equal(t, tt.wantMeasurement.Value, got.Value)
			assert.Equal(t, tt.wantMeasurement.Unit, got.Unit)
		})
	}
}
