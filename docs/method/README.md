# 方法手册说明（slice 集合）

本目录文档默认针对切片集合（包：`slice_collection`），示例基于泛型 API：
- 初始化：`collection.NewSliceCollect([]T{...})`
- 常用：`Filter`、`Map`、`MapFilter`、`Sum`、`Avg`、`Contains`、`Sort` 等
- `int`、`uint`、`float`、`string` 的比较在调用 `Sort`、`Max`、`Contains` 等无参方法时生成
- 结构体比较使用对应的 `By` 方法；`Sum`、`Avg`、`Median` 通过函数取出数值

Map 集合（包：`map_collection`）的使用说明请参考：
- `docs/guide/usage-map.md`
- pkg 文档：`github.com/ZHOUXING1997/collection/map_collection`
