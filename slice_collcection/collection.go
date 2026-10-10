// Package slice_collcection 提供针对泛型切片的链式集合操作。
//
// 提供过滤、映射、排序、查找、聚合、分组、分页、序列化等常见操作，
// 旨在在不追求极致性能的业务场景下提升开发效率。
package slice_collcection

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/ZHOUXING1997/collection/errorx"
)

// basicOrderedCmp 返回预声明有序类型的比较函数，并缓存在集合上。
// 结构体等非基础类型返回 NoComparableError，调用方必须使用 By 方法。
func (c *Collection[T]) basicOrderedCmp() (func(T, T) int, error) {
	if c.orderedCmp != nil {
		return c.orderedCmp, nil
	}
	fn := defaultCompare[T]()
	if fn == nil {
		return nil, errorx.NoComparableError
	}
	c.orderedCmp = fn
	return fn, nil
}

// // Err 返回Collection的错误信息
// func (c *Collection[T]) Err() error {
// 	return c.err
// }
//
// // SetErr 设置Collection的错误信息
// func (c *Collection[T]) SetErr(err error) *Collection[T] {
// 	c.err = err
// 	return c
// }

// Copy 复制一个新的Collection
func (c *Collection[T]) Copy() *Collection[T] {
	// 深拷贝切片,避免共享底层数组
	copied := make([]T, len(c.value))
	copy(copied, c.value)
	return NewCollection[T](copied)
}

// IsEmpty 判断是否为空
func (c *Collection[T]) IsEmpty() bool {
	return len(c.value) == 0
}

// IsNotEmpty 判断是否不为空
func (c *Collection[T]) IsNotEmpty() bool {
	return len(c.value) != 0
}

// Append 添加元素
func (c *Collection[T]) Append(item T) *Collection[T] {
	c.value = append(c.value, item)
	return c
}

// Remove 删除元素
func (c *Collection[T]) Remove(index int) *Collection[T] {
	// 添加边界检查
	if index < 0 || index >= len(c.value) {
		return c
	}
	c.value = append(c.value[:index], c.value[index+1:]...)
	return c
}

// Insert 插入元素
func (c *Collection[T]) Insert(index int, item T) *Collection[T] {
	// 添加边界检查
	if index < 0 || index > len(c.value) {
		return c
	}
	arr := make([]T, 0, len(c.value)+1)
	arr = append(arr, c.value[:index]...)
	arr = append(arr, item)
	arr = append(arr, c.value[index:]...)
	newColl := NewCollection[T](arr)
	return newColl
}

// Search 在预声明有序类型中查找元素，返回首次出现的下标。
func (c *Collection[T]) Search(item T) (int, error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return -1, err
	}
	return c.searchWith(item, cmpFn)
}

// SearchBy 用 cmpFn 查找元素。cmpFn 返回 0 表示相等。
func (c *Collection[T]) SearchBy(item T, cmpFn func(a, b T) int) (int, error) {
	if cmpFn == nil {
		return -1, errorx.NilFunc
	}
	return c.searchWith(item, cmpFn)
}

func (c *Collection[T]) searchWith(item T, cmpFn func(a, b T) int) (int, error) {
	for i, v := range c.value {
		if cmpFn(v, item) == 0 {
			return i, nil
		}
	}
	return -1, errorx.NotFoundError
}

// Unique 对预声明有序类型去重，保留首次出现的元素。
func (c *Collection[T]) Unique() (*Collection[T], error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return nil, err
	}
	return c.uniqueWith(cmpFn), nil
}

// UniqueBy 用 cmpFn 去重。cmpFn 返回 0 表示相等。保留首次出现的元素。
func (c *Collection[T]) UniqueBy(cmpFn func(a, b T) int) (*Collection[T], error) {
	if cmpFn == nil {
		return nil, errorx.NilFunc
	}
	return c.uniqueWith(cmpFn), nil
}

func (c *Collection[T]) uniqueWith(cmpFn func(a, b T) int) *Collection[T] {
	res := make([]T, 0, len(c.value))
	for _, v := range c.value {
		found := false
		for _, existed := range res {
			if cmpFn(existed, v) == 0 {
				found = true
				break
			}
		}
		if !found {
			res = append(res, v)
		}
	}
	return NewCollection(res)
}

