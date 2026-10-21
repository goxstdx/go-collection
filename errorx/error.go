package errorx

import (
	"errors"
	"fmt"
)

// NoComparableError means the current type has no default comparison function.
var NoComparableError = errors.New("collection is not comparable")

// NotFoundError means Search did not find the element.
var NotFoundError = fmt.Errorf("not found")

// EmptyError means Min, Max, or Mode, or a map MinBy or MaxBy, ran on an empty collection.
// The accompanying zero value is not a result.
var EmptyError = errors.New("collection is empty")

var NilFunc = errors.New("func param is nil")
