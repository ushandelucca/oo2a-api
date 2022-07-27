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

var readDeleteTestCases = []struct {
	description string
	valueObject Measurement
	expectError bool
}{
	{"simple", Measurement{Timestamp: "t1", Sensor: "s1", Value: 1.3, Unit: "%"}, false},
	{"simple", Measurement{Timestamp: "t2", Sensor: "s1", Value: 2.1, Unit: "%"}, false},
}

func TestReadMeasurement(t *testing.T) {
	var entity Measurement
	var actual Measurement
	var err error

	for _, tc := range readDeleteTestCases {
		entity, err = testStore.CreateMeasurement(tc.valueObject)
		require.NoError(t, err)

		actual, err = testStore.ReadMeasurement(entity.ID)

		if tc.expectError {
			assert.Error(t, err)
			assert.Equal(t, emptyMeasurement, entity)
		} else {
			assert.NoError(t, err)
			assert.NotZero(t, actual.ID)
			assert.Equal(t, actual.Timestamp, entity.Timestamp)
			assert.Equal(t, actual.Sensor, entity.Sensor)
			assert.Equal(t, actual.Value, entity.Value)
			assert.Equal(t, actual.Unit, entity.Unit)
		}
	}
}

func TestDeleteMeasurement(t *testing.T) {
	var entity Measurement
	var err error

	for _, tc := range readDeleteTestCases {
		entity, err = testStore.CreateMeasurement(tc.valueObject)
		require.NoError(t, err)
		require.NotZero(t, entity.ID)

		err = testStore.DeleteMeasurement(entity.ID)

		if tc.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

var updateMeasurementTestCases = []struct {
	description        string
	measurement        Measurement
	updatedMeasurement Measurement
	expectError        bool
}{
	{"simple", Measurement{Timestamp: "t1", Sensor: "s1", Value: 1.1, Unit: "%"}, Measurement{Timestamp: "t1", Sensor: "s1", Value: 1.2, Unit: "%"}, false},
	{"error when id", Measurement{Timestamp: "t1", Sensor: "s1", Value: 2.1, Unit: "%"}, Measurement{Timestamp: "t1.1", Sensor: "s1", Value: 2.2, Unit: "%"}, false},
}

func TestUpdateMeasurement(t *testing.T) {
	var entity Measurement
	var actual Measurement
	var err error

	for _, tc := range updateMeasurementTestCases {
		entity, err = testStore.CreateMeasurement(tc.measurement)
		require.NoError(t, err)

		tc.updatedMeasurement.ID = entity.ID

		actual, err = testStore.UpdateMeasurement(tc.updatedMeasurement)

		if tc.expectError {
			assert.Error(t, err)
			assert.Equal(t, emptyMeasurement, entity)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, actual.ID, tc.updatedMeasurement.ID)
			assert.Equal(t, actual.Timestamp, tc.updatedMeasurement.Timestamp)
			assert.Equal(t, actual.Sensor, tc.updatedMeasurement.Sensor)
			assert.Equal(t, actual.Value, tc.updatedMeasurement.Value)
			assert.Equal(t, actual.Unit, tc.updatedMeasurement.Unit)
		}
	}
}
