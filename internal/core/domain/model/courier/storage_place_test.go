package courier

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_StoragePlaceBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange

	// Act
	storagePlace, err := NewStoragePlace("bag", 10)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, storagePlace)
	assert.Equal(t, "bag", storagePlace.Name())
	assert.Nil(t, storagePlace.OrderID())

	assert.Equal(t, 10, storagePlace.TotalVolume())
}

func Test_StoragePlaceCanStore(t *testing.T) {
	// Arrange
	storagePlace, err := NewStoragePlace("bag", 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, storagePlace)

	// Act
	canStore, err := storagePlace.CanStore(5)
	assert.NoError(t, err)
	assert.True(t, canStore)

	canStore, err = storagePlace.CanStore(11)
	assert.Error(t, err)
	assert.False(t, canStore)
}

func Test_StoragePlaceStoreAndClear(t *testing.T) {
	// Arrange
	storagePlace, err := NewStoragePlace("bag", 10)
	assert.NoError(t, err)
	assert.NotEmpty(t, storagePlace)

	// Act
	canStore, err := storagePlace.CanStore(5)
	assert.NoError(t, err)
	assert.True(t, canStore)

	orderId := uuid.New()
	err = storagePlace.Store(orderId, 5)
	assert.NoError(t, err)
	assert.NotNil(t, storagePlace.OrderID())

	canStore, err = storagePlace.CanStore(5)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrOrderAlreadyStored)

	err = storagePlace.Clear(orderId)
	assert.NoError(t, err)
	canStore, err = storagePlace.CanStore(5)
	assert.NoError(t, err)
	assert.True(t, canStore)
}
