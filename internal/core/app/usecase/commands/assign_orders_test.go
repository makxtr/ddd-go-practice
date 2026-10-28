package commands_test

import (
	"context"
	"testing"

	"delivery/internal/core/app/usecase/commands"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/domain/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_AssignOrdersHandlerAssignsOrderToClosestFreeCourier(t *testing.T) {
	// Arrange
	env := newTestUoW()
	ctx := context.Background()
	far := mustCourier(t, "Далеко", 1, 10, 10)
	near := mustCourier(t, "Рядом", 1, 2, 2)
	o := mustOrder(t, 1, 1, 5)
	require.NoError(t, env.courierRepo.Add(ctx, far))
	require.NoError(t, env.courierRepo.Add(ctx, near))
	require.NoError(t, env.orderRepo.Add(ctx, o))
	handler, err := commands.NewAssignOrdersCommandHandler(env.uow, services.NewOrderDispatcher())
	require.NoError(t, err)

	// Act
	err = handler.Handle(ctx, commands.NewAssignOrdersCommand())

	// Assert
	require.NoError(t, err)
	gotOrder, _ := env.orderRepo.Get(ctx, o.ID())
	assert.Equal(t, order.StatusAssigned, gotOrder.Status())
	require.NotNil(t, gotOrder.CourierID())
	assert.Equal(t, near.ID(), *gotOrder.CourierID())
	gotNear, _ := env.courierRepo.Get(ctx, near.ID())
	assert.False(t, gotNear.IsFree())
	gotFar, _ := env.courierRepo.Get(ctx, far.ID())
	assert.True(t, gotFar.IsFree())
}

func Test_AssignOrdersHandlerDoesNothingWhenNoCreatedOrders(t *testing.T) {
	// Arrange
	env := newTestUoW()
	ctx := context.Background()
	require.NoError(t, env.courierRepo.Add(ctx, mustCourier(t, "Рядом", 1, 2, 2)))
	handler, _ := commands.NewAssignOrdersCommandHandler(env.uow, services.NewOrderDispatcher())

	// Act
	err := handler.Handle(ctx, commands.NewAssignOrdersCommand())

	// Assert
	assert.NoError(t, err)
}

func Test_AssignOrdersHandlerDoesNothingWhenNoFreeCouriers(t *testing.T) {
	// Arrange
	env := newTestUoW()
	ctx := context.Background()
	o := mustOrder(t, 1, 1, 5)
	require.NoError(t, env.orderRepo.Add(ctx, o))
	handler, _ := commands.NewAssignOrdersCommandHandler(env.uow, services.NewOrderDispatcher())

	// Act
	err := handler.Handle(ctx, commands.NewAssignOrdersCommand())

	// Assert
	require.NoError(t, err)
	gotOrder, _ := env.orderRepo.Get(ctx, o.ID())
	assert.Equal(t, order.StatusCreated, gotOrder.Status())
}

func Test_AssignOrdersHandlerDoesNothingWhenNoCourierCanTakeOrder(t *testing.T) {
	// Arrange: единственный свободный курьер имеет только сумку на 10, заказ на 15
	env := newTestUoW()
	ctx := context.Background()
	c := mustCourier(t, "Пеший", 1, 2, 2)
	o := mustOrder(t, 1, 1, 15)
	require.NoError(t, env.courierRepo.Add(ctx, c))
	require.NoError(t, env.orderRepo.Add(ctx, o))
	handler, _ := commands.NewAssignOrdersCommandHandler(env.uow, services.NewOrderDispatcher())

	// Act
	err := handler.Handle(ctx, commands.NewAssignOrdersCommand())

	// Assert
	require.NoError(t, err)
	gotOrder, _ := env.orderRepo.Get(ctx, o.ID())
	assert.Equal(t, order.StatusCreated, gotOrder.Status())
	gotCourier, _ := env.courierRepo.Get(ctx, c.ID())
	assert.True(t, gotCourier.IsFree())
}
