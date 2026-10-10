package slice_collection

import "cmp"

// NewCollection 创建并返回一个新的 Collection 实例。
//
// 不复制 values。集合和调用方传入的切片共用底层数组，容量足够时，Append、Remove、Insert 等会写到原切片。
// 比较函数不在这里决定。预声明数值类型和 string 在排序、查找等方法调用时生成。
// 结构体和定义类型必须向对应的 By 方法传入比较函数。
//
// 参数:
//   - values: 要初始化的切片
//
// 返回:
//   - *Collection[T]: 新创建的 Collection 实例
func NewCollection[T any](values []T) *Collection[T] {
	return &Collection[T]{value: values}
}

// cloneSlice 复制元素。结果容量等于长度，后续 append 不会写回 src。
func cloneSlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

// orderedFunc 把 cmp.Compare[N] 转成 func(T, T) int。
// 类型开关已经确认 T 就是 N，断言只在生成比较函数时发生一次。
func orderedFunc[T any, N cmp.Ordered]() func(T, T) int {
	return any(cmp.Compare[N]).(func(T, T) int)
}

// defaultCompare 为预声明有序类型生成比较函数。
// 结构体和 type MyInt int 这类定义类型返回 nil，调用方必须传入比较函数。
func defaultCompare[T any]() func(T, T) int {
	var zero T
	switch any(zero).(type) {
	case int:
		return orderedFunc[T, int]()
	case int8:
		return orderedFunc[T, int8]()
	case int16:
		return orderedFunc[T, int16]()
	case int32:
		return orderedFunc[T, int32]()
	case int64:
		return orderedFunc[T, int64]()
	case uint:
		return orderedFunc[T, uint]()
	case uint8:
		return orderedFunc[T, uint8]()
	case uint16:
		return orderedFunc[T, uint16]()
	case uint32:
		return orderedFunc[T, uint32]()
	case uint64:
		return orderedFunc[T, uint64]()
	case float32:
		return orderedFunc[T, float32]()
	case float64:
		return orderedFunc[T, float64]()
	case string:
		return orderedFunc[T, string]()
	default:
		return nil
	}
}

// NewEmptyCollection 返回一个空的Collection
func NewEmptyCollection[T any]() *Collection[T] {
	return NewCollection[T](nil)
}
