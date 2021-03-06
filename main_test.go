package degob

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type Inner struct {
	A float64
	B complex64
	C []byte
}

type Test struct {
	W Inner
	X int
	Y uint
	Z string
}

type ArrayInner struct {
	Float float64
	Int   int
}

type SliceInner struct {
	Uint uint
	Byte byte
}

type KeyType complex128
type ElemType struct {
	Complex complex64
	Float   float64
}

type AllPointers struct {
	X *int
	Y *string
	Z *bool
	Q *interface{}
}

type UserArray [2]int
type UserSlice []uint
type UserMap map[string]interface{}

type degobTestObject struct {
	fileName string
	item     interface{}
	expected *Gob
}

const (
	unpredictableId typeId = -100000
)

var testObjects = []degobTestObject{
	degobTestObject{
		fileName: "nestedstructfull.bin",
		item: Test{
			W: Inner{
				A: 3.14,
				B: 5 + 3i,
				C: []byte{1, 2, 3, 4, 5},
			},
			X: -10,
			Y: 10,
			Z: "Hello",
		},
		expected: &Gob{
			Value: &structValue{
				name: "Test",
				fields: structFields{
					structField{
						name: "W",
						value: &structValue{
							name: "Inner",
							fields: structFields{
								structField{name: "A", value: _float_type(3.14)},
								structField{name: "B", value: _complex_type(5 + 3i)},
								structField{name: "C", value: _bytes_type([]byte{1, 2, 3, 4, 5})},
							},
						},
					},
					structField{name: "X", value: _int_type(-10)},
					structField{name: "Y", value: _uint_type(10)},
					structField{name: "Z", value: _string_type("Hello")},
				},
			},
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "Test",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{Name: "W", Id: int(unpredictableId)},
							&FieldType{Name: "X", Id: int(_int_id)},
							&FieldType{Name: "Y", Id: int(_uint_id)},
							&FieldType{Name: "Z", Id: int(_string_id)},
						},
					},
				},
				unpredictableId - 1: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "Inner",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{Name: "A", Id: int(_float_id)},
							&FieldType{Name: "B", Id: int(_complex_id)},
							&FieldType{Name: "C", Id: int(_bytes_id)},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "nestedstructempty.bin",
		item:     &Test{},
		expected: &Gob{
			Value: &structValue{
				name: "Test",
				fields: structFields{
					structField{
						name: "W",
						value: &structValue{
							name: "Inner",
							fields: structFields{
								structField{name: "A", value: _float_type(0)},
								structField{name: "B", value: _complex_type(0)},
								structField{name: "C", value: _bytes_type([]byte{})},
							},
						},
					},
					structField{name: "X", value: _int_type(0)},
					structField{name: "Y", value: _uint_type(0)},
					structField{name: "Z", value: _string_type("")},
				},
			},
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "Test",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{Name: "W", Id: int(unpredictableId)},
							&FieldType{Name: "X", Id: int(_int_id)},
							&FieldType{Name: "Y", Id: int(_uint_id)},
							&FieldType{Name: "Z", Id: int(_string_id)},
						},
					},
				},
				unpredictableId - 1: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "Inner",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{Name: "A", Id: int(_float_id)},
							&FieldType{Name: "B", Id: int(_complex_id)},
							&FieldType{Name: "C", Id: int(_bytes_id)},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "arraybuiltin.bin",
		item:     [5]int{-2, -1, 0, 1, 2},
		expected: &Gob{
			Value: &arrayValue{
				length:   5,
				elemType: "int64",
				values: []Value{
					_int_type(-2),
					_int_type(-1),
					_int_type(0),
					_int_type(1),
					_int_type(2),
				},
			},
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					ArrayT: &ArrayType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Len:  5,
						Elem: _int_id,
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "slicebuiltin.bin",
		item: []string{
			"one", "two", "three",
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					SliceT: &SliceType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Elem: _string_id,
					},
				},
			},
			Value: &sliceValue{
				elemType: _string_id.name(),
				values: []Value{
					_string_type("one"),
					_string_type("two"),
					_string_type("three"),
				},
			},
		},
	},
	degobTestObject{
		fileName: "mapbuiltin.bin",
		item: map[string]float64{
			"one point two":           1.2,
			"negative ten point five": -10.5,
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					MapT: &MapType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Key:  _string_id,
						Elem: _float_id,
					},
				},
			},
			Value: &mapValue{
				keyType:  _string_id.name(),
				elemType: _float_id.name(),
				values: []mapEntry{
					mapEntry{
						key:  _string_type("one point two"),
						elem: _float_type(1.2),
					},
					mapEntry{
						key:  _string_type("negative ten point five"),
						elem: _float_type(-10.5),
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "arraystruct.bin",
		item: &[3]ArrayInner{
			ArrayInner{1.5, 10},
			ArrayInner{-1.5, -10},
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					ArrayT: &ArrayType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Len:  3,
						Elem: unpredictableId,
					},
				},
				unpredictableId - 1: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "Anon70",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{
								Name:       "Float",
								TypeString: "float64",
								Id:         int(_float_id),
							},
							&FieldType{
								Name:       "Int",
								TypeString: "int64",
								Id:         int(_int_id),
							},
						},
					},
				},
			},
			Value: &arrayValue{
				elemType: "Anon70",
				length:   3,
				values: []Value{
					&structValue{
						name: "Anon70",
						fields: structFields{
							structField{name: "Float", value: _float_type(1.5)},
							structField{name: "Int", value: _int_type(10)},
						},
					},
					&structValue{
						name: "Anon70",
						fields: structFields{
							structField{name: "Float", value: _float_type(-1.5)},
							structField{name: "Int", value: _int_type(-10)},
						},
					},
					&structValue{
						name: "Anon70",
						fields: structFields{
							structField{name: "Float", value: _float_type(0.0)},
							structField{name: "Int", value: _int_type(0)},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "slicestruct.bin",
		item: &[]SliceInner{
			SliceInner{0, 0x30},
			SliceInner{5, 0x35},
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "SliceInner",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{
								Name: "Uint",
								Id:   int(_uint_id),
							},
							&FieldType{
								Name: "Byte",
								Id:   int(_uint_id),
							},
						},
					},
				},
				unpredictableId - 1: &WireType{
					SliceT: &SliceType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Elem: unpredictableId,
					},
				},
			},
			Value: &sliceValue{
				elemType: "SliceInner",
				values: []Value{
					&structValue{
						name: "SliceInner",
						fields: structFields{
							structField{name: "Uint", value: _uint_type(0)},
							structField{name: "Byte", value: _uint_type(0x30)},
						},
					},
					&structValue{
						name: "SliceInner",
						fields: structFields{
							structField{name: "Uint", value: _uint_type(5)},
							structField{name: "Byte", value: _uint_type(0x35)},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "mapuserdefined.bin",
		item: &map[KeyType]ElemType{
			KeyType(5 - 2.1i):    ElemType{-2 + 3i, 10.2},
			KeyType(10.2 + 3.5i): ElemType{2 - 3i, -10.2},
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					MapT: &MapType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Key:  unpredictableId,
						Elem: unpredictableId,
					},
				},
				unpredictableId - 1: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "Anon74",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{
								Name: "Complex",
								Id:   int(_complex_id),
							},
							&FieldType{
								Name: "Float",
								Id:   int(_float_id),
							},
						},
					},
				},
			},
			Value: &mapValue{
				keyType:  "complex128",
				elemType: "Anon74",
				values: []mapEntry{
					mapEntry{
						key: _complex_type(5 - 2.1i),
						elem: &structValue{
							name: "Anon74",
							fields: structFields{
								structField{name: "Complex", value: _complex_type(-2 + 3i)},
								structField{name: "Float", value: _float_type(10.2)},
							},
						},
					},
					mapEntry{
						key: _complex_type(10.2 + 3.5i),
						elem: &structValue{
							name: "Anon74",
							fields: structFields{
								structField{name: "Complex", value: _complex_type(2 - 3i)},
								structField{name: "Float", value: _float_type(-10.2)},
							},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "string.bin",
		item:     "Hello there",
		expected: &Gob{
			Value: _string_type("Hello there"),
		},
	},
	degobTestObject{
		fileName: "complex.bin",
		item:     124.3 + 438.2i,
		expected: &Gob{
			Value: _complex_type(124.3 + 438.2i),
		},
	},
	degobTestObject{
		fileName: "float.bin",
		item:     12.5,
		expected: &Gob{
			Value: _float_type(12.5),
		},
	},
	degobTestObject{
		fileName: "ambiguousint.bin",
		item:     8,
		expected: &Gob{
			Value: _int_type(8),
		},
	},
	degobTestObject{
		fileName: "uint64.bin",
		item:     uint64(0xFFFFFFFFFFFFFFFF),
		expected: &Gob{
			Value: _uint_type(0xFFFFFFFFFFFFFFFF),
		},
	},
	degobTestObject{
		fileName: "booltrue.bin",
		item:     true,
		expected: &Gob{
			Value: _bool_type(true),
		},
	},
	degobTestObject{
		fileName: "boolfalse.bin",
		item:     false,
		expected: &Gob{
			Value: _bool_type(false),
		},
	},
	degobTestObject{
		fileName: "interfaceswithnil.bin",
		item: map[interface{}]interface{}{
			"StringToBool": false,
			"StringToInt":  12,
			1234:           "IntToString",
			5732:           nil,
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					MapT: &MapType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Key:            _interface_id,
						KeyTypeString:  "interface{}",
						Elem:           _interface_id,
						ElemTypeString: "interface{}",
					},
				},
			},
			Value: &mapValue{
				keyType:  "interface{}",
				elemType: "interface{}",
				values: []mapEntry{
					mapEntry{
						key: interfaceValue{
							name:  "string",
							value: _string_type("StringToBool"),
						},
						elem: interfaceValue{
							name:  "bool",
							value: _bool_type(false),
						},
					},
					mapEntry{
						key: interfaceValue{
							name:  "string",
							value: _string_type("StringToInt"),
						},
						elem: interfaceValue{
							name:  "int",
							value: _int_type(12),
						},
					},
					mapEntry{
						key: interfaceValue{
							name:  "int",
							value: _int_type(1234),
						},
						elem: interfaceValue{
							name:  "string",
							value: _string_type("IntToString"),
						},
					},
					mapEntry{
						key: interfaceValue{
							name:  "int",
							value: _int_type(5732),
						},
						elem: interfaceValue{
							value: _nil_value{},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "allpointers.bin",
		item:     newAllPointers(),
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "AllPointers",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{
								TypeString: "int64",
								Id:         int(_int_id),
								Name:       "X",
							},
							&FieldType{
								TypeString: "string",
								Id:         int(_string_id),
								Name:       "Y",
							},
							&FieldType{
								TypeString: "bool",
								Id:         int(_bool_id),
								Name:       "Z",
							},
							&FieldType{
								TypeString: "interface{}",
								Id:         int(_interface_id),
								Name:       "Q",
							},
						},
					},
				},
			},
			Value: &structValue{
				name: "AllPointers",
				fields: structFields{
					structField{name: "X", value: _int_type(10)},
					structField{name: "Y", value: _string_type("string pointer")},
					structField{name: "Z", value: _bool_type(true)},
					structField{
						name: "Q",
						value: interfaceValue{
							name:  "uint",
							value: _uint_type(80),
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "interfacemap.bin",
		item: map[interface{}]interface{}{
			"StringToBool":     false,
			"StringToInt":      12,
			1234:               "IntToString",
			ArrayInner{1.2, 1}: SliceInner{10, 0x04},
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					MapT: &MapType{
						CommonType: CommonType{
							Id: int(unpredictableId),
						},
						Key:            _interface_id,
						KeyTypeString:  "interface{}",
						Elem:           _interface_id,
						ElemTypeString: "interface{}",
					},
				},
				unpredictableId - 1: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "ArrayInner",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{
								Name:       "Float",
								TypeString: "float64",
								Id:         int(_float_id),
							},
							&FieldType{
								Name:       "Int",
								TypeString: "float64",
								Id:         int(_int_id),
							},
						},
					},
				},
				unpredictableId - 2: &WireType{
					StructT: &StructType{
						CommonType: CommonType{
							Name: "SliceInner",
							Id:   int(unpredictableId),
						},
						Field: []*FieldType{
							&FieldType{
								Name: "Uint",
								Id:   int(_uint_id),
							},
							&FieldType{
								Name: "Byte",
								Id:   int(_uint_id),
							},
						},
					},
				},
			},
			Value: &mapValue{
				keyType:  "interface{}",
				elemType: "interface{}",
				values: []mapEntry{
					mapEntry{
						key: interfaceValue{
							name:  "string",
							value: _string_type("StringToBool"),
						},
						elem: interfaceValue{
							name:  "bool",
							value: _bool_type(false),
						},
					},
					mapEntry{
						key: interfaceValue{
							name:  "string",
							value: _string_type("StringToInt"),
						},
						elem: interfaceValue{
							name:  "int",
							value: _int_type(12),
						},
					},
					mapEntry{
						key: interfaceValue{
							name:  "int",
							value: _int_type(1234),
						},
						elem: interfaceValue{
							name:  "string",
							value: _string_type("IntToString"),
						},
					},
					mapEntry{
						key: interfaceValue{
							name: "ArrayInner",
							value: &structValue{
								name: "ArrayInner",
								fields: structFields{
									structField{name: "Float", value: _float_type(1.2)},
									structField{name: "Int", value: _int_type(1)},
								},
							},
						},
						elem: interfaceValue{
							name: "SliceInner",
							value: &structValue{
								name: "SliceInner",
								fields: structFields{
									structField{name: "Uint", value: _uint_type(10)},
									structField{name: "Byte", value: _uint_type(0x04)},
								},
							},
						},
					},
				},
			},
		},
	},
	degobTestObject{
		fileName: "usermap.bin",
		item: UserMap{
			"hi":  int64(10),
			"bye": -10.4,
		},
		expected: &Gob{
			Types: map[typeId]*WireType{
				unpredictableId: &WireType{
					MapT: &MapType{
						CommonType: CommonType{
							Name: "UserMap",
							Id:   int(unpredictableId),
						},
						Key:            _string_id,
						KeyTypeString:  "string",
						Elem:           _interface_id,
						ElemTypeString: "interface{}",
					},
				},
			},
			Value: &mapValue{
				keyType:  "string",
				elemType: "interface{}",
				values: []mapEntry{
					mapEntry{
						key:  _string_type("hi"),
						elem: interfaceValue{name: "int64", value: _int_type(10)},
					},
					mapEntry{
						key:  _string_type("bye"),
						elem: interfaceValue{name: "float64", value: _float_type(-10.4)},
					},
				},
			},
		},
	},
}

