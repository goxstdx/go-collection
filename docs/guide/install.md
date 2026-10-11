# 安装与版本

本项目包含两类集合能力：
- 切片集合（包：`slice_collection`，入口函数：`collection.NewSliceCollect` 等）
- Map 集合（包：`map_collection`，入口函数：`collection.NewMapCollect` 等）

推荐使用 Go Modules，模块要求 Go 1.27：

```bash
go get github.com/goxstdx/go-collection@latest
```

如需固定版本，请将 `latest` 替换为具体 tag，例如：

```bash
go get github.com/goxstdx/go-collection@v0.1.1
```

快速开始（切片与 Map 示例）：

```go
package main

import (
    "fmt"
    "github.com/goxstdx/go-collection"
)

func main() {
    // 切片集合（slice_collection）
    c := collection.NewSliceCollect([]int{1, 2, 3})
    sum := c.Sum(func(n int) int { return n })
    fmt.Println(sum)

    // Map 集合（map_collection）
    m := map[string]int{"a": 1, "b": 2}
    mc := collection.NewMapCollect(m)
    fmt.Println(mc.Count())
}
```

更多版本信息请查看 `docs/guide/RELEASES.md` 或前往 pkg 站点查看：
- 入口包：[`github.com/goxstdx/go-collection`](https://pkg.go.dev/github.com/goxstdx/go-collection)
- 切片集合：[`.../slice_collection`](https://pkg.go.dev/github.com/goxstdx/go-collection/slice_collection)
- Map 集合：[`.../map_collection`](https://pkg.go.dev/github.com/goxstdx/go-collection/map_collection)