// Filter 过滤
func (c *Collection[T]) Filter(f func(item T, key int) bool) *Collection[T] {
	// 按照某个方法进行过滤, 保留符合的
	res := make([]T, 0, len(c.value))
	for i, v := range c.value {
		if f(v, i) {
			res = append(res, v)
		}
	}

	return NewCollection[T](res)
}

// Reject 过滤
func (c *Collection[T]) Reject(f func(item T, key int) bool) *Collection[T] {
	res := make([]T, 0, len(c.value))
	for i, v := range c.value {
		if !f(v, i) {
			res = append(res, v)
		}
	}

	return NewCollection[T](res)
}

// First 获取第一个元素
func (c *Collection[T]) First() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}

	return c.value[0]
}

// Last 获取最后一个元素
func (c *Collection[T]) Last() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}

	return c.value[len(c.value)-1]
}

// Slice 获取数组片段
func (c *Collection[T]) Slice(params ...int) (*Collection[T], error) {
	if len(params) == 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid params")
	}
	start := params[0]
	if start < 0 || start >= len(c.value) {
		return NewCollection[T](nil), fmt.Errorf("invalid start index")
	}
	if len(params) == 1 {
		return NewCollection(c.value[start:]), nil
	}
	end := params[1]
	if end < 0 || end > len(c.value) {
		return NewCollection[T](nil), fmt.Errorf("invalid end index")
	}
	if start > end {
		return NewCollection[T](nil), fmt.Errorf("start index should be less than end index")
	}

	return NewCollection(c.value[start:end]), nil
}

// Index 获取某个下标
func (c *Collection[T]) Index(i int) T {

	// 获取某个下标，对所有Collection生效
	var zero T
	if i < 0 || i >= len(c.value) {
		return zero
	}

	return c.value[i]
}

// SetIndex 设置某个下标
func (c *Collection[T]) SetIndex(i int, val T) *Collection[T] {
	// if c.err != nil {
	// 	return c
	// }

	// 设置数组的下标为某个值
	if i < 0 || i >= len(c.value) {
		return c
	}
	c.value[i] = val
	return c
}

// Count 获取数组长度
func (c *Collection[T]) Count() int {
	// 获取数组长度，对所有Collection生效
	return len(c.value)
}

// Merge 合并数组
func (c *Collection[T]) Merge(arr *Collection[T]) *Collection[T] {
	// 将两个数组进行合并
	if arr == nil {
		return c
	}

	res := c.Copy()
	for i := 0; i < arr.Count(); i++ {
		res.value = append(res.value, arr.Index(i))
	}

	return res
}

// Each 遍历
// 当 f 返回 false 时，终止
func (c *Collection[T]) Each(f func(item T, key int) bool) {
	for i, v := range c.value {
		if !f(v, i) {
			return
		}
	}
}

// Foreach 遍历，能保证每一个元素都被遍历到，内部决定是否跳过
func (c *Collection[T]) Foreach(f func(item T, key int)) {
	for i, v := range c.value {
		f(v, i)
	}
}

// Map 映射。R 由 f 的返回值决定，可以与 T 不同。
func (c *Collection[T]) Map[R any](f func(item T, key int) R) *Collection[R] {
	res := make([]R, 0, len(c.value))
	for i, v := range c.value {
		res = append(res, f(v, i))
	}

	return NewCollection(res)
}

// MapFilter 映射并过滤。f 返回 false 的元素会被丢弃。
func (c *Collection[T]) MapFilter[R any](f func(item T, key int) (R, bool)) *Collection[R] {
	res := make([]R, 0, len(c.value))
	for i, v := range c.value {
		r, ok := f(v, i)
		if ok {
			res = append(res, r)
		}
	}

	return NewCollection(res)
}

// Reduce 求和
func (c *Collection[T]) Reduce(f func(carry T, item T) T) T {
	var zero T
	if c.IsEmpty() {
		return zero
	}
	res := c.value[0]
	for i := 1; i < len(c.value); i++ {
		res = f(res, c.value[i])
	}

	return res
}

