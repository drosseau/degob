package degob

// At the moment this is pretty meh, but I don't particularly
// want to spend too much time on it right now

import (
	"fmt"
	"strings"
)

type style uint8

const (
	// SingleLine tells the Value to format its display as a Go type on a single line
	SingleLine = iota
	// Commented single line adds a // before SingleLine output. This is useful if
	// you'd like to stream Gobs to create a file of discovered types but don't
	// want to ignore values
	CommentedSingleLine
	//JSON tells the Value to format its display as JSON
	//
	// There are not always perfect mappings from Go types to JSON. Two examples
	// are maps that don't have string keys and complex numbers. Complex numbers
	// are displayed as a JSON object `{"Re": real(x), "Im": imag(x)}`. Non
	// string key maps are returned as an error object that looks like:
	// `{"error": STRING, "val": "SingleLine displayed val"}`.
	JSON
)

func (w *WireType) String() string {
	switch {
	case w.StructT != nil:
		return w.StructT.String()
	case w.SliceT != nil:
		return w.SliceT.String()
	case w.ArrayT != nil:
		return w.ArrayT.String()
	case w.MapT != nil:
		return w.MapT.String()
	default:
		return "Unset WireType"
	}
}

// Id will return -1 for all nil wire types
func (w *WireType) Id() int {
	switch {
	case w.StructT != nil:
		return w.StructT.Id
	case w.SliceT != nil:
		return w.SliceT.Id
	case w.ArrayT != nil:
		return w.ArrayT.Id
	case w.MapT != nil:
		return w.MapT.Id
	default:
		return -1
	}
}

func (a *ArrayType) String() string {
	if a.CommonType.Name != "" {
		return fmt.Sprintf("type %s [%d]%s", a.CommonType.Name, a.Len, a.ElemTypeString)
	}
	return fmt.Sprintf("// [%d]%s", a.Len, a.ElemTypeString)
}

func (s *SliceType) String() string {
	if s.CommonType.Name != "" {
		return fmt.Sprintf("type %s []%s", s.CommonType.Name, s.ElemTypeString)
	}
	return fmt.Sprintf("// []%s", s.ElemTypeString)
}

func (m *MapType) String() string {
	if m.CommonType.Name != "" {
		return fmt.Sprintf("type %s map[%s]%s", m.CommonType.Name, m.KeyTypeString, m.ElemTypeString)
	}
	return fmt.Sprintf("// map[%s]%s", m.KeyTypeString, m.ElemTypeString)
}

func (s *StructType) String() string {
	st := fmt.Sprintf("type %s struct {\n", s.CommonType.Name)
	nfields := len(s.Field)
	for i, f := range s.Field {
		if i < nfields-1 {
			st += fmt.Sprintf("\t%s %s\n", f.Name, f.TypeString)
		} else {
			st += fmt.Sprintf("\t%s %s\n}", f.Name, f.TypeString)
		}
	}
	return st
}

func (v sliceValue) valuesSep(sty style, sep string) string {
	var out []string
	for _, val := range v.values {
		out = append(out, val.Display(sty))
	}
	return strings.Join(out, sep)
}

func (v sliceValue) Display(sty style) string {
	switch sty {
	case JSON:
		return fmt.Sprintf("[%s]", v.valuesSep(sty, ", "))
	case SingleLine:
		return fmt.Sprintf("[]%s{%s}", v.elemType, v.valuesSep(SingleLine, ", "))
	case CommentedSingleLine:
		return fmt.Sprintf("//[]%s{%s}", v.elemType, v.valuesSep(SingleLine, ", "))
	default:
		panic("unimplemented style")
	}
}

func (v arrayValue) valuesSep(sty style, sep string) string {
	var out []string
	for _, val := range v.values {
		out = append(out, val.Display(sty))
	}
	return strings.Join(out, sep)
}

func (v arrayValue) Display(sty style) string {
	switch sty {
	case JSON:
		return fmt.Sprintf("[%s]", v.valuesSep(sty, ", "))
	case SingleLine:
		return fmt.Sprintf("[%d]%s{%s}", v.length, v.elemType, v.valuesSep(SingleLine, ", "))
	case CommentedSingleLine:
		return fmt.Sprintf("//[%d]%s{%s}", v.length, v.elemType, v.valuesSep(SingleLine, ", "))
	default:
		panic("unimplemented style")
	}
}

