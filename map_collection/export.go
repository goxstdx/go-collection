package map_collection

import (
	"sync"
)

// Collection is the core type for collection operations over many kinds of data.
// It is the core collection type for maps.
type Collection[K comparable, V any] struct {
	value map[K]V // map data
}

// NewCollection creates and returns a new Collection.
// It does not record key order. MinBy and MaxBy rescan on every call.
// EachSorted re-sorts on every call. None of those results are stored.
// A nil values map is replaced with an empty map so later writes such as Set and
// MergeInPlace do not panic on a nil map.
func NewCollection[T ~map[K]V, K comparable, V any](values T) *Collection[K, V] {
	if values == nil {
		values = make(T)
	}
	return &Collection[K, V]{
		value: values,
	}
}

// SafeCollection is a concurrency-safe wrapper around Collection.
// It guards the underlying Collection with a read-write mutex.
// Do not call the same SafeCollection again from a callback. Write methods wait on this lock,
// and read methods can also stall while a write is queued.
type SafeCollection[K comparable, V any] struct {
	mu   sync.RWMutex
	coll *Collection[K, V]
}

// NewSafeCollection creates a concurrency-safe Collection wrapper.
func NewSafeCollection[T ~map[K]V, K comparable, V any](values T) *SafeCollection[K, V] {
	return &SafeCollection[K, V]{
		coll: NewCollection(values),
	}
}
