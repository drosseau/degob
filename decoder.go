package degob

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

var anonTypes map[string]struct{}

func init() {
	anonTypes = make(map[string]struct{})
	rand.Seed(time.Now().UnixNano())
}

// This is a gob
// (byteCount (-type id, encoding of a wireType)* (type id, encoding of a value))*

const (
	uintByteSize       = 8
	smallestUserTypeId = 64
)

// Decoder reads and decodes gobs into our `Gob` type.
type Decoder struct {
	r            io.Reader // The base reader
	gobBuf       gobBuf    // Holds the current gob
	buf          [9]byte   // a buffer for reading uints
	seenTypes    map[typeId]*WireType
	decodedValue Value
	//values         map[typeId]*Value
	bytesProcessed uint64

	err *Error
}

func NewDecoder(r io.Reader) *Decoder {
	dec := new(Decoder)
	dec.r = bufio.NewReader(r)
	dec.seenTypes = make(map[typeId]*WireType)
	dec.decodedValue = nil
	return dec
}

// Decode reads all data on the input reader and decodes the Gobs returning
// a slice of `Gob`s.
func (dec *Decoder) Decode() ([]*Gob, error) {
	dec.bytesProcessed = 0
	dec.clearGob()
	gobs := make([]*Gob, 0, 5)
	for {
		dec.getGobPiece()
		if dec.err != nil {
			if dec.err.Err == io.EOF {
				break
			}
			return nil, dec.err
		}
		dec.decodeGobPiece()
		if dec.err != nil {
			return nil, dec.err
		}
		// we read the value of one of the input gobs and now we're on
		// to the next
		if dec.decodedValue != nil {
			g := new(Gob)
			dec.setGob(g)
			gobs = append(gobs, g)
			dec.clearGob()
		} else {
			dec.gobBuf.Reset()
		}
	}
	return gobs, nil
}

func (dec *Decoder) setGob(g *Gob) {
	g.Value = dec.decodedValue
	if len(dec.seenTypes) > 0 {
		for _, t := range dec.seenTypes {
			switch {
			case t.SliceT != nil:
				t.SliceT.ElemTypeString = dec.getName(t.SliceT.Elem)
			case t.ArrayT != nil:
				t.ArrayT.ElemTypeString = dec.getName(t.ArrayT.Elem)
			case t.MapT != nil:
				t.MapT.KeyTypeString = dec.getName(t.MapT.Key)
				t.MapT.ElemTypeString = dec.getName(t.MapT.Elem)
			case t.StructT != nil:
				for _, f := range t.StructT.Field {
					f.TypeString = dec.getName(typeId(f.Id))
				}
			}
		}
		g.Types = dec.seenTypes
	}
}

// Result is used for streaming
type Result struct {
	Gob *Gob
	Err *Error
}

// DecodeStream keeps reading from the underlying reader and returning
// on the returned channel. Errors do not stop the decoding. You can
// stop the decoding by closing the passed kill struct. If it is nil
// DecodeStream doesn't stop until EOF.
func (dec *Decoder) DecodeStream(kill <-chan struct{}) <-chan Result {
	dec.clearGob()
	c := make(chan Result)
	go func() {
		defer close(c)
		for {
			dec.getGobPiece()
			if dec.err != nil {
				if dec.err.Err == io.EOF {
					return
				}
				c <- Result{Err: dec.err}
				dec.err = nil
			}
			dec.decodeGobPiece()
			if dec.err != nil {
				c <- Result{Err: dec.err}
				dec.err = nil
			}
			if dec.decodedValue != nil {
				g := new(Gob)
				dec.setGob(g)
				c <- Result{Gob: g}
				select {
				case <-kill:
					return
				default:
					dec.clearGob()
				}
			} else {
				dec.gobBuf.Reset()
			}
		}
	}()
	return c
}

