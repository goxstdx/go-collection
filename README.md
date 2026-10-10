# Collection

Collection包目标是用于替换golang原生的Slice，使用场景是在大量不追求极致性能，追求业务开发效能的场景。

[![Go Reference](https://pkg.go.dev/badge/github.com/ZHOUXING1997/collection.svg)](https://pkg.go.dev/github.com/ZHOUXING1997/collection)

> 本项目 forked from [jianfengye/collection](https://github.com/jianfengye/collection)，并修改部分内容以适配自己的项目。

## 简介

Collection包封装了对切片（Slice）的各种操作，使其更符合业务开发的语义，提高开发效率。该库支持多种数据类型，包括基本类型和结构体类型，并提供了丰富的方法来操作和转换这些数据。

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
| `OrderByValueFunc(func(V) any, func(any, any) int)` | `OrderByValueFunc[R any](func(V) R, func(R, R) int)` |
| `SafeCollection.Reduce(any, func(any, V, K) any)` | `Reduce[R any](init R, fn func(R, V, K) R)` |

行为也有变化：

- `Min`、`Max`、`Mode` 在空集合上返回 `errorx.EmptyError`，不再把零值和 `nil` error 当作成功。
- `Unique`、`Diff` 等比较失败时返回 `nil` 和错误，不再把原集合和错误一起返回。
- `Insert`、`Prepend`、`Reverse`、`Shuffle` 改为直接修改当前集合。要保留原内容时，先 `Copy`。
- `Slice`、`ForPage`、`Split` 改为返回元素副本，不再和原集合共用底层数组。
- `NewSliceCollect` 仍不复制传入的切片，`Values` 仍返回内部切片，两者和调用方共用底层数组。
- `SafeCollection.All` 返回 map 副本。普通 `Collection.All` 仍返回内部 map 的引用。
- 已删除且没有替代的错误：`NoComputableError`、`ElementNoComputableError`、`KeyUnComparableError`、`InvalidTypeError`。

## 版本历史

| 版本 | 说明 |
| --- | --- |
| v0.1.0 | 去掉反射和预存比较函数；切片的增删改和重排直接修改原集合 |
| v0.0.6 | 重新打tag，无功能变更 |
| v0.0.5 | 添加 pkg.go.dev 文档支持，完善包和主要类型的文档注释 |
| v0.0.4 | 修复部分问题，优化代码结构 |
| v0.0.3 | 增加新的功能和测试用例 |
| v0.0.2 | 修复初始版本的问题，改进性能 |
| v0.0.1 | 初始版本，从 jianfengye/collection 项目分支而来 |

## 许可证

`collection` 使用 [Apache License 2.0](LICENSE) 许可证。
