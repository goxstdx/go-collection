// Package slice_collcection 提供对泛型切片的链式集合操作。
//
// int、uint、float32、float64、string 在 Sort、Max、Contains 等方法调用时生成比较函数。
// 结构体和定义类型必须把比较函数传给对应的 By 方法。
// Min、Max、Mode 在空集合上返回 errorx.EmptyError，零值只在 err 为 nil 时有效。
package slice_collcection
