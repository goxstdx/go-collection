# API 说明

本库提供两类集合，模块要求 Go 1.27。

- 切片集合：`collection.NewSliceCollect`，实现在 `slice_collcection`
- Map 集合：`collection.NewMapCollect`，实现在 `map_collection`

取值、排序键和数值聚合都由调用方传入函数完成，库本身不按字段名反射。

## 切片集合

```go
c := collection.NewSliceCollect([]int{3, 1, 2})
```

### 比较

`int`、`uint`（含各宽度）、`float32`、`float64`、`string` 调用无参方法。比较函数在第一次调用时生成，并缓存在当前集合上。

结构体和 `type MyInt int` 这类定义类型必须使用 `By` 方法，传入 `func(a, b T) int`：

- 小于 0：`a < b`
- 等于 0：相等
- 大于 0：`a > b`

传入 `nil` 返回 `errorx.NilFunc`。非基础类型调用无参方法返回 `errorx.NoComparableError`，集合类结果为 `nil`。

| 基础类型 | 其他类型 |
| --- | --- |
| `Max` `Min` `Mode` | `MaxBy` `MinBy` `ModeBy` |
| `Contains` `ContainsCount` | `ContainsBy` `ContainsCountBy` |
| `Search` `Unique` | `SearchBy` `UniqueBy` |
| `Diff` `Union` `Intersect` | `DiffBy` `UnionBy` `IntersectBy` |
| `Sort` `SortDesc` | `SortBy` `SortByDesc` `SortByFunc` |

`Min`、`Max`、`Mode` 在空集合上返回 `errorx.EmptyError`。返回值里的零值只有在 `err == nil` 时才是结果。`Mode` 在次数相同时保留先出现的元素。

```go
import "cmp"

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

### 排序

`Sort`、`SortDesc`、`SortBy`、`SortByDesc`、`SortByFunc` 都是稳定排序，并直接修改当前集合。连续调用时，后一次是主键；前一次的顺序只在本次比较结果为 0 时保留。

先按 `ID` 倒序，再按 `Age` 正序：

```go
users.SortByDesc(func(u User) int { return u.ID })
users.SortBy(func(u User) int { return u.Age })
```

上面这组数据的结果是 `id 3、2、4、1`：年龄 8 在前，同龄时保留刚才 id 倒序的相对顺序。字符串形式的数字先在函数里解析，再交给 `SortBy`。

`SortByFunc` 的函数返回 `true` 表示第一个元素应排在前面。

### 取值与聚合

```go
names := users.PluckField(func(u User) string { return u.Name })
byName := users.KeyBy(func(u User) string { return u.Name }) // 相同 key 时后者覆盖前者
sum := users.Sum(func(u User) int { return u.Age })
avg := users.Avg(func(u User) int { return u.Age })
median := users.Median(func(u User) int { return u.Age })
```

`PluckField` 返回新集合，元素类型由函数决定。`Sum`、`Avg`、`Median` 只接受预声明整数或浮点数，空集合返回 `0`，不返回 error。`Median` 不修改原集合；偶数个元素时取中间两个的平均值。

`Map` 和 `MapFilter` 的结果类型也由函数决定，可以和原元素不同。`MapFilter` 在函数返回 `false` 时丢弃该元素。`GroupBy` 的 key 必须可比较。

## Map 集合

```go
mc := collection.NewMapCollect(map[string]int{"b": 2, "a": 1})
```

Map 的 key、value 比较仍由调用方设置，创建时不会自动生成：

```go
mc.OrderKeyByFunc(cmp.Compare)
mc.SetValCompare(cmp.Compare)
mc.OrderValue()
```

`PluckField` 用函数从 value 取值。集合已有排序键时按该顺序返回，否则按 map 的遍历顺序。包级函数 `map_collection.Pluck` 同样接收取值函数，顺序不稳定。

```go
values := mc.PluckField(func(n int) int { return n })
labels := mc.Map(func(v int, k string) string { return k })
sum := mc.Reduce(0, func(acc int, v int, k string) int { return acc + v })
```

`Map` 保留 key 的比较函数和已排序的 key，不保留 value 的比较函数。`Reduce` 的累加类型由初始值和函数决定。

## 错误

| 错误 | 含义 |
| --- | --- |
| `errorx.NoComparableError` | 当前类型没有默认比较，需要改用 `By` 方法 |
| `errorx.EmptyError` | `Min`、`Max`、`Mode` 作用在空集合上 |
| `errorx.NilFunc` | 必填函数参数为 `nil` |
| `errorx.NotFoundError` | `Search` 没有找到元素，下标为 `-1` |
| `errorx.NotHaveKeyCompareFunc` | Map 按 key 排序时还没有比较函数 |
| `errorx.NotHaveValCompareFunc` | Map 按 value 排序时还没有比较函数 |
