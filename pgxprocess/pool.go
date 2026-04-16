package pgxprocess

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"go.pact.im/x/process"
)

var _ database = (*Pool)(nil)

// Pool wraps [pgxpool.Pool] to delay pool setup to application runtime.
type Pool struct {
	Config *pgxpool.Config

	handle[pgxpool.Pool, *pgxpool.Pool]
}

// Run implements the [process.Runner] interface.
func (p *Pool) Run(ctx context.Context, callback process.Callback) error {
	pool, err := pgxpool.NewWithConfig(ctx, p.Config)
	if err != nil {
		return err
	}

	p.pointer.Store(pool)

	callbackError := callback(ctx)

	p.pointer.CompareAndSwap(pool, nil)

	pool.Close()

	return callbackError
}
