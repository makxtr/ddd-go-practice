package commands_test

import (
	"context"
	"testing"

	"delivery/internal/core/app/usecase/commands"
	"delivery/internal/core/domain/model/order"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MoveCouriersHandlerMovesCourierTowardsOrder(t *testing.T) {
	// Arrange
	env := newTestUoW()
	ctx := context.Background()
	c := mustCourier(t, "Пеший", 1, 1, 1)
	o := mustOrder(t, 1, 3, 5)
	require.NoError(t, o.Assign(c.ID()))
	require.NoError(t, c.TakeOrder(o.ID(), o.Volume()))
	require.NoError(t, env.courierRepo.Add(ctx, c))
	require.NoError(t, env.orderRepo.Add(ctx, o))
	handler, err := commands.NewMoveCouriersCommandHandler(env.uow)
	require.NoError(t, err)

	// Act
	err = handler.Handle(ctx, commands.NewMoveCouriersCommand())

	// Assert
	require.NoError(t, err)
	gotCourier, _ := env.courierRepo.Get(ctx, c.ID())
	assert.Equal(t, mustLocation(t, 1, 2), gotCourier.Location())
	gotOrder, _ := env.orderRepo.Get(ctx, o.ID())
	assert.Equal(t, order.StatusAssigned, gotOrder.Status())
	assert.False(t, gotCourier.IsFree())
}

func Test_MoveCouriersHandlerCompletesOrderWhenCourierArrives(t *testing.T) {
	// Arrange
	env := newTestUoW()
	ctx := context.Background()
	c := mustCourier(t, "Авто", 3, 1, 1)
	o := mustOrder(t, 1, 3, 5)
	require.NoError(t, o.Assign(c.ID()))
	require.NoError(t, c.TakeOrder(o.ID(), o.Volume()))
	require.NoError(t, env.courierRepo.Add(ctx, c))
	require.NoError(t, env.orderRepo.Add(ctx, o))
	handler, _ := commands.NewMoveCouriersCommandHandler(env.uow)

	// Act
	err := handler.Handle(ctx, commands.NewMoveCouriersCommand())

	// Assert
	require.NoError(t, err)
	gotCourier, _ := env.courierRepo.Get(ctx, c.ID())
	assert.Equal(t, mustLocation(t, 1, 3), gotCourier.Location())
	assert.True(t, gotCourier.IsFree())
	gotOrder, _ := env.orderRepo.Get(ctx, o.ID())
	assert.Equal(t, order.StatusCompleted, gotOrder.Status())
}

func Test_MoveCouriersHandlerDoesNothingWhenNoAssignedOrders(t *testing.T) {
	// Arrange
	env := newTestUoW()
	ctx := context.Background()
	c := mustCourier(t, "Пеший", 1, 1, 1)
	require.NoError(t, env.courierRepo.Add(ctx, c))
	require.NoError(t, env.orderRepo.Add(ctx, mustOrder(t, 5, 5, 5)))
	handler, _ := commands.NewMoveCouriersCommandHandler(env.uow)

	// Act
	err := handler.Handle(ctx, commands.NewMoveCouriersCommand())

	// Assert
	require.NoError(t, err)
	gotCourier, _ := env.courierRepo.Get(ctx, c.ID())
	assert.Equal(t, mustLocation(t, 1, 1), gotCourier.Location())
}
