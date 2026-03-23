package ports

import "context"

type UnitOfWork interface {
	OrderRepository() OrderRepository
	CourierRepository() CourierRepository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}