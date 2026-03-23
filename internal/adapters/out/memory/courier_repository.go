package memory

import (
	"context"
	"sync"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CourierRepository struct {
	mu       sync.RWMutex
	couriers map[uuid.UUID]*courier.Courier
}

func NewCourierRepository() *CourierRepository {
	return &CourierRepository{
		couriers: make(map[uuid.UUID]*courier.Courier),
	}
}

func (r *CourierRepository) Add(_ context.Context, aggregate *courier.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.couriers[aggregate.ID()] = aggregate
	return nil
}

func (r *CourierRepository) Update(_ context.Context, aggregate *courier.Courier) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.couriers[aggregate.ID()]; !exists {
		return errs.NewObjectNotFoundError("courier", aggregate.ID())
	}

	r.couriers[aggregate.ID()] = aggregate
	return nil
}

func (r *CourierRepository) Get(_ context.Context, id uuid.UUID) (*courier.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c, exists := r.couriers[id]
	if !exists {
		return nil, errs.NewObjectNotFoundError("courier", id)
	}

	return c, nil
}

func (r *CourierRepository) GetAllFree(_ context.Context) ([]*courier.Courier, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*courier.Courier
	for _, c := range r.couriers {
		if c.IsFree() {
			result = append(result, c)
		}
	}

	return result, nil
}