package cmd

import (
	"log"

	"delivery/internal/adapters/out/memory"
	"delivery/internal/core/app/usecase/commands"
	"delivery/internal/core/app/usecase/queries"
	"delivery/internal/core/domain/services"
	"delivery/internal/core/ports"
)

type CompositionRoot struct {
	configs Config

	closers []Closer

	unitOfWork      ports.UnitOfWork
	orderDispatcher services.OrderDispatcher

	createOrderHandler     commands.CreateOrderCommandHandler
	createCourierHandler   commands.CreateCourierCommandHandler
	addStoragePlaceHandler commands.AddStoragePlaceCommandHandler
	moveCouriersHandler    commands.MoveCouriersCommandHandler
	assignOrdersHandler    commands.AssignOrdersCommandHandler

	getAllCouriersHandler        queries.GetAllCouriersQueryHandler
	getNotCompletedOrdersHandler queries.GetNotCompletedOrdersQueryHandler
}

func NewCompositionRoot(configs Config) *CompositionRoot {
	// Adapters (out)
	orderRepo := memory.NewOrderRepository()
	courierRepo := memory.NewCourierRepository()
	uow := memory.NewUnitOfWork(orderRepo, courierRepo)

	// Domain services
	orderDispatcher := services.NewOrderDispatcher()

	// Commands
	createOrderHandler, err := commands.NewCreateOrderCommandHandler(uow)
	exitOnError(err)
	createCourierHandler, err := commands.NewCreateCourierCommandHandler(uow)
	exitOnError(err)
	addStoragePlaceHandler, err := commands.NewAddStoragePlaceCommandHandler(uow)
	exitOnError(err)
	moveCouriersHandler, err := commands.NewMoveCouriersCommandHandler(uow)
	exitOnError(err)
	assignOrdersHandler, err := commands.NewAssignOrdersCommandHandler(uow, orderDispatcher)
	exitOnError(err)

	// Queries
	getAllCouriersHandler, err := queries.NewGetAllCouriersQueryHandler(courierRepo)
	exitOnError(err)
	getNotCompletedOrdersHandler, err := queries.NewGetNotCompletedOrdersQueryHandler(orderRepo)
	exitOnError(err)

	return &CompositionRoot{
		configs:         configs,
		unitOfWork:      uow,
		orderDispatcher: orderDispatcher,

		createOrderHandler:     createOrderHandler,
		createCourierHandler:   createCourierHandler,
		addStoragePlaceHandler: addStoragePlaceHandler,
		moveCouriersHandler:    moveCouriersHandler,
		assignOrdersHandler:    assignOrdersHandler,

		getAllCouriersHandler:        getAllCouriersHandler,
		getNotCompletedOrdersHandler: getNotCompletedOrdersHandler,
	}
}

func (cr *CompositionRoot) UnitOfWork() ports.UnitOfWork {
	return cr.unitOfWork
}

func (cr *CompositionRoot) OrderDispatcher() services.OrderDispatcher {
	return cr.orderDispatcher
}

func (cr *CompositionRoot) CreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	return cr.createOrderHandler
}

func (cr *CompositionRoot) CreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	return cr.createCourierHandler
}

func (cr *CompositionRoot) AddStoragePlaceCommandHandler() commands.AddStoragePlaceCommandHandler {
	return cr.addStoragePlaceHandler
}

func (cr *CompositionRoot) MoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	return cr.moveCouriersHandler
}

func (cr *CompositionRoot) AssignOrdersCommandHandler() commands.AssignOrdersCommandHandler {
	return cr.assignOrdersHandler
}

func (cr *CompositionRoot) GetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	return cr.getAllCouriersHandler
}

func (cr *CompositionRoot) GetNotCompletedOrdersQueryHandler() queries.GetNotCompletedOrdersQueryHandler {
	return cr.getNotCompletedOrdersHandler
}

func exitOnError(err error) {
	if err != nil {
		log.Fatalf("composition root: %v", err)
	}
}
