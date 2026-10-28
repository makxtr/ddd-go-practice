package commands_test

import (
	"context"
	"testing"

	"delivery/internal/core/app/usecase/commands"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCreateOrderCommandReturnsErrorWhenParamsAreInvalid(t *testing.T) {
	tests := map[string]struct {
		orderID uuid.UUID
		street  string
		volume  int
	}{
		"nil order id":  {orderID: uuid.Nil, street: "Ленина", volume: 5},
		"empty street":  {orderID: uuid.New(), street: "", volume: 5},
		"zero volume":   {orderID: uuid.New(), street: "Ленина", volume: 0},
		"negative vol.": {orderID: uuid.New(), street: "Ленина", volume: -1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := commands.NewCreateOrderCommand(tc.orderID, tc.street, tc.volume)
			assert.ErrorIs(t, err, errs.ErrValueIsRequired)
		})
	}
}

func Test_CreateOrderHandlerStoresCreatedOrderWithRandomLocation(t *testing.T) {
	// Arrange
	env := newTestUoW()
	handler, err := commands.NewCreateOrderCommandHandler(env.uow)
	require.NoError(t, err)
	orderID := uuid.New()
	cmd, err := commands.NewCreateOrderCommand(orderID, "Ленина", 5)
	require.NoError(t, err)

	// Act
	err = handler.Handle(context.Background(), cmd)

	// Assert
	require.NoError(t, err)
	got, err := env.orderRepo.Get(context.Background(), orderID)
	require.NoError(t, err)
	assert.Equal(t, order.StatusCreated, got.Status())
	assert.Equal(t, 5, got.Volume())
	assert.False(t, got.Location().IsEmpty())
}

func Test_CreateOrderHandlerIsIdempotentForSameOrderID(t *testing.T) {
	// Arrange
	env := newTestUoW()
	handler, _ := commands.NewCreateOrderCommandHandler(env.uow)
	orderID := uuid.New()
	cmd, _ := commands.NewCreateOrderCommand(orderID, "Ленина", 5)
	require.NoError(t, handler.Handle(context.Background(), cmd))

	// Act
	err := handler.Handle(context.Background(), cmd)

	// Assert
	require.NoError(t, err)
	all, _ := env.orderRepo.GetAllNotCompleted(context.Background())
	assert.Len(t, all, 1)
}

func Test_NewCreateOrderCommandHandlerReturnsErrorWhenUoWIsNil(t *testing.T) {
	_, err := commands.NewCreateOrderCommandHandler(nil)
	assert.ErrorIs(t, err, errs.ErrValueIsRequired)
}
