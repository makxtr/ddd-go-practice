package order

import (
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type Order struct {
	id        uuid.UUID
	courierID *uuid.UUID
	location  kernel.Location
	volume    int
	status    Status
}

func (o *Order) Complete() error {
	if o.status != StatusAssigned {
		return errs.NewValueIsRequiredError("order status must be Assigned")
	}
	o.status = StatusCompleted
	return nil
}

func (o *Order) Assign(courierID uuid.UUID) error {
	if o.status != StatusCreated {
		return errs.NewValueIsRequiredError("order status must be Created")
	}

	o.courierID = &courierID
	o.status = StatusAssigned
	return nil
}

func NewOrder(orderID uuid.UUID, location kernel.Location, volume int) (*Order, error) {
	if orderID == uuid.Nil {
		return nil, errs.NewValueIsRequiredError("orderID")
	}
	if location.IsEmpty() {
		return nil, errs.NewValueIsRequiredError("location")
	}
	if volume <= 0 {
		return nil, errs.NewValueIsRequiredError("volume")
	}

	return &Order{
		id:       orderID,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func (o *Order) ID() uuid.UUID {
	return o.id
}

func (o *Order) CourierID() *uuid.UUID {
	return o.courierID
}

func (o *Order) Location() kernel.Location {
	return o.location
}

func (o *Order) Status() Status {
	return o.status
}

func (o *Order) Volume() int {
	return o.volume
}
