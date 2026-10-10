# API 说明

本库提供两类集合，模块要求 Go 1.27。

- 切片集合：`collection.NewSliceCollect`，实现在 `slice_collection`
- Map 集合：`collection.NewMapCollect`，实现在 `map_collection`

取值和数值聚合都由调用方传入函数完成，库本身不按字段名反射。

## 切片集合

```go
c := collection.NewSliceCollect([]int{3, 1, 2})
```

### 比较

`int`、`uint`（含各宽度）、`float32`、`float64`、`string` 调用无参方法。比较函数在调用时按元素类型生成，不保存在集合上。

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

### 是否修改原集合

增删改和重新排序直接修改当前集合：`Append`、`Push`、`Pop`、`Remove`、`Insert`、`Prepend`、`SetIndex`、`FromJson`，以及 `Sort`、`SortDesc`、`SortBy`、`SortByDesc`、`SortByFunc`、`Reverse`、`Shuffle`。要保留原内容时，先 `Copy` 再调用。

查询和转换返回新集合。`Slice`、`ForPage`、`Split` 复制元素，之后对结果 `Push` 或 `SetIndex` 不会改到原集合。

`NewSliceCollect` 不复制传入的切片，`Values` 返回内部切片。集合和调用方共用底层数组；容量足够时，上面的增删改会写到最初传入的切片。

### 排序

`Sort`、`SortDesc`、`SortBy`、`SortByDesc`、`SortByFunc`、`Reverse`、`Shuffle` 都直接修改当前集合。排序是稳定的。连续调用时，后一次是主键；前一次的顺序只在本次比较结果为 0 时保留。

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

`PluckField` 返回新集合，元素类型由函数决定。`Sum`、`Avg`、`Median` 接受底层类型为整数或浮点数的值，因此 `type Money int64` 可以直接传入。空集合返回 `0`，不返回 error。`Median` 不修改原集合；偶数个元素时取中间两个的平均值。

`Map` 和 `MapFilter` 的结果类型也由函数决定，可以和原元素不同。`MapFilter` 在函数返回 `false` 时丢弃该元素。`GroupBy` 的 key 必须可比较。

## Map 集合

```go
mc := collection.NewMapCollect(map[string]int{"b": 2, "a": 1})
```

Map 不保存 key 顺序。`MinBy`、`MaxBy` 每次调用都重新扫描全部元素，复杂度随元素个数线性增长，不是排序。结果不保存，再调一次会再扫描一遍。比较函数同时收到 key 和 value，按 key 或按 value 都用它。比较结果相等时保留哪一个不固定；比较函数里再用 key 做最后比较，结果就固定。空集合返回 `errorx.EmptyError`。比较函数为 `nil` 时返回 `errorx.NilFunc`。伴随返回的零值不是结果。

带条件时先 `Filter`，再 `MinBy` 或 `MaxBy`。`Filter` 会多分配一个 map。

`Each`、`Foreach` 都不排序，顺序和直接遍历 map 一样，不确定。`Each` 的回调返回 `bool`，返回 `false` 时提前终止；`Foreach` 总是完整遍历。`EachSorted` 每次调用都按比较函数排序后再遍历，复杂度随元素个数按 n log n 增长。顺序不保存，再调一次会再排一次。比较结果相等时先后不固定，因为 key 的原始顺序来自 map 遍历。比较函数或遍历函数为 `nil` 时返回 `errorx.NilFunc`。空集合调用成功，遍历函数不会被调用。排序用的是调用开始时的键值快照：回调里删掉的 key 会跳过，新增的 key 不会进入这次遍历，回调读到的是当时 map 里的 value。含 NaN 的 key 无法查回，使用快照里的 value。`PluckField` 也不排序。

自己排 key 时，先取出 map 再交给标准库。`All` 返回内部 map 的引用，遍历期间不要修改它。

```go
k, v, err := mc.MinBy(func(a string, va int, b string, vb int) int {
    return cmp.Compare(a, b)
})
k, v, err = mc.MaxBy(func(a string, va int, b string, vb int) int {
    return cmp.Compare(va, vb)
})
err = mc.EachSorted(func(a string, va int, b string, vb int) int {
    return cmp.Compare(a, b)
}, func(v int, k string) {
    _ = v
    _ = k
})
m := mc.All()
for _, k := range slices.SortedFunc(maps.Keys(m), func(a, b string) int {
    return cmp.Compare(a, b)
}) {
    v := m[k]
    _ = v
}
```

`PluckField` 用函数从 value 取值，结果顺序不确定。包级函数 `map_collection.Pluck` 同样接收取值函数，顺序也不稳定。

```go
values := mc.PluckField(func(n int) int { return n })
labels := mc.Map(func(v int, k string) string { return k })
sum := mc.Reduce(0, func(acc int, v int, k string) int { return acc + v })
```

`Map` 的结果类型由函数返回值决定。`Reduce` 的累加类型由初始值和函数决定。`SafeCollection.Reduce` 同样是泛型方法。

`All` 返回底层 map 的引用，修改会影响集合。`SafeCollection.All` 返回副本，修改返回值不会影响集合。回调里不要再调用同一个 `SafeCollection`：写方法会等待这把锁，读方法在已有写操作等待时也可能卡住。

## 错误

| 错误 | 含义 |
| --- | --- |
| `errorx.NoComparableError` | 当前类型没有默认比较，需要改用 `By` 方法 |
| `errorx.EmptyError` | 切片的 `Min`、`Max`、`Mode`，以及 map 的 `MinBy`、`MaxBy` 作用在空集合上 |
| `errorx.NilFunc` | 必填函数参数为 `nil` |
| `errorx.NotFoundError` | `Search` 没有找到元素，下标为 `-1` |
