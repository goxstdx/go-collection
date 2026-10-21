package map_collection

// This file provides helpers for generic maps, modeled on common Laravel Collection
// operations and implemented with Go generics and the standard maps package.
//
// Notes:
//   - These functions are pure (except those whose names end in InPlace). They do not
//     modify the input map, and they return a new map or a result value.
//   - Go maps are unordered, so no order-related API is provided here.
//   - Prefer the Go 1.27+ standard maps package (Keys, Values, Clone, Equal, Copy).

import (
	"maps"
)

// Keys returns all keys of the map.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}

	return res
}

// Values returns all values of the map.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	res := make([]V, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

// Clone returns a copy of m.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	return maps.Clone(m)
}

// Equal reports whether two maps are equal (keys and values match one by one).
func Equal[M1 ~map[K]V, M2 ~map[K]V, K comparable, V comparable](m1 M1, m2 M2) bool {
	return maps.Equal(m1, m2)
}

// Has reports whether key is present in the map.
func Has[M ~map[K]V, K comparable, V any](m M, key K) bool {
	_, ok := m[key]
	return ok
}

// Get returns the value for key, or the zero value if it is absent.
func Get[M ~map[K]V, K comparable, V any](m M, key K) V {
	return m[key]
}

// GetOr returns the value for key, or def if it is absent.
func GetOr[M ~map[K]V, K comparable, V any](m M, key K, def V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return def
}

// Set returns a new map with key set to val. It does not modify the original map.
func Set[M ~map[K]V, K comparable, V any](m M, key K, val V) M {
	nm := maps.Clone(m)
	if nm == nil {
		nm = make(M)
	}
	nm[key] = val
	return nm
}

// SetInPlace sets key to val on the original map.
func SetInPlace[M ~map[K]V, K comparable, V any](m M, key K, val V) {
	m[key] = val
}

// Delete returns a new map with key removed. It does not modify the original map.
func Delete[M ~map[K]V, K comparable, V any](m M, key K) M {
	nm := maps.Clone(m)
	delete(nm, key)
	return nm
}

// DeleteInPlace deletes key from the original map.
func DeleteInPlace[M ~map[K]V, K comparable, V any](m M, key K) {
	delete(m, key)
}

// Merge returns a new map containing every entry from m and other.
// On a key conflict, other wins.
func Merge[M ~map[K]V, K comparable, V any](m M, other M) M {
	// Both nil: return an empty map rather than nil, so a later write does not panic.
	if m == nil && other == nil {
		return make(M)
	}
	nm := maps.Clone(m)
	if nm == nil {
		nm = make(M)
	}
	maps.Copy(nm, other)
	return nm
}

// MergeInPlace merges other's entries into m. On a key conflict, other wins.
func MergeInPlace[M ~map[K]V, K comparable, V any](m M, other M) {
	maps.Copy(m, other)
}

// Only keeps only the keys listed in keys and returns a new map.
func Only[M ~map[K]V, K comparable, V any](m M, keys []K) M {
	nm := make(M, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			nm[k] = v
		}
	}
	return nm
}

// Except drops the keys listed in keys and returns a new map. It does not modify the original map.
func Except[M ~map[K]V, K comparable, V any](m M, keys []K) M {
	// An empty input returns an empty map rather than nil, so a later write on the result does not panic.
	if len(m) == 0 {
		return make(M)
	}
	nm := maps.Clone(m)
	for _, k := range keys {
		delete(nm, k)
	}
	return nm
}

// MapValues maps values, keeps the original keys, and returns a new map.
func MapValues[M ~map[K]V, K comparable, V any, R any](m M, fn func(value V, key K) R) map[K]R {
	res := make(map[K]R, len(m))
	for k, v := range m {
		res[k] = fn(v, k)
	}
	return res
}

// MapKeys maps keys and may change the key type. It returns a new map.
// If two keys map to the same new key, the later one overwrites the earlier one.
func MapKeys[M ~map[K]V, K comparable, V any, NK comparable](m M, fn func(key K, value V) NK) map[NK]V {
	res := make(map[NK]V, len(m))
	for k, v := range m {
		res[fn(k, v)] = v
	}
	return res
}

// Filter filters key-value pairs and returns a new map.
func Filter[M ~map[K]V, K comparable, V any](m M, fn func(value V, key K) bool) M {
	nm := make(M)
	for k, v := range m {
		if fn(v, k) {
			nm[k] = v
		}
	}
	return nm
}

// Each walks the map in its own order and stops when fn returns false.
func Each[M ~map[K]V, K comparable, V any](m M, fn func(value V, key K) bool) {
	for k, v := range m {
		if !fn(v, k) {
			return
		}
	}
}

// Foreach calls fn once for every key-value pair. Order is not specified.
func Foreach[M ~map[K]V, K comparable, V any](m M, fn func(value V, key K)) {
	for k, v := range m {
		fn(v, k)
	}
}

// Reduce folds the map into a single result.
func Reduce[M ~map[K]V, K comparable, V any, R any](m M, init R, fn func(acc R, value V, key K) R) R {
	acc := init
	for k, v := range m {
		acc = fn(acc, v, k)
	}
	return acc
}

// Pluck uses fn to take one value from each map value. Iteration order is not stable.
func Pluck[M ~map[K]V, K comparable, V any, R any](m M, fn func(V) R) []R {
	result := make([]R, 0, len(m))
	for _, v := range m {
		result = append(result, fn(v))
	}

	return result
}
