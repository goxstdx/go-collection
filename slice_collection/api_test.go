package slice_collection

import (
	"cmp"
	"errors"
	"testing"

	"github.com/goxstdx/go-collection/errorx"
)

// TestBasicTypesCompareAtCall 验证预声明有序类型在调用时生成比较函数。
func TestBasicTypesCompareAtCall(t *testing.T) {
	t.Run(
		"int", func(t *testing.T) {
			assertBasicOrder(t, []int{3, 1, 2}, 1, 3)
		},
	)
	t.Run(
		"string", func(t *testing.T) {
			assertBasicOrder(t, []string{"b", "a", "c"}, "a", "c")
		},
	)
	t.Run(
		"float64", func(t *testing.T) {
			assertBasicOrder(t, []float64{1.5, 0.5, 2}, 0.5, 2)
		},
	)
	t.Run(
		"uint8", func(t *testing.T) {
			assertBasicOrder(t, []uint8{3, 1, 2}, 1, 3)
		},
	)
}

func assertBasicOrder[T cmp.Ordered](t *testing.T, in []T, wantMin, wantMax T) {
	t.Helper()
	c := NewCollection(in)

	gotMin, err := c.Min()
	if err != nil || gotMin != wantMin {
		t.Fatalf("Min() = %v, %v; want %v, nil", gotMin, err, wantMin)
	}
	gotMax, err := c.Max()
	if err != nil || gotMax != wantMax {
		t.Fatalf("Max() = %v, %v; want %v, nil", gotMax, err, wantMax)
	}

	ok, err := c.Contains(wantMin)
	if err != nil || !ok {
		t.Fatalf("Contains(%v) = %v, %v", wantMin, ok, err)
	}

	sorted, err := c.Copy().Sort()
	if err != nil {
		t.Fatalf("Sort() error: %v", err)
	}
	if sorted.Index(0) != wantMin || sorted.Index(sorted.Count()-1) != wantMax {
		t.Fatalf("Sort() = %v", sorted.Values())
	}
	// Sort 只修改副本，原集合顺序不变。
	if c.Index(0) != in[0] {
		t.Fatalf("Sort changed the original collection: %v", c.Values())
	}
}

// TestDefinedTypeRequiresCompareFunc 验证定义类型不能走无参比较。
func TestDefinedTypeRequiresCompareFunc(t *testing.T) {
	type score int
	c := NewCollection([]score{3, 1, 2})

	if _, err := c.Max(); !errors.Is(err, errorx.NoComparableError) {
		t.Fatalf("Max() err = %v, want NoComparableError", err)
	}

	got, err := c.MaxBy(func(a, b score) int { return cmp.Compare(a, b) })
	if err != nil || got != 3 {
		t.Fatalf("MaxBy() = %v, %v; want 3, nil", got, err)
	}
}

type record struct {
	ID   int
	Name string
	Age  int
}

func byAge(a, b record) int { return cmp.Compare(a.Age, b.Age) }

func sampleRecords() *Collection[record] {
	return NewCollection(
		[]record{
			{ID: 1, Name: "a", Age: 10},
			{ID: 2, Name: "b", Age: 8},
			{ID: 3, Name: "c", Age: 8},
			{ID: 4, Name: "d", Age: 10},
		},
	)
}

// TestStructRequiresByMethod 验证结构体无参比较返回错误，By 方法按传入函数比较。
func TestStructRequiresByMethod(t *testing.T) {
	c := sampleRecords()

	if _, err := c.Min(); !errors.Is(err, errorx.NoComparableError) {
		t.Fatalf("Min() err = %v, want NoComparableError", err)
	}
	if _, err := c.Sort(); !errors.Is(err, errorx.NoComparableError) {
		t.Fatalf("Sort() err = %v, want NoComparableError", err)
	}
	if got, err := c.Unique(); got != nil || !errors.Is(err, errorx.NoComparableError) {
		t.Fatalf("Unique() = %v, %v; want nil, NoComparableError", got, err)
	}

	youngest, err := c.MinBy(byAge)
	if err != nil || youngest.ID != 2 {
		t.Fatalf("MinBy() = %+v, %v; want id 2", youngest, err)
	}

	ok, err := c.ContainsBy(
		record{Name: "c"}, func(a, b record) int {
			return cmp.Compare(a.Name, b.Name)
		},
	)
	if err != nil || !ok {
		t.Fatalf("ContainsBy() = %v, %v", ok, err)
	}
}