// Every 判断是否所有元素都满足条件
func (c *Collection[T]) Every(f func(item T, key int) bool) bool {
	for i, v := range c.value {
		if !f(v, i) {
			return false
		}
	}

	return true
}

// ForPage 分页
func (c *Collection[T]) ForPage(page int, perPage int) (*Collection[T], error) {
	if page <= 0 || perPage <= 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid page or perPage")
	}
	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(c.value) {
		return NewCollection[T](nil), nil
	}
	if end > len(c.value) {
		end = len(c.value)
	}

	return &Collection[T]{value: c.value[start:end]}, nil
}

// Nth 每隔n个取一个
func (c *Collection[T]) Nth(n int, offset int) (*Collection[T], error) {
	if n <= 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid n")
	}
	res := make([]T, 0, len(c.value)/n+1)
	for i := offset; i < len(c.value); i += n {
		res = append(res, c.value[i])
	}

	return NewCollection[T](res), nil
}

// Pad 填充
func (c *Collection[T]) Pad(count int, def T) *Collection[T] {
	if count <= len(c.value) {
		return c.Copy()
	}
	res := make([]T, count)
	for i := 0; i < len(c.value); i++ {
		res[i] = c.value[i]
	}
	for i := len(c.value); i < count; i++ {
		res[i] = def
	}

	return NewCollection[T](res)
}

// Pop 弹出最后一个元素
func (c *Collection[T]) Pop() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}
	res := c.value[len(c.value)-1]
	c.value = c.value[:len(c.value)-1]
	return res
}

// Push 添加元素
func (c *Collection[T]) Push(item T) *Collection[T] {
	c.value = append(c.value, item)
	return c
}

// Prepend 添加元素到头部
func (c *Collection[T]) Prepend(item T) *Collection[T] {
	res := make([]T, 0, len(c.value)+1)
	res = append(res, item)
	res = append(res, c.value...)
	return NewCollection[T](res)
}

// Random 随机取一个元素
func (c *Collection[T]) Random() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return c.value[random.Intn(len(c.value))]
}

// Reverse 反转
func (c *Collection[T]) Reverse() *Collection[T] {
	res := make([]T, 0, len(c.value))
	for i := len(c.value) - 1; i >= 0; i-- {
		res = append(res, c.value[i])
	}

	return NewCollection(res)
}

// Shuffle 随机排序
func (c *Collection[T]) Shuffle() *Collection[T] {
	res := make([]T, len(c.value))
	perm := rand.Perm(len(c.value))
	for i, v := range perm {
		res[v] = c.value[i]
	}

	return NewCollection(res)
}

// GroupBy 按 f 的返回值分组。K 必须可比较。
func (c *Collection[T]) GroupBy[K comparable](f func(T, int) K) map[K]*Collection[T] {
	res := make(map[K]*Collection[T])
	for i, v := range c.value {
		key := f(v, i)
		if _, ok := res[key]; !ok {
			res[key] = NewCollection[T](nil)
		}
		res[key].Push(v)
	}

	return res
}

// Split 按照size个数进行分组
func (c *Collection[T]) Split(size int) []*Collection[T] {
	if size <= 0 {
		return []*Collection[T]{}
	}
	var res []*Collection[T]
	for i := 0; i < len(c.value); i += size {
		end := i + size
		if end > len(c.value) {
			end = len(c.value)
		}
		res = append(res, NewCollection(c.value[i:end]))
	}

	return res
}

// DD 打印出当前数组结构
func (c *Collection[T]) DD() {
	var zero T
	ret := fmt.Sprintf("Collection(%d, %T):{\n", c.Count(), zero)
	for k, v := range c.value {
		ret = ret + fmt.Sprintf("\t%d:\t%v\n", k, v)
	}
	ret = ret + "}\n"
	fmt.Print(ret)
}

// PluckField 用 fn 从每个元素取出一个值，返回新集合。
func (c *Collection[T]) PluckField[R any](fn func(T) R) *Collection[R] {
	res := make([]R, 0, len(c.value))
	for _, v := range c.value {
		res = append(res, fn(v))
	}

	return NewCollection(res)
}

