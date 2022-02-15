package process

import (
	"context"
)

// Table defines a table of key and value pairs that is potentially backed by a
// persistent and shared storage.
type Table[K comparable, V any] interface {
	// Get returns the value for the given key.
	Get(ctx context.Context, key K) (V, error)
	// Iter returns an iterator for key and value pairs in the table.
	Iter(ctx context.Context) (Iterator[K, V], error)
}

// Iterator iterates over key and value pairs in the Table. It follows the
// semantics of the standard sql.Rows type and does not necessarily correspond
// to any consistent snapshot of the Table’s contents.
type Iterator[K comparable, V any] interface {
	// Next prepares the value for the next iteration. It returns true on
	// success, and false if there is no next value. Consult Err to check
	// whether iterator successfully reached the end or an error occurred.
	Next() bool

	// Get returns the key and value for the current iteration.
	Get(ctx context.Context) (K, V, error)

	// Err returns any error that occurred during iteration.
	Err() error
	// Close closes the iterator.
	Close() error
}