func newAllPointers() AllPointers {
	x := new(int)
	*x = 10
	y := new(string)
	*y = "string pointer"
	z := new(bool)
	*z = true
	q := new(interface{})
	*q = uint(80)
	return AllPointers{x, y, z, q}
}

func compareGobs(expected *Gob, o *Gob, fname string, t *testing.T) {
	if expected.Types == nil {
		if o.Types != nil {
			t.Fatalf("expected nil `Types` but was non nil from gob in: `%s`\n%v", fname, o.Types)
		}
	} else {
		for _, wt := range expected.Types {
			found := false
			for _, owt := range o.Types {
				if compareWire(wt, owt) {
					found = true
					break
				}
			}
			if !found {
				if wt.SliceT != nil {
					t.Fatalf("expected Slice WireType not found, %v from gob in: `%s`\n", *wt.SliceT, fname)
				}
				if wt.StructT != nil {
					s := fmt.Sprintf("\nexpected WireType.StructT not found in %s", fname)
					s += fmt.Sprintf("%s\n\tName: %s\n\tFields: ", s, wt.StructT.CommonType.Name)
					for _, f := range wt.StructT.Field {
						s += fmt.Sprintf("%v ", *f)
					}
					s += "\n\nFound structs:\n"
					for _, owt := range o.Types {
						if owt.StructT != nil {
							s += fmt.Sprintf("\tName: %s\n\tFields:", owt.StructT.CommonType.Name)
							for _, f := range owt.StructT.Field {
								s += fmt.Sprintf("%v ", *f)
							}
							s += "\n"
						}
					}
					t.Fatal(s)
				}
				if wt.MapT != nil {
					t.Fatalf("expected Map WireType not found, %v from gob in: `%s`\n", *wt.MapT, fname)
				}
				if wt.ArrayT != nil {
					t.Fatalf("expected Array WireType not found, %v from gob in: `%s`\n", *wt.ArrayT, fname)
				}
			}
		}
	}
	if expected.Value == nil {
		if o.Value != nil {
			t.Fatalf("expected nil `Values` but was non nil from gob in: `%s`\n%v", fname, o.Value)
		}
	} else {
		v := expected.Value
		ov := o.Value
		if !v.Equal(ov) {
			s := fmt.Sprintf("expected Value not found for gob in %s:\n\t%v\n\t%s\nFound value:\n\t%v\n\t%s", fname, v, v.Display(SingleLine), ov, ov.Display(SingleLine))
			t.Fatal(s)
		}
	}
}

