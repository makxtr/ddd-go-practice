package cmd

import (
	"delivery/internal/adapters/out/memory"
	"delivery/internal/core/ports"
)

type CompositionRoot struct {
	configs Config

	closers []Closer

	unitOfWork ports.UnitOfWork
}

func NewCompositionRoot(configs Config) *CompositionRoot {
	orderRepo := memory.NewOrderRepository()
	courierRepo := memory.NewCourierRepository()
	uow := memory.NewUnitOfWork(orderRepo, courierRepo)

	return &CompositionRoot{
		configs:    configs,
		unitOfWork: uow,
	}
}

func (cr *CompositionRoot) UnitOfWork() ports.UnitOfWork {
	return cr.unitOfWork
}
