# Map 集合使用（map_collection）
## 初始化
```go
m := map[string]int{"a": 2, "b": 1}
mc := collection.NewMapCollect(m)
```
## 取最小和最大
集合不保存顺序。`MinBy`、`MaxBy` 每次调用都重新扫描全部元素，复杂度随元素个数线性增长，不是排序。结果不保存，再调一次就再扫描一遍。比较函数返回负数表示左边小于右边。比较结果相等时保留哪一个不固定；比较函数里再用 key 做最后比较，结果就固定。空集合返回 `errorx.EmptyError`。带条件时先 `Filter`，再取最小或最大。
```go
k, v, err := mc.MinBy(func(a string, va int, b string, vb int) int {
    if a < b { return -1 }
    if a > b { return 1 }
    return 0
})
k, v, err = mc.Filter(func(n int, _ string) bool {
    return n > 0
}).MinBy(func(a string, va int, b string, vb int) int {
    if va < vb { return -1 }
    if va > vb { return 1 }
    return 0
})
```
`Each`、`Foreach` 都不排序，顺序不确定；`Each` 的回调返回 `false` 时提前终止，`Foreach` 总是完整遍历。`EachSorted` 每次调用都按比较函数排序后再遍历，复杂度随元素个数按 n log n 增长，顺序不保存。比较结果相等时先后不固定。回调里删掉的 key 会跳过，新增的 key 不会进入这次遍历。`PluckField` 也不排序。自己排 key 时先取出 map：`m := mc.All()`，再用 `slices.SortedFunc(maps.Keys(m), cmp)`。`All` 返回内部 map 的引用，遍历期间不要修改它。
## 常用操作
- 提取：`Keys`、`Values`、`PluckField`
- 修改：`Set`（就地）、`Put`（返回新集合）、`Merge`/`MergeInPlace`
- 过滤：`Filter`、`Only`、`Except`
- 定位与聚合：`MinBy`、`MaxBy`、`EachSorted`、`Reduce`
- 序列化：`ToJSON`

`PluckField` 用函数从 value 取值，不再按字段名反射。结果顺序不确定。包级 `map_collection.Pluck` 同样接收取值函数，遍历顺序也不稳定。

`Map`、`Reduce` 的结果类型由函数返回值决定，可以和原 value 不同。

```go
k, v, err := mc.MinBy(func(a string, va int, b string, vb int) int {
    return cmp.Compare(va, vb)
})
jsonStr, _ := mc.ToJSON()
values := mc.PluckField(func(n int) int { return n })
_ = k
_ = v
_ = err
_ = jsonStr
_ = values
```
更多 API 说明见 pkg 文档：
- [`github.com/goxstdx/go-collection/map_collection`](https://pkg.go.dev/github.com/goxstdx/go-collection/map_collection)
