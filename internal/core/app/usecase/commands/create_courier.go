package commands

import (
	"context"

	"delivery/internal/core/domain/model/courier"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"
)

type CreateCourierCommand struct {
	name  string
	speed int
}

func NewCreateCourierCommand(name string, speed int) (CreateCourierCommand, error) {
	if name == "" {
		return CreateCourierCommand{}, errs.NewValueIsRequiredError("name")
	}
	if speed <= 0 {
		return CreateCourierCommand{}, errs.NewValueIsRequiredError("speed")
	}
	return CreateCourierCommand{name: name, speed: speed}, nil
}

func (c CreateCourierCommand) Name() string { return c.name }
func (c CreateCourierCommand) Speed() int   { return c.speed }

type CreateCourierCommandHandler interface {
	Handle(ctx context.Context, cmd CreateCourierCommand) error
}

var _ CreateCourierCommandHandler = &createCourierCommandHandler{}

type createCourierCommandHandler struct {
	uow ports.UnitOfWork
}

func NewCreateCourierCommandHandler(uow ports.UnitOfWork) (CreateCourierCommandHandler, error) {
	if uow == nil {
		return nil, errs.NewValueIsRequiredError("uow")
	}
	return &createCourierCommandHandler{uow: uow}, nil
}

func (h *createCourierCommandHandler) Handle(ctx context.Context, cmd CreateCourierCommand) error {
	return runInUnitOfWork(ctx, h.uow, func(ctx context.Context) error {
		c, err := courier.NewCourier(cmd.Name(), cmd.Speed(), kernel.NewRandomLocation())
		if err != nil {
			return err
		}
		return h.uow.CourierRepository().Add(ctx, c)
	})
}
