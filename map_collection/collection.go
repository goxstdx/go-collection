package map_collection

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/goxstdx/go-collection/errorx"
)

// Copy 复制一个新的 Collection
func (c *Collection[K, V]) Copy() *Collection[K, V] {
	return &Collection[K, V]{value: Clone(c.value)}
}

// IsEmpty 判断是否为空
func (c *Collection[K, V]) IsEmpty() bool {
	return len(c.value) == 0
}

// IsNotEmpty 判断是否不为空
func (c *Collection[K, V]) IsNotEmpty() bool {
	return len(c.value) != 0
}

// Count 返回 map 中键值对的数量
func (c *Collection[K, V]) Count() int {
	return len(c.value)
}

// Keys 返回所有的 key 组成的切片
func (c *Collection[K, V]) Keys() []K {
	return Keys(c.value)
}

// Values 返回所有的 value 组成的切片
func (c *Collection[K, V]) Values() []V {
	return Values(c.value)
}

// GetValue 获取指定 key 的值，如果不存在返回零值
func (c *Collection[K, V]) GetValue(key K) V {
	return Get(c.value, key)
}

// Get 获取指定 key 的值
// 返回值和是否存在的标志
func (c *Collection[K, V]) Get(key K) (V, bool) {
	val, ok := c.value[key]

	return val, ok
}

// GetOr 获取 key 对应的值；如果不存在，返回默认值 def
func (c *Collection[K, V]) GetOr(key K, def V) V {
	return GetOr(c.value, key, def)
}

// Has 判断是否存在指定的 key
func (c *Collection[K, V]) Has(key K) bool {
	return Has(c.value, key)
}

// Set 设置 key->val（直接修改当前 Collection，返回自身以支持链式调用）
func (c *Collection[K, V]) Set(key K, val V) *Collection[K, V] {
	c.value[key] = val
	return c
}

// Put 设置 key->val（不修改原 Collection，返回新的 Collection）
func (c *Collection[K, V]) Put(key K, val V) *Collection[K, V] {
	return &Collection[K, V]{value: Set(c.value, key, val)}
}

// Delete 删除指定的 key（不修改原 Collection，返回新的 Collection）
func (c *Collection[K, V]) Delete(key K) *Collection[K, V] {
	return &Collection[K, V]{value: Delete(c.value, key)}
}

// DeleteByFunc 删除满足条件的 key（不修改原 Collection，返回新的 Collection）
func (c *Collection[K, V]) DeleteByFunc(fn func(K, V) bool) *Collection[K, V] {
	newMap := make(map[K]V, len(c.value))
	// 仍遍历原始快照；重建保留项，避免 delete 无法删除含 NaN 的 key。
	for k, v := range Clone(c.value) {
		if !fn(k, v) {
			newMap[k] = v
		}
	}
	return &Collection[K, V]{value: newMap}
}

// Remove 删除指定的 key（直接修改当前 Collection，返回自身以支持链式调用）
func (c *Collection[K, V]) Remove(key K) *Collection[K, V] {
	delete(c.value, key)
	return c
}

// Merge 合并另一个 map 到当前 Collection（不修改原 Collection，返回新的 Collection）
// 键冲突时，以 other 为准
func (c *Collection[K, V]) Merge(other map[K]V) *Collection[K, V] {
	return &Collection[K, V]{value: Merge(c.value, other)}
}

// MergeCollection 合并另一个 Collection 到当前 Collection（不修改原 Collection，返回新的 Collection）
func (c *Collection[K, V]) MergeCollection(other *Collection[K, V]) *Collection[K, V] {
	if other == nil {
		return c.Copy()
	}
	return &Collection[K, V]{value: Merge(c.value, other.value)}
}

// MergeInPlace 将另一个 map 合并到当前 Collection（直接修改当前 Collection，返回自身以支持链式调用）
func (c *Collection[K, V]) MergeInPlace(other map[K]V) *Collection[K, V] {
	MergeInPlace(c.value, other)
	return c
}

// Only 仅保留指定的 keys（不修改原 Collection，返回新的 Collection）
func (c *Collection[K, V]) Only(keys []K) *Collection[K, V] {
	return &Collection[K, V]{value: Only(c.value, keys)}
}

// Except 排除指定的 keys（不修改原 Collection，返回新的 Collection）
func (c *Collection[K, V]) Except(keys []K) *Collection[K, V] {
	return &Collection[K, V]{value: Except(c.value, keys)}
}

