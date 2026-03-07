package services

import (
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"
	"math"
)

type OrderDispatcher interface {
	Dispatch(*order.Order, []*courier.Courier) (*courier.Courier, error)
}

var _ OrderDispatcher = &orderDispatcher{}

type orderDispatcher struct{}

func NewOrderDispatcher() OrderDispatcher {
	return &orderDispatcher{}
}

func (*orderDispatcher) Dispatch(ord *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if ord.Status() != order.StatusCreated {
		return nil, errs.NewValueIsRequiredError("order status must be Created")
	}

	var bestCourier *courier.Courier
	bestScore := math.MaxInt

	for _, c := range couriers {
		canTake, err := c.CanTakeOrder(ord.Volume())
		if err != nil || !canTake {
			continue
		}

		score, err := c.CalculateTimeToLocation(ord.Location())
		if err != nil {
			continue
		}
		if score < bestScore {
			bestScore = score
			bestCourier = c
		}
	}

	if bestCourier == nil {
		return nil, errs.NewObjectNotFoundError("courier", nil)
	}

	if err := ord.Assign(bestCourier.ID()); err != nil {
		return nil, err
	}

	if err := bestCourier.TakeOrder(ord.ID(), ord.Volume()); err != nil {
		return nil, err
	}

	return bestCourier, nil
}
