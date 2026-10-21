package map_collection

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/goxstdx/go-collection/errorx"
	"github.com/goxstdx/go-collection/map_collection"
)

func TestNewCollection(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	if c == nil {
		t.Error("NewCollection returned nil")
	}
	if c.Count() != 3 {
		t.Errorf("Expected count 3, got %d", c.Count())
	}
}

func TestCopy(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	copied := c.Copy()

	if !reflect.DeepEqual(copied.All(), c.All()) {
		t.Error("Copy did not copy the values correctly")
	}

	// Changing a copy must not affect the original collection.
	copied.Set("d", 4)
	if c.Count() == copied.Count() {
		t.Error("Copy should create independent collection")
	}
}

func TestIsEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})
	if !c.IsEmpty() {
		t.Error("Expected collection to be empty")
	}

	c.Set("a", 1)
	if c.IsEmpty() {
		t.Error("Expected collection to be not empty")
	}
}

func TestIsNotEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1})
	if !c.IsNotEmpty() {
		t.Error("Expected collection to be not empty")
	}

	c2 := map_collection.NewCollection(map[string]int{})
	if c2.IsNotEmpty() {
		t.Error("Expected collection to be empty")
	}
}

func TestCount(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	if c.Count() != 3 {
		t.Errorf("Expected count 3, got %d", c.Count())
	}
}

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	keys := c.Keys()

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	expectedKeys := map[string]bool{"a": true, "b": true, "c": true}
	for _, k := range keys {
		if !expectedKeys[k] {
			t.Errorf("Unexpected key: %s", k)
		}
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	values := c.Values()

	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}

	expectedValues := map[int]bool{1: true, 2: true, 3: true}
	for _, v := range values {
		if !expectedValues[v] {
			t.Errorf("Unexpected value: %d", v)
		}
	}
}

func TestGetValue(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	if c.GetValue("a") != 1 {
		t.Error("GetValue returned wrong value")
	}

	// A missing key should return the zero value.
	if c.GetValue("d") != 0 {
		t.Error("GetValue should return zero value for non-existent key")
	}
}

func TestGet(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	v, ok := c.Get("a")
	if !ok || v != 1 {
		t.Error("Get returned wrong value")
	}

	_, ok = c.Get("d")
	if ok {
		t.Error("Get should return false for non-existent key")
	}
}

func TestGetOr(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	if c.GetOr("a", 99) != 1 {
		t.Error("GetOr returned wrong value")
	}

	if c.GetOr("d", 99) != 99 {
		t.Error("GetOr should return default value for non-existent key")
	}
}

func TestHas(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	if !c.Has("a") {
		t.Error("Has should return true for existing key")
	}

	if c.Has("d") {
		t.Error("Has should return false for non-existent key")
	}
}

func TestSet(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	c := map_collection.NewCollection(m)

	c.Set("c", 3)
	if c.Count() != 3 {
		t.Error("Set did not add new key")
	}
	if c.GetValue("c") != 3 {
		t.Error("Set did not set correct value")
	}

	// Update a key that already exists.
	c.Set("a", 10)
	if c.GetValue("a") != 10 {
		t.Error("Set did not update existing key")
	}
}

func TestPut(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	c := map_collection.NewCollection(m)

	newC := c.Put("c", 3)

	// The original collection must not change.
	if c.Count() != 2 {
		t.Error("Put should not modify original collection")
	}

	// The new collection should contain the new value.
	if newC.Count() != 3 || newC.GetValue("c") != 3 {
		t.Error("Put did not create correct new collection")
	}
}

func TestDelete(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	newC := c.Delete("b")

	// The original collection must not change.
	if c.Count() != 3 {
		t.Error("Delete should not modify original collection")
	}

	// The new collection should have the given key removed.
	if newC.Count() != 2 || newC.Has("b") {
		t.Error("Delete did not remove key correctly")
	}
}

func TestDeleteByFunc(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	c := map_collection.NewCollection(m)

	// Delete entries whose value is even.
	newC := c.DeleteByFunc(func(k string, v int) bool {
		return v%2 == 0
	})

	// The original collection must not change.
	if c.Count() != 4 {
		t.Error("DeleteByFunc should not modify original collection")
	}

	// The new collection should keep only odd values.
	if newC.Count() != 2 || newC.Has("b") || newC.Has("d") {
		t.Error("DeleteByFunc did not delete correctly")
	}
}

func TestRemove(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	c.Remove("b")

	if c.Count() != 2 || c.Has("b") {
		t.Error("Remove did not remove key correctly")
	}
}

