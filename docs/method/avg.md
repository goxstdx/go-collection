# Avg

适用集合：slice_collcection

`Avg[N number](fn func(T) N) float64`

用 `fn` 从每个元素取出数值后求平均，返回 `float64`。空集合返回 `0`。

```go
c := collection.NewSliceCollect([]int{1, 2, 2, 3})
avg := c.Avg(func(n int) int { return n })
if avg != 2.0 {
    panic("Avg error")
}
```
