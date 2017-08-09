package degob

import (
	"bytes"
	"testing"
)

func readUintTest(t *testing.T, from []byte, width int, into []byte) uint64 {
	var read uint64
	val, w, err := readUint(bytes.NewBuffer(from), into, &read)
	if err != nil {
		t.Fatal(err)
	}
	if width != w {
		t.Fatalf("expected to read %d but read %d", width, w)
	}
	return val
}

func TestUintRead(t *testing.T) {
	// Thus 0 is transmitted as (00), 7 is transmitted as (07) and 256 is transmitted as (FE 01 00).
	var into [9]byte
	var from [9]byte
	from[0] = 0x07
	val := readUintTest(t, from[:], 1, into[:])
	if val != uint64(7) {
		t.Fatalf("expected 7 got %v\n", val)
	}
	from[0] = 0xfe
	from[1] = 0x01
	from[2] = 0x00
	val = readUintTest(t, from[:], 3, into[:])
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
	i64 := uintToInt(readUintTest(t, from[:], 3, into[:]))
	if i64 != int64(-129) {
		t.Fatalf("expected -129 got %v\n", i64)
	}
}

func TestBadUint(t *testing.T) {
	from := [11]byte{0xf6, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	var into [9]byte
	var read uint64
	_, w, err := readUint(bytes.NewBuffer(from[:]), into[:], &read)
	if err == nil {
		t.Fatal("should have errored")
	}
	if w != 11 {
		t.Fatal("should have had width 11 but was", w)
	}
	if err.Processed != 1 {
		t.Fatal("should have only processed the first byte")
	}
}
