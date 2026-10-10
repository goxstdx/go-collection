package errorx

import (
	"errors"
	"fmt"
)

// NoComparableError 表示当前类型没有默认比较函数。
var NoComparableError = errors.New("collection is not comparable")

// NotFoundError 表示 Search 没有找到元素。
var NotFoundError = fmt.Errorf("not found")

// EmptyError 表示 Min、Max、Mode 作用在空集合上。伴随返回的零值不是结果。
var EmptyError = errors.New("collection is empty")

var NotHaveKeyCompareFunc = errors.New("not have key compare func")
var NotHaveValCompareFunc = errors.New("not have value compare func")

var NilFunc = errors.New("func param is nil")
