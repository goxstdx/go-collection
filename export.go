package collection

import (
	"github.com/ZHOUXING1997/collection/map_collection"
	"github.com/ZHOUXING1997/collection/slice_collcection"
)

// NewSliceCollect 创建并返回一个新的切片 Collection。
//
// 比较不在创建时决定。int、uint、float、string 在 Sort、Max、Contains 等方法调用时生成比较函数。
// 结构体和定义类型必须把比较函数传给对应的 By 方法。
func NewSliceCollect[T any](values []T) *slice_collcection.Collection[T] {
	return slice_collcection.NewCollection[T](values)
}

// NewEmptyCollection 返回一个空的Collection
func NewEmptyCollection[T any]() *slice_collcection.Collection[T] {
	return slice_collcection.NewEmptyCollection[T]()
}

// NewMapCollect 创建并返回一个新的 map Collection 实例。
//
// 该函数接收一个 map 作为输入，并返回一个支持链式操作的 Collection 对象。
// 支持通过函数式选项配置比较函数。
//
// 参数:
//   - values: 要初始化的 map
//   - opts: 可选的配置选项，如 WithKeyCompare、WithValCompare
//
// 返回:
//   - *map_collection.Collection[K, V]: 新创建的 map Collection 实例
//
// 示例:
//   - NewMapCollect(m) // 默认不排序
//   - NewMapCollect(m, WithKeyCompare(func(a, b string) int { ... })) // 使用自定义key比较函数
func NewMapCollect[K comparable, V any](values map[K]V, opts ...map_collection.CollectionOption[K, V]) *map_collection.Collection[K, V] {
	return map_collection.NewCollection[map[K]V, K, V](values, opts...)
}

// NewEmptyMapCollection 返回一个空的 map Collection（默认排序）
func NewEmptyMapCollection[K comparable, V any]() *map_collection.Collection[K, V] {
	return map_collection.NewCollection[map[K]V, K, V](make(map[K]V))
}
