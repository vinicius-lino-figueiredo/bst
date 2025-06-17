package customutils

import (
	"fmt"
	"reflect"
	"time"
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
	if a == nil {
		if b == nil {
			return true
		}
		return false
	}
	if b == nil {
		return false
	}
	switch v := a.(type) {
	case int:
		return compareUnknownType(v, b)
	case int8:
		return compareUnknownType(v, b)
	case int16:
		return compareUnknownType(v, b)
	case int32:
		return compareUnknownType(v, b)
	case int64:
		return compareUnknownType(v, b)
	case uint:
		return compareUnknownType(v, b)
	case uint8:
		return compareUnknownType(v, b)
	case uint16:
		return compareUnknownType(v, b)
	case uint32:
		return compareUnknownType(v, b)
	case uint64:
		return compareUnknownType(v, b)
	case uintptr:
		return compareUnknownType(v, b)
	case float32:
		return compareUnknownType(v, b)
	case float64:
		return compareUnknownType(v, b)
	case string:
		return compareUnknownType(v, b)
	case time.Time:
		if b == nil {
			return false
		}
		if bt, ok := b.(time.Time); ok {
			return v.Equal(bt)
		}
		return false
	default:
		return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
	}
}

func compareUnknownType[T comparable](a T, b any) bool {
	if b == false {
		return false
	}
	if bt, ok := b.(T); ok {
		return a == bt
	}
	return false
}
