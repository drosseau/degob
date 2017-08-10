package degob

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

type mapValue struct {
	keyType  string
	elemType string
	values   map[Value]Value
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
	for k, v := range v.values {
		ov, ok := ov.values[k]
		if !ok {
			return false
		}
		if !v.Equal(ov) {
			return false
		}
	}
	return true
}

type structValue struct {
	name   string
	fields map[string]Value
}

func (s *structValue) Equal(o Value) bool {
	v, ok := o.(*structValue)
	if !ok {
		return false
	}
	if s.name != v.name {
		return false
	}
	for k, val := range s.fields {
		ov, ok := v.fields[k]
		if !ok {
			return false
		}
		if !val.Equal(ov) {
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
	eq := v.value.Equal(ov.value)
	if !eq {
	}
	return eq
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
		var v interfaceValue
		return v
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
