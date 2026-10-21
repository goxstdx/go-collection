// Package slice_collection provides generic, chainable slice collections with filtering,
// mapping, sorting, grouping, and aggregation.
//
// It covers common operations such as filtering, mapping, sorting, searching,
// aggregation, grouping, paging, and serialization, aimed at business code that
// does not need peak performance.
package slice_collection

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/goxstdx/go-collection/errorx"
)

// basicOrderedCmp returns a comparison function for predeclared ordered types.
// Structs and defined types return NoComparableError. The caller must use a By method.
// The function is not cached on the collection: read-only methods may run concurrently,
// and this is only a type switch.
func basicOrderedCmp[T any]() (func(T, T) int, error) {
	fn := defaultCompare[T]()
	if fn == nil {
		return nil, errorx.NoComparableError
	}
	return fn, nil
}

// // Err returns the Collection error.
// func (c *Collection[T]) Err() error {
// 	return c.err
// }
//
// // SetErr sets the Collection error.
// func (c *Collection[T]) SetErr(err error) *Collection[T] {
// 	c.err = err
// 	return c
// }

// Copy returns a new copy of the Collection.
func (c *Collection[T]) Copy() *Collection[T] {
	// Deep-copy the slice so the copy does not share the backing array.
	copied := make([]T, len(c.value))
	copy(copied, c.value)
	return NewCollection[T](copied)
}

// IsEmpty reports whether the collection is empty.
func (c *Collection[T]) IsEmpty() bool {
	return len(c.value) == 0
}

// IsNotEmpty reports whether the collection is not empty.
func (c *Collection[T]) IsNotEmpty() bool {
	return len(c.value) != 0
}

// Append adds an element at the end and mutates the current collection.
func (c *Collection[T]) Append(item T) *Collection[T] {
	c.value = append(c.value, item)
	return c
}

// Remove deletes the element at index and mutates the current collection.
// An out-of-range index leaves the collection unchanged.
func (c *Collection[T]) Remove(index int) *Collection[T] {
	// Bounds check.
	if index < 0 || index >= len(c.value) {
		return c
	}
	c.value = append(c.value[:index], c.value[index+1:]...)
	return c
}

// Insert inserts an element at index and mutates the current collection.
// An index equal to the length appends. An out-of-range index leaves the collection unchanged.
func (c *Collection[T]) Insert(index int, item T) *Collection[T] {
	if index < 0 || index > len(c.value) {
		return c
	}
	c.value = slices.Insert(c.value, index, item)
	return c
}

// Search looks up an element among predeclared ordered types and returns the index of the first match.
func (c *Collection[T]) Search(item T) (int, error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return -1, err
	}
	return c.searchWith(item, cmpFn)
}

// SearchBy looks up an element with cmpFn. cmpFn returns 0 when the values are equal.
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

// Unique removes duplicates from predeclared ordered types and keeps the first occurrence.
func (c *Collection[T]) Unique() (*Collection[T], error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return nil, err
	}
	return c.uniqueWith(cmpFn), nil
}

// UniqueBy removes duplicates with cmpFn. cmpFn returns 0 when the values are equal.
// The first occurrence is kept.
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

// Filter keeps elements that match.
func (c *Collection[T]) Filter(f func(item T, key int) bool) *Collection[T] {
	// Keep elements for which the function returns true.
	res := make([]T, 0, len(c.value))
	for i, v := range c.value {
		if f(v, i) {
			res = append(res, v)
		}
	}

	return NewCollection[T](res)
}

// Reject drops elements that match.
func (c *Collection[T]) Reject(f func(item T, key int) bool) *Collection[T] {
	res := make([]T, 0, len(c.value))
	for i, v := range c.value {
		if !f(v, i) {
			res = append(res, v)
		}
	}

	return NewCollection[T](res)
}

// First returns the first element.
func (c *Collection[T]) First() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}

	return c.value[0]
}