func TestMerge(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	c := map_collection.NewCollection(m1)

	m2 := map[string]int{"b": 20, "c": 3}
	newC := c.Merge(m2)

	// The original collection must not change.
	if c.Count() != 2 {
		t.Error("Merge should not modify original collection")
	}

	// The new collection should merge both maps. On a conflict, other wins.
	if newC.Count() != 3 || newC.GetValue("b") != 20 || newC.GetValue("c") != 3 {
		t.Error("Merge did not merge correctly")
	}
}

func TestMergeCollection(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	c1 := map_collection.NewCollection(m1)

	m2 := map[string]int{"b": 20, "c": 3}
	c2 := map_collection.NewCollection(m2)

	newC := c1.MergeCollection(c2)

	// The original collection must not change.
	if c1.Count() != 2 {
		t.Error("MergeCollection should not modify original collection")
	}

	// The new collection should merge both Collections.
	if newC.Count() != 3 || newC.GetValue("b") != 20 || newC.GetValue("c") != 3 {
		t.Error("MergeCollection did not merge correctly")
	}

	// Test the nil case.
	newC2 := c1.MergeCollection(nil)
	if newC2.Count() != c1.Count() {
		t.Error("MergeCollection with nil should return copy")
	}
}

func TestMergeInPlace(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	c := map_collection.NewCollection(m1)

	m2 := map[string]int{"b": 20, "c": 3}
	c.MergeInPlace(m2)

	// The original collection should be modified.
	if c.Count() != 3 || c.GetValue("b") != 20 || c.GetValue("c") != 3 {
		t.Error("MergeInPlace did not merge correctly")
	}
}

func TestOnly(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	c := map_collection.NewCollection(m)

	newC := c.Only([]string{"a", "c"})

	// The original collection must not change.
	if c.Count() != 4 {
		t.Error("Only should not modify original collection")
	}

	// The new collection should contain only the given keys.
	if newC.Count() != 2 || !newC.Has("a") || !newC.Has("c") || newC.Has("b") {
		t.Error("Only did not filter correctly")
	}
}

func TestExcept(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	c := map_collection.NewCollection(m)

	newC := c.Except([]string{"b", "d"})

	// The original collection must not change.
	if c.Count() != 4 {
		t.Error("Except should not modify original collection")
	}

	// The new collection should exclude the given keys.
	if newC.Count() != 2 || newC.Has("b") || newC.Has("d") {
		t.Error("Except did not filter correctly")
	}
}

func TestFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	c := map_collection.NewCollection(m)

	// Keep even values.
	newC := c.Filter(func(v int, k string) bool {
		return v%2 == 0
	})

	// The original collection must not change.
	if c.Count() != 4 {
		t.Error("Filter should not modify original collection")
	}

	// The new collection should contain only even values.
	if newC.Count() != 2 || !newC.Has("b") || !newC.Has("d") {
		t.Error("Filter did not filter correctly")
	}
}

func TestEach(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	sum := 0
	c.Each(func(v int, k string) bool {
		sum += v
		return true
	})

	if sum != 6 {
		t.Errorf("Each did not iterate correctly, expected sum 6, got %d", sum)
	}
}

func TestEachSorted(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	sum := 0
	err := c.EachSorted(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}, func(v int, k string) {
		sum += v
	})

	if err != nil || sum != 6 {
		t.Errorf("EachSorted did not iterate correctly, expected sum 6, got %d, err %v", sum, err)
	}
}

func TestMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	// Double every value.
	newC := c.Map(func(v int, k string) int {
		return v * 2
	})

	// The original collection must not change.
	if c.GetValue("a") != 1 {
		t.Error("Map should not modify original collection")
	}

	// Values in the new collection should be doubled.
	if newC.GetValue("a") != 2 || newC.GetValue("b") != 4 || newC.GetValue("c") != 6 {
		t.Error("Map did not transform values correctly")
	}
}

func TestReduce(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	// Sum.
	sum := c.Reduce(0, func(acc any, v int, k string) any {
		return acc.(int) + v
	})

	if sum.(int) != 6 {
		t.Errorf("Reduce did not calculate correctly, expected 6, got %d", sum)
	}
}

func TestMinBy(t *testing.T) {
	m := map[string]int{"a": 1}
	c := map_collection.NewCollection(m)

	k, v, err := c.MinBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	if err != nil || k != "a" || v != 1 {
		t.Error("MinBy did not return correct value")
	}

	emptyC := map_collection.NewCollection(map[string]int{})
	_, _, err = emptyC.MinBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("MinBy should return EmptyError for empty collection")
	}
}

