package degob

// At the moment this is pretty meh, but I don't particularly
// want to spend too much time on it right now

import (
	"fmt"
	"strings"
)

type style uint8

const (
	// CommentedPretty prints the values inside /* ... */ pretty printed
	//CommentedPretty style = iota
	// CommentedSingleLine prints values after // on a single line
	CommentedSingleLine = iota
	// Pretty prints values pretty uncommented
	//Pretty
	// SingleLine prints values in a single line
	SingleLine
	// JSON prints the values as JSON objects
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

func (a *ArrayType) String() string {
	return fmt.Sprintf("[%d]%s", a.Len, a.ElemTypeString)
}

func (s *SliceType) String() string {
	return fmt.Sprintf("[]%s", s.ElemTypeString)
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.KeyTypeString, m.ElemTypeString)
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

// pretty printed as []elemType{
//		... value.display() ...
//	}
type sliceValue struct {
	elemType string
	values   []Value
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

// pretty printed as [length]elemType{
//		... value.display() ...
//	}
type arrayValue struct {
	elemType string
	length   int
	values   []Value
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

// pretty printed as []map[key]elem{
//		.. value.display(): value.display() ..
// }
type mapValue struct {
	keyType  string
	elemType string
	values   map[Value]Value
}

func (v mapValue) getValues(sty style) string {
	var out string
	nval := len(v.values)
	i := 0
	for k, v := range v.values {
		if i < nval-1 {
			out += fmt.Sprintf("%s: %s,", k.Display(sty), v.Display(sty))
		} else {
			out += fmt.Sprintf("%s: %s", k.Display(sty), v.Display(sty))
		}
		i += 1
	}
	return out
}

func (v mapValue) Display(sty style) string {
	return fmt.Sprintf("map[%s]%s{%s}", v.keyType, v.elemType, v.getValues(sty))
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

// pretty printed as name {
//		... fields[name]: fields[value].display(), ...
//	}
type structValue struct {
	name   string
	fields map[string]Value
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

func (s *structValue) commentedSingleLine() string {
	return fmt.Sprintf("//%s", s.singleLine())
}

func (s *structValue) getFieldVals(newline bool, sty style) string {
	var out string
	nfields := len(s.fields)
	i := 0
	for k, v := range s.fields {
		out += fmt.Sprintf("%s: %s", k, v.Display(sty))
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
	return "{}"
}
