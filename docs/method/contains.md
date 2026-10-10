# Contains

适用集合：slice_collcection

`Contains(obj T) (bool, error)`

判断元素是否在集合中。`int`、`uint`、`float32`、`float64`、`string` 可直接调用。结构体和定义类型使用 `ContainsBy`，并传入比较函数。

```go
c := collection.NewSliceCollect([]int{1, 2, 2, 3})

ok, err := c.Contains(1)  // true, nil
ok2, err := c.Contains(5) // false, nil
_ = ok
_ = ok2
_ = err
```

```go
type User struct {
    Name string
}

users := collection.NewSliceCollect([]User{{Name: "a"}})
ok, err := users.ContainsBy(User{Name: "a"}, func(a, b User) int {
    return cmp.Compare(a.Name, b.Name)
})
```
