package order

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func validLocation() kernel.Location {
	loc, _ := kernel.NewLocation(3, 4)
	return loc
}

func validOrder() *Order {
	o, _ := NewOrder(uuid.New(), validLocation(), 5)
	return o
}

func Test_NewOrderBeCorrectWhenParamsAreCorrect(t *testing.T) {
	// Arrange
	id := uuid.New()
	loc := validLocation()

	// Act
	o, err := NewOrder(id, loc, 5)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, o)
	assert.Equal(t, id, o.ID())
	assert.Equal(t, loc, o.Location())
	assert.Equal(t, 5, o.Volume())
	assert.Equal(t, StatusCreated, o.Status())
	assert.Nil(t, o.CourierID())
}

func Test_NewOrderReturnErrorWhenParamsAreInvalid(t *testing.T) {
	loc := validLocation()

	tests := map[string]struct {
		id       uuid.UUID
		location kernel.Location
		volume   int
	}{
		"nil id": {
			id:       uuid.Nil,
			location: loc,
			volume:   5,
		},
		"empty location": {
			id:       uuid.New(),
			location: kernel.Location{},
			volume:   5,
		},
		"zero volume": {
			id:       uuid.New(),
			location: loc,
			volume:   0,
		},
		"negative volume": {
			id:       uuid.New(),
			location: loc,
			volume:   -1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			o, err := NewOrder(test.id, test.location, test.volume)

			// Assert
			assert.Nil(t, o)
			assert.Error(t, err)
			assert.ErrorIs(t, err, errs.ErrValueIsRequired)
		})
	}
}

func Test_AssignSetsStatusAndCourierID(t *testing.T) {
	// Arrange
	o := validOrder()
	courierID := uuid.New()

	// Act
	err := o.Assign(courierID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StatusAssigned, o.Status())
	assert.NotNil(t, o.CourierID())
	assert.Equal(t, courierID, *o.CourierID())
}

func Test_AssignReturnsErrorWhenNotCreated(t *testing.T) {
	// Arrange: заказ уже назначен
	o := validOrder()
	_ = o.Assign(uuid.New())

	// Act: повторный Assign
	err := o.Assign(uuid.New())

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func Test_CompleteSetsStatusToCompleted(t *testing.T) {
	// Arrange
	o := validOrder()
	_ = o.Assign(uuid.New())

	// Act
	err := o.Complete()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, StatusCompleted, o.Status())
}

func Test_CompleteReturnsErrorWhenNotAssigned(t *testing.T) {
	// Arrange: заказ в статусе Created
	o := validOrder()

	// Act
	err := o.Complete()

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func Test_CompleteReturnsErrorWhenAlreadyCompleted(t *testing.T) {
	// Arrange
	o := validOrder()
	_ = o.Assign(uuid.New())
	_ = o.Complete()

	// Act
	err := o.Complete()

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}