package map_collection

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/goxstdx/go-collection/errorx"
	"github.com/goxstdx/go-collection/map_collection"
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

	// Ten goroutines read concurrently.
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

	// Ten goroutines write concurrently.
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

	// Check that every value was written correctly.
	expectedCount := numGoroutines * itemsPerGoroutine
	if sc.Count() != expectedCount {
		t.Errorf("Expected count %d, got %d", expectedCount, sc.Count())
	}

	// Check data integrity.
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

	// Five goroutines read concurrently.
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

	// Five goroutines write concurrently.
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

	// Check that the original data is intact.
	if !sc.Has("a") || !sc.Has("b") || !sc.Has("c") {
		t.Error("Original keys were corrupted")
	}

	// Check that the written data is present.
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

	// Use immutable operations concurrently.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				// Put returns a new SafeCollection.
				newSC := sc.Put(fmt.Sprintf("k%d", id), id*100)
				if newSC.Count() <= sc.Count() {
					t.Errorf("Put should return new collection with more items")
				}

				// Delete returns a new SafeCollection.
				newSC2 := sc.Delete("a")
				if newSC2.Count() >= sc.Count() {
					t.Errorf("Delete should return new collection with fewer items")
				}

				// Filter returns a new SafeCollection.
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

	// Check that the original Collection was not modified.
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

	// Chain calls concurrently.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Chain Set and Remove.
			sc.Set(fmt.Sprintf("a%d", id), id).
				Set(fmt.Sprintf("b%d", id), id*10).
				Set(fmt.Sprintf("c%d", id), id*100)
		}(i)
	}

	wg.Wait()

	// Check that every value was written.
	if sc.Count() != 15 { // 5 goroutines × 3 keys
		t.Errorf("Expected count 15, got %d", sc.Count())
	}
}

func TestSafeCollectionConcurrentMerge(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"base": 0,
	})

	var wg sync.WaitGroup

	// Merge concurrently.
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

	// Check that every merged key exists.
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

	// Copy concurrently.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			copies[id] = sc.Copy()
		}(i)
	}

	wg.Wait()

	// Check that every copy is correct.
	for i, copy := range copies {
		if copy.Count() != 3 {
			t.Errorf("Copy %d has wrong count: %d", i, copy.Count())
		}
		if !copy.Has("a") || !copy.Has("b") || !copy.Has("c") {
			t.Errorf("Copy %d is missing keys", i)
		}
	}

	// Changing a copy must not affect the original collection.
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

// ========== Edge-case tests ==========

// Test concurrent operations on an empty SafeCollection.
func TestEmptySafeCollectionConcurrentOps(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// Phase 1: read the empty collection concurrently.
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

	// Wait for every read to finish.
	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}

	// Phase 2: write into the empty collection concurrently.
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

	// Check the final state.
	if sc.Count() != 20 {
		t.Errorf("After concurrent writes, expected count 20, got %d", sc.Count())
	}
}

// Test deleting every element concurrently.
func TestSafeCollectionConcurrentDeleteAll(t *testing.T) {
	// Create a collection that already has some elements.
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3, "d": 4, "e": 5,
	})

	var wg sync.WaitGroup

	// Delete every element concurrently.
	keys := []string{"a", "b", "c", "d", "e"}
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			sc.Remove(k)
		}(key)
	}

	wg.Wait()

	// Check that the collection is empty.
	if !sc.IsEmpty() {
		t.Error("After deleting all elements concurrently, collection should be empty")
	}
	if sc.Count() != 0 {
		t.Errorf("After deleting all elements, Count should be 0, got %d", sc.Count())
	}
}

// Test concurrent deletes and adds that finish empty.
func TestSafeCollectionConcurrentDeleteAndAddToEmpty(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2,
	})

	var wg sync.WaitGroup

	// Delete every existing element first.
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

	// Add a temporary element.
	wg.Add(1)
	go func() {
		defer wg.Done()
		sc.Set("temp", 99)
	}()

	wg.Wait()

	// Delete the temporary element.
	sc.Remove("temp")

	// Check that the collection ends empty.
	if !sc.IsEmpty() {
		t.Error("After all operations, collection should be empty")
	}
}

