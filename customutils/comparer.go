// Package customutils has some auxiliar types and methods for the main package.
package customutils

import (
	"cmp"
	"reflect"
	"time"
)

// BoolTrueNumericValue represents the value the comparison assumes to be the
// correct numeric representation for the boolean true. Since the correct order
// is arbitrary, user can change the value used by the package to fit their
// needs.
var BoolTrueNumericValue = 1

// NewCaster returns an implementation of caster by using an equivalent type
func NewCaster(v any) caster {
	if c, ok := v.(caster); ok {
		return c
	}
	switch v := v.(type) {
	case bool:
		return boolCaster(v)
	case uint8:
		return uint8Caster(v)
	case uint16:
		return uint16Caster(v)
	case uint32:
		return uint32Caster(v)
	case uint64:
		return uint64Caster(v)
	case int8:
		return int8Caster(v)
	case int16:
		return int16Caster(v)
	case int32:
		return int32Caster(v)
	case int64:
		return int64Caster(v)
	case float32:
		return float32Caster(v)
	case float64:
		return float64Caster(v)
	case string:
		return stringCaster(v)
	case int:
		return intCaster(v)
	case uint:
		return uintCaster(v)
	case time.Time:
		return timeCaster(v)
	default:
		return NewInterfaceCaster(v)
	}
}

// Retrieve returns the original value from any. If v implements caster,
// Retrieve returns the [caster.cast] value.
func Retrieve(v any) any {
	if v == nil {
		return nil
	}
	if c, ok := v.(caster); ok {
		return c.cast()
	}
	return v
}

// Comparer can be implemented to allow comparison for non-Comparable types.
type Comparer interface {
	Compare(b any) (int, bool)
}

// caster is a comparer that returns the original type.
type caster interface {
	Comparer
	cast() any
}

// boolCaster implements Comparer.
type boolCaster bool

// Compare implements Comparer.
func (b boolCaster) Compare(other any) (int, bool) { return compareBool(bool(b), other) }

// cast implements caster.
func (b boolCaster) cast() any { return bool(b) }

// uint8Caster implements Comparer.
type uint8Caster uint8

// Compare implements Comparer.
func (u uint8Caster) Compare(other any) (int, bool) { return compare(uint8(u), other) }

// cast implements caster.
func (u uint8Caster) cast() any { return uint8(u) }

// uint16Caster implements Comparer.
type uint16Caster uint16

// Compare implements Comparer.
func (u uint16Caster) Compare(other any) (int, bool) { return compare(uint16(u), other) }

// cast implements caster.
func (u uint16Caster) cast() any { return uint16(u) }

// uint32Caster implements Comparer.
type uint32Caster uint32

// Compare implements Comparer.
func (u uint32Caster) Compare(other any) (int, bool) { return compare(uint32(u), other) }

// cast implements caster.
func (u uint32Caster) cast() any { return uint32(u) }

// uint64Caster implements Comparer.
type uint64Caster uint64

// Compare implements Comparer.
func (u uint64Caster) Compare(other any) (int, bool) { return compare(uint64(u), other) }

// cast implements caster.
func (u uint64Caster) cast() any { return uint64(u) }

// int8Caster implements Comparer.
type int8Caster int8

// Compare implements Comparer.
func (i int8Caster) Compare(other any) (int, bool) { return compare(int8(i), other) }

// cast implements caster.
func (i int8Caster) cast() any { return int8(i) }

// int16Caster implements Comparer.
type int16Caster int16

// Compare implements Comparer.
func (i int16Caster) Compare(other any) (int, bool) { return compare(int16(i), other) }

// cast implements caster.
func (i int16Caster) cast() any { return int16(i) }

// int32Caster implements Comparer.
type int32Caster int32

// Compare implements Comparer.
func (i int32Caster) Compare(other any) (int, bool) { return compare(int32(i), other) }

// cast implements caster.
func (i int32Caster) cast() any { return int32(i) }

// int64Caster implements Comparer.
type int64Caster int64

// Compare implements Comparer.
func (i int64Caster) Compare(other any) (int, bool) { return compare(int64(i), other) }

