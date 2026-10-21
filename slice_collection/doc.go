// Package slice_collection provides generic, chainable slice collections with filtering,
// mapping, sorting, grouping, and aggregation.
//
// It provides chainable collection operations over generic slices.
//
// Inserts, deletes, updates, and reordering mutate the current collection: Append, Push, Pop,
// Remove, Insert, Prepend, SetIndex, FromJson, and Sort, SortDesc, SortBy, SortByDesc,
// SortByFunc, Reverse, Shuffle.
// Call Copy first when the original contents must be kept.
// Queries and transformations return a new collection. Slice, ForPage, and Split copy elements
// and do not share the backing array with the original collection.
// NewCollection does not copy the input slice, and Values returns the internal slice, so both
// share the backing array with the caller.
//
// int, uint, float32, float64, and string get a comparison function when methods such as Sort,
// Max, and Contains are called.
// Structs and defined types must pass a comparison function to the corresponding By method.
// Min, Max, and Mode return errorx.EmptyError on an empty collection. The zero value is valid
// only when err is nil.
package slice_collection
