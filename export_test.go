package collection

import (
	"errors"
	"testing"

	"github.com/ZHOUXING1997/collection/errorx"
)

// TestNewMapCollect_Basic 测试基本功能
func TestNewMapCollect_Basic(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := NewMapCollect(m)

	if c == nil {
		t.Error("NewMapCollect returned nil")
	}
	if c.Count() != 3 {
		t.Errorf("Expected count 3, got %d", c.Count())
	}
}

func TestNewMapCollect_MinByKey(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	c := NewMapCollect(m)
	k, v, err := c.MinBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if err != nil || k != "a" || v != 1 {
		t.Errorf("Expected min key 'a' with value 1, got key '%s' with value %d, err %v", k, v, err)
	}
}

func TestNewMapCollect_MinByValue(t *testing.T) {
	m := map[string]int{"a": 3, "b": 1, "c": 2}
	c := NewMapCollect(m)
	k, v, err := c.MinBy(func(_ string, va int, _ string, vb int) int {
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	})
	if err != nil || v != 1 || k != "b" {
		t.Errorf("Expected min key 'b' with value 1, got key '%s' with value %d, err %v", k, v, err)
	}
}

func TestNewMapCollect_MinByKeyAndValue(t *testing.T) {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	c := NewMapCollect(m)
	k, _, err := c.MinBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if err != nil || k != "a" {
		t.Errorf("Expected min key 'a', got '%s', err %v", k, err)
	}

	k, v, err := c.MinBy(func(_ string, va int, _ string, vb int) int {
		if va < vb {
			return -1
		} else if va > vb {
			return 1
		}
		return 0
	})
	if err != nil || v != 1 || k != "a" {
		t.Errorf("Expected min key 'a' with value 1, got key '%s' with value %d, err %v", k, v, err)
	}
}

// TestNewMapCollect_EmptyMap 测试空map
func TestNewMapCollect_EmptyMap(t *testing.T) {
	m := map[string]int{}
	c := NewMapCollect(m)

	if c.Count() != 0 {
		t.Errorf("Expected count 0, got %d", c.Count())
	}

	_, _, err := c.MinBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	if !errors.Is(err, errorx.EmptyError) {
		t.Error("MinBy on empty map should return EmptyError")
	}
}

// TestNewMapCollect_NilMap 测试nil map
func TestNewMapCollect_NilMap(t *testing.T) {
	var m map[string]int
	c := NewMapCollect(m)

	if c == nil {
		t.Error("NewMapCollect with nil map should not return nil")
	}

	if c.Count() != 0 {
		t.Errorf("Expected count 0 for nil map, got %d", c.Count())
	}
}

// TestNewMapCollect_LargeMap 测试大量数据
func TestNewMapCollect_LargeMap(t *testing.T) {
	m := make(map[int]int, 10000)
	for i := 0; i < 10000; i++ {
		m[i] = i * 2
	}

	c := NewMapCollect(m)
	cmpKey := func(a int, _ int, b int, _ int) int {
		return a - b
	}

	if c.Count() != 10000 {
		t.Errorf("Expected count 10000, got %d", c.Count())
	}

	k, v, err := c.MinBy(cmpKey)
	if err != nil || k != 0 || v != 0 {
		t.Errorf("Expected min key 0 with value 0, got key %d with value %d, err %v", k, v, err)
	}

	k, v, err = c.MaxBy(cmpKey)
	if err != nil || k != 9999 || v != 19998 {
		t.Errorf("Expected max key 9999 with value 19998, got key %d with value %d, err %v", k, v, err)
	}
}

// TestNewMapCollect_WithoutCompare 测试创建时不传入比较函数
func TestNewMapCollect_WithoutCompare(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	c := NewMapCollect(m)

	// 创建时不排序
	if c.Count() != 2 {
		t.Errorf("Expected count 2, got %d", c.Count())
	}
}

func TestNewMapCollect_DescendingOrder(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	c := NewMapCollect(m)
	// 比较函数反过来时，MinBy 取到最大的 key。
	k, v, err := c.MinBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return 1
		} else if a > b {
			return -1
		}
		return 0
	})
	if err != nil || k != "c" || v != 3 {
		t.Errorf("Expected min-by reversed cmp key 'c' with value 3, got key '%s' with value %d, err %v", k, v, err)
	}

	k, v, err = c.MaxBy(func(a string, _ int, b string, _ int) int {
		if a < b {
			return 1
		} else if a > b {
			return -1
		}
		return 0
	})
	if err != nil || k != "a" || v != 1 {
		t.Errorf("Expected max-by reversed cmp key 'a' with value 1, got key '%s' with value %d, err %v", k, v, err)
	}
}

// TestNewMapCollect_DuplicateKeys 测试重复的key(map本身不允许,但测试数据完整性)
func TestNewMapCollect_DuplicateKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	m["a"] = 10 // 覆盖

	c := NewMapCollect(m)

	if c.GetValue("a") != 10 {
		t.Errorf("Expected value 10 for key 'a', got %d", c.GetValue("a"))
	}
	if c.Count() != 2 {
		t.Errorf("Expected count 2, got %d", c.Count())
	}
}

// TestNewEmptyMapCollection 测试空集合创建
func TestNewEmptyMapCollection(t *testing.T) {
	c := NewEmptyMapCollection[string, int]()

	if c == nil {
		t.Error("NewEmptyMapCollection returned nil")
	}
	if c.Count() != 0 {
		t.Errorf("Expected count 0, got %d", c.Count())
	}
	if !c.IsEmpty() {
		t.Error("Expected empty collection")
	}
}

// TestNewEmptyMapCollection_AddData 测试空集合添加数据
func TestNewEmptyMapCollection_AddData(t *testing.T) {
	c := NewEmptyMapCollection[string, int]()
	c.Set("a", 1).Set("b", 2)

	if c.Count() != 2 {
		t.Errorf("Expected count 2, got %d", c.Count())
	}
	if c.GetValue("a") != 1 {
		t.Errorf("Expected value 1 for key 'a', got %d", c.GetValue("a"))
	}
}

// TestNewMapCollect_StructKey 测试结构体作为key
func TestNewMapCollect_StructKey(t *testing.T) {
	type Key struct {
		ID int
	}
	m := map[Key]string{
		{ID: 1}: "one",
		{ID: 2}: "two",
	}

	c := NewMapCollect(m)

	if c.Count() != 2 {
		t.Errorf("Expected count 2, got %d", c.Count())
	}

	k, v, err := c.MinBy(func(a Key, _ string, b Key, _ string) int {
		return a.ID - b.ID
	})
	if err != nil || k.ID != 1 || v != "one" {
		t.Errorf("Expected min key {ID:1} with value 'one', got key {ID:%d} with value '%s', err %v", k.ID, v, err)
	}
}

// TestNewMapCollect_PointerValue 测试指针类型的value
func TestNewMapCollect_PointerValue(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	p1 := &Person{Name: "Alice", Age: 30}
	p2 := &Person{Name: "Bob", Age: 25}

	m := map[string]*Person{
		"p1": p1,
		"p2": p2,
	}

	c := NewMapCollect(m)

	if c.Count() != 2 {
		t.Errorf("Expected count 2, got %d", c.Count())
	}

	v := c.GetValue("p1")
	if v == nil || v.Name != "Alice" {
		t.Error("Failed to get pointer value correctly")
	}
}
