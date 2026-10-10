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

	sc.mu.RLock()
	other.mu.RLock()
	newColl := sc.coll.MergeCollection(other.coll)
	other.mu.RUnlock()
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
	newColl := sc.coll.Map(fn)
	sc.mu.RUnlock()

	return &SafeCollection[K, R]{
		coll: newColl,
	}
}

// Each 遍历每个元素
func (sc *SafeCollection[K, V]) Each(fn func(V, K)) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.coll.Each(fn)
}

// Reduce 聚合操作
func (sc *SafeCollection[K, V]) Reduce(init any, fn func(any, V, K) any) any {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Reduce(init, fn)
}

// First 返回第一个元素（基于sortedKeys）
func (sc *SafeCollection[K, V]) First() (K, V, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.First()
}

// Last 返回最后一个元素（基于sortedKeys）
func (sc *SafeCollection[K, V]) Last() (K, V, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.Last()
}

// All 返回底层 map 的副本
func (sc *SafeCollection[K, V]) All() map[K]V {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.All()
}

// ToJSON 序列化为 JSON 字符串
func (sc *SafeCollection[K, V]) ToJSON() (string, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.coll.ToJSON()
}

// SetKeyCompare 设置 key 的比较函数
func (sc *SafeCollection[K, V]) SetKeyCompare(fn func(K, K) int) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.SetKeyCompare(fn)
	return sc
}

// OrderKey 根据 keyCompareFunc 排序
func (sc *SafeCollection[K, V]) OrderKey() (*SafeCollection[K, V], error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, err := sc.coll.OrderKey()
	return sc, err
}

// SetValCompare 设置 value 的比较函数
func (sc *SafeCollection[K, V]) SetValCompare(fn func(V, V) int) *SafeCollection[K, V] {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.coll.SetValCompare(fn)
	return sc
}

// OrderValue 按 value 排序
func (sc *SafeCollection[K, V]) OrderValue() (*SafeCollection[K, V], error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, err := sc.coll.OrderValue()
	return sc, err
}
