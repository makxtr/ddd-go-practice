package memory

import (
	"context"
	"sync"

	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*order.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[uuid.UUID]*order.Order),
	}
}

func (r *OrderRepository) Add(_ context.Context, aggregate *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[aggregate.ID()] = aggregate
	return nil
}

func (r *OrderRepository) Update(_ context.Context, aggregate *order.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[aggregate.ID()]; !exists {
		return errs.NewObjectNotFoundError("order", aggregate.ID())
	}

	r.orders[aggregate.ID()] = aggregate
	return nil
}

func (r *OrderRepository) Get(_ context.Context, id uuid.UUID) (*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	o, exists := r.orders[id]
	if !exists {
		return nil, errs.NewObjectNotFoundError("order", id)
	}

	return o, nil
}

func (r *OrderRepository) GetFirstCreated(_ context.Context) (*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, o := range r.orders {
		if o.Status() == order.StatusCreated {
			return o, nil
		}
	}

	return nil, errs.NewObjectNotFoundError("order", "status=Created")
}

func (r *OrderRepository) GetAllAssigned(_ context.Context) ([]*order.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*order.Order
	for _, o := range r.orders {
		if o.Status() == order.StatusAssigned {
			result = append(result, o)
		}
	}

	return result, nil
}