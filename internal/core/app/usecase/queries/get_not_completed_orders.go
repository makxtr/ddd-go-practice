package queries

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type GetNotCompletedOrdersQuery struct{}

func NewGetNotCompletedOrdersQuery() GetNotCompletedOrdersQuery {
	return GetNotCompletedOrdersQuery{}
}

type GetNotCompletedOrdersResponse struct {
	Orders []OrderResponse
}

type OrderResponse struct {
	ID       uuid.UUID
	Location LocationResponse
}

type GetNotCompletedOrdersQueryHandler interface {
	Handle(ctx context.Context, query GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error)
}

var _ GetNotCompletedOrdersQueryHandler = &getNotCompletedOrdersQueryHandler{}

type getNotCompletedOrdersQueryHandler struct {
	orderRepo ports.OrderRepository
}

func NewGetNotCompletedOrdersQueryHandler(orderRepo ports.OrderRepository) (GetNotCompletedOrdersQueryHandler, error) {
	if orderRepo == nil {
		return nil, errs.NewValueIsRequiredError("orderRepo")
	}
	return &getNotCompletedOrdersQueryHandler{orderRepo: orderRepo}, nil
}

func (h *getNotCompletedOrdersQueryHandler) Handle(ctx context.Context, _ GetNotCompletedOrdersQuery) (GetNotCompletedOrdersResponse, error) {
	orders, err := h.orderRepo.GetAllNotCompleted(ctx)
	if err != nil {
		return GetNotCompletedOrdersResponse{}, err
	}

	resp := GetNotCompletedOrdersResponse{Orders: make([]OrderResponse, 0, len(orders))}
	for _, o := range orders {
		resp.Orders = append(resp.Orders, OrderResponse{
			ID:       o.ID(),
			Location: toLocationResponse(o.Location()),
		})
	}
	return resp, nil
}
