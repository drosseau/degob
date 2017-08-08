package degob

import (
	"io"
	"math"
)

type typeId int32

// reads a uint from the reader
func readUint(r io.Reader, into []byte, read *uint64) (uint64, *Error) {
	var n int
	n, err := r.Read(into[:1])
	if err != nil {
		if err == io.EOF {
			return 0, genericError(io.EOF, *read, nil)
		}
		var gobBytes []byte
		if v, ok := r.(*gobBuf); ok {
			gobBytes = v.Data()
		}
		return 0, genericError(err, *read, gobBytes)
	}
	*read += uint64(n)
	b := into[0]
	// anything less than 0x7f is encoded as a single byte with
	// that value so we're done
	if b <= 0x7f {
		return uint64(b), nil
	}
	// FROM DOCS:
	// Otherwise it is sent as a minimal-length big-endian (high byte first)
	// byte stream holding the value, preceded by one byte holding the byte
	// count, negated.
	n = -int(int8(b))
	if n > uintByteSize {
		var gobBytes []byte
		if v, ok := r.(*gobBuf); ok {
			gobBytes = v.Data()
		}
		return 0, errUintTooBig(*read, gobBytes)
	}
	// now we read n bytes and that is our uint
	n, err = io.ReadFull(r, into[0:n])
	if err != nil {
		var gobBytes []byte
		if v, ok := r.(*gobBuf); ok {
			gobBytes = v.Data()
		}
		if err == io.EOF {
			return 0, genericError(io.ErrUnexpectedEOF, *read, gobBytes)
		}
		return 0, genericError(err, *read, gobBytes)
	}
	*read += uint64(n)
	var val uint64
	for _, b := range into[0:n] {
		val = val<<8 | uint64(b)
	}
	return val, nil
}

func uintToInt(x uint64) int64 {
	i := int64(x >> 1)
	if x&1 != 0 {
		i = ^i
	}
	return i
}

func uintToFloat(x uint64) float64 {
	var v uint64
	for i := 0; i < 8; i++ {
		v <<= 8
		v |= x & 0xFF
		x >>= 8
	}
	return math.Float64frombits(v)
}

func uintToComplex(r uint64, i uint64) complex128 {
	return complex(uintToFloat(r), uintToFloat(i))
}