func compareWire(expected *WireType, got *WireType) bool {
	if expected.StructT != nil {
		if compareStruct(expected.StructT, got.StructT) {
			return true
		}
		return false
	}
	if expected.ArrayT != nil {
		if compareArray(expected.ArrayT, got.ArrayT) {
			return true
		}
		return false
	}
	if expected.SliceT != nil {
		if compareSlice(expected.SliceT, got.SliceT) {
			return true
		}
		return false
	}
	if expected.MapT != nil {
		if compareMap(expected.MapT, got.MapT) {
			return true
		}
		return false
	}
	return false
}

func compareMap(expected *MapType, got *MapType) bool {
	if got == nil {
		return false
	}

	if expected.CommonType.Name != got.CommonType.Name {
		return false
	}

	if expected.Key > unpredictableId {
		if expected.Key != got.Key {
			return false
		}
	}
	if expected.Elem > unpredictableId {
		if expected.Elem != got.Elem {
			return false
		}
	}
	return true
}

func compareSlice(expected *SliceType, got *SliceType) bool {
	if got == nil {
		return false
	}

	if expected.CommonType.Name != got.CommonType.Name {
		return false
	}

	if expected.Elem > unpredictableId {
		if expected.Elem != got.Elem {
			return false
		}
	}
	return true
}

