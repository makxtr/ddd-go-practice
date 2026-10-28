package commands

import (
	"context"
	"errors"

	"delivery/internal/core/ports"
)

// runInUnitOfWork выполняет работу внутри Unit of Work:
// при ошибке откатывает изменения, при успехе фиксирует их.
func runInUnitOfWork(ctx context.Context, uow ports.UnitOfWork, work func(ctx context.Context) error) error {
	if err := work(ctx); err != nil {
		return errors.Join(err, uow.Rollback(ctx))
	}
	return uow.Commit(ctx)
}