// SortBy 按 fn 提取的值升序排列，直接修改当前集合。
// 使用稳定排序：本次比较相等的元素保持排序前的相对顺序，因此可以先排次要字段，再排主要字段。
func (c *Collection[T]) SortBy[R cmp.Ordered](fn func(T) R) *Collection[T] {
	slices.SortStableFunc(c.value, func(a, b T) int {
		return cmp.Compare(fn(a), fn(b))
	})

	return c
}

// SortByFunc 按自定义比较函数排序，直接修改当前集合。fn 返回 true 表示 v1 应排在 v2 前面。
// 使用稳定排序，比较结果相等时保留原相对顺序。
func (c *Collection[T]) SortByFunc(fn func(v1, v2 T) bool) (*Collection[T], error) {
	sort.SliceStable(c.value, func(i, j int) bool {
		return fn(c.value[i], c.value[j])
	})

	return c, nil
}

// SortByDesc 按 fn 提取的值降序排列，直接修改当前集合。
// 使用稳定排序：本次比较相等的元素保持排序前的相对顺序。
func (c *Collection[T]) SortByDesc[R cmp.Ordered](fn func(T) R) *Collection[T] {
	slices.SortStableFunc(c.value, func(a, b T) int {
		return cmp.Compare(fn(b), fn(a))
	})

	return c
}

// KeyBy 用 fn 的返回值作为 key。相同 key 时，后者覆盖前者。
func (c *Collection[T]) KeyBy[K comparable](fn func(T) K) map[K]T {
	res := make(map[K]T, len(c.value))
	for _, v := range c.value {
		res[fn(v)] = v
	}

	return res
}

// Max 返回预声明有序类型中的最大元素。空集合返回 EmptyError，不把零值当作结果。
func (c *Collection[T]) Max() (T, error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		var zero T
		return zero, err
	}
	return c.maxWith(cmpFn)
}

// MaxBy 用 cmpFn 取最大元素。cmpFn 返回正数表示 a 大于 b。
func (c *Collection[T]) MaxBy(cmpFn func(a, b T) int) (T, error) {
	if cmpFn == nil {
		var zero T
		return zero, errorx.NilFunc
	}
	return c.maxWith(cmpFn)
}

func (c *Collection[T]) maxWith(cmpFn func(a, b T) int) (T, error) {
	var zero T
	if c.IsEmpty() {
		return zero, errorx.EmptyError
	}
	maxElem := c.value[0]
	for _, v := range c.value[1:] {
		if cmpFn(v, maxElem) > 0 {
			maxElem = v
		}
	}
	return maxElem, nil
}

// Min 返回预声明有序类型中的最小元素。空集合返回 EmptyError，不把零值当作结果。
func (c *Collection[T]) Min() (T, error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		var zero T
		return zero, err
	}
	return c.minWith(cmpFn)
}

// MinBy 用 cmpFn 取最小元素。cmpFn 返回负数表示 a 小于 b。
func (c *Collection[T]) MinBy(cmpFn func(a, b T) int) (T, error) {
	if cmpFn == nil {
		var zero T
		return zero, errorx.NilFunc
	}
	return c.minWith(cmpFn)
}

func (c *Collection[T]) minWith(cmpFn func(a, b T) int) (T, error) {
	var zero T
	if c.IsEmpty() {
		return zero, errorx.EmptyError
	}
	minElem := c.value[0]
	for _, v := range c.value[1:] {
		if cmpFn(v, minElem) < 0 {
			minElem = v
		}
	}
	return minElem, nil
}

// Contains 判断预声明有序类型是否包含某个元素。
func (c *Collection[T]) Contains(obj T) (bool, error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return false, err
	}
	return c.containsWith(obj, cmpFn), nil
}

// ContainsBy 用 cmpFn 判断是否包含。cmpFn 返回 0 表示相等。
func (c *Collection[T]) ContainsBy(obj T, cmpFn func(a, b T) int) (bool, error) {
	if cmpFn == nil {
		return false, errorx.NilFunc
	}
	return c.containsWith(obj, cmpFn), nil
}

