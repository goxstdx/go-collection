package slice_collection

import "cmp"

// NewCollection creates and returns a new Collection.
//
// It does not copy values. The collection and the caller's slice share the same backing array.
// When capacity allows, Append, Remove, Insert, and similar methods write into that slice.
// The comparison function is not chosen here. Predeclared numeric types and string get one
// when sort or search methods are called.
// Structs and defined types must pass a comparison function to the corresponding By method.
//
// Parameters:
//   - values: the slice used to initialize the collection
//
// Returns:
//   - *Collection[T]: the new Collection
func NewCollection[T any](values []T) *Collection[T] {
	return &Collection[T]{value: values}
}

// cloneSlice copies elements. The result's capacity equals its length, so a later append
// does not write back into src.
func cloneSlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

// orderedFunc converts cmp.Compare[N] into func(T, T) int.
// The type switch has already established that T is N, so the assertion runs once,
// when the comparison function is created.
func orderedFunc[T any, N cmp.Ordered]() func(T, T) int {
	return any(cmp.Compare[N]).(func(T, T) int)
}

// defaultCompare builds a comparison function for predeclared ordered types.
// Structs and defined types such as type MyInt int return nil. The caller must pass
// a comparison function.
func defaultCompare[T any]() func(T, T) int {
	var zero T
	switch any(zero).(type) {
	case int:
		return orderedFunc[T, int]()
	case int8:
		return orderedFunc[T, int8]()
	case int16:
		return orderedFunc[T, int16]()
	case int32:
		return orderedFunc[T, int32]()
	case int64:
		return orderedFunc[T, int64]()
	case uint:
		return orderedFunc[T, uint]()
	case uint8:
		return orderedFunc[T, uint8]()
	case uint16:
		return orderedFunc[T, uint16]()
	case uint32:
		return orderedFunc[T, uint32]()
	case uint64:
		return orderedFunc[T, uint64]()
	case float32:
		return orderedFunc[T, float32]()
	case float64:
		return orderedFunc[T, float64]()
	case string:
		return orderedFunc[T, string]()
	default:
		return nil
	}
}

// NewEmptyCollection returns an empty Collection.
func NewEmptyCollection[T any]() *Collection[T] {
	return NewCollection[T](nil)
}