// TestEmptyMinMaxMode 验证空集合不把零值当作成功结果。
func TestEmptyMinMaxMode(t *testing.T) {
	c := NewCollection([]int{})

	if _, err := c.Min(); !errors.Is(err, errorx.EmptyError) {
		t.Fatalf("Min() err = %v, want EmptyError", err)
	}
	if _, err := c.Max(); !errors.Is(err, errorx.EmptyError) {
		t.Fatalf("Max() err = %v, want EmptyError", err)
	}
	if _, err := c.Mode(); !errors.Is(err, errorx.EmptyError) {
		t.Fatalf("Mode() err = %v, want EmptyError", err)
	}

	if _, err := NewCollection([]record{}).MinBy(byAge); !errors.Is(err, errorx.EmptyError) {
		t.Fatalf("MinBy() err = %v, want EmptyError", err)
	}
}

// TestNilCompareFunc 验证比较函数为 nil 时返回 NilFunc。
func TestNilCompareFunc(t *testing.T) {
	c := sampleRecords()

	if _, err := c.MaxBy(nil); !errors.Is(err, errorx.NilFunc) {
		t.Fatalf("MaxBy(nil) err = %v", err)
	}
	if _, err := c.ContainsBy(record{}, nil); !errors.Is(err, errorx.NilFunc) {
		t.Fatalf("ContainsBy(nil) err = %v", err)
	}
	if got, err := c.UniqueBy(nil); got != nil || !errors.Is(err, errorx.NilFunc) {
		t.Fatalf("UniqueBy(nil) = %v, %v", got, err)
	}
}

// TestStableSortKeepsEarlierOrder 验证连续稳定排序时，后一次是主键。
// 数据：id 1 a 10、id 2 b 8、id 3 c 8、id 4 d 10。
// 先按 id 倒序，再按 age 正序，相同 age 应保留 id 倒序后的相对顺序。
func TestStableSortKeepsEarlierOrder(t *testing.T) {
	c := sampleRecords()
	c.SortByDesc(func(r record) int { return r.ID })
	c.SortBy(func(r record) int { return r.Age })

	got := c.Values()
	want := []int{3, 2, 4, 1}
	for i, id := range want {
		if got[i].ID != id {
			t.Fatalf("after stable sorts, ids = %v, want %v", []int{got[0].ID, got[1].ID, got[2].ID, got[3].ID}, want)
		}
	}
}

// TestModeKeepsFirstOnTie 验证出现次数相同时保留先出现的元素。
func TestModeKeepsFirstOnTie(t *testing.T) {
	got, err := NewCollection([]int{1, 2, 2, 1}).Mode()
	if err != nil || got != 1 {
		t.Fatalf("Mode() = %v, %v; want 1, nil", got, err)
	}
}

// TestStructSetOps 验证结构体差集、并集、查找和去重使用传入的相等判断。
func TestStructSetOps(t *testing.T) {
	left := sampleRecords()
	right := NewCollection([]record{{ID: 2, Name: "b", Age: 8}, {ID: 9, Name: "z", Age: 1}})
	sameID := func(a, b record) int { return cmp.Compare(a.ID, b.ID) }

	diff, err := left.DiffBy(right, sameID)
	if err != nil || diff.Count() != 3 || diff.Index(0).ID != 1 {
		t.Fatalf("DiffBy() = %v, %v", diff.Values(), err)
	}

	union, err := left.UnionBy(right, sameID)
	if err != nil || union.Count() != 5 || union.Index(4).ID != 9 {
		t.Fatalf("UnionBy() = %v, %v", union.Values(), err)
	}

	idx, err := left.SearchBy(record{ID: 3}, sameID)
	if err != nil || idx != 2 {
		t.Fatalf("SearchBy() = %d, %v", idx, err)
	}

	unique, err := NewCollection([]record{{ID: 1, Age: 10}, {ID: 1, Age: 8}}).UniqueBy(sameID)
	if err != nil || unique.Count() != 1 || unique.Index(0).Age != 10 {
		t.Fatalf("UniqueBy() = %v, %v", unique.Values(), err)
	}
}