func (c *Collection[T]) containsWith(obj T, cmpFn func(a, b T) int) bool {
	for _, v := range c.value {
		if cmpFn(v, obj) == 0 {
			return true
		}
	}
	return false
}

// ContainsCount 返回预声明有序类型中等于 obj 的元素个数。
func (c *Collection[T]) ContainsCount(obj T) (int, error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return 0, err
	}
	return c.containsCountWith(obj, cmpFn), nil
}

// ContainsCountBy 用 cmpFn 统计相等元素的个数。cmpFn 返回 0 表示相等。
func (c *Collection[T]) ContainsCountBy(obj T, cmpFn func(a, b T) int) (int, error) {
	if cmpFn == nil {
		return 0, errorx.NilFunc
	}
	return c.containsCountWith(obj, cmpFn), nil
}

func (c *Collection[T]) containsCountWith(obj T, cmpFn func(a, b T) int) int {
	count := 0
	for _, v := range c.value {
		if cmpFn(v, obj) == 0 {
			count++
		}
	}
	return count
}

// Diff 返回存在于当前集合、但不存在于 arr 的元素。仅适用于预声明有序类型。
func (c *Collection[T]) Diff(arr *Collection[T]) (*Collection[T], error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return nil, err
	}
	return c.diffWith(arr, cmpFn), nil
}

// DiffBy 用 cmpFn 计算差集。cmpFn 返回 0 表示相等。
func (c *Collection[T]) DiffBy(arr *Collection[T], cmpFn func(a, b T) int) (*Collection[T], error) {
	if cmpFn == nil {
		return nil, errorx.NilFunc
	}
	return c.diffWith(arr, cmpFn), nil
}

func (c *Collection[T]) diffWith(arr *Collection[T], cmpFn func(a, b T) int) *Collection[T] {
	res := make([]T, 0, len(c.value))
	for _, v := range c.value {
		if arr == nil || !arr.containsWith(v, cmpFn) {
			res = append(res, v)
		}
	}
	return NewCollection(res)
}

// Sort 按元素自身升序排列，仅适用于预声明有序类型。直接修改当前集合。
// 稳定排序，比较结果相等时保留原相对顺序。
func (c *Collection[T]) Sort() (*Collection[T], error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(c.value, cmpFn)
	return c, nil
}

// SortDesc 按元素自身降序排列，仅适用于预声明有序类型。直接修改当前集合。
// 稳定排序，比较结果相等时保留原相对顺序。
func (c *Collection[T]) SortDesc() (*Collection[T], error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(c.value, func(a, b T) int {
		return cmpFn(b, a)
	})
	return c, nil
}

// Join 进行拼接
func (c *Collection[T]) Join(split string, format ...func(item interface{}) string) string {
	if c.IsEmpty() {
		return ""
	}

	// 使用 strings.Builder 优化字符串拼接性能
	var builder strings.Builder
	// 预估容量以减少内存分配
	builder.Grow(len(c.value) * (len(split) + 10))

	for i, v := range c.value {
		if len(format) > 0 {
			builder.WriteString(format[0](v))
		} else {
			builder.WriteString(fmt.Sprintf("%v", v))
		}
		if i != len(c.value)-1 {
			builder.WriteString(split)
		}
	}

	return builder.String()
}

// Union 返回两个预声明有序类型集合的并集，保留当前集合中的重复情况，再追加 arr 中没有的元素。
func (c *Collection[T]) Union(arr *Collection[T]) (*Collection[T], error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return nil, err
	}
	return c.unionWith(arr, cmpFn), nil
}

// UnionBy 用 cmpFn 计算并集。cmpFn 返回 0 表示相等。
func (c *Collection[T]) UnionBy(arr *Collection[T], cmpFn func(a, b T) int) (*Collection[T], error) {
	if cmpFn == nil {
		return nil, errorx.NilFunc
	}
	return c.unionWith(arr, cmpFn), nil
}

func (c *Collection[T]) unionWith(arr *Collection[T], cmpFn func(a, b T) int) *Collection[T] {
	res := make([]T, len(c.value), len(c.value)+len(arrValue(arr)))
	copy(res, c.value)
	if arr == nil {
		return NewCollection(res)
	}
	for _, v := range arr.value {
		if !c.containsWith(v, cmpFn) {
			res = append(res, v)
		}
	}
	return NewCollection(res)
}

