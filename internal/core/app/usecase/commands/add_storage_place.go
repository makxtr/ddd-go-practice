package commands

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type AddStoragePlaceCommand struct {
	courierID   uuid.UUID
	name        string
	totalVolume int
}

func NewAddStoragePlaceCommand(courierID uuid.UUID, name string, totalVolume int) (AddStoragePlaceCommand, error) {
	if courierID == uuid.Nil {
		return AddStoragePlaceCommand{}, errs.NewValueIsRequiredError("courierID")
	}
	if name == "" {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("name")
	}
	if totalVolume <= 0 {
		return AddStoragePlaceCommand{}, errs.NewValueIsInvalidError("totalVolume")
	}
	return AddStoragePlaceCommand{courierID: courierID, name: name, totalVolume: totalVolume}, nil
}

func (c AddStoragePlaceCommand) CourierID() uuid.UUID { return c.courierID }
func (c AddStoragePlaceCommand) Name() string         { return c.name }
func (c AddStoragePlaceCommand) TotalVolume() int     { return c.totalVolume }

type AddStoragePlaceCommandHandler interface {
	Handle(ctx context.Context, cmd AddStoragePlaceCommand) error
}

var _ AddStoragePlaceCommandHandler = &addStoragePlaceCommandHandler{}

type addStoragePlaceCommandHandler struct {
	uow ports.UnitOfWork
}

func NewAddStoragePlaceCommandHandler(uow ports.UnitOfWork) (AddStoragePlaceCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}
	return &addStoragePlaceCommandHandler{uow: uow}, nil
}

func (h *addStoragePlaceCommandHandler) Handle(ctx context.Context, cmd AddStoragePlaceCommand) error {
	return runInUnitOfWork(ctx, h.uow, func(ctx context.Context) error {
		courierRepo := h.uow.CourierRepository()

		c, err := courierRepo.Get(ctx, cmd.CourierID())
		if err != nil {
			return err
		}
		if err := c.AddStoragePlace(cmd.Name(), cmd.TotalVolume()); err != nil {
			return err
		}
		return courierRepo.Update(ctx, c)
	})
}
