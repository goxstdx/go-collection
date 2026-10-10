package map_collection

import (
	"reflect"
	"testing"

	"github.com/ZHOUXING1997/collection/map_collection"
)

func TestFuncKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := map_collection.Keys(m)

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

func TestFuncValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	values := map_collection.Values(m)

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

func TestFuncClone(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	cloned := map_collection.Clone(m)

	if !reflect.DeepEqual(m, cloned) {
		t.Error("Clone did not copy correctly")
	}

	// 修改克隆不应影响原map
	cloned["d"] = 4
	if len(m) == len(cloned) {
		t.Error("Clone should create independent map")
	}
}

func TestFuncEqual(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	m3 := map[string]int{"a": 1, "b": 2}

	if !map_collection.Equal(m1, m2) {
		t.Error("Equal should return true for equal maps")
	}

	if map_collection.Equal(m1, m3) {
		t.Error("Equal should return false for different maps")
	}
}

func TestFuncHas(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if !map_collection.Has(m, "a") {
		t.Error("Has should return true for existing key")
	}

	if map_collection.Has(m, "d") {
		t.Error("Has should return false for non-existent key")
	}
}

func TestFuncGet(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if map_collection.Get(m, "a") != 1 {
		t.Error("Get returned wrong value")
	}

	// 获取不存在的key应返回零值
	if map_collection.Get(m, "d") != 0 {
		t.Error("Get should return zero value for non-existent key")
	}
}

func TestFuncGetOr(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	if map_collection.GetOr(m, "a", 99) != 1 {
		t.Error("GetOr returned wrong value")
	}

	if map_collection.GetOr(m, "d", 99) != 99 {
		t.Error("GetOr should return default value for non-existent key")
	}
}

func TestFuncSet(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	newM := map_collection.Set(m, "c", 3)

	// 原map不应改变
	if len(m) != 2 {
		t.Error("Set should not modify original map")
	}

	// 新map应包含新值
	if len(newM) != 3 || newM["c"] != 3 {
		t.Error("Set did not create correct new map")
	}
}

func TestFuncSetInPlace(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	map_collection.SetInPlace(m, "c", 3)

	// 应修改原map
	if len(m) != 3 || m["c"] != 3 {
		t.Error("SetInPlace did not modify map correctly")
	}
}

func TestFuncDelete(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	newM := map_collection.Delete(m, "b")

	// 原map不应改变
	if len(m) != 3 {
		t.Error("Delete should not modify original map")
	}

	// 新map应删除指定key
	if len(newM) != 2 || newM["b"] != 0 {
		t.Error("Delete did not remove key correctly")
	}
}

func TestFuncDeleteInPlace(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	map_collection.DeleteInPlace(m, "b")

	// 应修改原map
	if len(m) != 2 || m["b"] != 0 {
		t.Error("DeleteInPlace did not remove key correctly")
	}
}

func TestFuncMerge(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 20, "c": 3}
	merged := map_collection.Merge(m1, m2)

	// 原map不应改变
	if len(m1) != 2 {
		t.Error("Merge should not modify original maps")
	}

	// 合并后的map应包含所有键，冲突时以m2为准
	if len(merged) != 3 || merged["b"] != 20 || merged["c"] != 3 {
		t.Error("Merge did not merge correctly")
	}
}

func TestFuncMergeInPlace(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 20, "c": 3}
	map_collection.MergeInPlace(m1, m2)

	// 应修改m1
	if len(m1) != 3 || m1["b"] != 20 || m1["c"] != 3 {
		t.Error("MergeInPlace did not merge correctly")
	}
}

func TestFuncOnly(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := map_collection.Only(m, []string{"a", "c"})

	// 原map不应改变
	if len(m) != 4 {
		t.Error("Only should not modify original map")
	}

	// 结果应只包含指定的key
	if len(result) != 2 || result["a"] != 1 || result["c"] != 3 {
		t.Error("Only did not filter correctly")
	}
}

