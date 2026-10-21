package slice_collection

// number is a type whose underlying type is an integer or a floating-point number.
// The ~ allows defined types such as type Money int64 to be used with Sum, Avg, and Median.
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Collection is the core type for collection operations over many kinds of data.
//
// The type parameter T supports basic types and struct types.
// Methods cover filtering, mapping, sorting, searching, aggregation, and similar transforms.
//
// Example:
//
//	intColl := NewCollection([]int{1, 2, 3, 4, 5})
//	filtered := intColl.Filter(func(item int, key int) bool {
//	    return item > 2
//	})
type Collection[T any] struct {
	value []T // slice
}
