package queries

import (
	"context"

	"delivery/internal/core/ports"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

type GetAllCouriersQuery struct{}

func NewGetAllCouriersQuery() GetAllCouriersQuery {
	return GetAllCouriersQuery{}
}

type GetAllCouriersResponse struct {
	Couriers []CourierResponse
}

type CourierResponse struct {
	ID       uuid.UUID
	Name     string
	Location LocationResponse
}

type GetAllCouriersQueryHandler interface {
	Handle(ctx context.Context, query GetAllCouriersQuery) (GetAllCouriersResponse, error)
}

var _ GetAllCouriersQueryHandler = &getAllCouriersQueryHandler{}

type getAllCouriersQueryHandler struct {
	courierRepo ports.CourierRepository
}

func NewGetAllCouriersQueryHandler(courierRepo ports.CourierRepository) (GetAllCouriersQueryHandler, error) {
	if courierRepo == nil {
		return nil, errs.NewValueIsRequiredError("courierRepo")
	}
	return &getAllCouriersQueryHandler{courierRepo: courierRepo}, nil
}

func (h *getAllCouriersQueryHandler) Handle(ctx context.Context, _ GetAllCouriersQuery) (GetAllCouriersResponse, error) {
	couriers, err := h.courierRepo.GetAll(ctx)
	if err != nil {
		return GetAllCouriersResponse{}, err
	}

	resp := GetAllCouriersResponse{Couriers: make([]CourierResponse, 0, len(couriers))}
	for _, c := range couriers {
		resp.Couriers = append(resp.Couriers, CourierResponse{
			ID:       c.ID(),
			Name:     c.Name(),
			Location: toLocationResponse(c.Location()),
		})
	}
	return resp, nil
}
