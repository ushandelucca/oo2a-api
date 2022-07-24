package storage

import (
	"testing"

	// using https://github.com/stretchr/testify library for brevity
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupMeasurements(t *testing.T) {

	// expecting existing table
	// drop table
	// -> select error
	// setup
	// -> select empty

}

func TestInsertMeasurement(t *testing.T) {
	store := &measurementStore{
		db: testDB.db,
	}

	object := Measurement{Timestamp: "t1", Sensor: "s1", Value: 1.3, Unit: "%"}

	err := store.SetupMeasurements()
	entity, err := store.CreateMeasurement(object)

	require.NoError(t, err)

	assert.Equal(t, "t1", entity.Timestamp)
	assert.Equal(t, "s1", entity.Sensor)
	assert.Equal(t, 1.3, entity.Value)
	assert.Equal(t, "%", entity.Unit)
	assert.NotZero(t, entity.ID)
}
