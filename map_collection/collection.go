package map_collection

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/goxstdx/go-collection/errorx"
)

// Copy returns a new copy of the Collection.
func (c *Collection[K, V]) Copy() *Collection[K, V] {
	return &Collection[K, V]{value: Clone(c.value)}
}

// IsEmpty reports whether the collection is empty.
func (c *Collection[K, V]) IsEmpty() bool {
	return len(c.value) == 0
}

// IsNotEmpty reports whether the collection is not empty.
func (c *Collection[K, V]) IsNotEmpty() bool {
	return len(c.value) != 0
}

// Count returns the number of key-value pairs in the map.
func (c *Collection[K, V]) Count() int {
	return len(c.value)
}

// Keys returns a slice of all keys.
func (c *Collection[K, V]) Keys() []K {
	return Keys(c.value)
}

// Values returns a slice of all values.
func (c *Collection[K, V]) Values() []V {
	return Values(c.value)
}

// GetValue returns the value for key, or the zero value if it is absent.
func (c *Collection[K, V]) GetValue(key K) V {
	return Get(c.value, key)
}

// Get returns the value for key and whether it exists.
func (c *Collection[K, V]) Get(key K) (V, bool) {
	val, ok := c.value[key]

	return val, ok
}

// GetOr returns the value for key, or def if it is absent.
func (c *Collection[K, V]) GetOr(key K, def V) V {
	return GetOr(c.value, key, def)
}

// Has reports whether key exists.
func (c *Collection[K, V]) Has(key K) bool {
	return Has(c.value, key)
}

// Set sets key to val. It mutates the current Collection and returns it for chaining.
func (c *Collection[K, V]) Set(key K, val V) *Collection[K, V] {
	c.value[key] = val
	return c
}

// Put sets key to val without modifying the original Collection, and returns a new Collection.
func (c *Collection[K, V]) Put(key K, val V) *Collection[K, V] {
	return &Collection[K, V]{value: Set(c.value, key, val)}
}

// Delete removes key without modifying the original Collection, and returns a new Collection.
func (c *Collection[K, V]) Delete(key K) *Collection[K, V] {
	return &Collection[K, V]{value: Delete(c.value, key)}
}

// DeleteByFunc removes keys that match the condition without modifying the original Collection,
// and returns a new Collection.
func (c *Collection[K, V]) DeleteByFunc(fn func(K, V) bool) *Collection[K, V] {
	newMap := make(map[K]V, len(c.value))
	// Still walks a snapshot of the original. Kept entries are rebuilt so a NaN key,
	// which delete cannot remove, is not left behind.
	for k, v := range Clone(c.value) {
		if !fn(k, v) {
			newMap[k] = v
		}
	}
	return &Collection[K, V]{value: newMap}
}

// Remove deletes key. It mutates the current Collection and returns it for chaining.
func (c *Collection[K, V]) Remove(key K) *Collection[K, V] {
	delete(c.value, key)
	return c
}

// Merge merges another map into a new Collection. It does not modify the original Collection.
// On a key conflict, other wins.
func (c *Collection[K, V]) Merge(other map[K]V) *Collection[K, V] {
	return &Collection[K, V]{value: Merge(c.value, other)}
}

// MergeCollection merges another Collection into a new Collection.
// It does not modify the original Collection.
func (c *Collection[K, V]) MergeCollection(other *Collection[K, V]) *Collection[K, V] {
	if other == nil {
		return c.Copy()
	}
	return &Collection[K, V]{value: Merge(c.value, other.value)}
}

// MergeInPlace merges another map into the current Collection.
// It mutates the current Collection and returns it for chaining.
func (c *Collection[K, V]) MergeInPlace(other map[K]V) *Collection[K, V] {
	MergeInPlace(c.value, other)
	return c
}

// Only keeps only the given keys. It does not modify the original Collection and returns a new one.
func (c *Collection[K, V]) Only(keys []K) *Collection[K, V] {
	return &Collection[K, V]{value: Only(c.value, keys)}
}

// Except drops the given keys. It does not modify the original Collection and returns a new one.
func (c *Collection[K, V]) Except(keys []K) *Collection[K, V] {
	return &Collection[K, V]{value: Except(c.value, keys)}
}

// Filter filters key-value pairs and returns a new Collection.
// Pairs for which fn returns true are kept.
func (c *Collection[K, V]) Filter(fn func(value V, key K) bool) *Collection[K, V] {
	return &Collection[K, V]{value: Filter(c.value, fn)}
}

// Each walks the map in its own order. Order is not specified. Iteration stops when fn returns false.
func (c *Collection[K, V]) Each(fn func(value V, key K) bool) *Collection[K, V] {
	Each(c.value, fn)

	return c
}

