package pgxprocess

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrNotConnected is an error that is returned if database handle is not
// connected.
var ErrNotConnected = errors.New("pgxprocess: not connected")

type errRow struct{}

func (errRow) Scan(_ ...any) error {
	return ErrNotConnected
}

type errRows struct{}

func (errRows) Close() {}

func (errRows) Err() error {
	return ErrNotConnected
}

func (errRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (errRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (errRows) Next() bool {
	return false
}

func (errRows) Scan(_ ...any) error {
	return ErrNotConnected
}

func (errRows) Values() ([]any, error) {
	return nil, ErrNotConnected
}

func (errRows) RawValues() [][]byte {
	return nil
}

func (errRows) Conn() *pgx.Conn {
	return nil
}

type errBatchResults struct{}

func (errBatchResults) Exec() (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, ErrNotConnected
}

func (errBatchResults) Query() (pgx.Rows, error) {
	return errRows{}, ErrNotConnected
}

func (errBatchResults) QueryRow() pgx.Row {
	return errRow{}
}

func (errBatchResults) Close() error {
	return ErrNotConnected
}
