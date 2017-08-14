package degob

import "testing"

func testEqual(v Value, eq Value, t *testing.T) {
	if !v.Equal(eq) {
		t.Fatalf("expected values to be equal but they weren't: %v != %v", v, eq)
	}
}

func testNotEqual(v Value, notEq Value, t *testing.T) {
	if v.Equal(notEq) {
		t.Fatalf("expected values to not be equal but they were: %v == %v", v, notEq)
	}
}

func TestEqualStructValue(t *testing.T) {
	v := structValue{
		name:   "TestStructValue",
		sorted: true,
		fields: structFields{
			structField{name: "Bar", value: _int_type(10)},
			structField{name: "Foo", value: _string_type("foo")},
		},
	}
	eq := structValue{
		name:   "TestStructValue",
		sorted: true,
		fields: structFields{
			structField{name: "Bar", value: _int_type(10)},
			structField{name: "Foo", value: _string_type("foo")},
		},
	}

	testEqual(&v, &eq, t)

	notEqName := structValue{
		name:   "NotEqualName",
		sorted: true,
		fields: structFields{
			structField{name: "Bar", value: _int_type(10)},
			structField{name: "Foo", value: _string_type("foo")},
		},
	}

	testNotEqual(&v, &notEqName, t)

	notEqFieldName := structValue{
		name:   "TestStructValue",
		sorted: true,
		fields: structFields{
			structField{name: "Baz", value: _int_type(10)},
			structField{name: "Foo", value: _string_type("foo")},
		},
	}

	testNotEqual(&v, &notEqFieldName, t)

	notEqFieldValue := structValue{
		name:   "TestStructValue",
		sorted: true,
		fields: structFields{
			structField{name: "Bar", value: _int_type(20)},
			structField{name: "Foo", value: _string_type("foo")},
		},
	}

	testNotEqual(&v, &notEqFieldValue, t)

	notEqFieldNum := structValue{
		name:   "TestStructValue",
		sorted: true,
		fields: structFields{
			structField{name: "Bar", value: _int_type(20)},
		},
	}

	testNotEqual(&v, &notEqFieldNum, t)

	notEqDiffType := _string_type("Hi")
	testNotEqual(&v, &notEqDiffType, t)
}

func TestEqualArrayValue(t *testing.T) {
	v := arrayValue{
		elemType: "complex128",
		length:   3,
		values: []Value{
			_complex_type(1 + 2i),
			_complex_type(1 - 2i),
			_complex_type(0),
		},
	}

	eq := v
	notEqType := arrayValue{
		elemType: "string",
		length:   3,
		values: []Value{
			_complex_type(1 + 2i),
			_complex_type(1 - 2i),
			_complex_type(0),
		},
	}
	notEqLen := arrayValue{
		elemType: "complex128",
		length:   2,
		values: []Value{
			_complex_type(1 + 2i),
			_complex_type(1 - 2i),
		},
	}
	notEqElems := arrayValue{
		elemType: "complex128",
		length:   3,
		values: []Value{
			_complex_type(1 - 2i),
			_complex_type(1 + 2i),
			_complex_type(0),
		},
	}

	notEqDiffType := _string_type("hi")

	testEqual(&v, &eq, t)
	testNotEqual(&v, &notEqElems, t)
	testNotEqual(&v, &notEqLen, t)
	testNotEqual(&v, &notEqDiffType, t)
	testNotEqual(&v, &notEqType, t)
}

