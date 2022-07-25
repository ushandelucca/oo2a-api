package storage

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
	description         string
	measurement         Measurement
	expectedMeasurement Measurement
	expectError         bool
}{
	{"simple", Measurement{ID: "", Timestamp: "t1", Sensor: "s1", Value: 1.3, Unit: "%"}, Measurement{Timestamp: "t1", Sensor: "s1", Value: 1.3, Unit: "%"}, false},
	{"error when id", Measurement{ID: "id1", Timestamp: "t1", Sensor: "s1", Value: 1.3, Unit: "%"}, Measurement{Timestamp: "t1", Sensor: "s1", Value: 1.3, Unit: "%"}, true},
}

func TestCreatMeasurement(t *testing.T) {
	var actual Measurement
	var err error

	for _, tc := range createMeasurementTestCases {
		actual, err = testStore.CreateMeasurement(tc.measurement)

		if tc.expectError {
			assert.Error(t, err)
			assert.Equal(t, emptyMeasurement, actual)
		} else {
			assert.NoError(t, err)
			assert.NotZero(t, actual.ID)
			assert.Equal(t, tc.expectedMeasurement.Timestamp, actual.Timestamp)
			assert.Equal(t, tc.expectedMeasurement.Sensor, actual.Sensor)
			assert.Equal(t, tc.expectedMeasurement.Value, actual.Value)
			assert.Equal(t, tc.expectedMeasurement.Unit, actual.Unit)
		}
	}
}
