package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupMeasurements(t *testing.T) {
	// expecting the existing table
	_, err := testStore.db.Exec("SELECT * FROM measurements LIMIT 1")
	require.NoError(t, err)
	err = testStore.SetupMeasurements()
	require.NoError(t, err)

	// drop table
	_, _ = testStore.db.Exec("DROP TABLE measurements")
	_, err = testStore.db.Exec("SELECT * FROM measurements LIMIT 1")
	require.Error(t, err)

	// create table
	err = testStore.SetupMeasurements()
	require.NoError(t, err)

	// check table
	_, err = testStore.db.Exec("SELECT * FROM measurements LIMIT 1")
	require.NoError(t, err)
}

var createMeasurementTestCases = []struct {
	name            string
	arg             Measurement
	wantMeasurement Measurement
	wantErr         bool
}{
	{"case 1", Measurement{"", "t1", "s1", 1.3, "%"}, Measurement{"id", "t1", "s1", 1.3, "%"}, false},
	{"error when id", Measurement{"id1", "t1", "s1", 1.3, "%"}, Measurement{}, true},
}

func TestCreateMeasurement(t *testing.T) {
	for _, tt := range createMeasurementTestCases {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testStore.CreateMeasurement(tt.arg)

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
	arg     Measurement
	wantErr bool
}{
	{"case 1", Measurement{"", "t1", "s1", 1.3, "%"}, false},
	{"case 2", Measurement{"", "t2", "s1", 2.1, "%"}, false},
}

func TestReadMeasurement(t *testing.T) {
	for _, tt := range readDeleteTestCases {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := testStore.CreateMeasurement(tt.arg)
			require.NoError(t, err)
			require.NotZero(t, entity.ID)

			got, err := testStore.ReadMeasurement(entity.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, got.ID, entity.ID)
			assert.Equal(t, got.Timestamp, entity.Timestamp)
			assert.Equal(t, got.Sensor, entity.Sensor)
			assert.Equal(t, got.Value, entity.Value)
			assert.Equal(t, got.Unit, entity.Unit)
		})
	}
}

func TestDeleteMeasurement(t *testing.T) {
	for _, tt := range readDeleteTestCases {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := testStore.CreateMeasurement(tt.arg)
			require.NoError(t, err)
			require.NotZero(t, entity.ID)

			err = testStore.DeleteMeasurement(entity.ID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

var updateMeasurementTestCases = []struct {
	name            string
	arg             Measurement
	wantMeasurement Measurement
	wantErr         bool
}{
	{"case 1", Measurement{"", "t1", "s1", 1.1, "%"}, Measurement{"", "t1", "s1", 1.2, "%"}, false},
	{"case 2", Measurement{"", "t1", "s1", 2.1, "%"}, Measurement{"", "t1.1", "s1", 2.2, "%"}, false},
}

func TestUpdateMeasurement(t *testing.T) {
	for _, tt := range updateMeasurementTestCases {
		t.Run(tt.name, func(t *testing.T) {
			entity, err := testStore.CreateMeasurement(tt.arg)
			require.NoError(t, err)
			require.NotZero(t, entity.ID)

			tt.wantMeasurement.ID = entity.ID
			got, err := testStore.UpdateMeasurement(tt.wantMeasurement)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMeasurement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.wantMeasurement.ID, got.ID)
			assert.Equal(t, tt.wantMeasurement.Timestamp, got.Timestamp)
			assert.Equal(t, tt.wantMeasurement.Sensor, got.Sensor)
			assert.Equal(t, tt.wantMeasurement.Value, got.Value)
			assert.Equal(t, tt.wantMeasurement.Unit, got.Unit)
		})
	}
}
