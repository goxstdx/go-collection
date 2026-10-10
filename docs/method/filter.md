# Filter

适用集合：slice_collection

`Filter(func(item T, key int) bool) *Collection[T]`

根据过滤函数返回保留的元素，生成新的集合。

```go
c := collection.NewSliceCollect([]int{1, 2, 2, 3})
res := c.Filter(func(item int, _ int) bool { return item == 2 })
res.DD()

/*
Collection(2, int):{
	0:	2
	1:	2
}
*/
```