// Last returns the last element.
func (c *Collection[T]) Last() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}

	return c.value[len(c.value)-1]
}

// Slice returns a new collection for [start, end). Elements are copied and do not
// share the backing array with the original collection.
// When only start is passed, the slice runs through the end.
func (c *Collection[T]) Slice(params ...int) (*Collection[T], error) {
	if len(params) == 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid params")
	}
	start := params[0]
	if start < 0 || start >= len(c.value) {
		return NewCollection[T](nil), fmt.Errorf("invalid start index")
	}
	if len(params) == 1 {
		return NewCollection(cloneSlice(c.value[start:])), nil
	}
	end := params[1]
	if end < 0 || end > len(c.value) {
		return NewCollection[T](nil), fmt.Errorf("invalid end index")
	}
	if start > end {
		return NewCollection[T](nil), fmt.Errorf("start index should be less than end index")
	}

	return NewCollection(cloneSlice(c.value[start:end])), nil
}

// Index returns the element at an index.
func (c *Collection[T]) Index(i int) T {

	// Index lookup. Applies to every Collection.
	var zero T
	if i < 0 || i >= len(c.value) {
		return zero
	}

	return c.value[i]
}

// SetIndex sets the element at an index and mutates the current collection.
// An out-of-range index leaves the collection unchanged.
func (c *Collection[T]) SetIndex(i int, val T) *Collection[T] {
	// if c.err != nil {
	// 	return c
	// }

	// Set the slice element at this index.
	if i < 0 || i >= len(c.value) {
		return c
	}
	c.value[i] = val
	return c
}

// Count returns the number of elements.
func (c *Collection[T]) Count() int {
	// Length of the slice. Applies to every Collection.
	return len(c.value)
}

// Merge merges another collection.
func (c *Collection[T]) Merge(arr *Collection[T]) *Collection[T] {
	// Merge the two slices.
	if arr == nil {
		return c
	}

	res := c.Copy()
	for i := 0; i < arr.Count(); i++ {
		res.value = append(res.value, arr.Index(i))
	}

	return res
}

// Each walks the elements.
// Iteration stops when f returns false.
func (c *Collection[T]) Each(f func(item T, key int) bool) {
	for i, v := range c.value {
		if !f(v, i) {
			return
		}
	}
}

// Foreach walks every element. The callback decides internally whether to skip work.
func (c *Collection[T]) Foreach(f func(item T, key int)) {
	for i, v := range c.value {
		f(v, i)
	}
}

// Map maps elements. R is decided by f's return value and may differ from T.
func (c *Collection[T]) Map[R any](f func(item T, key int) R) *Collection[R] {
	res := make([]R, 0, len(c.value))
	for i, v := range c.value {
		res = append(res, f(v, i))
	}

	return NewCollection(res)
}

// MapFilter maps and filters. Elements for which f returns false are dropped.
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

// Reduce folds the elements into one value.
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

// Every reports whether every element matches the condition.
func (c *Collection[T]) Every(f func(item T, key int) bool) bool {
	for i, v := range c.value {
		if !f(v, i) {
			return false
		}
	}

	return true
}

// ForPage returns a copy of the elements on that page. It does not share the backing array with the original collection.
func (c *Collection[T]) ForPage(page int, perPage int) (*Collection[T], error) {
	if page <= 0 || perPage <= 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid page or perPage")
	}
	// Check the page against the collection length first, so a huge page or perPage cannot overflow the multiplication.
	if len(c.value) == 0 || page-1 > (len(c.value)-1)/perPage {
		return NewCollection[T](nil), nil
	}
	start := (page - 1) * perPage
	// Cap the step at the remaining element count before adding, so the end index cannot overflow.
	end := start + min(perPage, len(c.value)-start)

	return NewCollection(cloneSlice(c.value[start:end])), nil
}

