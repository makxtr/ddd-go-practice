package commands

import (
	"context"
	"errors"

	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

// AssignOrdersCommand не имеет полей: команда берёт первый неназначенный заказ
// и назначает его на самого подходящего свободного курьера.
type AssignOrdersCommand struct{}

func NewAssignOrdersCommand() AssignOrdersCommand {
	return AssignOrdersCommand{}
}

type AssignOrdersCommandHandler interface {
	Handle(ctx context.Context, cmd AssignOrdersCommand) error
}

var _ AssignOrdersCommandHandler = &assignOrdersCommandHandler{}

type assignOrdersCommandHandler struct {
	uow        ports.UnitOfWork
	dispatcher services.OrderDispatcher
}

func NewAssignOrdersCommandHandler(uow ports.UnitOfWork, dispatcher services.OrderDispatcher) (AssignOrdersCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}
	if dispatcher == nil {
		return nil, errs.NewValueIsRequiredError("dispatcher")
	}
	return &assignOrdersCommandHandler{uow: uow, dispatcher: dispatcher}, nil
}

func (h *assignOrdersCommandHandler) Handle(ctx context.Context, _ AssignOrdersCommand) error {
	return runInUnitOfWork(ctx, h.uow, func(ctx context.Context) error {
		orderRepo := h.uow.OrderRepository()
		courierRepo := h.uow.CourierRepository()

		ord, err := orderRepo.GetFirstCreated(ctx)
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil // нет неназначенных заказов — делать нечего
		}
		if err != nil {
			return err
		}

		couriers, err := courierRepo.GetAllFree(ctx)
		if err != nil {
			return err
		}
		if len(couriers) == 0 {
			return nil // нет свободных курьеров — попробуем на следующем тике
		}

		c, err := h.dispatcher.Dispatch(ord, couriers)
		if errors.Is(err, errs.ErrObjectNotFound) {
			return nil // ни один свободный курьер не вмещает заказ — попробуем на следующем тике
		}
		if err != nil {
			return err
		}

		if err := orderRepo.Update(ctx, ord); err != nil {
			return err
		}
		return courierRepo.Update(ctx, c)
	})
}
