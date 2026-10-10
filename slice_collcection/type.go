package slice_collcection

// number 是可以累加为 float64 的预声明数值类型。
// 不用 ~，因为泛型函数体不能把整组联合类型直接转成 float64，
// 只能按具体类型做类型开关；定义类型请在传入的函数里先转成预声明类型。
type number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
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

	// orderedCmp 只缓存预声明有序类型的默认比较函数，在第一次需要比较时生成。
	orderedCmp func(T, T) int
}