// Foreach walks every entry in the map's own order. Order is not specified.
// Use EachSorted to walk in an order defined by a comparison function.
func (c *Collection[K, V]) Foreach(fn func(value V, key K)) *Collection[K, V] {
	Foreach(c.value, fn)

	return c
}

// EachSorted walks entries after sorting them with cmp. Every call sorts again and does not store the order.
// This call sorts, so the cost grows as n log n with the number of elements. Another call sorts again.
// The original key order comes from map iteration, so equal comparison results have no fixed order.
// A final key comparison inside the comparison function makes the order stable.
// Sorting uses a snapshot of keys and values taken when the call starts. Keys deleted in the callback
// are skipped. Keys added during the callback are not part of this walk.
// The callback reads the value currently in the map. A nil cmp or fn returns NilFunc.
// A key containing NaN cannot be looked up again, so the value snapshot from the start of the call is used.
func (c *Collection[K, V]) EachSorted(cmp func(a K, va V, b K, vb V) int, fn func(value V, key K)) error {
	if cmp == nil || fn == nil {
		return errorx.NilFunc
	}

	type item struct {
		key   K
		value V
	}
	items := make([]item, 0, len(c.value))
	for k, v := range c.value {
		items = append(items, item{key: k, value: v})
	}
	slices.SortFunc(items, func(a, b item) int {
		return cmp(a.key, a.value, b.key, b.value)
	})
	for _, it := range items {
		// A NaN key is not equal to itself, so it cannot be looked up or deleted by key.
		// Only the snapshot taken during iteration can be used.
		if it.key != it.key {
			fn(it.value, it.key)
			continue
		}
		// The callback may already have deleted this key. A missing map entry yields the zero
		// value, which must not be treated as still present.
		v, ok := c.value[it.key]
		if !ok {
			continue
		}
		fn(v, it.key)
	}
	return nil
}

// Map maps values. R is decided by fn's return value and may differ from V.
func (c *Collection[K, V]) Map[R any](fn func(value V, key K) R) *Collection[K, R] {
	return &Collection[K, R]{value: MapValues(c.value, fn)}
}

// Reduce folds the map into a single result.
func (c *Collection[K, V]) Reduce[R any](init R, fn func(acc R, value V, key K) R) R {
	return Reduce(c.value, init, fn)
}

// MinBy returns the key-value pair that is smallest according to cmp.
// Every call rescans every element. Cost grows linearly with the number of elements, and the result is not stored.
// Another call scans again. cmp returns a negative value when the left side is less than the right side.
// Which pair is kept when the comparison is equal is not fixed. A final key comparison inside the
// comparison function makes the result fixed.
// An empty collection returns EmptyError. A nil cmp returns NilFunc. The accompanying zero value is not a result.
func (c *Collection[K, V]) MinBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	return c.pick(cmp, true)
}

// MaxBy returns the key-value pair that is largest according to cmp.
// Every call rescans every element. Cost grows linearly with the number of elements, and the result is not stored.
// Another call scans again. cmp returns a positive value when the left side is greater than the right side.
// Which pair is kept when the comparison is equal is not fixed. A final key comparison inside the
// comparison function makes the result fixed.
// An empty collection returns EmptyError. A nil cmp returns NilFunc. The accompanying zero value is not a result.
func (c *Collection[K, V]) MaxBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	return c.pick(cmp, false)
}

// pick scans the map once. When wantMin is true it keeps the minimum, otherwise the maximum.
func (c *Collection[K, V]) pick(cmp func(a K, va V, b K, vb V) int, wantMin bool) (K, V, error) {
	var zk K
	var zv V
	if cmp == nil {
		return zk, zv, errorx.NilFunc
	}
	if len(c.value) == 0 {
		return zk, zv, errorx.EmptyError
	}

	var bestK K
	var bestV V
	first := true
	for k, v := range c.value {
		if first {
			bestK, bestV = k, v
			first = false
			continue
		}
		order := cmp(k, v, bestK, bestV)
		if (wantMin && order < 0) || (!wantMin && order > 0) {
			bestK, bestV = k, v
		}
	}
	return bestK, bestV, nil
}

// All returns the underlying map. This is a reference. Modifying it affects the Collection.
func (c *Collection[K, V]) All() map[K]V {
	return c.value
}

// ToJSON converts the Collection to a JSON string.
func (c *Collection[K, V]) ToJSON() (string, error) {
	data, err := json.Marshal(c.value)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DD prints the Collection contents for debugging.
func (c *Collection[K, V]) DD() *Collection[K, V] {
	fmt.Printf("Collection: %+v\n", c.value)

	return c
}

// PluckField uses fn to take one value from each value. Result order is not specified.
func (c *Collection[K, V]) PluckField[R any](fn func(V) R) []R {
	result := make([]R, 0, len(c.value))
	for _, v := range c.value {
		result = append(result, fn(v))
	}
	return result
}
