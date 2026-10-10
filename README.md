# Collection

Collection包目标是用于替换golang原生的Slice，使用场景是在大量不追求极致性能，追求业务开发效能的场景。

[![Go Reference](https://pkg.go.dev/badge/github.com/ZHOUXING1997/collection.svg)](https://pkg.go.dev/github.com/ZHOUXING1997/collection)

> 本项目 forked from [jianfengye/collection](https://github.com/jianfengye/collection)，并修改部分内容以适配自己的项目。

## 简介

Collection包封装了对切片（Slice）的各种操作，使其更符合业务开发的语义，提高开发效率。该库支持多种数据类型，包括基本类型和结构体类型，并提供了丰富的方法来操作和转换这些数据。

## 安装

```bash
go get github.com/ZHOUXING1997/collection@v0.0.5
```

## 支持的类型

元素类型是任意 `T`。比较和数值聚合另有约束：

- `Sort`、`Max`、`Min`、`Contains`、`Unique`、`Diff`、`Union`、`Intersect`、`Mode` 的无参形式只接受预声明的整数、无符号整数、`float32`、`float64`、`string`
- 结构体和 `type MyInt int` 这类定义类型，把比较函数传给对应的 `By` 方法
- `Sum`、`Avg`、`Median` 通过函数取出预声明整数或浮点数

## 使用示例

### 初始化 Collection

```go
intColl := collection.NewSliceCollect([]int{1, 2, 3, 4, 5})
strColl := collection.NewSliceCollect([]string{"a", "b", "c"})

type User struct {
    Name string
    Age  int
}
users := []User{{Name: "张三", Age: 18}, {Name: "李四", Age: 20}}
objColl := collection.NewSliceCollect(users)
```

### 常用方法示例

```go
filtered := intColl.Filter(func(item int, key int) bool {
    return item > 2
})

mapped := intColl.Map(func(item int, key int) int {
    return item * 2
})

sorted, err := intColl.Sort()
index, err := intColl.Search(3)
first := intColl.First()
contains, err := intColl.Contains(3)

oldest, err := objColl.MaxBy(func(a, b User) int {
    return cmp.Compare(a.Age, b.Age) // cmp 为标准库 cmp
})
```

## 更多示例

完整接口见 [API 说明](docs/guide/api.md)。

## 版本历史

| 版本     | 说明                                  |
|--------|-------------------------------------|
| v0.0.6 | 重新打tag，无功能变更                        |
| v0.0.5 | 添加 pkg.go.dev 文档支持，完善包和主要类型的文档注释    |
| v0.0.4 | 修复部分问题，优化代码结构                       |
| v0.0.3 | 增加新的功能和测试用例                         |
| v0.0.2 | 修复初始版本的问题，改进性能                      |
| v0.0.1 | 初始版本，从 jianfengye/collection 项目分支而来 |

## 许可证

`collection` 使用 [Apache License 2.0](LICENSE) 许可证。