func (dec *Decoder) genError(err error) *Error {
	return &Error{
		Processed: dec.bytesProcessed,
		Err:       err,
		RawGob:    dec.gobBuf.Data(),
	}
}

// Loads a gob from the reader into the current gob buffer
func (dec *Decoder) getGobPiece() {
	if dec.err != nil {
		return
	}
	size, err := readUint(dec.r, dec.buf[:], &dec.bytesProcessed)
	if err != nil {
		dec.err = err
		return
	}
	// read the entire gob into the gob buffer
	_, err_ := io.CopyN(&dec.gobBuf, dec.r, int64(size))
	if err_ != nil {
		if err_ == io.EOF {
			dec.err = dec.genError(io.ErrUnexpectedEOF)
			return
		}
		dec.err = dec.genError(err_)
		return
	}
}

// main decoding entrypoint, decodes the gob in the gobBuf
//
// this function consumes the entire gobBuf until EOF
func (dec *Decoder) decodeGobPiece() {
	if dec.err != nil {
		return
	}
	for dec.err == nil && dec.gobBuf.Len() > 0 {
		id := dec.readTypeId()
		if id >= 0 {
			dec.decodedValue = dec.valueForType(id)
			w, ok := dec.seenTypes[id]
			if !ok || w.StructT == nil {
				dec.consumeNextUint(0)
			}
			// each gob will have a value so after we read it
			// let's return and add it to the returned *Gob's
			dec.readValue(id, &dec.decodedValue)
			return
		}
		// we have a type definition
		dec.readType(-id)
	}
}

func (dec *Decoder) readTypeId() typeId {
	if dec.err != nil {
		return 0
	}
	n, err := readUint(&dec.gobBuf, dec.buf[:], &dec.bytesProcessed)
	if err != nil {
		dec.err = err
		return 0
	}
	return typeId(uintToInt(n))
}

// valueForWireType creates a new value for the wiretype set as
// the default
func (dec *Decoder) valueForWireType(w *WireType) Value {
	switch {
	case w.StructT != nil:
		v := new(structValue)
		v.name = w.StructT.CommonType.Name
		v.fields = make(map[string]Value)
		for _, f := range w.StructT.Field {
			v.fields[f.Name] = dec.valueForType(typeId(f.Id))
		}
		return v
	case w.SliceT != nil:
		v := new(sliceValue)
		v.elemType = dec.getName(w.SliceT.Elem)
		return v
	case w.ArrayT != nil:
		v := new(arrayValue)
		v.length = w.ArrayT.Len
		v.elemType = dec.getName(w.ArrayT.Elem)
		v.values = make([]Value, v.length)
		for i := 0; i < v.length; i++ {
			v.values[i] = dec.valueForType(w.ArrayT.Elem)
		}
		return v
	case w.MapT != nil:
		v := new(mapValue)
		v.keyType = dec.getName(w.MapT.Key)
		v.elemType = dec.getName(w.MapT.Elem)
		v.values = make(map[Value]Value)
		return v
	default:
		panic("all nil in wiretype")
	}
}

func (dec *Decoder) valueForType(id typeId) Value {
	if isBuiltin(id) {
		return valueFor(id)
	}
	if w, ok := dec.seenTypes[id]; ok {
		return dec.valueForWireType(w)
	}
	panic("asked for value type for unknown type")
}

func (dec *Decoder) readValue(id typeId, v *Value) {
	if dec.err != nil {
		return
	}
	wire, ok := dec.seenTypes[id]
	if !ok {
		if !isBuiltin(id) {
			dec.err = dec.genError(errors.New("unexpected type"))
			return
		} else {
			dec.readBuiltinValue(id, v)
		}
	} else {
		switch {
		case wire.StructT != nil:
			dec.readStructValue(wire, v)
		case wire.MapT != nil:
			dec.readMapValue(wire, v)
		case wire.SliceT != nil:
			dec.readSliceValue(wire, v)
		case wire.ArrayT != nil:
			dec.readArrayValue(wire, v)
		}
	}
}

