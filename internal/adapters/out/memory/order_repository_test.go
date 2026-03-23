package memory

import (
	"context"
	"testing"

	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestOrder(t *testing.T) *order.Order {
	t.Helper()
	loc, err := kernel.NewLocation(5, 5)
	require.NoError(t, err)
	o, err := order.NewOrder(uuid.New(), loc, 3)
	require.NoError(t, err)
	return o
}

func TestOrderRepository_Add_And_Get(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()
	o := newTestOrder(t)

	// Act
	err := repo.Add(ctx, o)

	// Assert
	assert.NoError(t, err)
	got, err := repo.Get(ctx, o.ID())
	assert.NoError(t, err)
	assert.Equal(t, o.ID(), got.ID())
}

func TestOrderRepository_Get_NotFound(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()

	// Act
	_, err := repo.Get(ctx, uuid.New())

	// Assert
	var target *errs.ObjectNotFoundError
	assert.ErrorAs(t, err, &target)
}

func TestOrderRepository_Update(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()
	o := newTestOrder(t)
	_ = repo.Add(ctx, o)
	courierID := uuid.New()
	_ = o.Assign(courierID)

	// Act
	err := repo.Update(ctx, o)

	// Assert
	assert.NoError(t, err)
	got, _ := repo.Get(ctx, o.ID())
	assert.Equal(t, order.StatusAssigned, got.Status())
}

func TestOrderRepository_Update_NotFound(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()
	o := newTestOrder(t)

	// Act
	err := repo.Update(ctx, o)

	// Assert
	var target *errs.ObjectNotFoundError
	assert.ErrorAs(t, err, &target)
}

func TestOrderRepository_GetFirstCreated(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()
	o1 := newTestOrder(t)
	o2 := newTestOrder(t)
	_ = repo.Add(ctx, o1)
	_ = repo.Add(ctx, o2)
	_ = o1.Assign(uuid.New())
	_ = repo.Update(ctx, o1)

	// Act
	got, err := repo.GetFirstCreated(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, order.StatusCreated, got.Status())
}

func TestOrderRepository_GetFirstCreated_NotFound(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()

	// Act
	_, err := repo.GetFirstCreated(ctx)

	// Assert
	var target *errs.ObjectNotFoundError
	assert.ErrorAs(t, err, &target)
}

func TestOrderRepository_GetAllAssigned(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()
	o1 := newTestOrder(t)
	o2 := newTestOrder(t)
	o3 := newTestOrder(t)
	_ = repo.Add(ctx, o1)
	_ = repo.Add(ctx, o2)
	_ = repo.Add(ctx, o3)
	_ = o1.Assign(uuid.New())
	_ = o2.Assign(uuid.New())

	// Act
	assigned, err := repo.GetAllAssigned(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, assigned, 2)
}

func TestOrderRepository_GetAllAssigned_Empty(t *testing.T) {
	// Arrange
	repo := NewOrderRepository()
	ctx := context.Background()
	_ = repo.Add(ctx, newTestOrder(t))

	// Act
	assigned, err := repo.GetAllAssigned(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, assigned)
}