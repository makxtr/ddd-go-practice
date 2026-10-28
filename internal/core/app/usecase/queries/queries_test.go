package queries_test

import (
	"context"
	"testing"

	"delivery/internal/adapters/out/memory"
	"delivery/internal/core/app/usecase/queries"
	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustLocation(t *testing.T, x, y int) kernel.Location {
	t.Helper()
	loc, err := kernel.NewLocation(x, y)
	require.NoError(t, err)
	return loc
}

func Test_GetAllCouriersReturnsEveryCourierWithLocation(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := memory.NewCourierRepository()
	c, err := courier.NewCourier("Вело", 2, mustLocation(t, 3, 4))
	require.NoError(t, err)
	require.NoError(t, repo.Add(ctx, c))
	handler, err := queries.NewGetAllCouriersQueryHandler(repo)
	require.NoError(t, err)

	// Act
	resp, err := handler.Handle(ctx, queries.NewGetAllCouriersQuery())

	// Assert
	require.NoError(t, err)
	require.Len(t, resp.Couriers, 1)
	assert.Equal(t, queries.CourierResponse{
		ID:       c.ID(),
		Name:     "Вело",
		Location: queries.LocationResponse{X: 3, Y: 4},
	}, resp.Couriers[0])
}

func Test_NewGetAllCouriersQueryHandlerReturnsErrorWhenRepoIsNil(t *testing.T) {
	_, err := queries.NewGetAllCouriersQueryHandler(nil)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}

func Test_GetNotCompletedOrdersReturnsCreatedAndAssignedOnly(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repo := memory.NewOrderRepository()
	created, _ := order.NewOrder(uuid.New(), mustLocation(t, 1, 2), 5)
	assigned, _ := order.NewOrder(uuid.New(), mustLocation(t, 3, 4), 5)
	require.NoError(t, assigned.Assign(uuid.New()))
	completed, _ := order.NewOrder(uuid.New(), mustLocation(t, 5, 6), 5)
	require.NoError(t, completed.Assign(uuid.New()))
	require.NoError(t, completed.Complete())
	for _, o := range []*order.Order{created, assigned, completed} {
		require.NoError(t, repo.Add(ctx, o))
	}
	handler, err := queries.NewGetNotCompletedOrdersQueryHandler(repo)
	require.NoError(t, err)

	// Act
	resp, err := handler.Handle(ctx, queries.NewGetNotCompletedOrdersQuery())

	// Assert
	require.NoError(t, err)
	assert.ElementsMatch(t, []queries.OrderResponse{
		{ID: created.ID(), Location: queries.LocationResponse{X: 1, Y: 2}},
		{ID: assigned.ID(), Location: queries.LocationResponse{X: 3, Y: 4}},
	}, resp.Orders)
}
