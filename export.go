package collection

import (
	"github.com/goxstdx/go-collection/map_collection"
	"github.com/goxstdx/go-collection/slice_collection"
)

// NewSliceCollect creates and returns a new slice Collection.
//
// It does not copy values. The collection and the caller's slice share the same backing array.
// Comparison is not decided at creation time. int, uint, float, and string get a comparison
// function when methods such as Sort, Max, and Contains are called.
// Structs and defined types must pass a comparison function to the corresponding By method.
func NewSliceCollect[T any](values []T) *slice_collection.Collection[T] {
	return slice_collection.NewCollection[T](values)
}

// NewEmptyCollection returns an empty Collection.
func NewEmptyCollection[T any]() *slice_collection.Collection[T] {
	return slice_collection.NewEmptyCollection[T]()
}

// NewMapCollect creates and returns a new map Collection.
//
// It does not record key order. MinBy and MaxBy rescan on every call.
// EachSorted re-sorts on every call.
func NewMapCollect[K comparable, V any](values map[K]V) *map_collection.Collection[K, V] {
	return map_collection.NewCollection[map[K]V, K, V](values)
}

// NewEmptyMapCollection returns an empty map Collection.
func NewEmptyMapCollection[K comparable, V any]() *map_collection.Collection[K, V] {
	return map_collection.NewCollection[map[K]V, K, V](make(map[K]V))
}
