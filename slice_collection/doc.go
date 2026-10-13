// Package slice_collection provides generic, chainable slice collections with filtering,
// mapping, sorting, grouping, and aggregation.
//
// 提供对泛型切片的链式集合操作。
//
// 增删改和重新排序直接修改当前集合：Append、Push、Pop、Remove、Insert、Prepend、SetIndex、FromJson，以及 Sort、SortDesc、SortBy、SortByDesc、SortByFunc、Reverse、Shuffle。
// 需要保留原内容时，先 Copy 再调用。
// 查询和转换返回新集合。Slice、ForPage、Split 复制元素，不和原集合共用底层数组。
// NewCollection 不复制传入的切片，Values 返回内部切片，两者和调用方共用底层数组。
//
// int、uint、float32、float64、string 在 Sort、Max、Contains 等方法调用时生成比较函数。
// 结构体和定义类型必须把比较函数传给对应的 By 方法。
// Min、Max、Mode 在空集合上返回 errorx.EmptyError，零值只在 err 为 nil 时有效。
package slice_collection
