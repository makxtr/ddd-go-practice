package commands

import (
	"context"
	"errors"

	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/domain/model/order"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type CreateOrderCommand struct {
	orderID uuid.UUID
	street  string
	volume  int
}

func NewCreateOrderCommand(orderID uuid.UUID, street string, volume int) (CreateOrderCommand, error) {
	if orderID == uuid.Nil {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("orderID")
	}
	if street == "" {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("street")
	}
	if volume <= 0 {
		return CreateOrderCommand{}, errs.NewValueIsRequiredError("volume")
	}
	return CreateOrderCommand{orderID: orderID, street: street, volume: volume}, nil
}

func (c CreateOrderCommand) OrderID() uuid.UUID { return c.orderID }
func (c CreateOrderCommand) Street() string     { return c.street }
func (c CreateOrderCommand) Volume() int        { return c.volume }

type CreateOrderCommandHandler interface {
	Handle(ctx context.Context, cmd CreateOrderCommand) error
}

var _ CreateOrderCommandHandler = &createOrderCommandHandler{}

type createOrderCommandHandler struct {
	uow ports.UnitOfWork
}

func NewCreateOrderCommandHandler(uow ports.UnitOfWork) (CreateOrderCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}
	return &createOrderCommandHandler{uow: uow}, nil
}

func (h *createOrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCommand) error {
	return runInUnitOfWork(ctx, h.uow, func(ctx context.Context) error {
		orderRepo := h.uow.OrderRepository()

		// Сообщение из Kafka может прийти повторно — заказ с таким ID уже создан.
		existing, err := orderRepo.Get(ctx, cmd.OrderID())
		if err != nil && !errors.Is(err, errs.ErrObjectNotFound) {
			return err
		}
		if existing != nil {
			return nil
		}

		// TODO: получать Location по Street из сервиса Geo (следующие модули).
		location := kernel.NewRandomLocation()

		ord, err := order.NewOrder(cmd.OrderID(), location, cmd.Volume())
		if err != nil {
			return err
		}

		return orderRepo.Add(ctx, ord)
	})
}