// Nth takes every n-th element.
func (c *Collection[T]) Nth(n int, offset int) (*Collection[T], error) {
	if n <= 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid n")
	}
	if offset < 0 {
		return NewCollection[T](nil), fmt.Errorf("invalid offset")
	}
	if offset >= len(c.value) {
		return NewCollection(make([]T, 0)), nil
	}
	res := make([]T, 0, (len(c.value)-1-offset)/n+1)
	for i := offset; i < len(c.value); {
		res = append(res, c.value[i])
		// Stop when the next step would pass the end, so i+n cannot overflow into a negative index.
		if n >= len(c.value)-i {
			break
		}
		i += n
	}

	return NewCollection[T](res), nil
}

// Pad pads the collection up to count with def.
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

// Pop removes and returns the last element and mutates the current collection.
// An empty collection returns the zero value.
func (c *Collection[T]) Pop() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}
	res := c.value[len(c.value)-1]
	c.value = c.value[:len(c.value)-1]
	return res
}

// Push adds an element at the end and mutates the current collection.
func (c *Collection[T]) Push(item T) *Collection[T] {
	c.value = append(c.value, item)
	return c
}

// Prepend inserts an element at the front and mutates the current collection.
func (c *Collection[T]) Prepend(item T) *Collection[T] {
	c.value = slices.Insert(c.value, 0, item)
	return c
}

// Random returns one element at random.
func (c *Collection[T]) Random() T {
	var zero T
	if c.IsEmpty() {
		return zero
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return c.value[random.Intn(len(c.value))]
}

// Reverse reverses element order and mutates the current collection.
func (c *Collection[T]) Reverse() *Collection[T] {
	slices.Reverse(c.value)
	return c
}

// Shuffle shuffles element order and mutates the current collection.
func (c *Collection[T]) Shuffle() *Collection[T] {
	rand.Shuffle(
		len(c.value), func(i, j int) {
			c.value[i], c.value[j] = c.value[j], c.value[i]
		},
	)
	return c
}

// GroupBy groups elements by f's return value. K must be comparable.
func (c *Collection[T]) GroupBy[K comparable](f func(T, int) K) map[K]*Collection[T] {
	res := make(map[K]*Collection[T])
	for i, v := range c.value {
		key := f(v, i)
		group, ok := res[key]
		if !ok {
			group = NewCollection[T](nil)
			res[key] = group
		}
		// A key containing NaN cannot be looked up again, so the group just created is updated through the local pointer.
		group.Push(v)
	}

	return res
}

// Split groups elements by size. Each group is a copy and does not share the backing array with the original collection.
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
		res = append(res, NewCollection(cloneSlice(c.value[i:end])))
	}

	return res
}

// DD prints the current slice structure.
func (c *Collection[T]) DD() {
	var zero T
	ret := fmt.Sprintf("Collection(%d, %T):{\n", c.Count(), zero)
	for k, v := range c.value {
		ret = ret + fmt.Sprintf("\t%d:\t%v\n", k, v)
	}
	ret = ret + "}\n"
	fmt.Print(ret)
}

// PluckField uses fn to take one value from each element and returns a new collection.
func (c *Collection[T]) PluckField[R any](fn func(T) R) *Collection[R] {
	res := make([]R, 0, len(c.value))
	for _, v := range c.value {
		res = append(res, fn(v))
	}

	return NewCollection(res)
}

// SortBy sorts ascending by the value extracted by fn and mutates the current collection.
// The sort is stable: elements that compare equal keep their previous relative order,
// so a secondary field can be sorted before the primary field.
func (c *Collection[T]) SortBy[R cmp.Ordered](fn func(T) R) *Collection[T] {
	slices.SortStableFunc(
		c.value, func(a, b T) int {
			return cmp.Compare(fn(a), fn(b))
		},
	)

	return c
}