func (v arrayValue) getValues(sty style) string {
	s := ""
	for _, val := range v.values {
		s += val.Display(sty)
		// TODO: This is dependent on other stuff
		s += ",\n"
	}
	return s
}

func (v mapValue) getValues(sty style) string {
	var out string
	nval := len(v.values)
	i := 0
	for _, v := range v.values {
		if i < nval-1 {
			out += fmt.Sprintf("%s: %s,", v.key.Display(sty), v.elem.Display(sty))
		} else {
			out += fmt.Sprintf("%s: %s", v.key.Display(sty), v.elem.Display(sty))
		}
		i += 1
	}
	return out
}

func (v mapValue) Display(sty style) string {
	switch sty {
	case JSON:
		return v.displayJSON()
	case CommentedSingleLine:
		return fmt.Sprintf("//map[%s]%s{%s}", v.keyType, v.elemType, v.getValues(sty))
	case SingleLine:
		return fmt.Sprintf("map[%s]%s{%s}", v.keyType, v.elemType, v.getValues(sty))
	default:
		panic("unknown style requested")
	}
}

func (v mapValue) displayJSON() string {
	if v.keyType != "string" {
		return `{
	"error": "cannot display map type with non key strings as JSON"
	"val": "` + v.Display(SingleLine) + `"
}`
	}
	s := "{"
	end := len(v.values)
	for i, v := range v.values {
		s += fmt.Sprintf("%s: %s", v.key.Display(JSON), v.elem.Display(JSON))
		if i+1 < end {
			s += ", "
		}
	}
	s += "}"
	return s
}

func (s *structValue) Display(sty style) string {
	switch sty {
	case CommentedSingleLine:
		return s.commentedSingleLine()
	case SingleLine:
		return s.singleLine()
	case JSON:
		return s.json()
	default:
		panic("unknown style requested")
	}
}

func (s *structValue) commentedSingleLine() string {
	return fmt.Sprintf("//%s", s.singleLine())
}

func (s *structValue) getFieldVals(newline bool, sty style) string {
	var out string
	nfields := len(s.fields)
	i := 0
	for _, v := range s.fields {
		out += fmt.Sprintf("%s: %s", v.name, v.value.Display(sty))
		if i < nfields-1 {
			if newline {
				out += ",\n"
			} else {
				out += ", "
			}
		} else {
			if newline {
				out += ",\n"
			}
		}
		i += 1
	}
	return out
}

func (s *structValue) singleLine() string {
	return fmt.Sprintf("%s{%s}", s.name, s.getFieldVals(false, SingleLine))
}

func (s *structValue) json() string {
	str := "{"
	end := len(s.fields)
	for i, v := range s.fields {
		str += fmt.Sprintf("\"%s\": %s", v.name, v.value.Display(JSON))
		if i+1 < end {
			str += ", "
		}
	}
	str += "}"
	return str
}

// Base types

func (v _bool_type) Display(sty style) string {
	return fmt.Sprintf("%v", bool(v))
}
func (v _int_type) Display(sty style) string {
	return fmt.Sprintf("%v", int64(v))
}
func (v _uint_type) Display(sty style) string {
	return fmt.Sprintf("%v", uint64(v))
}
func (v _float_type) Display(sty style) string {
	return fmt.Sprintf("%v", float64(v))
}
func (v _bytes_type) Display(sty style) string {
	if sty == JSON {
		s := fmt.Sprintf("[% #02x]", v)
		return strings.Join(strings.Split(s, " "), ", ")
	}
	return fmt.Sprintf("%#v", []byte(v))
}
func (v _string_type) Display(sty style) string {
	return fmt.Sprintf("\"%v\"", string(v))
}
func (v _complex_type) Display(sty style) string {
	if sty == JSON {
		return fmt.Sprintf(`{"Re": %f, "Im": %f}`, real(v), imag(v))
	}
	return fmt.Sprintf("%#v", complex128(v))
}
func (v interfaceValue) Display(sty style) string {
	return fmt.Sprintf("%v", v.value.Display(sty))
}
func (v _nil_value) Display(sty style) string {
	return "nil"
}