func compareArray(expected *ArrayType, got *ArrayType) bool {
	if got == nil {
		return false
	}

	if expected.CommonType.Name != got.CommonType.Name {
		return false
	}

	if expected.Len != got.Len {
		return false
	}

	if expected.Elem > unpredictableId {
		if expected.Elem != got.Elem {
			return false
		}
	}
	return true
}

func compareStruct(expected *StructType, got *StructType) bool {
	if got == nil {
		return false
	}
	if expected.CommonType.Name != got.CommonType.Name {
		return false
	}
	if len(expected.Field) != len(got.Field) {
		return false
	}
	fieldsMatch := true
	for _, f := range expected.Field {
		foundField := false
		for _, of := range got.Field {
			if f.Name != of.Name {
				continue
			}
			if f.Id > int(unpredictableId) {
				if f.Id != of.Id {
					continue
				}
			}
			foundField = true
			break
		}
		if !foundField {
			fieldsMatch = false
			break
		}
	}
	return fieldsMatch
}

func TestMain(m *testing.M) {
	gob.Register(Inner{})
	gob.Register(Test{})
	gob.Register(ArrayInner{})
	gob.Register(SliceInner{})
	gob.Register(ElemType{})
	gob.Register(AllPointers{})
	gob.Register(KeyType(0))
	gob.Register(UserArray{})
	gob.Register([]UserSlice{})
	gob.Register(UserMap{})
	gob.Register(map[interface{}]interface{}{})
	for _, obj := range testObjects {
		fname := filepath.Join("test_examples", obj.fileName)
		_, err := os.Stat(fname)
		if err == nil || !os.IsNotExist(err) {
			continue
		}
		f, err := os.Create(fname)
		if err != nil {
			panic(err)
		}
		err = gob.NewEncoder(f).Encode(obj.item)
		if err != nil {
			panic(err)
		}
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func openFileTest(fname string, t *testing.T) *os.File {
	f, err := os.Open(filepath.Join("test_examples", fname))
	if err != nil {
		t.Fatalf("err: %v opening file: %s", err, fname)
	}
	return f
}

func fileToBufferTest(fname string, buf *bytes.Buffer, t *testing.T) {
	f := openFileTest(fname, t)
	_, err := buf.ReadFrom(f)
	f.Close()
	if err != nil {
		t.Fatalf("err: %v reading file: %s", err, fname)
	}
}
