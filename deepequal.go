// Deep equality test via reflection
// Based on https://golang.org/src/reflect/deepequal.go

package deepequal

import (
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

// During deepValueEqual, must keep track of checks that are
// in progress. The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueEqual(v1, v2 reflect.Value, visited map[visit]bool, depth int) (bool, error) {
	if !v1.IsValid() || !v2.IsValid() {
		if v1.IsValid() == v2.IsValid() {
			return true, nil
		}
		return false, fmt.Errorf("Only one is valid")
	}
	if v1.Type() != v2.Type() {
		return false, fmt.Errorf("Types for %v (type %v) and %v (type %v) do not match", v1, v1.Type(), v2, v2.Type())
	}

	// We want to avoid putting more in the visited map than we need to.
	// For any possible reference cycle that might be encountered,
	// hard(t) needs to return true for at least one of the types in the cycle.
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true, nil
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Float32, reflect.Float64:
		// Handle special NaN values:
		// Go treats math.Nan() == math.Nan() as false so we have to check for this
		if math.IsNaN(v1.Float()) && math.IsNaN(v2.Float()) {
			return true, nil
		}
		// Will continue with normal value comparison in the default case
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if equal, err := deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1); !equal {
				return false, fmt.Errorf("Error in array %v: %s", v1, err.Error())
			}
		}
		return true, nil
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false, fmt.Errorf("One of the slices is nil (%v vs %v)", v1, v2)
		}
		if v1.Len() != v2.Len() {
			return false, fmt.Errorf("One of the slices has a different length (%v (len %d) vs %v (len %d))", v1, v1.Len(), v2, v2.Len())
		}
		if v1.Pointer() == v2.Pointer() {
			return true, fmt.Errorf("Both slices have the same pointer address (%v (ptr %d) vs %v (ptr %d))", v1, v1.Pointer(), v2, v2.Pointer())
		}
		for i := 0; i < v1.Len(); i++ {
			if equal, err := deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1); !equal {
				return false, fmt.Errorf("Error in array %v: %s", v1, err.Error())
			}
		}
		return true, nil
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			if v1.IsNil() == v2.IsNil() {
				return true, nil
			}
			return false, fmt.Errorf("One interface is nil")
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1)
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true, nil
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if equal, err := deepValueEqual(v1.Field(i), v2.Field(i), visited, depth+1); !equal {
				fieldName := v1.Type().Field(i).Name
				return false, fmt.Errorf("Error in struct field %v: %s", fieldName, err.Error())
			}
		}
		return true, nil
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false, fmt.Errorf("One map is nil")
		}
		if v1.Len() != v2.Len() {
			return false, fmt.Errorf("Length of the maps differs")
		}
		if v1.Pointer() == v2.Pointer() {
			return true, nil
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			equal, err := deepValueEqual(val1, val2, visited, depth+1)
			if !val1.IsValid() || !val2.IsValid() || !equal {
				return false, fmt.Errorf("Difference in values of %v: %s", k, err.Error())
			}
		}
		return true, nil
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true, nil
		}
		// Can't do better than this:
		return false, fmt.Errorf("Cannot compare functions")
	default:
		// Do nothing and let normal comparison follow
	}

	// Can't do better than this as unexported fields want to remain hidden
	if !v1.CanInterface() || !v2.CanInterface() {
		if v1.CanInterface() == v2.CanInterface() {
			return true, nil
		}
		return false, fmt.Errorf("Only one value is not interfaceable (maybe unexported) (%v (%v) vs %v (%v))", v1, v1.Type(), v2, v2.Type())
	}

	if v1.Interface() == v2.Interface() {
		return true, nil
	}

	return false, fmt.Errorf("Value interface differs (%v (%v) vs %v (%v))", v1, v1.Type(), v2, v2.Type())
}

// DeepEqual reports whether x and y are “deeply equal,” defined as follows.
// Two values of identical type are deeply equal if one of the following cases applies.
// Values of distinct types are never deeply equal.
//
// Array values are deeply equal when their corresponding elements are deeply equal.
//
// Struct values are deeply equal if their corresponding fields,
// both exported and unexported, are deeply equal.
//
// Func values are deeply equal if both are nil; otherwise they are not deeply equal.
//
// Interface values are deeply equal if they hold deeply equal concrete values.
//
// Map values are deeply equal when all of the following are true:
// they are both nil or both non-nil, they have the same length,
// and either they are the same map object or their corresponding keys
// (matched using Go equality) map to deeply equal values.
//
// Pointer values are deeply equal if they are equal using Go's == operator
// or if they point to deeply equal values.
//
// Slice values are deeply equal when all of the following are true:
// they are both nil or both non-nil, they have the same length,
// and either they point to the same initial entry of the same underlying array
// (that is, &x[0] == &y[0]) or their corresponding elements (up to length) are deeply equal.
// Note that a non-nil empty slice and a nil slice (for example, []byte{} and []byte(nil))
// are not deeply equal.
//
// Other values - numbers, bools, strings, and channels - are deeply equal
// if they are equal using Go's == operator.
//
// In general DeepEqual is a recursive relaxation of Go's == operator.
// However, this idea is impossible to implement without some inconsistency.
// Specifically, it is possible for a value to be unequal to itself,
// either because it is of func type (uncomparable in general)
// or because it is a floating-point NaN value (not equal to itself in floating-point comparison),
// or because it is an array, struct, or interface containing
// such a value.
// On the other hand, pointer values are always equal to themselves,
// even if they point at or contain such problematic values,
// because they compare equal using Go's == operator, and that
// is a sufficient condition to be deeply equal, regardless of content.
// DeepEqual has been defined so that the same short-cut applies
// to slices and maps: if x and y are the same slice or the same map,
// they are deeply equal regardless of content.
//
// As DeepEqual traverses the data values it may find a cycle. The
// second and subsequent times that DeepEqual compares two pointer
// values that have been compared before, it treats the values as
// equal rather than examining the values to which they point.
// This ensures that DeepEqual terminates.
func DeepEqual(x, y interface{}) (bool, error) {
	if x == nil || y == nil {
		if x == y {
			return true, nil
		}
		return false, fmt.Errorf("Only one value is nil (%v vs %v)", x, y)
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false, fmt.Errorf("Types %v and %v do not match", v1.Type(), v2.Type())
	}
	return deepValueEqual(v1, v2, make(map[visit]bool), 0)
}
