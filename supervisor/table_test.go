package supervisor

import (
	"context"
	"errors"
	"sync"
	"testing"

	"gotest.tools/v3/assert"
)

// errNotFound is an error that map table implementation return when an element
// does not exist for the given key.
var errNotFound = errors.New("process: map table element was not found")

// mapTable is a Table implementation that uses builtin map type.
type mapTable[K comparable, V any] struct {
	m typedMap[K, V]
}

// Get implements the Table interface.
func (m *mapTable[K, V]) Get(_ context.Context, key K) (V, error) {
	v, ok := m.m.Load(key)
	if !ok {
		return v, errNotFound
	}
	return v, nil
}

// Iter implements the Table interface.
func (m *mapTable[K, V]) Iter(_ context.Context) (Iterator[K, V], error) {
	it := &mapIterator[K, V]{
		stop: make(chan struct{}),
		next: make(chan pair[K, V]),
	}
	go func() {
		m.m.Range(func(k K, v V) bool {
			select {
			case it.next <- pair[K, V]{k, v}:
			case <-it.stop:
				return false
			}
			return true
		})
		close(it.next)
	}()
	return it, nil
}

// mapIterator is an Iterator implementation for the mapTable type.
type mapIterator[K comparable, V any] struct {
	next chan pair[K, V]

	once sync.Once
	stop chan struct{}

	exists  bool
	current pair[K, V]
}

// Next implements the Iterator interface.
func (it *mapIterator[K, V]) Next() bool {
	it.current, it.exists = <-it.next
	return it.exists
}

// Get implements the Iterator interface.
func (it *mapIterator[K, V]) Get(_ context.Context) (K, V, error) {
	k := it.current.key
	v := it.current.val
	if !it.exists {
		return k, v, errNotFound
	}
	return k, v, nil
}

// Err implements the Iterator interface.
func (it *mapIterator[K, V]) Err() error {
	return nil
}

// Close implements the Iterator interface.
func (it *mapIterator[K, V]) Close() error {
	it.once.Do(func() {
		close(it.stop)
	})
	return nil
}

// pair is the key and value pair.
type pair[K any, V any] struct {
	key K
	val V
}

func TestMapTable(t *testing.T) {
	ctx := context.Background()

	var tab mapTable[uint64, string]
	tab.m.Store(1, "foo")
	tab.m.Store(2, "bar")
	tab.m.Store(2, "baz")

	v, err := tab.Get(ctx, 1)
	assert.NilError(t, err)
	assert.Equal(t, v, "foo")

	v, err = tab.Get(ctx, 2)
	assert.NilError(t, err)
	assert.Equal(t, v, "baz")

	it, err := tab.Iter(ctx)
	assert.NilError(t, err)
	t.Cleanup(func() {
		assert.NilError(t, it.Close())
	})

	m := map[uint64]string{}
	for it.Next() {
		k, v, err := it.Get(ctx)
		assert.NilError(t, err)

		_, exists := m[k]
		assert.Assert(t, !exists)

		m[k] = v
	}
	assert.NilError(t, it.Err())
	assert.DeepEqual(t, map[uint64]string{
		1: "foo",
		2: "baz",
	}, m)
}
