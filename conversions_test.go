package degob

import (
	"bytes"
	"testing"
)

func readUintTest(t *testing.T, from []byte, into []byte) uint64 {
	var read uint64
	val, err := readUint(bytes.NewBuffer(from), into, &read)
	if err != nil {
		t.Fatal(err)
	}
	return val
}

func TestUintRead(t *testing.T) {
	// Thus 0 is transmitted as (00), 7 is transmitted as (07) and 256 is transmitted as (FE 01 00).
	var into [9]byte
	var from [9]byte
	from[0] = 0x07
	val := readUintTest(t, from[:], into[:])
	if val != uint64(7) {
		t.Fatalf("expected 7 got %v\n", val)
	}
	from[0] = 0xfe
	from[1] = 0x01
	from[2] = 0x00
	val = readUintTest(t, from[:], into[:])
	if val != uint64(256) {
		t.Fatalf("expected 256 got %v\n", val)
	}
}

func TestUintToInt(t *testing.T) {
	var into [9]byte
	var from [9]byte
	from[0] = 0xfe
	from[1] = 0x01
	from[2] = 0x01
	i64 := uintToInt(readUintTest(t, from[:], into[:]))
	if i64 != int64(-129) {
		t.Fatalf("expected -129 got %v\n", i64)
	}
}
