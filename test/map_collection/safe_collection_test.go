package map_collection

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/ZHOUXING1997/collection/errorx"
	"github.com/ZHOUXING1997/collection/map_collection"
)

func TestSafeCollectionBasicOperations(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{"a": 1, "b": 2})

	if sc.Count() != 2 {
		t.Errorf("Expected count 2, got %d", sc.Count())
	}

	sc.Set("c", 3)
	if sc.Count() != 3 {
		t.Error("Set did not add new key")
	}

	if !sc.Has("c") {
		t.Error("Has should return true for added key")
	}
}

func TestSafeCollectionConcurrentRead(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
	})

	var wg sync.WaitGroup
	errors := make(chan error, 50)

	// 10个goroutine并发读取
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				v := sc.GetValue("a")
				if v != 1 {
					errors <- fmt.Errorf("GetValue returned wrong value: %d", v)
					return
				}

				if !sc.Has("b") {
					errors <- fmt.Errorf("Has returned false for existing key")
					return
				}

				count := sc.Count()
				if count != 5 {
					errors <- fmt.Errorf("Count returned wrong value: %d", count)
					return
				}
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func TestSafeCollectionConcurrentWrite(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 100

	// 10个goroutine并发写入
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				key := fmt.Sprintf("k%d_%d", id, j)
				sc.Set(key, id*1000+j)
			}
		}(i)
	}

	wg.Wait()

	// 验证所有数据都被正确写入
	expectedCount := numGoroutines * itemsPerGoroutine
	if sc.Count() != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, sc.Count())
	}

	// 验证数据完整性
	for i := 0; i < numGoroutines; i++ {
		for j := 0; j < itemsPerGoroutine; j++ {
			key := fmt.Sprintf("k%d_%d", i, j)
			expected := i*1000 + j
			if v := sc.GetValue(key); v != expected {
				t.Errorf("Key %s: expected %d, got %d", key, expected, v)
			}
		}
	}
}

func TestSafeCollectionConcurrentReadWrite(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3,
	})

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// 5个goroutine并发读
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = sc.GetValue("a")
				_ = sc.Count()
				_ = sc.Has("b")
			}
		}()
	}

	// 5个goroutine并发写
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				key := fmt.Sprintf("w%d", id)
				sc.Set(key, id*100+j)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}

	// 验证原始数据未损坏
	if !sc.Has("a") || !sc.Has("b") || !sc.Has("c") {
		t.Error("Original keys were corrupted")
	}

	// 验证写入的数据存在
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("w%d", i)
		if !sc.Has(key) {
			t.Errorf("Written key %s not found", key)
		}
	}
}

func TestSafeCollectionImmutableOperations(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3,
	})

	var wg sync.WaitGroup

	// 并发使用不可变操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				// Put返回新的SafeCollection
				newSC := sc.Put(fmt.Sprintf("k%d", id), id*100)
				if newSC.Count() <= sc.Count() {
					t.Errorf("Put should return new collection with more items")
				}

				// Delete返回新的SafeCollection
				newSC2 := sc.Delete("a")
				if newSC2.Count() >= sc.Count() {
					t.Errorf("Delete should return new collection with fewer items")
				}

				// Filter返回新的SafeCollection
				newSC3 := sc.Filter(func(v int, k string) bool {
					return v > 1
				})
				if newSC3.Count() >= sc.Count() {
					t.Errorf("Filter should return filtered collection")
				}
			}
		}(i)
	}

	wg.Wait()

	// 验证原Collection未被修改
	if sc.Count() != 3 {
		t.Errorf("Original collection should not be modified, expected count 3, got %d", sc.Count())
	}
	if !sc.Has("a") || !sc.Has("b") || !sc.Has("c") {
		t.Error("Original collection keys were modified")
	}
}

func TestSafeCollectionChainedOperations(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup

	// 并发链式调用
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 链式调用Set和Remove
			sc.Set(fmt.Sprintf("a%d", id), id).
				Set(fmt.Sprintf("b%d", id), id*10).
				Set(fmt.Sprintf("c%d", id), id*100)
		}(i)
	}

	wg.Wait()

	// 验证所有数据都被写入
	if sc.Count() != 15 { // 5个goroutine × 3个key
		t.Errorf("Expected count 15, got %d", sc.Count())
	}
}