func TestEqualMapValue(t *testing.T) {
	v := mapValue{
		elemType: "float64",
		keyType:  "string",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("foo"),
				elem: _float_type(12.4),
			},
			mapEntry{
				key:  _string_type("bar"),
				elem: _float_type(-3.14),
			},
		},
	}
	eq := v

	notEqLen := mapValue{
		elemType: "float64",
		keyType:  "string",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("bar"),
				elem: _float_type(-3.14),
			},
		},
	}

	notEqElems := mapValue{
		elemType: "float64",
		keyType:  "string",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("foo"),
				elem: _float_type(12.4),
			},
			mapEntry{
				key:  _string_type("bar"),
				elem: _float_type(3.14),
			},
		},
	}
	notEqKeys := mapValue{
		elemType: "float64",
		keyType:  "string",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("baz"),
				elem: _float_type(12.4),
			},
			mapEntry{
				key:  _string_type("bar"),
				elem: _float_type(3.14),
			},
		},
	}
	notEqElemType := mapValue{
		elemType: "complex128",
		keyType:  "string",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("foo"),
				elem: _float_type(12.4),
			},
			mapEntry{
				key:  _string_type("bar"),
				elem: _float_type(3.14),
			},
		},
	}
	notEqKeyType := mapValue{
		elemType: "float64",
		keyType:  "uint64",
		values: []mapEntry{
			mapEntry{
				key:  _string_type("foo"),
				elem: _float_type(12.4),
			},
			mapEntry{
				key:  _string_type("bar"),
				elem: _float_type(3.14),
			},
		},
	}

	testEqual(&v, &eq, t)
	testNotEqual(&v, &notEqElems, t)
	testNotEqual(&v, &notEqKeys, t)
	testNotEqual(&v, &notEqLen, t)
	testNotEqual(&v, &notEqElemType, t)
	testNotEqual(&v, &notEqKeyType, t)
	testNotEqual(&v, _bool_type(false), t)
}

func TestEqualSliceValue(t *testing.T) {

	v := sliceValue{
		elemType: "Foo",
		values: []Value{
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(10),
					},
					structField{
						name:  "Y",
						value: _bool_type(false),
					},
				},
			},
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(100),
					},
					structField{
						name:  "Y",
						value: _bool_type(true),
					},
				},
			},
		},
	}
	eq := v
	notEqElemTypeString := sliceValue{
		elemType: "BadType",
		values: []Value{
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(10),
					},
					structField{
						name:  "Y",
						value: _bool_type(false),
					},
				},
			},
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(100),
					},
					structField{
						name:  "Y",
						value: _bool_type(true),
					},
				},
			},
		},
	}
	notEqLen := sliceValue{
		elemType: "Foo",
		values: []Value{
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(10),
					},
					structField{
						name:  "Y",
						value: _bool_type(false),
					},
				},
			},
		},
	}
	notEqElems := sliceValue{
		elemType: "Foo",
		values: []Value{
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(20),
					},
					structField{
						name:  "Y",
						value: _bool_type(false),
					},
				},
			},
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(100),
					},
					structField{
						name:  "Y",
						value: _bool_type(true),
					},
				},
			},
		},
	}
	notEqBadDefinition1 := sliceValue{
		elemType: "Foo",
		values: []Value{
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "Y",
						value: _uint_type(10),
					},
					structField{
						name:  "Z",
						value: _bool_type(false),
					},
				},
			},
			&structValue{
				name: "Foo",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(100),
					},
					structField{
						name:  "Y",
						value: _bool_type(true),
					},
				},
			},
		},
	}
	notEqBadDefinition2 := sliceValue{
		elemType: "Foo",
		values: []Value{
			&structValue{
				name: "Baz",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(10),
					},
					structField{
						name:  "Y",
						value: _bool_type(false),
					},
				},
			},
			&structValue{
				name: "Baz",
				fields: structFields{
					structField{
						name:  "X",
						value: _uint_type(100),
					},
					structField{
						name:  "Y",
						value: _bool_type(true),
					},
				},
			},
		},
	}

	testEqual(&v, &eq, t)
	testNotEqual(&v, &notEqElems, t)
	testNotEqual(&v, &notEqLen, t)
	testNotEqual(&v, &notEqBadDefinition1, t)
	testNotEqual(&v, &notEqBadDefinition2, t)
	testNotEqual(&v, &notEqElemTypeString, t)
	testNotEqual(&v, _complex_type(1), t)
}