func (dec *Decoder) readBuiltinValue(id typeId, val *Value) {
	if dec.err != nil {
		return
	}
	switch id {
	case _bool_id:
		b := dec.nextUint()
		// should I check that it is 0 or 1?
		if b == 0 {
			*val = _bool_type(false)
		} else {
			*val = _bool_type(true)
		}
	case _int_id:
		v := dec.nextUint()
		if v&1 != 0 {
			*val = _int_type(^int64(v >> 1))
		} else {
			*val = _int_type(int64(v >> 1))
		}
	case _uint_id:
		v := dec.nextUint()
		*val = _uint_type(v)
	case _float_id:
		v := dec.nextUint()
		*val = _float_type(uintToFloat(v))
	case _complex_id:
		r := dec.nextUint()
		i := dec.nextUint()
		*val = _complex_type(uintToComplex(r, i))
	case _bytes_id:
		l := int(dec.nextUint())
		b := make([]byte, l)
		dec.gobBuf.Read(b)
		*val = _bytes_type(b)
	case _string_id:
		l := int(dec.nextUint())
		b := make([]byte, l)
		dec.gobBuf.Read(b)
		*val = _string_type(b)
	case _interface_id:
		nameLen := int(dec.nextUint())
		if nameLen == 0 {
			dec.readNilInterface(val)
		} else {
			dec.readNonNilInterface(val, nameLen)
		}
	default:
		panic("id was not a builtin id")
	}
}

func (dec *Decoder) readNilInterface(v *Value) {
	panic("not implemented yet")
}

func (dec *Decoder) readNonNilInterface(v *Value, nl int) {
	if dec.err != nil {
		return
	}
	var into interfaceValue
	nameB := make([]byte, nl)
	dec.gobBuf.Read(nameB)
	into.name = string(nameB)
	for {
		id := dec.readTypeId()
		if id < 0 {
			dec.readType(-id)
		} else {
			// hmm
			// TODO: What is this next uint telling me?
			// I think it is a length of the next block to read
			_ = dec.nextUint()
			w, ok := dec.seenTypes[id]
			if !ok || w.StructT == nil {
				dec.consumeNextUint(0)
			}
			dec.readValue(id, &into.value)
			break
		}
	}
	*v = into
}

func (dec *Decoder) readMapValue(wire *WireType, val *Value) {
	if dec.err != nil {
		return
	}
	into := dec.valueForWireType(wire).(*mapValue)
	length := int(dec.nextUint())
	for i := 0; i < length; i++ {
		kVal := new(Value)
		eVal := new(Value)
		dec.readValue(wire.MapT.Key, kVal)
		dec.readValue(wire.MapT.Elem, eVal)
		into.values[*kVal] = *eVal
	}
	*val = into
}

func (dec *Decoder) readSliceValue(wire *WireType, val *Value) {
	if dec.err != nil {
		return
	}
	length := int(dec.nextUint())
	into := dec.valueForWireType(wire).(*sliceValue)
	into.values = make([]Value, length)
	// TODO: should I set all of these to the default?
	for i := 0; i < length; i++ {
		dec.readValue(wire.SliceT.Elem, &into.values[i])
	}
	*val = into
}

func (dec *Decoder) readArrayValue(wire *WireType, val *Value) {
	if dec.err != nil {
		return
	}
	into := dec.valueForWireType(wire).(*arrayValue)
	length := int(dec.nextUint())
	for i := 0; i < length; i++ {
		dec.readValue(wire.ArrayT.Elem, &into.values[i])
	}
	*val = into
}