func TestSafeCollectionConcurrentMerge(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"base": 0,
	})

	var wg sync.WaitGroup

	// 并发Merge操作
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			other := map[string]int{
				fmt.Sprintf("m%d", id): id,
			}
			sc.MergeInPlace(other)
		}(i)
	}

	wg.Wait()

	// 验证所有merge的key都存在
	if !sc.Has("base") {
		t.Error("Base key was lost")
	}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("m%d", i)
		if !sc.Has(key) {
			t.Errorf("Merged key %s not found", key)
		}
	}
}

func TestSafeCollectionConcurrentCopy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3,
	})

	var wg sync.WaitGroup
	copies := make([]*map_collection.SafeCollection[string, int], 10)

	// 并发Copy
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			copies[id] = sc.Copy()
		}(i)
	}

	wg.Wait()

	// 验证所有副本都正确
	for i, copy := range copies {
		if copy.Count() != 3 {
			t.Errorf("Copy %d has wrong count: %d", i, copy.Count())
		}
		if !copy.Has("a") || !copy.Has("b") || !copy.Has("c") {
			t.Errorf("Copy %d is missing keys", i)
		}
	}

	// 修改副本不应影响原集合
	copies[0].Set("d", 4)
	if sc.Has("d") {
		t.Error("Modifying copy affected original collection")
	}
}

func TestSafeCollectionConcurrentMinBy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"c": 3, "a": 1, "b": 2,
	})

	var wg sync.WaitGroup
	errors := make(chan error, 20)
	cmpKey := func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			k, v, err := sc.MinBy(cmpKey)
			if err != nil || k != "a" || v != 1 {
				errors <- fmt.Errorf("MinBy = %s %d err=%v", k, v, err)
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// ========== 极限场景测试 ==========

// 测试空SafeCollection的并发操作
func TestEmptySafeCollectionConcurrentOps(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// 阶段1: 并发读取空集合
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !sc.IsEmpty() {
				errors <- fmt.Errorf("Empty collection should return true for IsEmpty")
			}
			if sc.Count() != 0 {
				errors <- fmt.Errorf("Empty collection Count should be 0, got %d", sc.Count())
			}
		}()
	}

	// 等待所有读取完成
	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}

	// 阶段2: 并发向空集合写入
	var wg2 sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg2.Add(1)
		go func(id int) {
			defer wg2.Done()
			key := fmt.Sprintf("key%d", id)
			sc.Set(key, id)
		}(i)
	}

	wg2.Wait()

	// 验证最终状态
	if sc.Count() != 20 {
		t.Errorf("After concurrent writes, expected count 20, got %d", sc.Count())
	}
}

// 测试并发删除所有元素
func TestSafeCollectionConcurrentDeleteAll(t *testing.T) {
	// 创建包含一些元素的集合
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
	})

	var wg sync.WaitGroup

	// 并发删除所有元素
	keys := []string{"a", "b", "c", "d", "e"}
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			sc.Remove(k)
		}(key)
	}

	wg.Wait()

	// 验证集合为空
	if !sc.IsEmpty() {
		t.Error("After deleting all elements concurrently, collection should be empty")
	}
	if sc.Count() != 0 {
		t.Errorf("After deleting all elements, Count should be 0, got %d", sc.Count())
	}
}

// 测试并发删除和添加,最终为空
func TestSafeCollectionConcurrentDeleteAndAddToEmpty(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2,
	})

	var wg sync.WaitGroup

	// 先删除所有现有元素
	wg.Add(1)
	go func() {
		defer wg.Done()
		sc.Remove("a")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sc.Remove("b")
	}()

	// 添加临时元素
	wg.Add(1)
	go func() {
		defer wg.Done()
		sc.Set("temp", 99)
	}()

	wg.Wait()

	// 删除临时元素
	sc.Remove("temp")

	// 验证最终为空
	if !sc.IsEmpty() {
		t.Error("After all operations, collection should be empty")
	}
}

// 测试空集合的并发 MinBy / MaxBy
func TestEmptySafeCollectionConcurrentMinByMaxBy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	errs := make(chan error, 40)
	cmpKey := func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := sc.MinBy(cmpKey)
			if !errors.Is(err, errorx.EmptyError) {
				errs <- fmt.Errorf("MinBy on empty collection err=%v", err)
			}
		}()
	}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := sc.MaxBy(cmpKey)
			if !errors.Is(err, errorx.EmptyError) {
				errs <- fmt.Errorf("MaxBy on empty collection err=%v", err)
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Error(err)
	}
}

