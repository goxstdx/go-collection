# Map 集合使用（map_collection）
## 初始化
```go
m := map[string]int{"a": 2, "b": 1}
mc := collection.NewMapCollect(m)
```
## Key/Value 排序与遍历
- 按 key 排序遍历：
```go
mc.SetKeyCompare(func(a, b string) int {
    if a < b { return -1 }
    if a > b { return 1 }
    return 0
}).OrderKey()
mc.Foreach(func(v int, k string) { /* ... */ })
```
- 按 value 排序（更新有序键列表）：
```go
mc.SetValCompare(func(a, b int) int {
    switch {
    case a < b: return -1
    case a > b: return 1
    default: return 0
    }
}).OrderValue()
```
## 常用操作
- 提取：`Keys`、`Values`、`PluckField`
- 修改：`Set`（就地）、`Put`（返回新集合）、`Merge`/`MergeInPlace`
- 过滤：`Filter`、`Only`、`Except`
- 定位与聚合：`First`、`Last`、`FirstWhere`、`LastWhere`、`Reduce`
- 序列化：`ToJSON`

`PluckField` 用函数从 value 取值，不再按字段名反射。集合上已有排序键时按该顺序返回。包级 `map_collection.Pluck` 同样接收取值函数，遍历顺序不稳定。

`Map`、`Reduce` 的结果类型由函数返回值决定，可以和原 value 不同。

```go
k, v, ok := mc.First()
jsonStr, _ := mc.ToJSON()
values := mc.PluckField(func(n int) int { return n })
_ = k
_ = v
_ = ok
_ = jsonStr
_ = values
```
更多 API 说明见 pkg 文档：
- [`github.com/ZHOUXING1997/collection/map_collection`](https://pkg.go.dev/github.com/ZHOUXING1997/collection/map_collection)