func (dec *Decoder) readStructValue(wire *WireType, val *Value) {
	if dec.err != nil {
		return
	}
	into := dec.valueForWireType(wire).(*structValue)
	fields := wire.StructT.Field
	fieldNum := -1
	for {
		delta := int(dec.nextUint())
		if delta == 0 || dec.err != nil {
			break
		}
		fieldNum += delta
		if fieldNum < 0 || fieldNum >= len(fields) {
			dec.err = dec.genError(errors.New("bad fieldnum"))
			return
		}
		id := typeId(fields[fieldNum].Id)
		var v Value
		if isBuiltin(id) {
			v = valueFor(id)
		}
		dec.readValue(id, &v)
		into.fields[fields[fieldNum].Name] = v
	}
	*val = into
}

// reads newly defined types. These will always come as WireType structs
func (dec *Decoder) readType(id typeId) {
	if dec.err != nil {
		return
	}
	if id < smallestUserTypeId || dec.seenTypes[id] != nil {
		dec.err = errDuplicateType(dec.bytesProcessed, nil)
		return
	}
	wire := new(WireType)
	dec.decodeType(id, wire)
	// Every type definition will be followed by two null bytes
	dec.consumeNextUint(0)
	dec.consumeNextUint(0)
	dec.seenTypes[id] = wire
}

// reads the gobBuf and stores the read WireType only operates one
// WireType at a time
func (dec *Decoder) decodeType(id typeId, w *WireType) {
	if dec.err != nil {
		return
	}
	delta := int(dec.nextUint())
	fieldNum := delta - 1
	dec.consumeNextUint(1)
	switch fieldNum {
	case 0:
		dec.decodeArray(id, w)
	case 1:
		dec.decodeSlice(id, w)
	case 2:
		dec.decodeStruct(id, w)
	case 3:
		dec.decodeMap(id, w)
	default:
		dec.err = errUnknownDelta(dec.bytesProcessed, dec.gobBuf.Bytes())
		return
	}
}

func (dec *Decoder) consumeNextUint(expected int) {
	if dec.err != nil {
		return
	}
	delta := int(dec.nextUint())
	if delta != expected {
		dec.err = dec.genError(fmt.Errorf("expected delta %d but got %d", expected, delta))
		return
	}
}

func (dec *Decoder) nextUint() uint64 {
	if dec.err != nil {
		return 0
	}
	delta, err := readUint(&dec.gobBuf, dec.buf[:], &dec.bytesProcessed)
	if err != nil {
		dec.err = err
		return 0
	}
	return delta
}

func (dec *Decoder) decodeArray(id typeId, w *WireType) {
	if dec.err != nil {
		return
	}
	common := dec.decodeCommon()
	dec.consumeNextUint(1)
	elemId := dec.readTypeId()
	dec.consumeNextUint(1)
	l := int(uintToInt(dec.nextUint()))
	w.ArrayT = &ArrayType{
		CommonType: common,
		Elem:       elemId,
		Len:        l,
	}
}

func (dec *Decoder) decodeSlice(id typeId, w *WireType) {
	if dec.err != nil {
		return
	}
	common := dec.decodeCommon()
	dec.consumeNextUint(1)
	elemId := dec.readTypeId()
	w.SliceT = &SliceType{
		CommonType: common,
		Elem:       elemId,
	}
}

func (dec *Decoder) decodeStruct(id typeId, w *WireType) {
	if dec.err != nil {
		return
	}
	common := dec.decodeCommon()
	fields := dec.decodeFields()
	w.StructT = &StructType{
		CommonType: common,
		Field:      fields,
	}
	// set the name if it is anonymous
	if common.Name == "" {
		dec.anonymousStructTypeName(w)
	}
}

func (dec *Decoder) decodeMap(id typeId, w *WireType) {
	if dec.err != nil {
		return
	}
	common := dec.decodeCommon()
	dec.consumeNextUint(1)
	keyId := dec.readTypeId()
	dec.consumeNextUint(1)
	elemId := dec.readTypeId()
	w.MapT = &MapType{
		CommonType: common,
		Key:        keyId,
		Elem:       elemId,
	}
}

