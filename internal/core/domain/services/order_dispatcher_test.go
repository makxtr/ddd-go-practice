package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func getLocation() kernel.Location {
	loc, _ := kernel.NewLocation(3, 4)
	return loc
}

func getCouriers() []*courier.Courier {
	loc1, _ := kernel.NewLocation(1, 1)
	loc2, _ := kernel.NewLocation(5, 5)
	loc3, _ := kernel.NewLocation(2, 3)

	c1, _ := courier.NewCourier("Courier1", 2, loc1)
	c2, _ := courier.NewCourier("Courier2", 3, loc2)
	c3, _ := courier.NewCourier("Courier3", 1, loc3)

	return []*courier.Courier{c1, c2, c3}
}

func getOrder() *order.Order {
	o, _ := order.NewOrder(uuid.New(), getLocation(), 5)
	return o
}

func Test_ReturnsErrorWhenOrderInNotCreatedStatus(t *testing.T) {
	// Arrange: заказ в статусе Created
	o := getOrder()
	_ = o.Assign(uuid.New())

	// Act
	_, err := NewOrderDispatcher().Dispatch(o, getCouriers())

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func Test_DispatchesToClosestCourier(t *testing.T) {
	// Arrange
	o := getOrder() // location (3,4), volume 5

	// Act
	result, err := NewOrderDispatcher().Dispatch(o, getCouriers())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, order.StatusAssigned, o.Status())
	assert.NotNil(t, o.CourierID())
	assert.Equal(t, result.ID(), *o.CourierID())
}

func Test_ReturnsErrorWhenNoCouriersAvailable(t *testing.T) {
	// Arrange
	o := getOrder()

	// Act
	_, err := NewOrderDispatcher().Dispatch(o, []*courier.Courier{})

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, errs.ErrObjectNotFound)
}

func Test_SkipsCourierWithNoStorage(t *testing.T) {
	// Arrange: один курьер близко но без места, другой далеко но с местом
	o, _ := order.NewOrder(uuid.New(), getLocation(), 15) // volume 15 > default bag capacity 10

	loc1, _ := kernel.NewLocation(3, 3)
	loc2, _ := kernel.NewLocation(10, 10)

	c1, _ := courier.NewCourier("Close", 5, loc1)   // близко, но сумка = 10 < 15
	c2, _ := courier.NewCourier("Far", 2, loc2)      // далеко
	_ = c2.AddStoragePlace("Trunk", 20)               // есть место на 20

	couriers := []*courier.Courier{c1, c2}

	// Act
	result, err := NewOrderDispatcher().Dispatch(o, couriers)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Far", result.Name())
}
