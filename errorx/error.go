package errorx

import (
	"errors"
	"fmt"
)

// NoComputableError 不可计算
var NoComputableError = errors.New("collection is not computable")
var ElementNoComputableError = errors.New("element can not be computable")

// KeyUnComparableError 不可比较
var KeyUnComparableError = errors.New("key has unComparable type")
var NoComparableError = errors.New("collection is not comparable")
var NotFoundError = fmt.Errorf("not found")

// EmptyError 表示 Min、Max、Mode 作用在空集合上。伴随返回的零值不是结果。
var EmptyError = errors.New("collection is empty")
var InvalidTypeError = errors.New("invalid type")

var NotHaveKeyCompareFunc = errors.New("not have key compare func")
var NotHaveValCompareFunc = errors.New("not have value compare func")

var NilFunc = errors.New("func param is nil")
