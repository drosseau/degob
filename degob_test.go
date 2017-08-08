package degob

import (
	"bytes"
	"os"
	"testing"
)

func TestDegob(t *testing.T) {
	for _, obj := range testObjects {
		f, err := os.Open(obj.fileName)
		if err != nil {
			t.Fatalf("err: %v opening file: %s", err, obj.fileName)
		}
		var buf bytes.Buffer
		_, err = buf.ReadFrom(f)
		f.Close()
		if err != nil {
			t.Fatalf("err: %v reading file: %s", err, obj.fileName)
		}

		dec := NewDecoder(&buf)
		gobs, derr := dec.Decode()
		if derr != nil {
			t.Fatalf("err: %v decoding gob in file: %s", derr, obj.fileName)
		}
		compareGobs(obj.expected, gobs[0], obj.fileName, t)
	}
}

func TestGobStream(t *testing.T) {
	var buf bytes.Buffer
	for _, obj := range testObjects {
		f, err := os.Open(obj.fileName)
		if err != nil {
			t.Fatalf("err: %v opening file: %s", err, obj.fileName)
		}
		var buf bytes.Buffer
		_, err = buf.ReadFrom(f)
		f.Close()
		if err != nil {
			t.Fatalf("err: %v reading file: %s", err, obj.fileName)
		}
	}
	d := NewDecoder(&buf)
	i := 0
	for g := range d.DecodeStream(nil) {
		if g.Err != nil {
			t.Fatalf("err: %v decoding gob in file: %s", g.Err, testObjects[i].fileName)
		}
		compareGobs(testObjects[i].expected, g.Gob, testObjects[i].fileName, t)
		i += 1
	}
}