func TestFilterThenMinBy(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	cmpValue := func(_ string, va int, _ string, vb int) int {
		if va < vb {
			return -1
		}
		if va > vb {
			return 1
		}
		return 0
	}

	_, v, err := c.Filter(func(v int, k string) bool {
		return v > 1
	}).MinBy(cmpValue)
	if err != nil || v <= 1 {
		t.Error("Filter then MinBy did not find correct value")
	}

	_, _, err = c.Filter(func(v int, k string) bool {
		return v > 10
	}).MinBy(cmpValue)
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("Filter then MinBy should return EmptyError when no match found")
	}
}

func TestMaxBy(t *testing.T) {
	m := map[string]int{"a": 1}
	c := map_collection.NewCollection(m)

	k, v, err := c.MaxBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	if err != nil || k != "a" || v != 1 {
		t.Error("MaxBy did not return correct value")
	}

	emptyC := map_collection.NewCollection(map[string]int{})
	_, _, err = emptyC.MaxBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("MaxBy should return EmptyError for empty collection")
	}
}

func TestFilterThenMaxBy(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)
	cmpValue := func(_ string, va int, _ string, vb int) int {
		if va < vb {
			return -1
		}
		if va > vb {
			return 1
		}
		return 0
	}

	_, v, err := c.Filter(func(v int, k string) bool {
		return v < 3
	}).MaxBy(cmpValue)
	if err != nil || v >= 3 {
		t.Error("Filter then MaxBy did not find correct value")
	}

	_, _, err = c.Filter(func(v int, k string) bool {
		return v > 10
	}).MaxBy(cmpValue)
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("Filter then MaxBy should return EmptyError when no match found")
	}
}

func TestAll(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	all := c.All()
	if len(all) != 3 || all["a"] != 1 || all["b"] != 2 || all["c"] != 3 {
		t.Error("All did not return correct map")
	}
}

func TestToJSON(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := map_collection.NewCollection(m)

	jsonStr, err := c.ToJSON()
	if err != nil {
		t.Errorf("ToJSON returned error: %v", err)
	}

	// Check unmarshalling.
	var result map[string]int
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if !reflect.DeepEqual(result, m) {
		t.Error("ToJSON did not produce correct JSON")
	}
}

func TestDD(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	c := map_collection.NewCollection(m)

	// DD should return itself.
	result := c.DD()
	if result != c {
		t.Error("DD should return itself")
	}
}

type Person struct {
	Name string
	Age  int
}

func TestPluckFieldString(t *testing.T) {
	m := map[string]Person{
		"p1": {Name: "Alice", Age: 30},
		"p2": {Name: "Bob", Age: 25},
		"p3": {Name: "Charlie", Age: 35},
	}
	c := map_collection.NewCollection(m)

	names := c.PluckField(func(p Person) string { return p.Name })
	if len(names) != 3 {
		t.Errorf("PluckField returned wrong number of items: %d", len(names))
	}

	// Check that every name is present.
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	if !nameMap["Alice"] || !nameMap["Bob"] || !nameMap["Charlie"] {
		t.Error("PluckField did not extract all names correctly")
	}
}

func TestPluckFieldInt(t *testing.T) {
	m := map[string]Person{
		"p1": {Name: "Alice", Age: 30},
		"p2": {Name: "Bob", Age: 25},
		"p3": {Name: "Charlie", Age: 35},
	}
	c := map_collection.NewCollection(m)

	ages := c.PluckField(func(p Person) int { return p.Age })

	if len(ages) != 3 {
		t.Errorf("PluckField returned wrong number of items: %d", len(ages))
	}

	// Check that every age is present.
	ageMap := make(map[int]bool)
	for _, age := range ages {
		ageMap[age] = true
	}

	if !ageMap[30] || !ageMap[25] || !ageMap[35] {
		t.Error("PluckField did not extract all ages correctly")
	}
}

// ========== Edge-case tests ==========

// Test basic operations on an empty collection.
func TestEmptyCollectionOperations(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	// Test Count.
	if c.Count() != 0 {
		t.Errorf("Empty collection Count should be 0, got %d", c.Count())
	}

	// Test IsEmpty.
	if !c.IsEmpty() {
		t.Error("Empty collection should return true for IsEmpty")
	}

	// Test Keys.
	keys := c.Keys()
	if len(keys) != 0 {
		t.Errorf("Empty collection Keys should return empty slice, got %d keys", len(keys))
	}

	// Test Values.
	values := c.Values()
	if len(values) != 0 {
		t.Errorf("Empty collection Values should return empty slice, got %d values", len(values))
	}

	// Test All.
	all := c.All()
	if len(all) != 0 {
		t.Errorf("Empty collection All should return empty map, got %d items", len(all))
	}
}

