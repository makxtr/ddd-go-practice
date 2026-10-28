package commands

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

// MoveCouriersCommand не имеет полей: команда перемещает всех курьеров,
// у которых есть назначенный заказ, на один шаг к заказу.
type MoveCouriersCommand struct{}

func NewMoveCouriersCommand() MoveCouriersCommand {
	return MoveCouriersCommand{}
}

type MoveCouriersCommandHandler interface {
	Handle(ctx context.Context, cmd MoveCouriersCommand) error
}

var _ MoveCouriersCommandHandler = &moveCouriersCommandHandler{}

type moveCouriersCommandHandler struct {
	uow ports.UnitOfWork
}

func NewMoveCouriersCommandHandler(uow ports.UnitOfWork) (MoveCouriersCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}
	return &moveCouriersCommandHandler{uow: uow}, nil
}

func (h *moveCouriersCommandHandler) Handle(ctx context.Context, _ MoveCouriersCommand) error {
	return runInUnitOfWork(ctx, h.uow, func(ctx context.Context) error {
		orderRepo := h.uow.OrderRepository()
		courierRepo := h.uow.CourierRepository()

		assignedOrders, err := orderRepo.GetAllAssigned(ctx)
		if err != nil {
			return err
		}

		for _, ord := range assignedOrders {
			if ord.CourierID() == nil {
				return errs.NewValueIsRequiredError("courierID")
			}

			c, err := courierRepo.Get(ctx, *ord.CourierID())
			if err != nil {
				return err
			}

			if err := c.Move(ord.Location()); err != nil {
				return err
			}

			if c.Location().Equals(ord.Location()) {
				if err := ord.Complete(); err != nil {
					return err
				}
				if err := c.CompleteOrder(ord.ID()); err != nil {
					return err
				}
				if err := orderRepo.Update(ctx, ord); err != nil {
					return err
				}
			}

			if err := courierRepo.Update(ctx, c); err != nil {
				return err
			}
		}

		return nil
	})
}
