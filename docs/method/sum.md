# Sum

适用集合：slice_collection

`Sum[N number](fn func(T) N) float64`

用 `fn` 从每个元素取出数值后求和，结果为 `float64`。`N` 的底层类型必须是整数或浮点数，`type Money int64` 可以直接使用。空集合返回 `0`。

```go
c := collection.NewSliceCollect([]int{1, 2, 2, 3})
sum := c.Sum(func(n int) int { return n })
if sum != 8 {
    panic("sum 错误")
}
```

```go
type User struct {
    Age int
}

users := collection.NewSliceCollect([]User{{Age: 18}, {Age: 20}})
total := users.Sum(func(u User) int { return u.Age })
_ = total
```