// SortByFunc sorts with a custom comparison function and mutates the current collection.
// fn returns true when v1 should come before v2.
// The sort is stable: equal elements keep their previous relative order.
func (c *Collection[T]) SortByFunc(fn func(v1, v2 T) bool) (*Collection[T], error) {
	if fn == nil {
		return nil, errorx.NilFunc
	}

	sort.SliceStable(
		c.value, func(i, j int) bool {
			return fn(c.value[i], c.value[j])
		},
	)

	return c, nil
}

// SortByDesc sorts descending by the value extracted by fn and mutates the current collection.
// The sort is stable: elements that compare equal keep their previous relative order.
func (c *Collection[T]) SortByDesc[R cmp.Ordered](fn func(T) R) *Collection[T] {
	slices.SortStableFunc(
		c.value, func(a, b T) int {
			return cmp.Compare(fn(b), fn(a))
		},
	)

	return c
}

// KeyBy uses fn's return value as the key. When keys collide, the later element overwrites the earlier one.
func (c *Collection[T]) KeyBy[K comparable](fn func(T) K) map[K]T {
	res := make(map[K]T, len(c.value))
	for _, v := range c.value {
		res[fn(v)] = v
	}

	return res
}

// Max returns the maximum element among predeclared ordered types.
// An empty collection returns EmptyError. The zero value is not treated as a result.
func (c *Collection[T]) Max() (T, error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		var zero T
		return zero, err
	}
	return c.maxWith(cmpFn)
}

// MaxBy returns the maximum element according to cmpFn.
// cmpFn returns a positive value when a is greater than b.
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

// Min returns the minimum element among predeclared ordered types.
// An empty collection returns EmptyError. The zero value is not treated as a result.
func (c *Collection[T]) Min() (T, error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		var zero T
		return zero, err
	}
	return c.minWith(cmpFn)
}

// MinBy returns the minimum element according to cmpFn.
// cmpFn returns a negative value when a is less than b.
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

// Contains reports whether a predeclared ordered type contains the element.
func (c *Collection[T]) Contains(obj T) (bool, error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return false, err
	}
	return c.containsWith(obj, cmpFn), nil
}

// ContainsBy reports whether the element is present according to cmpFn.
// cmpFn returns 0 when the values are equal.
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

// ContainsCount returns how many elements of a predeclared ordered type equal obj.
func (c *Collection[T]) ContainsCount(obj T) (int, error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return 0, err
	}
	return c.containsCountWith(obj, cmpFn), nil
}

// ContainsCountBy counts equal elements with cmpFn. cmpFn returns 0 when the values are equal.
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

// Diff returns elements that are in the current collection but not in arr.
// It applies only to predeclared ordered types.
func (c *Collection[T]) Diff(arr *Collection[T]) (*Collection[T], error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return nil, err
	}
	return c.diffWith(arr, cmpFn), nil
}

// DiffBy computes the difference with cmpFn. cmpFn returns 0 when the values are equal.
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

// Sort sorts elements ascending by themselves. It applies only to predeclared ordered types
// and mutates the current collection.
// The sort is stable: equal elements keep their previous relative order.
func (c *Collection[T]) Sort() (*Collection[T], error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(c.value, cmpFn)
	return c, nil
}

// SortDesc sorts elements descending by themselves. It applies only to predeclared ordered types
// and mutates the current collection.
// The sort is stable: equal elements keep their previous relative order.
func (c *Collection[T]) SortDesc() (*Collection[T], error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return nil, err
	}

	slices.SortStableFunc(
		c.value, func(a, b T) int {
			return cmpFn(b, a)
		},
	)
	return c, nil
}

// Join joins elements with split. format customizes the text of each element.
func (c *Collection[T]) Join(split string, format ...func(item T) string) string {
	if c.IsEmpty() {
		return ""
	}

	// strings.Builder avoids repeated string allocation.
	var builder strings.Builder
	// Estimate capacity to reduce allocations.
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

// Union returns the union of two collections of a predeclared ordered type.
// Duplicates already in the current collection are kept, then elements from arr that are absent are appended.
func (c *Collection[T]) Union(arr *Collection[T]) (*Collection[T], error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return nil, err
	}
	return c.unionWith(arr, cmpFn), nil
}

