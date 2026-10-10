package map_collection

import (
	"errors"
	"testing"

	"github.com/ZHOUXING1997/collection/errorx"
	"github.com/ZHOUXING1997/collection/map_collection"
)

func cmpStringKey(a string, _ int, b string, _ int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func cmpIntValue(_ string, va int, _ string, vb int) int {
	if va < vb {
		return -1
	}
	if va > vb {
		return 1
	}
	return 0
}

func TestMinByKey(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"c": 3, "a": 1, "b": 2})

	k, v, err := c.MinBy(cmpStringKey)
	if err != nil || k != "a" || v != 1 {
		t.Errorf("MinBy() = %s %d %v", k, v, err)
	}

	k, v, err = c.MaxBy(cmpStringKey)
	if err != nil || k != "c" || v != 3 {
		t.Errorf("MaxBy() = %s %d %v", k, v, err)
	}
}

func TestMinByValue(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 3, "b": 1, "c": 2})

	k, v, err := c.MinBy(cmpIntValue)
	if err != nil || k != "b" || v != 1 {
		t.Errorf("MinBy() = %s %d %v", k, v, err)
	}

	k, v, err = c.MaxBy(cmpIntValue)
	if err != nil || k != "a" || v != 3 {
		t.Errorf("MaxBy() = %s %d %v", k, v, err)
	}
}

func TestMinByNilCompare(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1})

	_, _, err := c.MinBy(nil)
	if !errors.Is(err, errorx.NilFunc) {
		t.Errorf("MinBy(nil) err = %v", err)
	}
}

func TestMinByOnEmpty(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{})

	_, _, err := c.MinBy(cmpStringKey)
	if !errors.Is(err, errorx.EmptyError) {
		t.Errorf("MinBy empty err = %v", err)
	}
	_, _, err = c.MaxBy(cmpStringKey)
	if !errors.Is(err, errorx.EmptyError) {
		t.Errorf("MaxBy empty err = %v", err)
	}
}

func TestFilterThenMinByValue(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 3, "b": 1, "c": 2})

	k, v, err := c.Filter(func(v int, _ string) bool { return v > 1 }).MinBy(cmpIntValue)
	if err != nil || k != "c" || v != 2 {
		t.Errorf("Filter MinBy() = %s %d %v", k, v, err)
	}

	k, v, err = c.Filter(func(v int, _ string) bool { return v < 3 }).MaxBy(cmpIntValue)
	if err != nil || k != "c" || v != 2 {
		t.Errorf("Filter MaxBy() = %s %d %v", k, v, err)
	}

	_, _, err = c.Filter(func(v int, _ string) bool { return v > 10 }).MinBy(cmpIntValue)
	if !errors.Is(err, errorx.EmptyError) {
		t.Errorf("Filter no match err = %v", err)
	}
}

func TestMinByStructField(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	c := map_collection.NewCollection(map[string]Person{
		"p1": {Name: "Alice", Age: 30},
		"p2": {Name: "Bob", Age: 25},
		"p3": {Name: "Charlie", Age: 35},
	})
	cmpAge := func(ak string, a Person, bk string, b Person) int {
		if a.Age != b.Age {
			if a.Age < b.Age {
				return -1
			}
			return 1
		}
		if ak < bk {
			return -1
		}
		if ak > bk {
			return 1
		}
		return 0
	}

	k, v, err := c.MinBy(cmpAge)
	if err != nil || k != "p2" || v.Age != 25 {
		t.Errorf("MinBy() = %s age %d %v", k, v.Age, err)
	}
	k, v, err = c.MaxBy(cmpAge)
	if err != nil || k != "p3" || v.Age != 35 {
		t.Errorf("MaxBy() = %s age %d %v", k, v.Age, err)
	}
}

func TestEachSortedSortsEachCall(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"c": 3, "a": 1, "b": 2})

	var keys []string
	err := c.EachSorted(cmpStringKey, func(_ int, k string) {
		keys = append(keys, k)
	})
	if err != nil {
		t.Fatalf("EachSorted() err = %v", err)
	}
	if len(keys) != 3 || keys[0] != "a" || keys[1] != "b" || keys[2] != "c" {
		t.Fatalf("EachSorted keys = %v", keys)
	}

	err = c.EachSorted(nil, func(int, string) {})
	if !errors.Is(err, errorx.NilFunc) {
		t.Errorf("EachSorted nil cmp err = %v", err)
	}
}

func TestEachSortedSkipsRemovedKey(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"c": 3, "a": 1, "b": 2})

	var keys []string
	var vals []int
	err := c.EachSorted(cmpStringKey, func(v int, k string) {
		if k == "a" {
			c.Remove("b")
			c.Set("c", 9)
			c.Set("d", 4)
		}
		keys = append(keys, k)
		vals = append(vals, v)
	})
	if err != nil {
		t.Fatalf("EachSorted() err = %v", err)
	}
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "c" || vals[0] != 1 || vals[1] != 9 {
		t.Fatalf("EachSorted keys=%v vals=%v", keys, vals)
	}
	if !c.Has("d") || c.GetValue("d") != 4 {
		t.Fatalf("added key d = %d", c.GetValue("d"))
	}
}
