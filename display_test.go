package degob

import "testing"

func cmp(got string, ex string, t *testing.T) {
	if got != ex {
		t.Fatalf("expected %s got %s", ex, got)
	}
}

func TestDisplayArray(t *testing.T) {
	v := WireType{
		ArrayT: &ArrayType{
			CommonType: CommonType{
				Name: "",
				Id:   12345,
			},
			Elem:           _string_id,
			ElemTypeString: "string",
			Len:            10,
		},
	}
	out := v.String()
	cmp(out, "// [10]string", t)

	v.ArrayT.CommonType.Name = "Foo"
	out = v.String()
	cmp(out, "type Foo [10]string", t)
}

func TestDisplaySlice(t *testing.T) {
	v := WireType{
		SliceT: &SliceType{
			CommonType: CommonType{
				Name: "",
				Id:   12345,
			},
			Elem:           _string_id,
			ElemTypeString: "string",
		},
	}
	out := v.String()
	cmp(out, "// []string", t)

	v.SliceT.CommonType.Name = "Foo"
	out = v.String()
	cmp(out, "type Foo []string", t)
}

func TestDisplayMap(t *testing.T) {
	w := WireType{
		MapT: &MapType{
			ElemTypeString: "string",
			KeyTypeString:  "complex128",
		},
	}
	out := w.String()
	cmp(out, "// map[complex128]string", t)

	w.MapT.CommonType.Name = "Foo"
	out = w.String()
	cmp(out, "type Foo map[complex128]string", t)
}

func TestDisplayStruct(t *testing.T) {
	w := WireType{
		StructT: &StructType{
			CommonType: CommonType{
				Name: "Foo",
			},
			Field: []*FieldType{
				&FieldType{
					Name:       "X",
					TypeString: "int64",
				},
				&FieldType{
					Name:       "Y",
					TypeString: "string",
				},
			},
		},
	}
	out := w.String()
	expected := "type Foo struct {\n\tX int64\n\tY string\n}"
	cmp(out, expected, t)
}

func TestDisplayStructVal(t *testing.T) {
	v := structValue{
		name:   "Foo",
		sorted: true,
		fields: structFields{
			structField{name: "Complex", value: _complex_type(1 + 2i)},
			structField{name: "String", value: _string_type("1 + 2i")},
		},
	}

	out := v.Display(SingleLine)
	cmp(out, "Foo{Complex: (1+2i), String: \"1 + 2i\"}", t)
	out = v.Display(CommentedSingleLine)
	cmp(out, "//Foo{Complex: (1+2i), String: \"1 + 2i\"}", t)
}

func TestDisplayMapVal(t *testing.T) {
	v := mapValue{
		keyType:  "string",
		elemType: "int64",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("foo"),
				elem: _int_type(12),
			},
			mapEntry{
				key:  _string_type("bar"),
				elem: _int_type(-10),
			},
		},
	}
	out := v.Display(SingleLine)
	cmp(out, "map[string]int64{\"foo\": 12,\"bar\": -10}", t)
	out = v.Display(CommentedSingleLine)
	cmp(out, "//map[string]int64{\"foo\": 12,\"bar\": -10}", t)
}

func TestDisplayArrayVal(t *testing.T) {
	v := arrayValue{
		length:   2,
		elemType: "string",
		values: []Value{
			_string_type("one"),
			_string_type("two"),
		},
	}

	out := v.Display(SingleLine)
	cmp(out, "[2]string{\"one\", \"two\"}", t)
	out = v.Display(CommentedSingleLine)
	cmp(out, "//[2]string{\"one\", \"two\"}", t)
}

func TestDisplaySliceValue(t *testing.T) {
	v := sliceValue{
		elemType: "[]byte",
		values: []Value{
			_bytes_type([]byte{0x30, 0x31}),
			_bytes_type([]byte{0x32, 0x33, 0x34}),
		},
	}

	out := v.Display(SingleLine)
	cmp(out, "[][]byte{[]byte{0x30, 0x31}, []byte{0x32, 0x33, 0x34}}", t)
	out = v.Display(CommentedSingleLine)
	cmp(out, "//[][]byte{[]byte{0x30, 0x31}, []byte{0x32, 0x33, 0x34}}", t)
}
