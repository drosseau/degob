package degob

import (
	"fmt"
	"sort"
)

// WireType is basically the same as the `encoding/gob` wiretype with a few
// fields added and exported. It is also a `Stringer` so the types can be
// printed for better viewing.
type WireType struct {
	ArrayT  *ArrayType
	SliceT  *SliceType
	StructT *StructType
	MapT    *MapType
}

// ArrayType represents a fixed size array
type ArrayType struct {
	CommonType
	Elem           typeId
	ElemTypeString string
	Len            int
}

// CommonType has information common to all base types
type CommonType struct {
	Name string
	Id   int
}

// SliceType represents a slice
type SliceType struct {
	CommonType
	Elem           typeId
	ElemTypeString string
}
type StructType struct {
	CommonType
	Field []*FieldType
}

// FieldType is the information for a struct field
type FieldType struct {
	Name       string
	TypeString string
	Id         int
}

// MapType is the information for a map
type MapType struct {
	CommonType
	Key            typeId
	KeyTypeString  string
	Elem           typeId
	ElemTypeString string
}

const (
	// builtin types
	_bool_id      typeId = 1
	_int_id       typeId = 2
	_uint_id      typeId = 3
	_float_id     typeId = 4
	_bytes_id     typeId = 5
	_string_id    typeId = 6
	_complex_id   typeId = 7
	_interface_id typeId = 8
	// reserved types
	_reserved1_id typeId = 9
	_reserved2_id typeId = 10
	_reserved3_id typeId = 11
	_reserved4_id typeId = 12
	_reserved5_id typeId = 13
	_reserved6_id typeId = 14
	_reserved7_id typeId = 15
)

func (t typeId) String() string {
	switch t {
	case _bool_id:
		return "1 (bool)"
	case _int_id:
		return "2 (int)"
	case _uint_id:
		return "3 (uint)"
	case _float_id:
		return "4 (float)"
	case _bytes_id:
		return "5 (bytes)"
	case _string_id:
		return "6 (string)"
	case _complex_id:
		return "7 (complex)"
	case _interface_id:
		return "8 (interface)"
	case _reserved1_id:
		return "9 (reserved)"
	case _reserved2_id:
		return "10 (reserved)"
	case _reserved3_id:
		return "11 (reserved)"
	case _reserved4_id:
		return "12 (reserved)"
	case _reserved5_id:
		return "13 (reserved)"
	case _reserved6_id:
		return "14 (reserved)"
	case _reserved7_id:
		return "15 (reserved)"
	default:
		return fmt.Sprintf("%d (user defined)", t)
	}
}

// Value is a reprsentation of a Go value. You can test for equality and
// Display them stylized
type Value interface {
	// Equal tests for equality to another value
	Equal(Value) bool
	// Display shows the value accoring to the chosen style
	Display(sty style) string
}

type sliceValue struct {
	elemType string
	values   []Value
}

func (v sliceValue) Equal(o Value) bool {
	ov, ok := o.(*sliceValue)
	if !ok {
		return false
	}
	if v.elemType != ov.elemType {
		return false
	}
	if len(v.values) != len(ov.values) {
		return false
	}
	for i, v := range v.values {
		if !v.Equal(ov.values[i]) {
			return false
		}
	}
	return true
}

type arrayValue struct {
	elemType string
	length   int
	values   []Value
}

func (v arrayValue) Equal(o Value) bool {
	ov, ok := o.(*arrayValue)
	if !ok {
		return false
	}

	if v.elemType != ov.elemType {
		return false
	}
	if v.length != ov.length {
		return false
	}
	for i, v := range v.values {
		if !v.Equal(ov.values[i]) {
			return false
		}
	}
	return true
}

type mapEntry struct {
	key  Value
	elem Value
}

type mapValue struct {
	keyType  string
	elemType string
	values   []mapEntry
}

