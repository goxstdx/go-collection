package map_collection

import (
	"slices"

	"github.com/ZHOUXING1997/collection/errorx"
)

func (c *Collection[K, V]) sortKey(fn func(K, K) int) {
	if c.sortedKeys == nil {
		c.sortedKeys = Keys(c.value)
	}

	slices.SortFunc(c.sortedKeys, fn)
}

// SetKeyCompare 设置 key 的比较函数（不立即排序，需要调用 Order() 或使用 First/Last 时才排序）
// -1：小于，0：等于，1：大于
func (c *Collection[K, V]) SetKeyCompare(fn func(K, K) int) *Collection[K, V] {
	c.keyCompareFunc = fn
	return c
}

// OrderKey 根据当前的 keyCompareFunc 对 sortedKeys 重新排序
func (c *Collection[K, V]) OrderKey() (*Collection[K, V], error) {
	if c.keyCompareFunc == nil {
		return c, errorx.NotHaveKeyCompareFunc
	}

	c.sortKey(c.keyCompareFunc)

	return c, nil
}

// OrderKeyByFunc 设置 key 的比较函数并立即排序
// -1：小于，0：等于，1：大于
func (c *Collection[K, V]) OrderKeyByFunc(fn func(K, K) int) (*Collection[K, V], error) {
	if fn == nil {
		return c, errorx.NilFunc
	}

	c.keyCompareFunc = fn

	c.sortKey(c.keyCompareFunc)

	return c, nil
}

// SetValCompare 设置 value 的比较函数
// -1：小于，0：等于，1：大于
func (c *Collection[K, V]) SetValCompare(compareFunc func(V, V) int) *Collection[K, V] {
	c.valCompareFunc = compareFunc

	return c
}

// OrderValue 按 value 排序，更新 sortedKeys
// -1：小于，0：等于，1：大于
func (c *Collection[K, V]) OrderValue() (*Collection[K, V], error) {
	if c.valCompareFunc == nil {
		return c, errorx.NotHaveValCompareFunc
	}

	c.sortKey(func(a, b K) int {
		return c.valCompareFunc(c.value[a], c.value[b])
	})

	return c, nil
}

// OrderByValueFunc 按 value 的提取结果排序。
// extractFunc: 从 value 中提取用于比较的值
// compareFunc: 比较函数，-1：小于，0：等于，1：大于
func (c *Collection[K, V]) OrderByValueFunc[R any](extractFunc func(V) R, compareFunc func(R, R) int) (*Collection[K, V], error) {
	if extractFunc == nil || compareFunc == nil {
		return c, errorx.NilFunc
	}

	c.sortKey(func(a, b K) int {
		return compareFunc(extractFunc(c.value[a]), extractFunc(c.value[b]))
	})

	return c, nil
}
