// Package map_collection 提供对泛型 map 的链式集合操作。
//
// 集合不保存 key 顺序。MinBy、MaxBy 每次调用都重新扫描全部元素。EachSorted 每次调用都重新排序。结果都不保存。
// PluckField、Map、Reduce 的结果类型由传入函数决定，不按字段名反射取值。
package map_collection
