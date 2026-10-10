# Append

适用集合：slice_collection

`Append(item T) *Collection[T]`

向集合尾部追加一个元素，返回当前集合以便链式调用。

```go
c := collection.NewSliceCollect([]int{1, 2})
c.Append(3)
c.DD()

/*
Collection(3, int):{
	0:	1
	1:	2
	2:	3
}
*/
```