package supervisor

import (
	"sync"
)

// typedMap is a type-safe wrapper for sync.typedMap.
type typedMap[K any, V any] struct {
	m sync.Map
}

// Delete is a type-safe wrapper for sync.Map’s Delete method.
func (m *typedMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Load is a type-safe wrapper for sync.Map’s Load method.
func (m *typedMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if ok {
		value, _ = v.(V)
	}
	return
}

// LoadAndDelete is a type-safe wrapper for sync.Map’s LoadAndDelete method.
func (m *typedMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.m.LoadAndDelete(key)
	if loaded {
		value, _ = v.(V)
	}
	return
}

// LoadOrStore is a type-safe wrapper for sync.Map’s LoadOrStore method.
func (m *typedMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.m.LoadOrStore(key, value)
	if loaded {
		actual, _ = v.(V)
	} else {
		actual = value
	}
	return
}

// Range is a type-safe wrapper for sync.Map’s Range method.
func (m *typedMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(k, v any) bool {
		key, _ := k.(K)
		val, _ := v.(V)
		return f(key, val)
	})
}

// Store is a type-safe wrapper for sync.Map’s Store method.
func (m *typedMap[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}
