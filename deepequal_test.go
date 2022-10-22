package deepequal

import "testing"

type SimpleNestedStruct struct {
	Str         string
	Float32     float32
	hiddenField int64
}

type SimpleStruct struct {
	Int32          int32
	Int64          int64
	Nested         SimpleNestedStruct
	NestedPtr      *SimpleNestedStruct
	NestedSlice    []SimpleNestedStruct
	NestedPtrSlice []*SimpleNestedStruct
	hiddenField    float32
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if equal, err := DeepEqual(a, b); !equal {
		t.Errorf("not equal: %v", err)
	}
}

func TestSimpleStruct(t *testing.T) {
	assertEqual(t, SimpleStruct{}, SimpleStruct{})
	assertEqual(t, SimpleStruct{Nested: SimpleNestedStruct{Float32: 42}}, SimpleStruct{Nested: SimpleNestedStruct{Float32: 42}})
	assertEqual(t, SimpleStruct{hiddenField: -1}, SimpleStruct{})
}