// Test concurrent MinBy and MaxBy on an empty collection.
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

// Test concurrent iteration of an empty collection.
func TestEmptySafeCollectionConcurrentEach(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	var iterCount int
	var mu sync.Mutex

	// Call Each concurrently.
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

// Test concurrent Filter on an empty collection.
func TestEmptySafeCollectionConcurrentFilter(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup

	// Filter the empty collection concurrently.
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

// Test reading after concurrent deletes have emptied the collection.
func TestSafeCollectionDeleteToEmptyThenRead(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 2, "c": 3,
	})

	var wg sync.WaitGroup

	// Delete every element concurrently.
	keys := []string{"a", "b", "c"}
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			sc.Remove(k)
		}(key)
	}

	wg.Wait()

	// Read the empty collection concurrently.
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

// Test concurrent delete on a one-element SafeCollection.
func TestSingleElementSafeCollectionConcurrentDelete(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"only": 1,
	})

	var wg sync.WaitGroup

	// Several goroutines try to delete the same element.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sc.Remove("only")
		}()
	}

	wg.Wait()

	// Check that the collection is empty.
	if !sc.IsEmpty() {
		t.Error("After concurrent delete of single element, collection should be empty")
	}
	if sc.Count() != 0 {
		t.Errorf("After deleting single element, Count should be 0, got %d", sc.Count())
	}
}

// Test concurrent Copy of an empty collection.
func TestEmptySafeCollectionConcurrentCopy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup
	copies := make([]*map_collection.SafeCollection[string, int], 10)

	// Copy the empty collection concurrently.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			copies[id] = sc.Copy()
		}(i)
	}

	wg.Wait()

	// Check that every copy is empty.
	for i, copy := range copies {
		if !copy.IsEmpty() {
			t.Errorf("Copy %d of empty collection should be empty", i)
		}
		if copy.Count() != 0 {
			t.Errorf("Copy %d should have count 0, got %d", i, copy.Count())
		}
	}
}

// Test concurrent Merge of an empty collection.
func TestEmptySafeCollectionConcurrentMergeEmpty(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{})

	var wg sync.WaitGroup

	// Merge an empty map into an empty collection concurrently.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				// Half of the calls merge an empty map.
				sc.MergeInPlace(map[string]int{})
			} else {
				// Half of the calls merge a non-empty map.
				sc.MergeInPlace(map[string]int{
					fmt.Sprintf("key%d", id): id,
				})
			}
		}(i)
	}

	wg.Wait()

	// Check that only keys with an odd id were added.
	if sc.IsEmpty() {
		t.Error("After merging some non-empty maps, collection should not be empty")
	}
}

// Test concurrent Filter whose result is empty.
func TestSafeCollectionConcurrentFilterToEmpty(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{
		"a": 1, "b": 3, "c": 5,
	})

	var wg sync.WaitGroup

	// Filter concurrently. Every result is empty (keep evens, but there are no evens).
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

	// The original collection must not change.
	if sc.Count() != 3 {
		t.Error("Filter should not modify original collection")
	}
}

// TestSafeCollectionAllReturnsCopy checks that changing the map returned by All outside the lock does not write back into the collection.
func TestSafeCollectionAllReturnsCopy(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{"a": 1})
	all := sc.All()
	all["a"] = 2
	if sc.GetValue("a") != 1 {
		t.Fatalf("All() leaked the internal map, value = %d", sc.GetValue("a"))
	}
}

// TestSafeCollectionReduceTyped checks that Reduce uses the caller's accumulator type.
func TestSafeCollectionReduceTyped(t *testing.T) {
	sc := map_collection.NewSafeCollection(map[string]int{"a": 1, "b": 2})
	sum := sc.Reduce(0, func(acc int, v int, _ string) int { return acc + v })
	if sum != 3 {
		t.Fatalf("Reduce() = %d, want 3", sum)
	}
}
