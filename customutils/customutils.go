package customutils

import (
	"fmt"
)

// DefaultCompareKeysFunction will work for numbers, strings and dates
func DefaultCompareKeysFunction(a any, b any) int {
	acomp := NewCaster(a)
	c, ok := acomp.Compare(Retrieve(b))
	if !ok {
		bcomp := NewCaster(b)
		c, ok := bcomp.Compare(Retrieve(a))
		if !ok {
			panic(fmt.Errorf("Couldn't compare elements %T and %T", a, b))
		}
		return -c
	}
	return c
}

// DefaultCheckValueEquality checks whether two values are equal (used in
// non-unique deletion)
func DefaultCheckValueEquality(a, b any) bool {
	return a == b
}