// Test iteration of an empty collection.
func TestEmptyCollectionIteration(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	// Test Each: it should not iterate at all.
	count := 0
	c.Each(func(v int, k string) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("Each on empty collection should not iterate, but iterated %d times", count)
	}

	// Test EachSorted.
	count = 0
	err := c.EachSorted(func(a string, _ int, b string, _ int) int {
		return 0
	}, func(v int, k string) {
		count++
	})
	if err != nil || count != 0 {
		t.Errorf("EachSorted on empty collection should not iterate, but iterated %d times, err %v", count, err)
	}

	// Test Map: it should return an empty collection.
	mapped := c.Map(func(v int, k string) int {
		return v * 2
	})
	if !mapped.IsEmpty() {
		t.Error("Map on empty collection should return empty collection")
	}

	// Test Reduce: it should return the initial value.
	result := c.Reduce(100, func(acc any, v int, k string) any {
		return acc.(int) + v
	})
	if result.(int) != 100 {
		t.Errorf("Reduce on empty collection should return initial value 100, got %d", result)
	}
}

// Test filtering an empty collection.
func TestEmptyCollectionFiltering(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	// Test Filter.
	filtered := c.Filter(func(v int, k string) bool {
		return v > 0
	})
	if !filtered.IsEmpty() {
		t.Error("Filter on empty collection should return empty collection")
	}

	// Test Only.
	only := c.Only([]string{"a", "b"})
	if !only.IsEmpty() {
		t.Error("Only on empty collection should return empty collection")
	}

	// Test Except.
	except := c.Except([]string{"a", "b"})
	if !except.IsEmpty() {
		t.Error("Except on empty collection should return empty collection")
	}
}

// Test lookup on an empty collection.
func TestEmptyCollectionSearch(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	// Filtering an empty collection and then taking the minimum should return EmptyError.
	_, _, err := c.Filter(func(v int, k string) bool {
		return v > 0
	}).MinBy(func(_ string, va int, _ string, vb int) int {
		return va - vb
	})
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("MinBy on empty collection should return EmptyError")
	}

	_, _, err = c.Filter(func(v int, k string) bool {
		return v > 0
	}).MaxBy(func(_ string, va int, _ string, vb int) int {
		return va - vb
	})
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("MaxBy on empty collection should return EmptyError")
	}
}

// Test serialization of an empty collection.
func TestEmptyCollectionToJSON(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	jsonStr, err := c.ToJSON()
	if err != nil {
		t.Errorf("ToJSON on empty collection returned error: %v", err)
	}

	// It should be an empty object {}.
	var result map[string]int
	err = json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		t.Errorf("Failed to unmarshal empty collection JSON: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Empty collection JSON should produce empty map, got %d items", len(result))
	}
}

// Test copy and merge of an empty collection.
func TestEmptyCollectionCopyAndMerge(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	// Test Copy.
	copied := c.Copy()
	if !copied.IsEmpty() {
		t.Error("Copy of empty collection should be empty")
	}

	// Test Merge with an empty map.
	merged := c.Merge(map[string]int{})
	if !merged.IsEmpty() {
		t.Error("Merge empty with empty should return empty")
	}

	// Test Merge with a non-empty map.
	merged2 := c.Merge(map[string]int{"a": 1})
	if merged2.Count() != 1 {
		t.Error("Merge empty with non-empty should return non-empty")
	}

	// Test MergeCollection with nil.
	merged3 := c.MergeCollection(nil)
	if !merged3.IsEmpty() {
		t.Error("MergeCollection with nil should return empty")
	}
}

// Test deleting every element, the core case.
func TestDeleteAllElements(t *testing.T) {
	// Test Delete removing every element.
	c := map_collection.NewCollection(map[string]int{"a": 1, "b": 2})
	c = c.Delete("a")
	c = c.Delete("b")

	if !c.IsEmpty() {
		t.Error("After deleting all elements, collection should be empty")
	}
	if c.Count() != 0 {
		t.Errorf("After deleting all elements, Count should be 0, got %d", c.Count())
	}

	// Test Remove removing every element.
	c2 := map_collection.NewCollection(map[string]int{"a": 1, "b": 2})
	c2.Remove("a")
	c2.Remove("b")

	if !c2.IsEmpty() {
		t.Error("After removing all elements, collection should be empty")
	}
}

