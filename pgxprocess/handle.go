package pgxprocess

import (
	"context"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// database is a common interface for [pgx.Conn] and [pgxpool.Pool].
type database interface {
	Ping(context.Context) error
	Begin(context.Context) (pgx.Tx, error)
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

type handle[T any, P interface {
	*T
	database
}] struct {
	pointer atomic.Pointer[T]
}

// Handle returns the underlying database handle.
func (h *handle[T, P]) Handle() P {
	return P(h.pointer.Load())
}

func (h *handle[T, P]) Ping(ctx context.Context) error {
	p := P(h.pointer.Load())
	if p == nil {
		return ErrNotConnected
	}
	return p.Ping(ctx)
}

func (h *handle[T, P]) Begin(ctx context.Context) (pgx.Tx, error) {
	p := P(h.pointer.Load())
	if p == nil {
		return nil, ErrNotConnected
	}
	return p.Begin(ctx)
}

func (h *handle[T, P]) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	p := P(h.pointer.Load())
	if p == nil {
		return nil, ErrNotConnected
	}
	return p.BeginTx(ctx, txOptions)
}

func (h *handle[T, P]) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	p := P(h.pointer.Load())
	if p == nil {
		return pgconn.CommandTag{}, ErrNotConnected
	}
	return p.Exec(ctx, sql, args...)
}

func (h *handle[T, P]) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	p := P(h.pointer.Load())
	if p == nil {
		return nil, ErrNotConnected
	}
	return p.Query(ctx, sql, args...)
}

func (h *handle[T, P]) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	p := P(h.pointer.Load())
	if p == nil {
		return errRow{}
	}
	return p.QueryRow(ctx, sql, args...)
}

func (h *handle[T, P]) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	p := P(h.pointer.Load())
	if p == nil {
		return 0, ErrNotConnected
	}
	return p.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (h *handle[T, P]) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	p := P(h.pointer.Load())
	if p == nil {
		return errBatchResults{}
	}
	return p.SendBatch(ctx, b)
}