func arrValue[T any](arr *Collection[T]) []T {
	if arr == nil {
		return nil
	}
	return arr.value
}

// Intersect 返回两个预声明有序类型集合的交集，保留当前集合中的出现顺序和重复。
func (c *Collection[T]) Intersect(arr *Collection[T]) (*Collection[T], error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		return nil, err
	}
	return c.intersectWith(arr, cmpFn), nil
}

// IntersectBy 用 cmpFn 计算交集。cmpFn 返回 0 表示相等。
func (c *Collection[T]) IntersectBy(arr *Collection[T], cmpFn func(a, b T) int) (*Collection[T], error) {
	if cmpFn == nil {
		return nil, errorx.NilFunc
	}
	return c.intersectWith(arr, cmpFn), nil
}

func (c *Collection[T]) intersectWith(arr *Collection[T], cmpFn func(a, b T) int) *Collection[T] {
	res := make([]T, 0, len(c.value))
	for _, v := range c.value {
		if arr != nil && arr.containsWith(v, cmpFn) {
			res = append(res, v)
		}
	}
	return NewCollection(res)
}

// Avg 对 fn 取出的数值求平均。空集合返回 0。
func (c *Collection[T]) Avg[N number](fn func(T) N) float64 {
	if c.IsEmpty() {
		return 0
	}

	return c.Sum(fn) / float64(len(c.value))
}

// Median 对 fn 取出的数值求中位数。偶数个时取中间两个的平均值。空集合返回 0。
// 不修改原集合。
func (c *Collection[T]) Median[N number](fn func(T) N) float64 {
	if c.IsEmpty() {
		return 0
	}

	nums := make([]float64, len(c.value))
	for i, item := range c.value {
		nums[i] = toFloat(fn(item))
	}
	slices.Sort(nums)

	n := len(nums)
	if n%2 == 0 {
		return (nums[n/2-1] + nums[n/2]) / 2
	}

	return nums[n/2]
}

// Mode 返回预声明有序类型中出现次数最多的元素。空集合返回 EmptyError。
// 次数相同时，保留先出现的那个。
func (c *Collection[T]) Mode() (T, error) {
	cmpFn, err := c.basicOrderedCmp()
	if err != nil {
		var zero T
		return zero, err
	}
	return c.modeWith(cmpFn)
}

// ModeBy 用 cmpFn 统计众数。cmpFn 返回 0 表示相等。空集合返回 EmptyError。
func (c *Collection[T]) ModeBy(cmpFn func(a, b T) int) (T, error) {
	if cmpFn == nil {
		var zero T
		return zero, errorx.NilFunc
	}
	return c.modeWith(cmpFn)
}

func (c *Collection[T]) modeWith(cmpFn func(a, b T) int) (T, error) {
	var zero T
	if c.IsEmpty() {
		return zero, errorx.EmptyError
	}

	type bucket struct {
		item  T
		count int
	}
	buckets := make([]bucket, 0, len(c.value))
	for _, item := range c.value {
		found := false
		for i := range buckets {
			if cmpFn(buckets[i].item, item) == 0 {
				buckets[i].count++
				found = true
				break
			}
		}
		if !found {
			buckets = append(buckets, bucket{item: item, count: 1})
		}
	}

	best := buckets[0]
	for _, b := range buckets[1:] {
		if b.count > best.count {
			best = b
		}
	}
	return best.item, nil
}

// Sum 对 fn 取出的数值求和。空集合返回 0。
func (c *Collection[T]) Sum[N number](fn func(T) N) float64 {
	var sum float64
	for _, item := range c.value {
		sum += toFloat(fn(item))
	}

	return sum
}

// Values 获取值
func (c *Collection[T]) Values() []T {
	return c.value
}

// ToJson 获取json
func (c *Collection[T]) ToJson() ([]byte, error) {
	return json.Marshal(c.value)
}

// FromJson 从json中获取数据
func (c *Collection[T]) FromJson(data []byte) error {
	return json.Unmarshal(data, &c.value)
}
