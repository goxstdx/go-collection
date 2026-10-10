package slice_collection

// number 是底层类型为整数或浮点数的类型。
// 带 ~，type Money int64 这类定义类型也可以交给 Sum、Avg、Median。
type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Collection 是集合操作的核心结构，实现了对各种数据类型的集合操作。
//
// Collection 使用泛型参数 T 来支持不同类型的数据，包括基本类型和结构体类型。
// 它提供了丰富的方法来操作和转换这些数据，如过滤、映射、排序、查找、聚合等。
//
// 使用示例：
//
//	intColl := NewCollection([]int{1, 2, 3, 4, 5})
//	filtered := intColl.Filter(func(item int, key int) bool {
//	    return item > 2
//	})
type Collection[T any] struct {
	value []T // 数组
}
