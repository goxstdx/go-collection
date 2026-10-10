# Map

适用集合：slice_collection

`Map[R any](func(item T, key int) R) *Collection[R]`

对集合中的每个元素进行函数映射。结果元素类型由函数返回值决定，可以和原元素不同。若需映射并过滤，请使用 `MapFilter`。

```go
c := collection.NewSliceCollect([]int{1, 2, 3, 4})
newC := c.Map(func(item int, _ int) int { return item * 2 })
newC.DD()

// 结果：2,4,6,8
```