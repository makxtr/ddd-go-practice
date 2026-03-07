package courier

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"
	"math"

	"github.com/google/uuid"
)

type Courier struct {
	id            uuid.UUID
	name          string
	location      kernel.Location
	speed         int
	storagePlaces []StoragePlace
}

func NewCourier(name string, speed int, location kernel.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequiredError("name")
	}
	if speed <= 0 {
		return nil, errs.NewValueIsRequiredError("speed")
	}
	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}

	bag, err := NewStoragePlace("Сумка", 10)
	if err != nil {
		return nil, err
	}

	return &Courier{
		id:            uuid.New(),
		name:          name,
		speed:         speed,
		location:      location,
		storagePlaces: []StoragePlace{*bag},
	}, nil
}

func (c *Courier) ID() uuid.UUID {
	return c.id
}

func (c *Courier) Name() string {
	return c.name
}

func (c *Courier) Move(target kernel.Location) error {
	if !target.IsValid() {
		return errs.NewValueIsRequiredError("target")
	}

	dx := float64(target.X() - c.location.X())
	dy := float64(target.Y() - c.location.Y())
	remainingRange := float64(c.speed)

	if math.Abs(dx) > remainingRange {
		dx = math.Copysign(remainingRange, dx)
	}
	remainingRange -= math.Abs(dx)

	if math.Abs(dy) > remainingRange {
		dy = math.Copysign(remainingRange, dy)
	}

	newX := c.location.X() + int(dx)
	newY := c.location.Y() + int(dy)

	newLocation, err := kernel.NewLocation(newX, newY)
	if err != nil {
		return err
	}
	c.location = newLocation
	return nil
}

func (c *Courier) AddStoragePlace(name string, volume int) error {
	sp, err := NewStoragePlace(name, volume)
	if err != nil {
		return err
	}
	c.storagePlaces = append(c.storagePlaces, *sp)
	return nil
}

func (c *Courier) CanTakeOrder(volume int) (bool, error) {
	for i := range c.storagePlaces {
		canStore, _ := c.storagePlaces[i].CanStore(volume)
		if canStore {
			return true, nil
		}
	}
	return false, nil
}

func (c *Courier) TakeOrder(orderID uuid.UUID, volume int) error {
	for i := range c.storagePlaces {
		canStore, _ := c.storagePlaces[i].CanStore(volume)
		if canStore {
			return c.storagePlaces[i].Store(orderID, volume)
		}
	}
	return errs.NewValueIsRequiredError("no available storage place")
}

func (c *Courier) CompleteOrder(orderID uuid.UUID) error {
	sp, err := c.findStoragePlaceByOrderId(orderID)
	if err != nil {
		return err
	}
	return sp.Clear(orderID)
}

func (c *Courier) CalculateTimeToLocation(target kernel.Location) (int, error) {
	distance, err := c.location.DistanceTo(target)
	if err != nil {
		return 0, err
	}
	steps := (distance + c.speed - 1) / c.speed // округление вверх
	return steps, nil
}

func (c *Courier) findStoragePlaceByOrderId(orderID uuid.UUID) (*StoragePlace, error) {
	for i := range c.storagePlaces {
		if c.storagePlaces[i].orderId != nil && *c.storagePlaces[i].orderId == orderID {
			return &c.storagePlaces[i], nil
		}
	}
	return nil, errs.NewObjectNotFoundError("storagePlace", orderID)
}
