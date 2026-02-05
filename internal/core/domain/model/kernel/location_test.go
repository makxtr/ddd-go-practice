package kernel

import (
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LocationBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange

	// Act
	location, err := NewLocation(3, 4)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, location)
	assert.Equal(t, 3, location.X())
	assert.Equal(t, 4, location.Y())
}

func Test_LocationBeCorrectWhenCreateRandom(t *testing.T) {
	// Arrange

	// Act
	location := NewRandomLocation()

	// Assert
	assert.True(t, location.X() >= 1 && location.X() <= 10)
	assert.True(t, location.Y() >= 1 && location.Y() <= 10)
	assert.False(t, location.IsEmpty())
}

func Test_LocationCorrectEquals(t *testing.T) {
	// Arrange

	// Act
	location1, _ := NewLocation(3, 4)
	location2, _ := NewLocation(3, 4)

	// Assert
	assert.True(t, location1.Equals(location2))
}

func Test_LocationDistances(t *testing.T) {
	// Arrange
	startLocation, _ := NewLocation(3, 4)
	endLocation, _ := NewLocation(10, 4)
	expectedDistance := 7

	// Act
	actualDistance, err := startLocation.DistanceTo(endLocation)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedDistance, actualDistance)
}

func Test_LocationDistanceToEmptyReturnsError(t *testing.T) {
	// Arrange
	startLocation, _ := NewLocation(1, 1)
	var emptyLocation Location

	// Act
	_, err := startLocation.DistanceTo(emptyLocation)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func Test_LocationReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	// Arrange
	tests := map[string]struct {
		x        int
		y        int
		expected error
	}{
		"-1": {
			x:        -1,
			y:        3,
			expected: errs.NewValueIsOutOfRangeError("x", -1, 1, 10),
		},
		"30": {
			x:        1,
			y:        30,
			expected: errs.NewValueIsOutOfRangeError("y", 30, 1, 10),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewLocation(test.x, test.y)

			// Assert
			assert.EqualError(t, err, test.expected.Error())
		})
	}
}
