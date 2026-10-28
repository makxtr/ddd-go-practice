package commands_test

import (
	"context"
	"testing"

	"delivery/internal/core/app/usecase/commands"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewAddStoragePlaceCommandReturnsErrorWhenParamsAreInvalid(t *testing.T) {
	tests := map[string]struct {
		courierID uuid.UUID
		name      string
		volume    int
		wantErr   error
	}{
		"nil courier id": {courierID: uuid.Nil, name: "Багажник", volume: 30, wantErr: errs.ErrValueIsRequired},
		"empty name":     {courierID: uuid.New(), name: "", volume: 30, wantErr: errs.ErrValueIsInvalid},
		"zero volume":    {courierID: uuid.New(), name: "Багажник", volume: 0, wantErr: errs.ErrValueIsInvalid},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := commands.NewAddStoragePlaceCommand(tc.courierID, tc.name, tc.volume)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func Test_AddStoragePlaceHandlerAddsPlaceToCourier(t *testing.T) {
	// Arrange
	env := newTestUoW()
	c := mustCourier(t, "Вело", 2, 1, 1)
	require.NoError(t, env.courierRepo.Add(context.Background(), c))
	handler, err := commands.NewAddStoragePlaceCommandHandler(env.uow)
	require.NoError(t, err)
	cmd, err := commands.NewAddStoragePlaceCommand(c.ID(), "Багажник", 30)
	require.NoError(t, err)

	// Act
	err = handler.Handle(context.Background(), cmd)

	// Assert
	require.NoError(t, err)
	got, _ := env.courierRepo.Get(context.Background(), c.ID())
	canTake, _ := got.CanTakeOrder(30) // default bag is 10, so only the new place fits
	assert.True(t, canTake)
}

func Test_AddStoragePlaceHandlerReturnsErrorWhenCourierNotFound(t *testing.T) {
	// Arrange
	env := newTestUoW()
	handler, _ := commands.NewAddStoragePlaceCommandHandler(env.uow)
	cmd, _ := commands.NewAddStoragePlaceCommand(uuid.New(), "Багажник", 30)

	// Act
	err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.ErrorIs(t, err, errs.ErrObjectNotFound)
}
