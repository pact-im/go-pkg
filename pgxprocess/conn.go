package pgxprocess

import (
	"context"

	"github.com/jackc/pgx/v5"

	"go.pact.im/x/process"
)

var _ database = (*Conn)(nil)

// Conn wraps [pgx.Conn] to delay connection setup to application runtime.
type Conn struct {
	Config *pgx.ConnConfig

	handle[pgx.Conn, *pgx.Conn]
}

// Run implements the [process.Runner] interface.
func (p *Conn) Run(ctx context.Context, callback process.Callback) error {
	conn, err := pgx.ConnectConfig(ctx, p.Config)
	if err != nil {
		return err
	}

	p.pointer.Store(conn)

	callbackError := callback(ctx)

	p.pointer.CompareAndSwap(conn, nil)

	_ = conn.Close(ctx)

	return callbackError
}