func (v mapValue) Equal(o Value) bool {
	ov, ok := o.(*mapValue)
	if !ok {
		return false
	}
	if v.keyType != ov.keyType {
		return false
	}
	if v.elemType != ov.elemType {
		return false
	}
	if len(v.values) != len(ov.values) {
		return false
	}
	// TODO: At some point it'd be nice to sort `values` for each one to make
	// this comparison more efficient. The problem is how I would define `Less`
	// for Values that aren't the same type. I'd probably have to resort
	// to comparing TypeIds but then I'd have to give each Value a TypeId method
	// so I'm saving it for some other time.
	for _, vval := range v.values {
		found := false
		for _, ovval := range ov.values {
			if vval.key.Equal(ovval.key) && vval.elem.Equal(ovval.elem) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type structField struct {
	name  string
	value Value
}

type structFields []structField

func (s structFields) Len() int           { return len(s) }
func (s structFields) Less(i, j int) bool { return s[i].name < s[j].name }
func (s structFields) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type structValue struct {
	name   string
	sorted bool
	fields structFields
}

func (s *structValue) Equal(o Value) bool {
	v, ok := o.(*structValue)
	if !ok {
		return false
	}

	if s.name != v.name {
		return false
	}

	if len(s.fields) != len(v.fields) {
		return false
	}

	if !s.sorted {
		sort.Sort(s.fields)
		s.sorted = true
	}
	if !v.sorted {
		sort.Sort(v.fields)
		v.sorted = true
	}

	for fno := range s.fields {
		sf := s.fields[fno]
		vf := v.fields[fno]
		if sf.name != vf.name {
			return false
		}
		if !sf.value.Equal(vf.value) {
			return false
		}
	}
	return true
}

type _bool_type bool

func (v _bool_type) Equal(o Value) bool {
	ov, ok := o.(_bool_type)
	if !ok {
		return false
	}
	return v == ov
}

type _int_type int64

func (v _int_type) Equal(o Value) bool {
	ov, ok := o.(_int_type)
	if !ok {
		return false
	}
	return v == ov
}

type _uint_type uint64

func (v _uint_type) Equal(o Value) bool {
	ov, ok := o.(_uint_type)
	if !ok {
		return false
	}
	return v == ov
}

type _float_type float64

func (v _float_type) Equal(o Value) bool {
	ov, ok := o.(_float_type)
	if !ok {
		return false
	}
	return v == ov
}

type _bytes_type []byte

func (v _bytes_type) Equal(o Value) bool {
	ov, ok := o.(_bytes_type)
	if !ok {
		return false
	}
	if len(v) != len(ov) {
		return false
	}

	for i, b := range v {
		if b != ov[i] {
			return false
		}
	}
	return true
}

type _string_type string

func (v _string_type) Equal(o Value) bool {
	ov, ok := o.(_string_type)
	if !ok {
		return false
	}
	return v == ov
}

type _complex_type complex128

func (v _complex_type) Equal(o Value) bool {
	ov, ok := o.(_complex_type)
	if !ok {
		return false
	}
	return v == ov
}

type _nil_value struct{}

func (v _nil_value) Equal(o Value) bool {
	_, ok := o.(_nil_value)
	return ok
}

type interfaceValue struct {
	name  string
	value Value
}

func (v interfaceValue) Equal(o Value) bool {
	ov, ok := o.(interfaceValue)
	if !ok {
		return false
	}
	if v.name != ov.name {
		return false
	}
	return v.value.Equal(ov.value)
}

func valueFor(id typeId) Value {
	switch id {
	case _bool_id:
		var v _bool_type
		return v
	case _int_id:
		var v _int_type
		return v
	case _uint_id:
		var v _uint_type
		return v
	case _float_id:
		var v _float_type
		return v
	case _bytes_id:
		var v _bytes_type
		return v
	case _string_id:
		var v _string_type
		return v
	case _complex_id:
		var v _complex_type
		return v
	case _interface_id:
		return interfaceValue{value: _nil_value{}}
	default:
		panic("unknown id")
	}
}

func isBuiltin(id typeId) bool {
	return (id >= 1) && (id <= 8)
}

func (t typeId) name() string {
	switch t {
	case _bool_id:
		return "bool"
	case _int_id:
		return "int64"
	case _uint_id:
		return "uint64"
	case _float_id:
		return "float64"
	case _bytes_id:
		return "[]byte"
	case _string_id:
		return "string"
	case _complex_id:
		return "complex128"
	case _interface_id:
		return "interface{}"
	default:
		panic("unknown id")
	}
}