// TestMapAndAggregateUseCallerFunc 验证映射结果类型和数值聚合都由传入函数决定。
func TestMapAndAggregateUseCallerFunc(t *testing.T) {
	c := sampleRecords()

	names := c.Map(func(r record, _ int) string { return r.Name })
	if names.Count() != 4 || names.Index(0) != "a" || names.Index(3) != "d" {
		t.Fatalf("Map() = %v", names.Values())
	}

	ages := c.MapFilter(
		func(r record, _ int) (int, bool) {
			return r.Age, r.Age < 10
		},
	)
	if ages.Count() != 2 || ages.Index(0) != 8 || ages.Index(1) != 8 {
		t.Fatalf("MapFilter() = %v", ages.Values())
	}

	if sum := c.Sum(func(r record) int { return r.Age }); sum != 36 {
		t.Fatalf("Sum() = %v, want 36", sum)
	}
	if avg := c.Avg(func(r record) int { return r.Age }); avg != 9 {
		t.Fatalf("Avg() = %v, want 9", avg)
	}

	nums := NewCollection([]int{3, 1, 2})
	if median := nums.Median(func(n int) int { return n }); median != 2 {
		t.Fatalf("Median() = %v, want 2", median)
	}
	if nums.Index(0) != 3 || nums.Index(1) != 1 || nums.Index(2) != 2 {
		t.Fatalf("Median changed the original collection: %v", nums.Values())
	}
	if NewCollection([]int{}).Sum(func(n int) int { return n }) != 0 {
		t.Fatal("Sum() on empty collection should be 0")
	}
}

// TestDefinedNumericTypeCanAggregate 验证定义数值类型可以直接求和。
func TestDefinedNumericTypeCanAggregate(t *testing.T) {
	type money int64
	got := NewCollection([]money{1, 2, 3}).Sum(func(m money) money { return m })
	if got != 6 {
		t.Fatalf("Sum() = %v, want 6", got)
	}
}

// TestConcurrentReadDoesNotRace 验证只读比较不会写集合字段。
func TestConcurrentReadDoesNotRace(t *testing.T) {
	c := NewCollection([]int{1, 2, 3})
	done := make(chan struct{})
	for i := 0; i < 8; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				if _, err := c.Max(); err != nil {
					t.Errorf("Max() error: %v", err)
				}
				if _, err := c.Contains(2); err != nil {
					t.Errorf("Contains() error: %v", err)
				}
			}
			done <- struct{}{}
		}()
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}

// TestViewMethodsReturnCopies 验证分页和切分的结果不和原集合共用底层数组。
func TestViewMethodsReturnCopies(t *testing.T) {
	c := NewCollection([]int{1, 2, 3, 4})

	part, err := c.Slice(0, 2)
	if err != nil {
		t.Fatalf("Slice() error: %v", err)
	}
	part.Push(9)
	if c.Index(2) != 3 {
		t.Fatalf("Slice result changed the original: %v", c.Values())
	}

	page, err := c.ForPage(1, 2)
	if err != nil {
		t.Fatalf("ForPage() error: %v", err)
	}
	page.SetIndex(0, 8)
	if c.Index(0) != 1 {
		t.Fatalf("ForPage result changed the original: %v", c.Values())
	}

	parts := c.Split(2)
	parts[0].SetIndex(0, 7)
	if c.Index(0) != 1 {
		t.Fatalf("Split result changed the original: %v", c.Values())
	}
}

// TestContainerAndReorderMutate 验证插入、前置、反转、打乱都改当前集合。
func TestContainerAndReorderMutate(t *testing.T) {
	c := NewCollection([]int{1, 2, 3})
	if c.Insert(0, 0) != c || c.Values()[0] != 0 {
		t.Fatalf("Insert() = %v", c.Values())
	}
	if c.Prepend(9) != c || c.Values()[0] != 9 {
		t.Fatalf("Prepend() = %v", c.Values())
	}
	if c.Reverse() != c || c.Values()[0] != 3 {
		t.Fatalf("Reverse() = %v", c.Values())
	}
	if c.Shuffle() != c || c.Count() != 5 {
		t.Fatalf("Shuffle() = %v", c.Values())
	}
}
