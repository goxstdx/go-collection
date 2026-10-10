# 切片集合使用（slice_collcection）

需要 Go 1.27。比较、取值和聚合都通过泛型方法完成，不再按字段名做反射。

## 初始化

```go
c := collection.NewSliceCollect([]int{3, 1, 2})
empty := collection.NewEmptyCollection[string]()
```

## 比较

`int`、`uint`（含各宽度）、`float32`、`float64`、`string` 直接调用无参方法。比较函数在第一次调用时生成，并缓存在当前集合上。

`type MyInt int` 这类定义类型，以及结构体，必须使用 `By` 方法并传入 `func(a, b T) int`。返回值小于 0 表示 `a < b`，等于 0 表示相等，大于 0 表示 `a > b`。传入 `nil` 会得到 `errorx.NilFunc`。

| 基础类型 | 其他类型 |
| --- | --- |
| `Max` / `Min` / `Mode` | `MaxBy` / `MinBy` / `ModeBy` |
| `Contains` / `ContainsCount` | `ContainsBy` / `ContainsCountBy` |
| `Search` / `Unique` | `SearchBy` / `UniqueBy` |
| `Diff` / `Union` / `Intersect` | `DiffBy` / `UnionBy` / `IntersectBy` |
| `Sort` / `SortDesc` | `SortBy` / `SortByDesc` / `SortByFunc` |

非基础类型调用无参方法时返回 `errorx.NoComparableError`。`Min`、`Max`、`Mode` 遇到空集合返回 `errorx.EmptyError`；伴随返回的零值只有在 `err == nil` 时才是结果。

```go
max, err := c.Max()

type User struct {
    ID   int
    Name string
    Age  int
}
users := collection.NewSliceCollect([]User{
    {ID: 1, Name: "a", Age: 10},
    {ID: 2, Name: "b", Age: 8},
})
oldest, err := users.MaxBy(func(a, b User) int {
    return cmp.Compare(a.Age, b.Age)
})
```

## 常用操作

- 过滤：`Filter`、`Reject`
- 查找：`Search`、`First`、`Last`
- 排序：`Sort`、`SortDesc`、`SortBy`、`SortByDesc`、`SortByFunc`
- 取值：`PluckField`、`KeyBy`
- 聚合：`Sum`、`Avg`、`Median`、`Mode`
- 其它：`Map`、`MapFilter`、`GroupBy`、`Split`、`ForPage`、`Nth`

```go
filtered := c.Filter(func(item int, _ int) bool { return item > 1 })
max, err := c.Max()
joined := c.Join(",")
_ = filtered
_ = max
_ = err
_ = joined
```

`Sort`、`SortDesc`、`SortBy`、`SortByDesc` 都是稳定排序，并直接修改当前集合。连续调用时，后一次排序是主键，前一次只在本次比较结果为 0 时保留。字符串形式的数字先在传入的函数里解析，再交给 `SortBy`。

```go
users.SortByDesc(func(u User) int { return u.ID })
users.SortBy(func(u User) int { return u.Age })

names := users.PluckField(func(u User) string { return u.Name })
byName := users.KeyBy(func(u User) string { return u.Name })
sum := users.Sum(func(u User) int { return u.Age })
_ = names
_ = byName
_ = sum
```

`Sum`、`Avg`、`Median` 接收一个取出数值的函数，空集合返回 `0`，不返回 error。`Map` 的返回元素类型由函数返回值决定，可以和原元素不同。

更多 API 说明见 pkg 文档：
- [`github.com/ZHOUXING1997/collection/slice_collcection`](https://pkg.go.dev/github.com/ZHOUXING1997/collection/slice_collcection)
