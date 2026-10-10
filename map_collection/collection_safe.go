package map_collection

// Copy 返回一个新的线程安全的 Collection 副本
func (sc *SafeCollection[K, V]) Copy() *SafeCollection[K, V] {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: sc.coll.Copy(),
	}
}

// IsEmpty 判断是否为空
func (sc *SafeCollection[K, V]) IsEmpty() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.IsEmpty()
}

// IsNotEmpty 判断是否不为空
func (sc *SafeCollection[K, V]) IsNotEmpty() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.IsNotEmpty()
}

// Count 返回元素数量
func (sc *SafeCollection[K, V]) Count() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Count()
}

// Keys 返回所有 key 的切片
func (sc *SafeCollection[K, V]) Keys() []K {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Keys()
}

// Values 返回所有 value 的切片
func (sc *SafeCollection[K, V]) Values() []V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Values()
}

// GetValue 获取指定 key 的值，如果不存在返回零值
func (sc *SafeCollection[K, V]) GetValue(key K) V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.GetValue(key)
}

// Get 获取指定 key 的值，返回值和是否存在的标志
func (sc *SafeCollection[K, V]) Get(key K) (V, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Get(key)
}

// GetOr 获取 key 对应的值；如果不存在，返回默认值
func (sc *SafeCollection[K, V]) GetOr(key K, def V) V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.GetOr(key, def)
}

// Has 判断是否存在指定的 key
func (sc *SafeCollection[K, V]) Has(key K) bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Has(key)
}

// Set 设置 key->val（直接修改当前 Collection）
func (sc *SafeCollection[K, V]) Set(key K, val V) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.Set(key, val)
	return sc
}

// Put 设置 key->val（返回新的线程安全 Collection）
func (sc *SafeCollection[K, V]) Put(key K, val V) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Put(key, val)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// Delete 删除指定的 key（返回新的线程安全 Collection）
func (sc *SafeCollection[K, V]) Delete(key K) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Delete(key)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// DeleteByFunc 删除满足条件的 key（返回新的线程安全 Collection）
func (sc *SafeCollection[K, V]) DeleteByFunc(fn func(K, V) bool) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.DeleteByFunc(fn)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// Remove 删除指定的 key（直接修改当前 Collection）
func (sc *SafeCollection[K, V]) Remove(key K) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.Remove(key)
	return sc
}

// Merge 合并另一个 map（返回新的线程安全 Collection）
func (sc *SafeCollection[K, V]) Merge(other map[K]V) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Merge(other)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// MergeCollection 合并另一个 SafeCollection（返回新的线程安全 Collection）
func (sc *SafeCollection[K, V]) MergeCollection(other *SafeCollection[K, V]) *SafeCollection[K, V] {
	if other == nil {
		return sc.Copy()
	}

	// 先取 other 的副本再锁自己，同一时刻只持有一把锁。
	// 同时持有两把读锁时，a、b 互相合并或自己合并自己，遇到排队的写操作会死锁。
	otherMap := other.All()

	sc.mu.RLock()
	newColl := sc.coll.Merge(otherMap)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// MergeInPlace 将另一个 map 合并到当前 Collection（直接修改）
func (sc *SafeCollection[K, V]) MergeInPlace(other map[K]V) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.MergeInPlace(other)
	return sc
}

// Filter 过滤元素（返回新的线程安全 Collection）
func (sc *SafeCollection[K, V]) Filter(fn func(V, K) bool) *SafeCollection[K, V] {
	sc.mu.RLock()
	newColl := sc.coll.Filter(fn)
	sc.mu.RUnlock()

	return &SafeCollection[K, V]{
		coll: newColl,
	}
}

// Map 映射转换，返回新的线程安全 Collection。R 可以与 V 不同。
func (sc *SafeCollection[K, V]) Map[R any](fn func(V, K) R) *SafeCollection[K, R] {
	sc.mu.RLock()
	// 回调 panic 时也必须释放读锁，避免 recover 后的写操作一直阻塞。
	defer sc.mu.RUnlock()
	newColl := sc.coll.Map(fn)

	return &SafeCollection[K, R]{
		coll: newColl,
	}
}

// Each 按 map 自身顺序遍历，顺序不确定。fn 返回 false 时终止遍历。
func (sc *SafeCollection[K, V]) Each(fn func(V, K) bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.coll.Each(fn)
}

// Foreach 按 map 自身顺序完整遍历，顺序不确定。要按比较函数的顺序遍历，用 EachSorted。
func (sc *SafeCollection[K, V]) Foreach(fn func(V, K)) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.coll.Foreach(fn)
}

// Reduce 将 map 折叠为一个结果。R 由 init 和 fn 决定。
func (sc *SafeCollection[K, V]) Reduce[R any](init R, fn func(R, V, K) R) R {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Reduce(init, fn)
}

// MinBy 返回按 cmp 最小的键值对。每次调用都重新扫描全部元素。
func (sc *SafeCollection[K, V]) MinBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.MinBy(cmp)
}

// MaxBy 返回按 cmp 最大的键值对。每次调用都重新扫描全部元素。
func (sc *SafeCollection[K, V]) MaxBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.MaxBy(cmp)
}

// EachSorted 按 cmp 排序后遍历。每次调用都重新排序，不保存顺序。
func (sc *SafeCollection[K, V]) EachSorted(cmp func(a K, va V, b K, vb V) int, fn func(V, K)) error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.EachSorted(cmp, fn)
}

// All 返回底层 map 的副本。修改返回值不会影响 Collection。
func (sc *SafeCollection[K, V]) All() map[K]V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return Clone(sc.coll.All())
}

// ToJSON 序列化为 JSON 字符串
func (sc *SafeCollection[K, V]) ToJSON() (string, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.ToJSON()
}
