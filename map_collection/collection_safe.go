package map_collection

// Copy returns a new concurrency-safe copy of the Collection.
func (sc *SafeCollection[K, V]) Copy() *SafeCollection[K, V] {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: sc.coll.Copy(),
	}
}

// IsEmpty reports whether the collection is empty.
func (sc *SafeCollection[K, V]) IsEmpty() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.IsEmpty()
}

// IsNotEmpty reports whether the collection is not empty.
func (sc *SafeCollection[K, V]) IsNotEmpty() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.IsNotEmpty()
}

// Count returns the number of elements.
func (sc *SafeCollection[K, V]) Count() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Count()
}

// Keys returns a slice of all keys.
func (sc *SafeCollection[K, V]) Keys() []K {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Keys()
}

// Values returns a slice of all values.
func (sc *SafeCollection[K, V]) Values() []V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Values()
}

// GetValue returns the value for key, or the zero value if it is absent.
func (sc *SafeCollection[K, V]) GetValue(key K) V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.GetValue(key)
}

// Get returns the value for key and whether it exists.
func (sc *SafeCollection[K, V]) Get(key K) (V, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Get(key)
}

// GetOr returns the value for key, or the default if it is absent.
func (sc *SafeCollection[K, V]) GetOr(key K, def V) V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.GetOr(key, def)
}

// Has reports whether key exists.
func (sc *SafeCollection[K, V]) Has(key K) bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Has(key)
}

// Set sets key to val. It mutates the current Collection.
func (sc *SafeCollection[K, V]) Set(key K, val V) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.Set(key, val)
	return sc
}

// Put sets key to val and returns a new concurrency-safe Collection.
func (sc *SafeCollection[K, V]) Put(key K, val V) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Put(key, val)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// Delete removes key and returns a new concurrency-safe Collection.
func (sc *SafeCollection[K, V]) Delete(key K) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Delete(key)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// DeleteByFunc removes keys that match the condition and returns a new concurrency-safe Collection.
func (sc *SafeCollection[K, V]) DeleteByFunc(fn func(K, V) bool) *SafeCollection[K, V] {
	sc.mu.RLock()
	// Release the read lock even if the callback panics, so a later write is not blocked after recover.
	defer sc.mu.RUnlock()
	newColl := sc.coll.DeleteByFunc(fn)

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// Remove deletes key. It mutates the current Collection.
func (sc *SafeCollection[K, V]) Remove(key K) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.Remove(key)
	return sc
}

// Merge merges another map and returns a new concurrency-safe Collection.
func (sc *SafeCollection[K, V]) Merge(other map[K]V) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Merge(other)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// MergeCollection merges another SafeCollection and returns a new concurrency-safe Collection.
func (sc *SafeCollection[K, V]) MergeCollection(other *SafeCollection[K, V]) *SafeCollection[K, V] {
	if other == nil {
		return sc.Copy()
	}

	// Copy other before locking this collection, so only one lock is held at a time.
	// Holding both read locks can deadlock when a and b merge each other, or a collection
	// merges itself, if a write is already queued.
	otherMap := other.All()

	sc.mu.RLock()
	newColl := sc.coll.Merge(otherMap)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// MergeInPlace merges another map into the current Collection. It mutates the collection.
func (sc *SafeCollection[K, V]) MergeInPlace(other map[K]V) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.MergeInPlace(other)
	return sc
}

// Filter filters elements and returns a new concurrency-safe Collection.
func (sc *SafeCollection[K, V]) Filter(fn func(V, K) bool) *SafeCollection[K, V] {
	sc.mu.RLock()
	// Release the read lock even if the callback panics, so a later write is not blocked after recover.
	defer sc.mu.RUnlock()
	newColl := sc.coll.Filter(fn)

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// Map maps values and returns a new concurrency-safe Collection. R may differ from V.
func (sc *SafeCollection[K, V]) Map[R any](fn func(V, K) R) *SafeCollection[K, R] {
	sc.mu.RLock()
	// The read lock must be released even if the callback panics, so a later write is not blocked after recover.
	defer sc.mu.RUnlock()
	newColl := sc.coll.Map(fn)

	return &SafeCollection[K, R]{
		coll: newColl,
	}
}

// Each walks the map in its own order. Order is not specified. Iteration stops when fn returns false.
func (sc *SafeCollection[K, V]) Each(fn func(V, K) bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.coll.Each(fn)
}

// Foreach walks every entry in the map's own order. Order is not specified.
// Use EachSorted to walk in an order defined by a comparison function.
func (sc *SafeCollection[K, V]) Foreach(fn func(V, K)) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.coll.Foreach(fn)
}

// Reduce folds the map into a single result. R is decided by init and fn.
func (sc *SafeCollection[K, V]) Reduce[R any](init R, fn func(R, V, K) R) R {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Reduce(init, fn)
}

// MinBy returns the key-value pair that is smallest according to cmp.
// Every call rescans every element.
func (sc *SafeCollection[K, V]) MinBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.MinBy(cmp)
}

// MaxBy returns the key-value pair that is largest according to cmp.
// Every call rescans every element.
func (sc *SafeCollection[K, V]) MaxBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.MaxBy(cmp)
}

// EachSorted walks entries after sorting them with cmp.
// Every call sorts again and does not store the order.
func (sc *SafeCollection[K, V]) EachSorted(cmp func(a K, va V, b K, vb V) int, fn func(V, K)) error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.EachSorted(cmp, fn)
}

// All returns a copy of the underlying map. Modifying the result does not affect the Collection.
func (sc *SafeCollection[K, V]) All() map[K]V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return Clone(sc.coll.All())
}

// ToJSON serializes the collection to a JSON string.
func (sc *SafeCollection[K, V]) ToJSON() (string, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.ToJSON()
}
