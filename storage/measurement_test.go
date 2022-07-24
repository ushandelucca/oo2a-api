// Package storage implements the db access for the CRUD operations.
package storage

import (
	"context"
	"testing"

	// using https://github.com/stretchr/testify library for brevity
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertMeasurement(t *testing.T) {
	store := &measurementStore{
		db: testDB.db,
	}

	object := Measurement{ID: "", Timestamp: ""}

	entity, err := store.InsertMeasurement(context.TODO(), object)

	require.NoError(t, err)

	assert.Equal(t, "t1", entity.Timestamp)
	assert.Equal(t, "s1", entity.Sensor)
	assert.Equal(t, 1.3, entity.Value)
	assert.Equal(t, "%", entity.Unit)
	assert.NotZero(t, entity.ID)
}
