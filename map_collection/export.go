package map_collection

import (
	"sync"
)

// Collection 是集合操作的核心结构，实现了对各种数据类型的集合操作。
// Collection 是针对map类型的集合操作核心结构
type Collection[K comparable, V any] struct {
	value map[K]V // map数据

	// 比较函数，用于值的比较
	valCompareFunc func(V, V) int
	// key的比较函数，用于key的排序等操作
	keyCompareFunc func(K, K) int

	// 缓存排序后的keys,用于稳定的First/Last操作
	sortedKeys []K
}

// WithKeyCompare 设置 key 的比较函数（用于排序）
// compareFunc: 比较函数，返回 -1(小于)、0(等于)、1(大于)
func WithKeyCompare[K comparable, V any](compareFunc func(K, K) int) CollectionOption[K, V] {
	return func(c *Collection[K, V]) {
		c.keyCompareFunc = compareFunc
	}
}

// WithValCompare 设置 value 的比较函数（用于按值排序）
// compareFunc: 比较函数，返回 -1(小于)、0(等于)、1(大于)
func WithValCompare[K comparable, V any](compareFunc func(V, V) int) CollectionOption[K, V] {
	return func(c *Collection[K, V]) {
		c.valCompareFunc = compareFunc
	}
}

// NewCollection 创建并返回一个新的 Collection 实例
// 支持函数式选项模式配置比较函数
func NewCollection[T ~map[K]V, K comparable, V any](values T, opts ...CollectionOption[K, V]) *Collection[K, V] {
	coll := &Collection[K, V]{
		value: values,
	}

	// 比较函数在选项里设置，排好序后再把集合交给调用方。
	// First、Last 只读 sortedKeys，避免在 SafeCollection 的读锁里写入。
	for _, opt := range opts {
		opt(coll)
	}
	coll.initSortedKeys()

	return coll
}

// SafeCollection 是 Collection 的线程安全包装器
// 使用读写锁(RWMutex)保护底层Collection，支持并发读取和安全的写入
type SafeCollection[K comparable, V any] struct {
	mu   sync.RWMutex
	coll *Collection[K, V]
}

// NewSafeCollection 创建一个线程安全的 Collection 包装器
// 支持函数式选项模式配置比较函数
func NewSafeCollection[T ~map[K]V, K comparable, V any](values T, opts ...CollectionOption[K, V]) *SafeCollection[K, V] {
	return &SafeCollection[K, V]{
		coll: NewCollection(values, opts...),
	}
}
