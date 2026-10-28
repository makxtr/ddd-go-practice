package commands_test

import (
	"context"
	"testing"

	"delivery/internal/core/app/usecase/commands"
	"delivery/internal/pkg/errs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCreateCourierCommandReturnsErrorWhenParamsAreInvalid(t *testing.T) {
	tests := map[string]struct {
		name  string
		speed int
	}{
		"empty name": {name: "", speed: 2},
		"zero speed": {name: "Вело", speed: 0},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := commands.NewCreateCourierCommand(tc.name, tc.speed)
			assert.ErrorIs(t, err, errs.ErrValueIsRequired)
		})
	}
}

func Test_CreateCourierHandlerStoresCourier(t *testing.T) {
	// Arrange
	env := newTestUoW()
	handler, err := commands.NewCreateCourierCommandHandler(env.uow)
	require.NoError(t, err)
	cmd, err := commands.NewCreateCourierCommand("Вело", 2)
	require.NoError(t, err)

	// Act
	err = handler.Handle(context.Background(), cmd)

	// Assert
	require.NoError(t, err)
	all, _ := env.courierRepo.GetAll(context.Background())
	require.Len(t, all, 1)
	assert.Equal(t, "Вело", all[0].Name())
	assert.False(t, all[0].Location().IsEmpty())
}
