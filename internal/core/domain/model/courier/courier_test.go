package courier

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func validLocation() kernel.Location {
	loc, _ := kernel.NewLocation(1, 1)
	return loc
}

func validCourier() *Courier {
	c, _ := NewCourier("Иван", 2, validLocation())
	return c
}

func Test_NewCourierBeCorrectWhenParamsAreCorrect(t *testing.T) {
	// Arrange
	loc, _ := kernel.NewLocation(3, 4)

	// Act
	c, err := NewCourier("Иван", 2, loc)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "Иван", c.name)
	assert.Equal(t, 2, c.speed)
	assert.Equal(t, loc, c.location)
	assert.Len(t, c.storagePlaces, 1)
	assert.Equal(t, "Сумка", c.storagePlaces[0].Name())
	assert.Equal(t, 10, c.storagePlaces[0].TotalVolume())
}

func Test_NewCourierReturnErrorWhenParamsAreInvalid(t *testing.T) {
	loc, _ := kernel.NewLocation(1, 1)

	tests := map[string]struct {
		name     string
		speed    int
		location kernel.Location
	}{
		"empty name": {
			name:     "",
			speed:    1,
			location: loc,
		},
		"zero speed": {
			name:     "Иван",
			speed:    0,
			location: loc,
		},
		"empty location": {
			name:     "Иван",
			speed:    1,
			location: kernel.Location{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			c, err := NewCourier(test.name, test.speed, test.location)

			// Assert
			assert.Nil(t, c)
			assert.Error(t, err)
			assert.ErrorIs(t, err, errs.ErrValueIsRequired)
		})
	}
}

func Test_MoveMovesTowardTarget(t *testing.T) {
	// Arrange: курьер в (1,1), скорость 2, цель (5,5)
	c := validCourier()

	target, _ := kernel.NewLocation(5, 5)

	// Act
	err := c.Move(target)

	// Assert: скорость 2, сначала по X +2, остаток 0 на Y => (3,1)
	assert.NoError(t, err)
	assert.Equal(t, 3, c.location.X())
	assert.Equal(t, 1, c.location.Y())
}

func Test_MoveReturnsErrorOnInvalidTarget(t *testing.T) {
	// Arrange
	c := validCourier()

	// Act
	err := c.Move(kernel.Location{})

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func Test_AddStoragePlaceAddsNewPlace(t *testing.T) {
	// Arrange
	c := validCourier()
	assert.Len(t, c.storagePlaces, 1)

	// Act
	err := c.AddStoragePlace("Рюкзак", 20)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, c.storagePlaces, 2)
	assert.Equal(t, "Рюкзак", c.storagePlaces[1].Name())
	assert.Equal(t, 20, c.storagePlaces[1].TotalVolume())
}

func Test_AddStoragePlaceReturnsErrorOnInvalidParams(t *testing.T) {
	// Arrange
	c := validCourier()

	// Act
	err := c.AddStoragePlace("", 10)

	// Assert
	assert.Error(t, err)
	assert.Len(t, c.storagePlaces, 1)
}

func Test_CanTakeOrderReturnsTrueWhenVolumesFit(t *testing.T) {
	// Arrange
	c := validCourier() // сумка на 10

	// Act
	canTake, err := c.CanTakeOrder(5)

	// Assert
	assert.NoError(t, err)
	assert.True(t, canTake)
}

func Test_CanTakeOrderReturnsFalseWhenVolumeTooLarge(t *testing.T) {
	// Arrange
	c := validCourier() // сумка на 10

	// Act
	canTake, _ := c.CanTakeOrder(15)

	// Assert
	assert.False(t, canTake)
}

func Test_CanTakeOrderReturnsFalseWhenAllPlacesOccupied(t *testing.T) {
	// Arrange
	c := validCourier()
	orderID := uuid.New()
	_ = c.TakeOrder(orderID, 5)

	// Act
	canTake, _ := c.CanTakeOrder(5)

	// Assert
	assert.False(t, canTake)
}

func Test_TakeOrderStoresOrderInStoragePlace(t *testing.T) {
	// Arrange
	c := validCourier()
	orderID := uuid.New()

	// Act
	err := c.TakeOrder(orderID, 5)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, c.storagePlaces[0].OrderID())
	assert.Equal(t, orderID, *c.storagePlaces[0].OrderID())
}

func Test_TakeOrderReturnsErrorWhenNoSpace(t *testing.T) {
	// Arrange
	c := validCourier()
	_ = c.TakeOrder(uuid.New(), 5) // занимаем единственное место

	// Act
	err := c.TakeOrder(uuid.New(), 5)

	// Assert
	assert.Error(t, err)
}

func Test_CompleteOrderClearsStoragePlace(t *testing.T) {
	// Arrange
	c := validCourier()
	orderID := uuid.New()
	_ = c.TakeOrder(orderID, 5)

	// Act
	err := c.CompleteOrder(orderID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, c.storagePlaces[0].OrderID())
}

func Test_CompleteOrderReturnsErrorWhenOrderNotFound(t *testing.T) {
	// Arrange
	c := validCourier()

	// Act
	err := c.CompleteOrder(uuid.New())

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrObjectNotFound)
}

func Test_CalculateTimeToLocation(t *testing.T) {
	// Arrange: курьер в (1,1), скорость 2, цель (5,5)
	// дистанция = 4+4 = 8, время = 8/2 = 4 такта
	c := validCourier()
	target, _ := kernel.NewLocation(5, 5)

	// Act
	steps, err := c.CalculateTimeToLocation(target)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 4, steps)
}

func Test_CalculateTimeToLocationRoundsUp(t *testing.T) {
	// Arrange: курьер в (1,1), скорость 2, цель (2,1)
	// дистанция = 1, время = ceil(1/2) = 1 такт
	c := validCourier()
	target, _ := kernel.NewLocation(2, 1)

	// Act
	steps, err := c.CalculateTimeToLocation(target)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, steps)
}

func Test_CalculateTimeToLocationReturnsErrorOnEmptyTarget(t *testing.T) {
	// Arrange
	c := validCourier()

	// Act
	_, err := c.CalculateTimeToLocation(kernel.Location{})

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}