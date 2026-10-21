package map_collection

import (
	"testing"

	"github.com/goxstdx/go-collection/map_collection"
)

// TestMapResultType checks that Map's result type is decided by the function return value and that the original collection is unchanged.
func TestMapResultType(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1, "b": 2})

	texts := c.Map(func(v int, k string) string { return k })
	if texts.GetValue("a") != "a" || texts.GetValue("b") != "b" {
		t.Fatalf("Map() values = a:%q b:%q", texts.GetValue("a"), texts.GetValue("b"))
	}
	if c.GetValue("a") != 1 {
		t.Fatalf("Map changed the original value: %d", c.GetValue("a"))
	}
}

// TestReduceKeepsAccumulatorType checks that Reduce uses the caller's accumulator type directly.
func TestReduceKeepsAccumulatorType(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"a": 1, "b": 2, "c": 3})

	sum := c.Reduce(0, func(acc int, v int, _ string) int { return acc + v })
	if sum != 6 {
		t.Fatalf("Reduce() = %d, want 6", sum)
	}
}

// TestPluckFieldReturnsValues checks that PluckField takes every value. Order is not specified.
func TestPluckFieldReturnsValues(t *testing.T) {
	c := map_collection.NewCollection(map[string]int{"b": 2, "a": 1, "c": 3})

	got := c.PluckField(func(v int) int { return v })
	seen := map[int]int{}
	for _, n := range got {
		seen[n]++
	}
	if len(got) != 3 || seen[1] != 1 || seen[2] != 1 || seen[3] != 1 {
		t.Fatalf("PluckField() = %v", got)
	}
}
