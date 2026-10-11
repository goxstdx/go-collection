# Collection

Collection 是一个 Go 泛型集合库，为切片（Slice）和 Map 提供常用的集合操作，适用于侧重业务开发效率的场景。

[![Go Reference](https://pkg.go.dev/badge/github.com/ZHOUXING1997/collection.svg)](https://pkg.go.dev/github.com/ZHOUXING1997/collection)

## 简介

Collection 提供过滤、映射、查找、排序、分组和数值聚合等操作，支持基本类型和结构体类型。通过泛型和调用方传入的函数完成类型转换、比较及取值，并提供线程安全的 Map 集合。

## 安装

```bash
go get github.com/ZHOUXING1997/collection@v0.1.0
```

## 支持的类型

元素类型是任意 `T`。比较和数值聚合另有约束：

- `Sort`、`Max`、`Min`、`Contains`、`Unique`、`Diff`、`Union`、`Intersect`、`Mode` 的无参形式只接受预声明的整数、无符号整数、`float32`、`float64`、`string`
- 结构体和 `type MyInt int` 这类定义类型，把比较函数传给对应的 `By` 方法
- `Sum`、`Avg`、`Median` 通过函数取出底层类型为整数或浮点数的值

## 使用示例

### 初始化 Collection

```go
intColl := collection.NewSliceCollect([]int{1, 2, 3, 4, 5})
strColl := collection.NewSliceCollect([]string{"a", "b", "c"})

type User struct {
    Name string
    Age  int
}
users := []User{{Name: "张三", Age: 18}, {Name: "李四", Age: 20}}
objColl := collection.NewSliceCollect(users)
```

### 常用方法示例

```go
filtered := intColl.Filter(func(item int, key int) bool {
    return item > 2
})

mapped := intColl.Map(func(item int, key int) int {
    return item * 2
})

sorted, err := intColl.Sort()
index, err := intColl.Search(3)
first := intColl.First()
contains, err := intColl.Contains(3)

oldest, err := objColl.MaxBy(func(a, b User) int {
    return cmp.Compare(a.Age, b.Age) // cmp 为标准库 cmp
})
```

## 更多示例

完整接口见 [API 说明](docs/guide/api.md)。

## 从 v0.0.6 迁移到 v0.1.0

v0.1.0 删除了按字段名反射取值，也删除了预先保存的比较函数。Go 版本要求 1.27。

| v0.0.6 | v0.1.0 |
| --- | --- |
| `SetCompare(func(any, any) int)` | 删除。`int`、`uint`、`float32`、`float64`、`string` 直接调用；其他类型使用 `By` 方法 |
| `PluckString` / `PluckInt64` / `PluckFloat64` / `PluckUint64` / `PluckBool` | `PluckField(func(T) R)` |
| `SortBy(field)` / `SortByDesc(field)` / `SortFloatBy(field)` | `SortBy(func(T) R)` / `SortByDesc(func(T) R)`。字符串数字先在函数里解析 |
| `KeyByStrField(field)` | `KeyBy(func(T) K)` |
| `Sum()` / `Avg()` / `Median()`，返回 `(float64, error)` | `Sum(fn)` / `Avg(fn)` / `Median(fn)`，返回 `float64`。空集合为 `0` |
| `Map(func(T, int) T)` | `Map[R any](func(T, int) R)` |
| `MapFilter(func(T, int) (T, bool))` | `MapFilter[R any](func(T, int) (R, bool))` |
| `GroupBy(func(T, int) interface{})` | `GroupBy[K comparable](func(T, int) K)` |
| `Join(split, func(interface{}) string)` | `Join(split, func(T) string)` |
| `map.Pluck(field)` / `map.PluckFunc` | `PluckField(func(V) R)` |
| `map_collection.Pluck(m, fieldName)` | `map_collection.Pluck(m, func(V) R)`。nil 指针由函数自己处理 |
| `Map(func(V, K) V)` / `Reduce(any, func(...) any)` | 结果类型由函数返回值决定 |
| `OrderByValueFunc(func(V) any, func(any, any) int)` | 删除。按 value 取最小或最大时，把比较写进 `MinBy` / `MaxBy` 的 `func(a K, va V, b K, vb V) int` |
| `WithKeyCompare` / `SetKeyCompare` / `OrderKey` / `OrderKeyByFunc` | 删除。`MinBy`、`MaxBy` 每次调用都重新扫描全部元素，取最小或最大，不保存结果 |
| `First` / `Last` / `FirstWhere` / `LastWhere` | `MinBy` / `MaxBy`，返回 `(K, V, error)`。带条件时先 `Filter`。空集合返回 `errorx.EmptyError` |
| `Foreach(fn)`（原来有序） | `Foreach` 仍在，但不再有序，按 map 自身顺序完整遍历。需要有序时用 `EachSorted(cmp, fn) error` |
| `Each(func(V, K))` | `Each(func(V, K) bool)`，回调返回 `false` 时提前终止 |
| `WithValCompare` / `SetValCompare` / `OrderValue` | 删除。比较函数里比较 value |
| `SafeCollection.Reduce(any, func(any, V, K) any)` | `Reduce[R any](init R, fn func(R, V, K) R)` |

行为也有变化：

- `Min`、`Max`、`Mode` 在空集合上返回 `errorx.EmptyError`，不再把零值和 `nil` error 当作成功。
- `Unique`、`Diff` 等比较失败时返回 `nil` 和错误，不再把原集合和错误一起返回。
- `Insert`、`Prepend`、`Reverse`、`Shuffle` 改为直接修改当前集合。要保留原内容时，先 `Copy`。
- `Slice`、`ForPage`、`Split` 改为返回元素副本，不再和原集合共用底层数组。
- `NewSliceCollect` 仍不复制传入的切片，`Values` 仍返回内部切片，两者和调用方共用底层数组。
- `SafeCollection.All` 返回 map 副本。普通 `Collection.All` 仍返回内部 map 的引用。
- Map 的 `MinBy`、`MaxBy` 每次调用都重新扫描全部元素，是线性扫描，不是排序。结果不保存，连续调用会各扫描一遍。空集合返回 `errorx.EmptyError`。
- Map 的 `EachSorted` 每次调用都重新排序，复杂度随元素个数按 n log n 增长，只返回 error。`Each`、`Foreach` 都按 map 自身顺序遍历，顺序不确定；`Each` 的回调返回 `false` 时提前终止，`Foreach` 总是完整遍历。
- Map 的内部 map 不会为 nil。`NewMapCollect(nil)`、空集合的 `Except`、两个 nil map 的 `Merge` 都返回空 map，后续 `Set`、`MergeInPlace` 等写操作不会 panic。
- 已删除且没有替代的错误：`NoComputableError`、`ElementNoComputableError`、`KeyUnComparableError`、`InvalidTypeError`、`NotHaveKeyCompareFunc`、`NotHaveValCompareFunc`。

## 版本历史

| 版本 | 说明 |
| --- | --- |
| v0.1.0 | 去掉反射和预存比较函数；切片的增删改和重排直接修改原集合 |
| v0.0.6 | 重新打tag，无功能变更 |
| v0.0.5 | 添加 pkg.go.dev 文档支持，完善包和主要类型的文档注释 |
| v0.0.4 | 修复部分问题，优化代码结构 |
| v0.0.3 | 增加新的功能和测试用例 |
| v0.0.2 | 修复初始版本的问题，改进性能 |
| v0.0.1 | 初始版本 |

## 许可证

`collection` 使用 [Apache License 2.0](LICENSE) 许可证。
