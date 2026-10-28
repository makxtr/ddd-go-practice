package queries

import "delivery/internal/core/domain/model/kernel"

type LocationResponse struct {
	X int
	Y int
}

func toLocationResponse(location kernel.Location) LocationResponse {
	return LocationResponse{X: location.X(), Y: location.Y()}
}