// Filter 过滤键值对，返回新的 Collection
// fn 函数返回 true 的键值对会被保留
func (c *Collection[K, V]) Filter(fn func(value V, key K) bool) *Collection[K, V] {
	return &Collection[K, V]{value: Filter(c.value, fn)}
}

// Each 按 map 自身顺序遍历，顺序不确定。fn 返回 false 时终止遍历。
func (c *Collection[K, V]) Each(fn func(value V, key K) bool) *Collection[K, V] {
	Each(c.value, fn)

	return c
}

// Foreach 按 map 自身顺序完整遍历，顺序不确定。要按比较函数的顺序遍历，用 EachSorted。
func (c *Collection[K, V]) Foreach(fn func(value V, key K)) *Collection[K, V] {
	Foreach(c.value, fn)

	return c
}

// EachSorted 按 cmp 排序后遍历。每次调用都重新排序，不保存顺序。
// 这一次是排序，复杂度随元素个数按 n log n 增长。再调一次会再排一次。
// key 的原始顺序来自 map 遍历，比较结果相等时先后不固定。比较函数里再用 key 做最后的比较，顺序就固定。
// 排序用的是调用开始时的键值快照。回调里删掉的 key 会跳过；新增的 key 不会进入这次遍历。
// 回调读到的是当时 map 里的 value。cmp 或 fn 为 nil 时返回 NilFunc。
// 含 NaN 的 key 无法查回，使用调用开始时的 value 快照。
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
		// 含 NaN 的 key 不等于自身，无法查回或按 key 删除，只能使用遍历时的快照。
		if it.key != it.key {
			fn(it.value, it.key)
			continue
		}
		// 回调可能已经删掉这个 key。map 查不到时得到的是零值，不能当成还在。
		v, ok := c.value[it.key]
		if !ok {
			continue
		}
		fn(v, it.key)
	}
	return nil
}

// Map 对 value 做映射。R 由 fn 的返回值决定，可以与 V 不同。
func (c *Collection[K, V]) Map[R any](fn func(value V, key K) R) *Collection[K, R] {
	return &Collection[K, R]{value: MapValues(c.value, fn)}
}

// Reduce 将 map 折叠为一个结果。
func (c *Collection[K, V]) Reduce[R any](init R, fn func(acc R, value V, key K) R) R {
	return Reduce(c.value, init, fn)
}

// MinBy 返回按 cmp 最小的键值对。每次调用都重新扫描全部元素，复杂度随元素个数线性增长，不保存结果。
// 再调一次会再扫描一遍。cmp 返回负数表示左边小于右边。
// 比较结果相等时保留哪一个不固定。比较函数里再用 key 做最后的比较，结果就固定。
// 空集合返回 EmptyError。cmp 为 nil 时返回 NilFunc。伴随返回的零值不是结果。
func (c *Collection[K, V]) MinBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	return c.pick(cmp, true)
}

// MaxBy 返回按 cmp 最大的键值对。每次调用都重新扫描全部元素，复杂度随元素个数线性增长，不保存结果。
// 再调一次会再扫描一遍。cmp 返回正数表示左边大于右边。
// 比较结果相等时保留哪一个不固定。比较函数里再用 key 做最后的比较，结果就固定。
// 空集合返回 EmptyError。cmp 为 nil 时返回 NilFunc。伴随返回的零值不是结果。
func (c *Collection[K, V]) MaxBy(cmp func(a K, va V, b K, vb V) int) (K, V, error) {
	return c.pick(cmp, false)
}

// pick 扫一遍 map。wantMin 为真时取最小，否则取最大。
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

// All 返回底层 map。这是引用，修改会影响 Collection。
func (c *Collection[K, V]) All() map[K]V {
	return c.value
}

// ToJSON 将 Collection 转换为 JSON 字符串
func (c *Collection[K, V]) ToJSON() (string, error) {
	data, err := json.Marshal(c.value)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// DD 打印 Collection 的内容（用于调试）
func (c *Collection[K, V]) DD() *Collection[K, V] {
	fmt.Printf("Collection: %+v\n", c.value)

	return c
}

// PluckField 用 fn 从每个 value 取出一个值。结果顺序不确定。
func (c *Collection[K, V]) PluckField[R any](fn func(V) R) []R {
	result := make([]R, 0, len(c.value))
	for _, v := range c.value {
		result = append(result, fn(v))
	}
	return result
}
