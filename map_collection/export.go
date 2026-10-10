package map_collection

import (
	"sync"
)

// Collection 是集合操作的核心结构，实现了对各种数据类型的集合操作。
// Collection 是针对map类型的集合操作核心结构
type Collection[K comparable, V any] struct {
	value map[K]V // map数据
}

// NewCollection 创建并返回一个新的 Collection 实例。
// 不记录 key 顺序。MinBy、MaxBy 每次调用都重新扫描。EachSorted 每次调用都重新排序。结果都不保存。
// values 为 nil 时改用空 map，保证后续 Set、MergeInPlace 等写操作不会往 nil map 写入而 panic。
func NewCollection[T ~map[K]V, K comparable, V any](values T) *Collection[K, V] {
	if values == nil {
		values = make(T)
	}
	return &Collection[K, V]{
		value: values,
	}
}

// SafeCollection 是 Collection 的线程安全包装器。
// 使用读写锁保护底层 Collection。
// 回调里不要再调用同一个 SafeCollection。写方法会等待这把锁，读方法在已有写操作等待时也可能卡住。
type SafeCollection[K comparable, V any] struct {
	mu   sync.RWMutex
	coll *Collection[K, V]
}

// NewSafeCollection 创建一个线程安全的 Collection 包装器。
func NewSafeCollection[T ~map[K]V, K comparable, V any](values T) *SafeCollection[K, V] {
	return &SafeCollection[K, V]{
		coll: NewCollection(values),
	}
}
