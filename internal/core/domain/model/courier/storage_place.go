package courier

import (
	"delivery/internal/pkg/errs"
	"errors"

	"github.com/google/uuid"
)

type StoragePlace struct {
	id          uuid.UUID
	name        string
	totalVolume int
	orderId     *uuid.UUID
}

var (
	ErrOrderAlreadyStored      = errors.New("order already stored")
	ErrOrderNotHasEnoughVolume = errors.New("order not has enough volume")
)

func (sp *StoragePlace) CanStore(volume int) (bool, error) {
	if sp.isOccupied() {
		return false, ErrOrderAlreadyStored
	}
	if sp.totalVolume < volume {
		return false, ErrOrderNotHasEnoughVolume
	}
	return true, nil
}

func (sp *StoragePlace) Store(orderID uuid.UUID, volume int) error {
	if _, err := sp.CanStore(volume); err != nil {
		return err
	}

	sp.orderId = &orderID
	return nil
}

func (sp *StoragePlace) Clear(orderID uuid.UUID) error {
	if sp.orderId != nil && *sp.orderId != orderID {
		return errs.NewObjectNotFoundError("orderID", orderID)
	}
	sp.orderId = nil
	return nil
}

func (sp *StoragePlace) isOccupied() bool {
	return sp.orderId != nil
}

func NewStoragePlace(name string, totalVolume int) (*StoragePlace, error) {
	if name == "" {
		return nil, errs.NewValueIsInvalidError("name")
	}
	if totalVolume <= 0 {
		return nil, errs.NewValueIsInvalidError("totalVolume")
	}

	return &StoragePlace{
		id:          uuid.New(),
		name:        name,
		totalVolume: totalVolume,
	}, nil
}

func (sp *StoragePlace) ID() uuid.UUID {
	return sp.id
}

func (sp *StoragePlace) Name() string {
	return sp.name
}

func (sp *StoragePlace) TotalVolume() int {
	return sp.totalVolume
}

func (sp *StoragePlace) OrderID() *uuid.UUID {
	return sp.orderId
}

func (sp *StoragePlace) Equals(other *StoragePlace) bool {
	return sp.id == other.id
}
