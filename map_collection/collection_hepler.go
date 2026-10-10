package map_collection

import (
	"slices"
	"sort"
)

// initSortedKeys 初始化排序后的 keys
func (c *Collection[K, V]) initSortedKeys() {
	if c.sortedKeys == nil {
		c.sortedKeys = Keys(c.value)
	}
	if c.keyCompareFunc != nil {
		slices.SortFunc(c.sortedKeys, func(a, b K) int {
			return c.keyCompareFunc(a, b)
		})
	}
}

// insertKeyInOrder 将 key 插入到 sortedKeys 的正确位置（增量更新）
func (c *Collection[K, V]) insertKeyInOrder(key K) {
	if c.sortedKeys == nil {
		// 如果 sortedKeys 未初始化，初始化它
		c.initSortedKeys()
		return
	}

	if c.keyCompareFunc == nil {
		// 无排序函数，直接追加
		c.sortedKeys = append(c.sortedKeys, key)
		return
	}

	// 二分查找插入位置
	idx := sort.Search(len(c.sortedKeys), func(i int) bool {
		return c.keyCompareFunc(c.sortedKeys[i], key) >= 0
	})
	c.sortedKeys = slices.Insert(c.sortedKeys, idx, key)
}

// removeKeyFromSorted 从 sortedKeys 中移除指定的 key
func (c *Collection[K, V]) removeKeyFromSorted(keys ...K) {
	if c.sortedKeys == nil || len(keys) == 0 {
		return
	}

	// 创建keys的查找map
	keyMap := make(map[K]bool, len(keys))
	for _, k := range keys {
		keyMap[k] = true
	}

	// 过滤sortedKeys,保留不在keys中的元素
	filtered := make([]K, 0, len(c.sortedKeys))
	for _, k := range c.sortedKeys {
		if !keyMap[k] {
			filtered = append(filtered, k)
		}
	}
	c.sortedKeys = filtered
}

// cloneWithSortedKeys 内部方法：克隆 Collection 并继承排序信息
func (c *Collection[K, V]) cloneWithSortedKeys(newMap map[K]V) *Collection[K, V] {
	newColl := &Collection[K, V]{
		value:          newMap,
		valCompareFunc: c.valCompareFunc,
		keyCompareFunc: c.keyCompareFunc,
	}

	if c.sortedKeys != nil {
		newColl.sortedKeys = slices.Clone(c.sortedKeys)
	}

	return newColl
}

// filterSortedKeys 从 sortedKeys 中过滤出在 newMap 中存在的 key，保持顺序
func (c *Collection[K, V]) filterSortedKeys(newMap map[K]V) []K {
	if c.sortedKeys == nil {
		return nil
	}
	filtered := make([]K, 0, len(newMap))
	for _, k := range c.sortedKeys {
		if _, ok := newMap[k]; ok {
			filtered = append(filtered, k)
		}
	}
	return filtered
}

// mergeKeys 归并两个已排序的 key 列表（用于 Merge 操作）
func (c *Collection[K, V]) mergeKeys(otherKeys []K) []K {
	if c.keyCompareFunc == nil {
		// 无排序函数，直接合并（去重）
		result := slices.Clone(c.sortedKeys)
		for _, k := range otherKeys {
			found := false
			for _, existing := range c.sortedKeys {
				if existing == k {
					found = true
					break
				}
			}
			if !found {
				result = append(result, k)
			}
		}
		return result
	}

	// 归并排序两个已排序的列表
	result := make([]K, 0, len(c.sortedKeys)+len(otherKeys))
	i, j := 0, 0

	for i < len(c.sortedKeys) && j < len(otherKeys) {
		cmp := c.keyCompareFunc(c.sortedKeys[i], otherKeys[j])
		if cmp < 0 {
			result = append(result, c.sortedKeys[i])
			i++
		} else if cmp > 0 {
			result = append(result, otherKeys[j])
			j++
		} else {
			// 键相同，只添加一次（合并时以 other 为准）
			result = append(result, otherKeys[j])
			i++
			j++
		}
	}

	// 添加剩余元素
	for i < len(c.sortedKeys) {
		result = append(result, c.sortedKeys[i])
		i++
	}
	for j < len(otherKeys) {
		result = append(result, otherKeys[j])
		j++
	}

	return result
}