// 测试空集合的并发迭代
func TestEmptySafeCollectionConcurrentEach(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	var iterCount int
	var mu sync.Mutex

	// 并发调用Each
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Each(func(v int, k string) bool {
				mu.Lock()
				iterCount++
				mu.Unlock()
				return true
			})
		}()
	}

	wg.Wait()

	if iterCount != 0 {
		t.Errorf("Each on empty collection should not iterate, but iterated %d times", iterCount)
	}
}

// 测试空集合的并发Filter
func TestEmptySafeCollectionConcurrentFilter(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup

	// 并发Filter空集合
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filtered := sc.Filter(func(v int, k string) bool {
				return v > 0
			})
			if !filtered.IsEmpty() {
				t.Error("Filter on empty collection should return empty collection")
			}
		}()
	}

	wg.Wait()
}

// 测试并发删除直到为空后再读取
func TestSafeCollectionDeleteToEmptyThenRead(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3,
	})

	var wg sync.WaitGroup

	// 并发删除所有元素
	keys := []string{"a", "b", "c"}
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			sc.Remove(k)
		}(key)
	}

	wg.Wait()

	// 并发读取空集合
	errors := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !sc.IsEmpty() {
				errors <- fmt.Errorf("Collection should be empty after deleting all elements")
			}
			keys := sc.Keys()
			if len(keys) != 0 {
				errors <- fmt.Errorf("Keys should return empty slice, got %d keys", len(keys))
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

// 测试单元素SafeCollection并发删除
func TestSingleElementSafeCollectionConcurrentDelete(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"only": 1,
	})

	var wg sync.WaitGroup

	// 多个goroutine尝试删除同一个元素
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Remove("only")
		}()
	}

	wg.Wait()

	// 验证集合为空
	if !sc.IsEmpty() {
		t.Error("After concurrent delete of single element, collection should be empty")
	}
	if sc.Count() != 0 {
		t.Errorf("After deleting single element, Count should be 0, got %d", sc.Count())
	}
}

// 测试并发Copy空集合
func TestEmptySafeCollectionConcurrentCopy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	copies := make([]*map_collection.SafeCollection[string, int], 10)

	// 并发Copy空集合
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			copies[id] = sc.Copy()
		}(i)
	}

	wg.Wait()

	// 验证所有副本都是空的
	for i, copy := range copies {
		if !copy.IsEmpty() {
			t.Errorf("Copy %d of empty collection should be empty", i)
		}
		if copy.Count() != 0 {
			t.Errorf("Copy %d should have count 0, got %d", i, copy.Count())
		}
	}
}

// 测试并发Merge空集合
func TestEmptySafeCollectionConcurrentMergeEmpty(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup

	// 并发Merge空map到空集合
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				// 一半merge空map
				sc.MergeInPlace(map[string]int{})
			} else {
				// 一半merge非空map
				sc.MergeInPlace(map[string]int{
					fmt.Sprintf("key%d", id): id,
				})
			}
		}(i)
	}

	wg.Wait()

	// 验证只有奇数id的key被添加
	if sc.IsEmpty() {
		t.Error("After merging some non-empty maps, collection should not be empty")
	}
}

// 测试并发Filter结果为空
func TestSafeCollectionConcurrentFilterToEmpty(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 3, "c": 5,
	})

	var wg sync.WaitGroup

	// 并发Filter,结果都为空(过滤偶数,但没有偶数)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			filtered := sc.Filter(func(v int, k string) bool {
				return v%2 == 0
			})
			if !filtered.IsEmpty() {
				t.Error("Filter with no matches should return empty collection")
			}
		}()
	}

	wg.Wait()

	// 原集合不应改变
	if sc.Count() != 3 {
		t.Error("Filter should not modify original collection")
	}
}

// TestSafeCollectionAllReturnsCopy 验证 All 返回的 map 在锁外修改不会写回集合。
func TestSafeCollectionAllReturnsCopy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{"a": 1})
	all := sc.All()
	all["a"] = 2
	if sc.GetValue("a") != 1 {
		t.Fatalf("All() leaked the internal map, value = %d", sc.GetValue("a"))
	}
}

// TestSafeCollectionReduceTyped 验证 Reduce 使用调用方的累加类型。
func TestSafeCollectionReduceTyped(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{"a": 1, "b": 2})
	sum := sc.Reduce(0, func(acc int, v int, _ string) int { return acc + v })
	if sum != 3 {
		t.Fatalf("Reduce() = %d, want 3", sum)
	}
}