// Test DeleteByFunc removing every element.
func TestDeleteByFuncAllElements(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 2, "b": 4, "c": 6})

	// Delete every even value (which deletes every element).
	result := c.DeleteByFunc(func(k string, v int) bool {
		return v%2 == 0
	})

	if !result.IsEmpty() {
		t.Error("DeleteByFunc removing all elements should return empty collection")
	}
	if result.Count() != 0 {
		t.Errorf("DeleteByFunc removing all elements should have Count 0, got %d", result.Count())
	}

	// The original collection must not change.
	if c.Count() != 3 {
		t.Error("DeleteByFunc should not modify original collection")
	}
}

// Test Filter whose result is empty.
func TestFilterResultEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1, "b": 3, "c": 5})

	// Keep evens, but there are no evens.
	result := c.Filter(func(v int, k string) bool {
		return v%2 == 0
	})

	if !result.IsEmpty() {
		t.Error("Filter with no matches should return empty collection")
	}
	if result.Count() != 0 {
		t.Errorf("Filter with no matches should have Count 0, got %d", result.Count())
	}
}

// Test Only whose result is empty.
func TestOnlyResultEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1, "b": 2, "c": 3})

	// Keep keys that do not exist.
	result := c.Only([]string{"x", "y", "z"})

	if !result.IsEmpty() {
		t.Error("Only with non-existent keys should return empty collection")
	}
	if result.Count() != 0 {
		t.Errorf("Only with non-existent keys should have Count 0, got %d", result.Count())
	}

	// Test Only with an empty slice.
	result2 := c.Only([]string{})
	if !result2.IsEmpty() {
		t.Error("Only with empty keys slice should return empty collection")
	}
}

// Test Except whose result is empty.
func TestExceptResultEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1, "b": 2, "c": 3})

	// Exclude every key.
	result := c.Except([]string{"a", "b", "c"})

	if !result.IsEmpty() {
		t.Error("Except excluding all keys should return empty collection")
	}
	if result.Count() != 0 {
		t.Errorf("Except excluding all keys should have Count 0, got %d", result.Count())
	}

	// Test Except with an empty slice.
	result2 := c.Except([]string{})
	if result2.Count() != 3 {
		t.Error("Except with empty keys slice should return all elements")
	}
}

// Test operations on a one-element collection.
func TestSingleElementCollection(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1})

	// Test basic operations.
	if c.Count() != 1 {
		t.Error("Single element collection should have Count 1")
	}
	if c.IsEmpty() {
		t.Error("Single element collection should not be empty")
	}

	// Test Filter keeping the element.
	filtered := c.Filter(func(v int, k string) bool {
		return v == 1
	})
	if filtered.Count() != 1 {
		t.Error("Filter keeping single element should return collection with 1 element")
	}

	// Test Filter dropping the element.
	filtered2 := c.Filter(func(v int, k string) bool {
		return v == 2
	})
	if !filtered2.IsEmpty() {
		t.Error("Filter removing single element should return empty collection")
	}

	// Test Map.
	mapped := c.Map(func(v int, k string) int {
		return v * 2
	})
	if mapped.Count() != 1 {
		t.Error("Map on single element should return collection with 1 element")
	}
	if mapped.GetValue("a") != 2 {
		t.Error("Map should transform the single element")
	}

	// Test Reduce.
	result := c.Reduce(10, func(acc any, v int, k string) any {
		return acc.(int) + v
	})
	if result.(int) != 11 {
		t.Errorf("Reduce on single element should return 11, got %d", result)
	}

	// Test that deleting the only element leaves the collection empty.
	deleted := c.Delete("a")
	if !deleted.IsEmpty() {
		t.Error("Delete single element should return empty collection")
	}
}

// Test that a one-element collection is empty after delete.
func TestSingleElementDeleteToEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1})

	c.Remove("a")

	if !c.IsEmpty() {
		t.Error("After removing single element, collection should be empty")
	}
	if c.Count() != 0 {
		t.Errorf("After removing single element, Count should be 0, got %d", c.Count())
	}
}

// Test Pluck on an empty collection.
func TestPluckFieldStringOnEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]Person{})

	names := c.PluckField(func(p Person) string { return p.Name })
	if len(names) != 0 {
		t.Errorf("Pluck on empty collection should return empty slice, got %d items", len(names))
	}
}

// Test PluckField taking integers from an empty collection.
func TestPluckFieldIntOnEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]Person{})

	ages := c.PluckField(func(p Person) int { return p.Age })
	if len(ages) != 0 {
		t.Errorf("PluckField on empty collection should return empty slice, got %d items", len(ages))
	}
}