func (dec *Decoder) decodeFields() []*FieldType {
	if dec.err != nil {
		return nil
	}
	// TODO: potential bug here for empty structs?
	dec.consumeNextUint(1)
	nfields := int(dec.nextUint())
	fields := make([]*FieldType, nfields)
	for i := 0; i < nfields; i++ {
		fields[i] = new(FieldType)
		dec.consumeNextUint(1)
		dec.decodeString(&fields[i].Name)
		dec.consumeNextUint(1)
		fields[i].Id = int(dec.readTypeId())
		dec.consumeNextUint(0)
	}
	return fields
}

func (dec *Decoder) decodeCommon() CommonType {
	var c CommonType
	if dec.err != nil {
		return c
	}
	fieldNum := -1
	for {
		delta := int(dec.nextUint())
		// the end is noted with delta 0
		if delta == 0 {
			break
		}
		fieldNum += delta
		switch fieldNum {
		case 0:
			dec.decodeString(&c.Name)
		case 1:
			c.Id = int(dec.readTypeId())
		default:
			dec.err = errCorruptCommonType(dec.bytesProcessed, dec.gobBuf.Bytes())
			return c
		}
	}
	return c
}

func (dec *Decoder) decodeString(into *string) {
	if dec.err != nil {
		return
	}
	l := dec.nextUint()
	b := make([]byte, l)
	r, err_ := dec.gobBuf.Read(b)
	if err_ != nil {
		if err_ == io.EOF {
			dec.err = dec.genError(io.ErrUnexpectedEOF)
			return
		}
		dec.err = dec.genError(err_)
		return
	}
	dec.bytesProcessed += uint64(r)
	if uint64(r) != l {
		dec.err = errBadString(dec.bytesProcessed, dec.gobBuf.Bytes())
		return
	}
	*into = string(b)
}

// clears any gob that is already seenTypes
func (dec *Decoder) clearGob() {
	dec.gobBuf.Reset()
	dec.seenTypes = make(map[typeId]*WireType)
	dec.decodedValue = nil
}

func (dec *Decoder) getName(id typeId) string {
	if isBuiltin(id) {
		return strings.TrimSpace(id.name())
	}
	if v, ok := dec.seenTypes[id]; ok {
		switch {
		case v.StructT != nil:
			if v.StructT.CommonType.Name == "" {
				return strings.TrimSpace(dec.anonymousStructTypeName(v))
			}
			return strings.TrimSpace(v.StructT.CommonType.Name)
		case v.SliceT != nil:
			return strings.TrimSpace(v.SliceT.CommonType.Name)
		case v.MapT != nil:
			return strings.TrimSpace(v.MapT.CommonType.Name)
		case v.ArrayT != nil:
			return strings.TrimSpace(v.ArrayT.CommonType.Name)
		default:
			panic("empty wiretype at chosen id")
		}
	} else {
		panic("something has gone very wrong")
	}
}

func (dec *Decoder) anonymousStructTypeName(w *WireType) string {
	follow := make([]byte, 4)
	_, _ = rand.Read(follow)
	var followString string
	for {
		followString = hex.EncodeToString(follow)
		if _, ok := anonTypes[followString]; ok {
			// that is annoying..
			_, _ = rand.Read(follow)
		} else {
			anonTypes[followString] = struct{}{}
			break
		}
	}
	s := fmt.Sprintf("Anon%d_%s", w.StructT.Id, followString)
	/*
		s := "struct {"
		nfields := len(w.StructT.Field)
		for l, f := range w.StructT.Field {
			if l < nfields-1 {
				s += fmt.Sprintf("%s %s; ", f.Name, dec.getName(typeId(f.Id)))
			} else {
				s += fmt.Sprintf("%s %s}", f.Name, dec.getName(typeId(f.Id)))
			}
		}
	*/
	w.StructT.CommonType.Name = s
	return s
}