// UnionBy computes the union with cmpFn. cmpFn returns 0 when the values are equal.
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

// Intersect returns the intersection of two collections of a predeclared ordered type,
// keeping the current collection's order and duplicates.
func (c *Collection[T]) Intersect(arr *Collection[T]) (*Collection[T], error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		return nil, err
	}
	return c.intersectWith(arr, cmpFn), nil
}

// IntersectBy computes the intersection with cmpFn. cmpFn returns 0 when the values are equal.
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

// Avg returns the average of the numbers taken by fn. An empty collection returns 0.
func (c *Collection[T]) Avg[N number](fn func(T) N) float64 {
	if c.IsEmpty() {
		return 0
	}

	var sum, mean float64
	allFinite := true
	for i, item := range c.value {
		value := float64(fn(item))
		sum += value
		if math.IsInf(value, 0) || math.IsNaN(value) {
			allFinite = false
		}
		if allFinite {
			count := float64(i + 1)
			delta := value - mean
			// Subtracting large values of opposite sign can also overflow, so share the value first, then update the mean.
			if math.IsInf(delta, 0) {
				mean = mean - mean/count + value/count
			} else {
				mean += delta / count
			}
		}
	}
	// Use the running mean only when a sum of finite inputs overflows, keeping the original
	// precision and Inf/NaN behavior.
	// Computing both at once avoids calling a side-effecting fn twice.
	if allFinite && math.IsInf(sum, 0) {
		return mean
	}
	return sum / float64(len(c.value))
}

// Median returns the median of the numbers taken by fn.
// With an even count, it returns the average of the two middle values. An empty collection returns 0.
// The original collection is not modified.
func (c *Collection[T]) Median[N number](fn func(T) N) float64 {
	if c.IsEmpty() {
		return 0
	}

	nums := make([]float64, len(c.value))
	for i, item := range c.value {
		nums[i] = float64(fn(item))
	}
	slices.Sort(nums)

	n := len(nums)
	if n%2 == 0 {
		a, b := nums[n/2-1], nums[n/2]
		sum := a + b
		// When adding finite large numbers overflows, halve first. Otherwise keep the original
		// calculation so halving a small number does not underflow.
		if math.IsInf(sum, 0) && !math.IsInf(a, 0) && !math.IsInf(b, 0) {
			return a/2 + b/2
		}
		return sum / 2
	}

	return nums[n/2]
}

// Mode returns the most frequent element among predeclared ordered types.
// An empty collection returns EmptyError.
// When frequencies tie, the element that appeared first is kept.
func (c *Collection[T]) Mode() (T, error) {
	cmpFn, err := basicOrderedCmp[T]()
	if err != nil {
		var zero T
		return zero, err
	}
	return c.modeWith(cmpFn)
}

// ModeBy finds the mode with cmpFn. cmpFn returns 0 when the values are equal.
// An empty collection returns EmptyError.
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

// Sum returns the sum of the numbers taken by fn. An empty collection returns 0.
func (c *Collection[T]) Sum[N number](fn func(T) N) float64 {
	var sum float64
	for _, item := range c.value {
		sum += float64(fn(item))
	}

	return sum
}

// Values returns the internal slice. This is a reference. Modifying it affects the current
// collection and the slice passed to NewCollection.
func (c *Collection[T]) Values() []T {
	return c.value
}

// ToJson returns the JSON encoding.
func (c *Collection[T]) ToJson() ([]byte, error) {
	return json.Marshal(c.value)
}

// FromJson replaces the current elements with a JSON array and mutates the current collection.
func (c *Collection[T]) FromJson(data []byte) error {
	return json.Unmarshal(data, &c.value)
}
