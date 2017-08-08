package degob

import "fmt"

type WireType struct {
	ArrayT  *ArrayType
	SliceT  *SliceType
	StructT *StructType
	MapT    *MapType
}

type ArrayType struct {
	CommonType
	Elem           typeId
	ElemTypeString string
	Len            int
}
type CommonType struct {
	Name string // the name of the struct type
	Id   int    // the id of the type, repeated so it's inside the type
}
type SliceType struct {
	CommonType
	Elem           typeId
	ElemTypeString string
}
type StructType struct {
	CommonType
	Field []*FieldType // the fields of the struct.
}
type FieldType struct {
	Name       string // the name of the field.
	TypeString string
	Id         int // the type id of the field, which must be already defined
}
type MapType struct {
	CommonType
	Key            typeId
	KeyTypeString  string
	Elem           typeId
	ElemTypeString string
}

const (
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
	// gob specific types
)

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

type _bool_type bool

func (v _bool_type) Equal(o Value) bool {
	ov, ok := o.(_bool_type)
	if !ok {
		return false
	}
	return v == ov
}

func (v _bool_type) Display(sty style) string {
	return fmt.Sprintf("%v", bool(v))
}

type _int_type int64

func (v _int_type) Equal(o Value) bool {
	ov, ok := o.(_int_type)
	if !ok {
		return false
	}
	return v == ov
}

func (v _int_type) Display(sty style) string {
	return fmt.Sprintf("%v", int64(v))
}

type _uint_type uint64

func (v _uint_type) Equal(o Value) bool {
	ov, ok := o.(_uint_type)
	if !ok {
		return false
	}
	return v == ov
}

func (v _uint_type) Display(sty style) string {
	return fmt.Sprintf("%v", uint64(v))
}

type _float_type float64

func (v _float_type) Equal(o Value) bool {
	ov, ok := o.(_float_type)
	if !ok {
		return false
	}
	return v == ov
}

func (v _float_type) Display(sty style) string {
	return fmt.Sprintf("%v", float64(v))
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

func (v _bytes_type) Display(sty style) string {
	return fmt.Sprintf("%#v", []byte(v))
}

type _string_type string

func (v _string_type) Equal(o Value) bool {
	ov, ok := o.(_string_type)
	if !ok {
		return false
	}
	return v == ov
}

func (v _string_type) Display(sty style) string {
	return fmt.Sprintf("\"%v\"", string(v))
}

type _complex_type complex128

func (v _complex_type) Equal(o Value) bool {
	ov, ok := o.(_complex_type)
	if !ok {
		return false
	}
	return v == ov
}

func (v _complex_type) Display(sty style) string {
	return fmt.Sprintf("%#v", complex128(v))
}

type _interface_type struct {
	name  string
	value Value
}

func (v _interface_type) Equal(o Value) bool {
	ov, ok := o.(_interface_type)
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

func (v _interface_type) Display(sty style) string {
	return fmt.Sprintf("%v", v.value.Display(sty))
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
		var v _interface_type
		return v
	default:
		panic("unknown id")
	}
}

func isBuiltin(id typeId) bool {
	return (id >= 1) && (id <= 8)
}

// ???
type interfaceValue struct{}

func (v interfaceValue) Display(sty style) string {
	return ""
}

// Value is any displayable value
type Value interface {
	Equal(Value) bool
	Display(sty style) string
}
