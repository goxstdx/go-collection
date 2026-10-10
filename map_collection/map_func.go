package map_collection

// 本文件提供一组针对泛型 map 的辅助方法，参考 Laravel Collection 的常见能力，
// 但以 Go 原生泛型与标准库 maps 包为基础实现，保持简洁与实用。
//
// 说明：
// - 这些函数均为纯函数（除 InPlace 结尾者外），不会修改传入的原始 map，
//   会返回一个新的 map 或结果值。
// - Go 的 map 无序，这里未提供与顺序相关的 API。
// - 优先使用 Go 1.27+ 的标准库 maps 包（Keys/Values/Clone/Equal/Copy）。

import (
	"maps"
)

// Keys 返回 map 的所有键。
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	res := make([]K, 0, len(m))
	for k := range m {
		res = append(res, k)
	}

	return res
}

// Values 返回 map 的所有值。
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	res := make([]V, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

// Clone 返回一个 m 的拷贝。
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	return maps.Clone(m)
}

// Equal 比较两个 map 是否相等（键和值逐一相等）。
func Equal[M1 ~map[K]V, M2 ~map[K]V, K comparable, V comparable](m1 M1, m2 M2) bool {
	return maps.Equal(m1, m2)
}

// Has 判断 key 是否在 map 中。
func Has[M ~map[K]V, K comparable, V any](m M, key K) bool {
	_, ok := m[key]
	return ok
}

// Get 获取 key 对应的值；如果不存在，则返回零值。
func Get[M ~map[K]V, K comparable, V any](m M, key K) V {
	return m[key]
}

// GetOr 获取 key 对应的值；如果不存在，返回默认值 def。
func GetOr[M ~map[K]V, K comparable, V any](m M, key K, def V) V {
	if v, ok := m[key]; ok {
		return v
	}
	return def
}

// Set 返回设置了 key->val 后的新 map（不修改原 map）。
func Set[M ~map[K]V, K comparable, V any](m M, key K, val V) M {
	nm := maps.Clone(m)
	nm[key] = val
	return nm
}

// SetInPlace 直接在原 map 上设置 key->val。
func SetInPlace[M ~map[K]V, K comparable, V any](m M, key K, val V) {
	m[key] = val
}

// Delete 返回删除 key 后的新 map（不修改原 map）。
func Delete[M ~map[K]V, K comparable, V any](m M, key K) M {
	nm := maps.Clone(m)
	delete(nm, key)
	return nm
}

// DeleteInPlace 直接在原 map 上删除 key。
func DeleteInPlace[M ~map[K]V, K comparable, V any](m M, key K) {
	delete(m, key)
}

// Merge 返回一个新 map，包含 m 与 other 的所有键值，若键冲突则以 other 为准。
func Merge[M ~map[K]V, K comparable, V any](m M, other M) M {
	// 两边都为 nil 时返回空 map 而非 nil，避免调用方后续写入 panic。
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

// MergeInPlace 将 other 的键值合并进 m，键冲突以 other 为准。
func MergeInPlace[M ~map[K]V, K comparable, V any](m M, other M) {
	maps.Copy(m, other)
}

// Only 仅保留 keys 中出现的键，返回新 map。
func Only[M ~map[K]V, K comparable, V any](m M, keys []K) M {
	nm := make(M, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			nm[k] = v
		}
	}
	return nm
}

// Except 排除 keys 中出现的键，返回新 map（不修改原 map）。
func Except[M ~map[K]V, K comparable, V any](m M, keys []K) M {
	// 输入为空时返回空 map 而非 nil，避免调用方在结果上写入 panic。
	if len(m) == 0 {
		return make(M)
	}
	nm := maps.Clone(m)
	for _, k := range keys {
		delete(nm, k)
	}
	return nm
}

// MapValues 对值进行映射，保留原键，返回新 map。
func MapValues[M ~map[K]V, K comparable, V any, R any](m M, fn func(value V, key K) R) map[K]R {
	res := make(map[K]R, len(m))
	for k, v := range m {
		res[k] = fn(v, k)
	}
	return res
}

// MapKeys 对键进行映射（可能改变键类型），返回新 map；若映射过程中出现键冲突，后者覆盖前者。
func MapKeys[M ~map[K]V, K comparable, V any, NK comparable](m M, fn func(key K, value V) NK) map[NK]V {
	res := make(map[NK]V, len(m))
	for k, v := range m {
		res[fn(k, v)] = v
	}
	return res
}

// Filter 过滤键值对，返回新 map。
func Filter[M ~map[K]V, K comparable, V any](m M, fn func(value V, key K) bool) M {
	nm := make(M)
	for k, v := range m {
		if fn(v, k) {
			nm[k] = v
		}
	}
	return nm
}

// Each 按 map 自身顺序遍历，fn 返回 false 时终止遍历。
func Each[M ~map[K]V, K comparable, V any](m M, fn func(value V, key K) bool) {
	for k, v := range m {
		if !fn(v, k) {
			return
		}
	}
}

// Foreach 对每个键值对执行一次回调，顺序不确定。
func Foreach[M ~map[K]V, K comparable, V any](m M, fn func(value V, key K)) {
	for k, v := range m {
		fn(v, k)
	}
}

// Reduce 聚合：将 map 折叠为一个结果。
func Reduce[M ~map[K]V, K comparable, V any, R any](m M, init R, fn func(acc R, value V, key K) R) R {
	acc := init
	for k, v := range m {
		acc = fn(acc, v, k)
	}
	return acc
}

// Pluck 用 fn 从 map 的每个 value 取出一个值。遍历顺序不稳定。
func Pluck[M ~map[K]V, K comparable, V any, R any](m M, fn func(V) R) []R {
	result := make([]R, 0, len(m))
	for _, v := range m {
		result = append(result, fn(v))
	}

	return result
}