func TestFuncExcept(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := map_collection.Except(m, []string{"b", "d"})

	// 原map不应改变
	if len(m) != 4 {
		t.Error("Except should not modify original map")
	}

	// 结果应排除指定的key
	if len(result) != 2 || result["a"] != 1 || result["c"] != 3 {
		t.Error("Except did not filter correctly")
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := map_collection.MapValues(m, func(v int, k string) int {
		return v * 2
	})

	// 原map不应改变
	if m["a"] != 1 {
		t.Error("MapValues should not modify original map")
	}

	// 结果的值应翻倍
	if result["a"] != 2 || result["b"] != 4 || result["c"] != 6 {
		t.Error("MapValues did not transform values correctly")
	}
}

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	result := map_collection.MapKeys(m, func(k string, v int) string {
		return k + "_new"
	})

	// 原map不应改变
	if len(m) != 3 {
		t.Error("MapKeys should not modify original map")
	}

	// 结果应有新的key
	if result["a_new"] != 1 || result["b_new"] != 2 || result["c_new"] != 3 {
		t.Error("MapKeys did not transform keys correctly")
	}
}

func TestFuncFilter(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	result := map_collection.Filter(m, func(v int, k string) bool {
		return v%2 == 0
	})

	// 原map不应改变
	if len(m) != 4 {
		t.Error("Filter should not modify original map")
	}

	// 结果应只包含偶数值
	if len(result) != 2 || result["b"] != 2 || result["d"] != 4 {
		t.Error("Filter did not filter correctly")
	}
}

func TestFuncEach(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	sum := 0

	map_collection.Each(m, func(v int, k string) bool {
		sum += v
		return true
	})

	if sum != 6 {
		t.Errorf("Each did not iterate correctly, expected sum 6, got %d", sum)
	}
}

func TestFuncReduce(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}

	sum := map_collection.Reduce(m, 0, func(acc int, v int, k string) int {
		return acc + v
	})

	if sum != 6 {
		t.Errorf("Reduce did not calculate correctly, expected 6, got %d", sum)
	}

	// 测试字符串拼接
	m2 := map[int]string{1: "a", 2: "b", 3: "c"}
	concat := map_collection.Reduce(m2, "", func(acc string, v string, k int) string {
		return acc + v
	})

	// 结果应包含所有字符（顺序不确定）
	if len(concat) != 3 {
		t.Errorf("Reduce string concat failed, expected length 3, got %d", len(concat))
	}
}

type TestPerson struct {
	Name string
	Age  int
}

func TestFuncPluck(t *testing.T) {
	m := map[string]TestPerson{
		"p1": {Name: "Alice", Age: 30},
		"p2": {Name: "Bob", Age: 25},
		"p3": {Name: "Charlie", Age: 35},
	}

	names := map_collection.Pluck(m, func(p TestPerson) string { return p.Name })

	if len(names) != 3 {
		t.Errorf("Pluck returned wrong number of items: %d", len(names))
	}

	// 验证是否包含所有名字
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	if !nameMap["Alice"] || !nameMap["Bob"] || !nameMap["Charlie"] {
		t.Error("Pluck did not extract all names correctly")
	}

	// 测试提取Age
	ages := map_collection.Pluck(m, func(p TestPerson) int { return p.Age })
	if len(ages) != 3 {
		t.Errorf("Pluck Age returned wrong number of items: %d", len(ages))
	}

	ageMap := make(map[int]bool)
	for _, age := range ages {
		ageMap[age] = true
	}

	if !ageMap[30] || !ageMap[25] || !ageMap[35] {
		t.Error("Pluck did not extract all ages correctly")
	}
}

func TestPluckWithPointer(t *testing.T) {
	m := map[string]*TestPerson{
		"p1": {Name: "Alice", Age: 30},
		"p2": {Name: "Bob", Age: 25},
		"p3": nil, // 测试nil指针
	}

	names := map_collection.Pluck(m, func(p *TestPerson) string {
		if p == nil {
			return ""
		}
		return p.Name
	})

	// 应该包含2个有效名字
	validNames := 0
	for _, name := range names {
		if name != "" {
			validNames++
		}
	}

	if validNames != 2 {
		t.Errorf("Pluck with pointer expected 2 valid names, got %d", validNames)
	}
}
