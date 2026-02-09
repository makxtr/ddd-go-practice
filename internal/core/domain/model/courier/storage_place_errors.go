package courier

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrOrderAlreadyStored      = errors.New("order already stored")
	ErrOrderNotHasEnoughVolume = errors.New("order not has enough volume")
)

type OrderAlreadyStoredError struct {
	OrderID uuid.UUID
}

func NewOrderAlreadyStoredError(orderID uuid.UUID) *OrderAlreadyStoredError {
	return &OrderAlreadyStoredError{OrderID: orderID}
}

func (e *OrderAlreadyStoredError) Error() string {
	return fmt.Sprintf("%s: %s", ErrOrderAlreadyStored, e.OrderID)
}

func (e *OrderAlreadyStoredError) Unwrap() error {
	return ErrOrderAlreadyStored
}

type OrderNotHasEnoughVolumeError struct {
	RequiredVolume  int
	AvailableVolume int
}

func NewOrderNotHasEnoughVolumeError(requiredVolume, availableVolume int) *OrderNotHasEnoughVolumeError {
	return &OrderNotHasEnoughVolumeError{
		RequiredVolume:  requiredVolume,
		AvailableVolume: availableVolume,
	}
}

func (e *OrderNotHasEnoughVolumeError) Error() string {
	return fmt.Sprintf("%s: required %d, available %d", ErrOrderNotHasEnoughVolume, e.RequiredVolume, e.AvailableVolume)
}

func (e *OrderNotHasEnoughVolumeError) Unwrap() error {
	return ErrOrderNotHasEnoughVolume
}
