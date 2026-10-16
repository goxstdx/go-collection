# go-collection

[English](README.md) | [简体中文](README.zh-CN.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/goxstdx/go-collection.svg)](https://pkg.go.dev/github.com/goxstdx/go-collection)

## Overview

go-collection is a Go library for generic, chainable slice and map operations: filtering, mapping, sorting, grouping, and aggregation. Callbacks provide typed comparisons and value extraction without field-name reflection. Mapping can change the element type, and the `map_collection` package also provides a concurrency-safe `SafeCollection`.

## Installation

Requires Go 1.27 or later.

```bash
go get github.com/goxstdx/go-collection@latest
```

## Quick start

```go
package main

import (
    "fmt"

    "github.com/goxstdx/go-collection"
)

func main() {
    numbers := collection.NewSliceCollect([]int{1, 2, 3, 4, 5})
    labels := numbers.Filter(func(n int, _ int) bool {
        return n > 2
    }).Map(func(n int, _ int) string {
        return fmt.Sprintf("item-%d", n)
    })
    fmt.Println(labels.Values()) // [item-3 item-4 item-5]

    scores := collection.NewMapCollect(map[string]int{"alice": 80, "bob": 95})
    passed := scores.Filter(func(score int, _ string) bool {
        return score >= 90
    })
    fmt.Println(passed.Count()) // 1
}
```

## API documentation

The root package provides constructors. Collection methods are documented in the subpackages that define their types.

| Package | Documentation |
| --- | --- |
| `collection` | [Constructors](https://pkg.go.dev/github.com/goxstdx/go-collection) |
| `slice_collection` | [Slice Collection methods](https://pkg.go.dev/github.com/goxstdx/go-collection/slice_collection#Collection) |
| `map_collection` | [Map Collection methods](https://pkg.go.dev/github.com/goxstdx/go-collection/map_collection#Collection) |
| `map_collection` | [SafeCollection methods](https://pkg.go.dev/github.com/goxstdx/go-collection/map_collection#SafeCollection) |

The [API guide](docs/guide/api.md) and [migration notes](README.zh-CN.md#从-v006-迁移到-v010) are available in Chinese.

## Supported types

- Collection elements can be any type `T`.
- Comparison methods such as `Sort`, `Max`, `Min`, `Contains`, `Unique`, `Diff`, `Union`, `Intersect`, and `Mode` support predeclared integer, unsigned integer, floating-point, and string types without a callback.
- For structs and defined types such as `type MyInt int`, use the corresponding `By` methods with a comparison callback.
- `Sum`, `Avg`, and `Median` accept callbacks returning values whose underlying types are integers or floating-point numbers.

## Mutation and iteration

- Slice updates and reordering operations such as `Append`, `Insert`, `Remove`, `Sort`, `Reverse`, and `Shuffle` modify the current collection. Call `Copy` first to preserve it.
- Slice `Filter` and `Map` return new collections. `Slice`, `ForPage`, and `Split` copy elements into separate backing arrays.
- `NewSliceCollect` and `Values` share the backing array with the caller; they do not copy it.
- Map iteration order is unspecified. `Each` stops when its callback returns `false`; `Foreach` visits every entry. Use `EachSorted` for ordered traversal.
- `SafeCollection.All` returns a map copy; ordinary map `Collection.All` returns the internal map.

## License

Licensed under the [Apache License 2.0](LICENSE).
