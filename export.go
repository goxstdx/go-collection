package collection

import (
	"github.com/ZHOUXING1997/collection/map_collection"
	"github.com/ZHOUXING1997/collection/slice_collection"
)

// NewSliceCollect 创建并返回一个新的切片 Collection。
//
// 不复制 values。集合和调用方传入的切片共用底层数组。
// 比较不在创建时决定。int、uint、float、string 在 Sort、Max、Contains 等方法调用时生成比较函数。
// 结构体和定义类型必须把比较函数传给对应的 By 方法。
func NewSliceCollect[T any](values []T) *slice_collection.Collection[T] {
	return slice_collection.NewCollection[T](values)
}

// NewEmptyCollection 返回一个空的Collection
func NewEmptyCollection[T any]() *slice_collection.Collection[T] {
	return slice_collection.NewEmptyCollection[T]()
}

// NewMapCollect 创建并返回一个新的 map Collection 实例。
//
// 不记录 key 顺序。MinBy、MaxBy 每次调用都重新扫描。EachSorted 每次调用都重新排序。
func NewMapCollect[K comparable, V any](values map[K]V) *map_collection.Collection[K, V] {
	return map_collection.NewCollection[map[K]V, K, V](values)
}

// NewEmptyMapCollection 返回一个空的 map Collection。
func NewEmptyMapCollection[K comparable, V any]() *map_collection.Collection[K, V] {
	return map_collection.NewCollection[map[K]V, K, V](make(map[K]V))
}
