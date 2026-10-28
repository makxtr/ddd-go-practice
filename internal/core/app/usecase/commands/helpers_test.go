package commands_test

import (
	"testing"

	"delivery/internal/adapters/out/memory"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type testUoW struct {
	uow         *memory.UnitOfWork
	orderRepo   *memory.OrderRepository
	courierRepo *memory.CourierRepository
}

func newTestUoW() testUoW {
	orderRepo := memory.NewOrderRepository()
	courierRepo := memory.NewCourierRepository()
	return testUoW{
		uow:         memory.NewUnitOfWork(orderRepo, courierRepo),
		orderRepo:   orderRepo,
		courierRepo: courierRepo,
	}
}

func mustLocation(t *testing.T, x, y int) kernel.Location {
	t.Helper()
	loc, err := kernel.NewLocation(x, y)
	require.NoError(t, err)
	return loc
}

func mustCourier(t *testing.T, name string, speed int, x, y int) *courier.Courier {
	t.Helper()
	c, err := courier.NewCourier(name, speed, mustLocation(t, x, y))
	require.NoError(t, err)
	return c
}

func mustOrder(t *testing.T, x, y int, volume int) *order.Order {
	t.Helper()
	o, err := order.NewOrder(uuid.New(), mustLocation(t, x, y), volume)
	require.NoError(t, err)
	return o
}
