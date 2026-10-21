// Package map_collection provides generic, chainable map collections with filtering,
// mapping, reduction, and a concurrency-safe collection type.
//
// It provides chainable collection operations over generic maps.
//
// The collection does not store key order. MinBy and MaxBy rescan every element on each call.
// EachSorted re-sorts on each call. None of those results are stored.
// The result types of PluckField, Map, and Reduce are decided by the functions passed in.
// Values are not read by reflecting on field names.
package map_collection
