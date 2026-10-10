package slice_collcection

import "cmp"

// NewCollection 创建并返回一个新的 Collection 实例。
//
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

// cmpAs 按具体有序类型比较两个 T。T 在构造时已经确认是 N。
func cmpAs[T any, N cmp.Ordered](a, b T) int {
	return cmp.Compare(any(a).(N), any(b).(N))
}

// defaultCompare 为预声明有序类型生成比较函数。
// 结构体和 type MyInt int 这类定义类型返回 nil，调用方必须传入比较函数。
func defaultCompare[T any]() func(T, T) int {
	var zero T
	switch any(zero).(type) {
	case int:
		return cmpAs[T, int]
	case int8:
		return cmpAs[T, int8]
	case int16:
		return cmpAs[T, int16]
	case int32:
		return cmpAs[T, int32]
	case int64:
		return cmpAs[T, int64]
	case uint:
		return cmpAs[T, uint]
	case uint8:
		return cmpAs[T, uint8]
	case uint16:
		return cmpAs[T, uint16]
	case uint32:
		return cmpAs[T, uint32]
	case uint64:
		return cmpAs[T, uint64]
	case float32:
		return cmpAs[T, float32]
	case float64:
		return cmpAs[T, float64]
	case string:
		return cmpAs[T, string]
	default:
		return nil
	}
}

// toFloat 把预声明数值转成 float64。
// 泛型联合类型不能直接转换，所以按具体类型分支。
func toFloat[N number](n N) float64 {
	switch v := any(n).(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		return 0
	}
}

// NewEmptyCollection 返回一个空的Collection
func NewEmptyCollection[T any]() *Collection[T] {
	return NewCollection[T](nil)
}
