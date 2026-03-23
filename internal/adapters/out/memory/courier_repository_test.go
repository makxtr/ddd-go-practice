package memory

import (
	"context"
	"testing"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCourier(t *testing.T) *courier.Courier {
	t.Helper()
	loc, err := kernel.NewLocation(3, 3)
	require.NoError(t, err)
	c, err := courier.NewCourier("Иван", 2, loc)
	require.NoError(t, err)
	return c
}

func TestCourierRepository_Add_And_Get(t *testing.T) {
	// Arrange
	repo := NewCourierRepository()
	ctx := context.Background()
	c := newTestCourier(t)

	// Act
	err := repo.Add(ctx, c)

	// Assert
	assert.NoError(t, err)
	got, err := repo.Get(ctx, c.ID())
	assert.NoError(t, err)
	assert.Equal(t, c.ID(), got.ID())
}

func TestCourierRepository_Get_NotFound(t *testing.T) {
	// Arrange
	repo := NewCourierRepository()
	ctx := context.Background()

	// Act
	_, err := repo.Get(ctx, uuid.New())

	// Assert
	var target *errs.ObjectNotFoundError
	assert.ErrorAs(t, err, &target)
}

func TestCourierRepository_Update(t *testing.T) {
	// Arrange
	repo := NewCourierRepository()
	ctx := context.Background()
	c := newTestCourier(t)
	_ = repo.Add(ctx, c)

	target, _ := kernel.NewLocation(5, 5)
	_ = c.Move(target)

	// Act
	err := repo.Update(ctx, c)

	// Assert
	assert.NoError(t, err)
}

func TestCourierRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo := NewCourierRepository()
	ctx := context.Background()
	c := newTestCourier(t)

	// Act
	err := repo.Update(ctx, c)

	// Assert
	var target *errs.ObjectNotFoundError
	assert.ErrorAs(t, err, &target)
}

func TestCourierRepository_GetAllFree(t *testing.T) {
	// Arrange
	repo := NewCourierRepository()
	ctx := context.Background()
	c1 := newTestCourier(t)
	c2 := newTestCourier(t)
	_ = repo.Add(ctx, c1)
	_ = repo.Add(ctx, c2)

	// Занимаем одного курьера заказом
	_ = c1.TakeOrder(uuid.New(), 5)

	// Act
	free, err := repo.GetAllFree(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, free, 1)
	assert.Equal(t, c2.ID(), free[0].ID())
}

func TestCourierRepository_GetAllFree_AllBusy(t *testing.T) {
	// Arrange
	repo := NewCourierRepository()
	ctx := context.Background()
	c := newTestCourier(t)
	_ = repo.Add(ctx, c)
	_ = c.TakeOrder(uuid.New(), 5)

	// Act
	free, err := repo.GetAllFree(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, free)
}