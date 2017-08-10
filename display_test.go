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

// since these are a map the order isn't deterministic so this
// can fail
/*
func TestDisplayStructVal(t *testing.T) {
	v := structValue{
		name: "Foo",
		fields: map[string]Value{
			"Complex": _complex_type(1 + 2i),
			"String":  _string_type("1 + 2i"),
		},
	}

	out := v.Display(SingleLine)
	cmp(out, "Foo{Complex: (1+2i), String: \"1 + 2i\"}", t)
}
*/