// cast implements caster.
func (i int64Caster) cast() any { return int64(i) }

// float32Caster implements Comparer.
type float32Caster float32

// Compare implements Comparer.
func (f float32Caster) Compare(other any) (int, bool) { return compare(float32(f), other) }

// cast implements caster.
func (f float32Caster) cast() any { return float32(f) }

// float64Caster implements Comparer.
type float64Caster float64

// Compare implements Comparer.
func (f float64Caster) Compare(other any) (int, bool) { return compare(float64(f), other) }

// cast implements caster.
func (f float64Caster) cast() any { return float64(f) }

// stringCaster implements Comparer.
type stringCaster string

// Compare implements Comparer.
func (s stringCaster) Compare(other any) (int, bool) { return compare(string(s), other) }

// cast implements caster.
func (s stringCaster) cast() any { return string(s) }

// intCaster implements Comparer.
type intCaster int

// Compare implements Comparer.
func (i intCaster) Compare(other any) (int, bool) { return compare(int(i), other) }

// cast implements caster.
func (i intCaster) cast() any { return int(i) }

// uintCaster implements Comparer.
type uintCaster uint

// Compare implements Comparer.
func (u uintCaster) Compare(other any) (int, bool) { return compare(uint(u), other) }

// cast implements caster.
func (u uintCaster) cast() any { return uint(u) }

// timeCaster implements Comparer.
type timeCaster time.Time

// Compare implements Comparer.
func (t timeCaster) Compare(other any) (int, bool) { return compareTime(time.Time(t), other) }

// cast implements caster.
func (t timeCaster) cast() any { return time.Time(t) }

// NewInterfaceCaster returns a new valid instance of Any.
func NewInterfaceCaster(v any) *Any {
	if v == nil {
		return new(Any)
	}
	reflectV := reflect.ValueOf(v)
	if !reflectV.IsValid() || reflectV.IsNil() || reflectV.Comparable() {
		return &Any{value: v}
	}
	return &Any{
		value:        v,
		reflectValue: reflectV,
		Type:         reflectV.Type(),
		Comparable:   true,
	}
}

// Any implements Comparer.
type Any struct {
	value        any
	reflectValue reflect.Value
	Type         reflect.Type
	Comparable   bool
}

// Compare implements Comparer.
func (a Any) Compare(other any) (int, bool) {
	if other == nil {
		return 1, true
	}
	otherV := reflect.ValueOf(other)
	if !otherV.IsValid() {
		return 0, false
	}
	if otherV.IsNil() {
		return 1, true
	}
	if otherV.Type() != a.Type {
		return 0, false
	}
	if !a.Comparable {
		return 0, false
	}
	if !otherV.Comparable() {
		return 0, false
	}
	if reflect.DeepEqual(a.value, otherV) {
		return 0, true
	}
	return 0, false

}

// cast implements caster.
func (a Any) cast() any { return a.value }

func compareTime(a time.Time, b any) (int, bool) {
	if b == nil { // if no value, a is greater
		return 1, true
	}
	bt, ok := b.(time.Time)
	if !ok {
		return 0, false
	}
	return a.Compare(bt), true
}

// compareBool compares two values assuming they are booleans. Since comparing
// booleans might be arbitrary, the value order can be set.
func compareBool(a bool, b any) (int, bool) {
	if b == nil { // if no value, a is greater
		return 1, true
	}
	bb, ok := b.(bool)
	if !ok {
		return 0, false
	}

	ai, bi := 0, 0
	if a {
		ai = BoolTrueNumericValue
	}
	if bb {
		bi = BoolTrueNumericValue
	}
	return ai - bi, true
}

func compare[T cmp.Ordered](a T, b any) (int, bool) {
	if b == nil { // if no value, a is greater
		return 1, true
	}
	switch b := b.(type) {
	case T:
		return cmp.Compare(a, b), true
	default:
		if bc, ok := b.(Comparer); ok {
			c, ok := bc.Compare(a)
			return -c, ok
		}
		return 0, false
	}
}
