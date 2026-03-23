package memory

import (
	"context"
	"delivery/internal/core/ports"
)

type UnitOfWork struct {
	orderRepo   *OrderRepository
	courierRepo *CourierRepository
}

func NewUnitOfWork(orderRepo *OrderRepository, courierRepo *CourierRepository) *UnitOfWork {
	return &UnitOfWork{
		orderRepo:   orderRepo,
		courierRepo: courierRepo,
	}
}

func (u *UnitOfWork) OrderRepository() ports.OrderRepository {
	return u.orderRepo
}

func (u *UnitOfWork) CourierRepository() ports.CourierRepository {
	return u.courierRepo
}

func (u *UnitOfWork) Commit(_ context.Context) error {
	return nil
}

func (u *UnitOfWork) Rollback(_ context.Context) error {
	return nil
}