package degob

import (
	"bytes"
	"io"
	"testing"
)

func TestDegob(t *testing.T) {
	for _, obj := range testObjects {
		var buf bytes.Buffer
		fileToBufferTest(obj.fileName, &buf, t)
		dec := NewDecoder(&buf)
		gobs, derr := dec.Decode()
		if derr != nil {
			t.Fatalf("err: %v decoding gob in file: %s", derr, obj.fileName)
		}
		if len(gobs) == 0 {
			t.Fatalf("no error but empty gobs for file %s", obj.fileName)
		}
		compareGobs(obj.expected, gobs[0], obj.fileName, t)
	}
}

func TestGobStream(t *testing.T) {
	var buf bytes.Buffer
	for i, obj := range testObjects {
		fileToBufferTest(obj.fileName, &buf, t)
		if i == 2 {
			break
		}
	}
	d := NewDecoder(&buf)
	i := 0
	for g := range d.DecodeStream(nil, 0) {
		if g.Err != nil {
			t.Fatalf("err: %v decoding gob in file: %s", g.Err, testObjects[i].fileName)
		}
		compareGobs(testObjects[i].expected, g.Gob, testObjects[i].fileName, t)
		i += 1
	}
}

func TestKillGobStream(t *testing.T) {
	kill := make(chan struct{})
	var buf bytes.Buffer
	tot := 0
	for i, obj := range testObjects {
		fileToBufferTest(obj.fileName, &buf, t)
		tot += 1
		if i == 2 {
			break
		}
	}
	d := NewDecoder(&buf)
	out := d.DecodeStream(kill, 0)
	close(kill)
	i := 0
	for g := range out {
		if g.Err != nil {
			t.Fatalf("err: %v decoding gob in file: %s", g.Err, testObjects[i].fileName)
		}
		compareGobs(testObjects[i].expected, g.Gob, testObjects[i].fileName, t)
		i += 1
	}
	if !(i < tot) {
		t.Fatal("expected to exit streaming early")
	}
}

func TestDuplicateDefinition(t *testing.T) {
	obj := testObjects[0]
	var buf bytes.Buffer
	fileToBufferTest(obj.fileName, &buf, t)
	b := buf.Bytes()
	b[46] = 0x81
	buf.Reset()
	_, err := buf.Write(b)
	if err != nil {
		t.Fatal("writing to buffer", err)
	}
	d := NewDecoder(&buf)
	_, err = d.Decode()
	if err == nil {
		t.Fatal("should have errored")
	}
	derr := err.(*Error)
	if derr.Processed != 47 {
		t.Fatal("should have been at byte 47 but was at", derr.Processed)
	}
}

func TestNonEOFErrorStream(t *testing.T) {
	var buf bytes.Buffer
	for i, obj := range testObjects {
		if i == 0 {
			var tmpbuf bytes.Buffer
			fileToBufferTest(obj.fileName, &tmpbuf, t)
			b := tmpbuf.Bytes()
			b[46] = 0x81
			_, err := buf.Write(b)
			if err != nil {
				t.Fatal("writing to buffer", err)
			}
			continue
		} else {
			fileToBufferTest(obj.fileName, &buf, t)
		}
		if i == 2 {
			break
		}
	}
	d := NewDecoder(&buf)
	i := 0
	for g := range d.DecodeStream(nil, 0) {
		if i != 0 {
			if g.Err != nil {
				t.Fatalf("err: %v decoding gob in file: %s", g.Err, testObjects[i].fileName)
			}
			compareGobs(testObjects[i].expected, g.Gob, testObjects[i].fileName, t)
		} else {
			if g.Err == nil {
				t.Fatalf("expected error")
			}
		}
		i += 1
	}
	if i != 3 {
		t.Fatal("didn't process all gobs")
	}
}

func TestUnexpectedEOFStream(t *testing.T) {
	var buf bytes.Buffer
	for i, obj := range testObjects {
		if i == 1 {
			// get only part of it and break
			f := openFileTest(obj.fileName, t)
			b := make([]byte, 5)
			_, err := io.ReadFull(f, b)
			if err != nil {
				t.Fatal("error reading to buf %v", err)
			}
			_, err = buf.Write(b)
			if err != nil {
				t.Fatal("error writing to buf %v", err)
			}
			break
		}
		fileToBufferTest(obj.fileName, &buf, t)
	}
	d := NewDecoder(&buf)
	out := d.DecodeStream(nil, 1)
	i := 0
	g := <-out
	if g.Err != nil {
		t.Fatalf("err: %v decoding gob in file: %s", g.Err, testObjects[i].fileName)
	}
	compareGobs(testObjects[0].expected, g.Gob, testObjects[i].fileName, t)
	g = <-out
	if g.Err == nil {
		t.Fatal("expected an error")
	}
	if g.Err.Err != io.ErrUnexpectedEOF {
		t.Fatal("expected io.ErrUnexpectedEOF")
	}

	for _ = range out {
		t.Fatal("channel should be closed")
	}
}

func TestUnexpectedEOF(t *testing.T) {
	obj := testObjects[0]
	f := openFileTest(obj.fileName, t)
	buf := make([]byte, 10)
	_, err := io.ReadFull(f, buf)
	if err != nil {
		t.Fatalf("err: %v ReadFull from %s", err, obj.fileName)
	}
	b := bytes.NewBuffer(buf)
	d := NewDecoder(b)
	_, err = d.Decode()
	if err == nil {
		t.Fatal("expected an ErrUnexpectedEOF error but got nil")
	}
	derr, ok := err.(*Error)
	if !ok {
		t.Fatal("expected library error type")
	}

	if derr.Err != io.ErrUnexpectedEOF {
		t.Fatalf("expected an ErrUnexpectedEOF error but was %v", derr.Err)
	}
}

func TestBadGobUnexpectedTypeId(t *testing.T) {
	obj := testObjects[3]
	var buf bytes.Buffer
	fileToBufferTest(obj.fileName, &buf, t)
	b := buf.Bytes()
	b[2] = 0x80
	buf.Reset()
	_, err := buf.Write(b)
	if err != nil {
		t.Fatal("writing to buffer", err)
	}
	d := NewDecoder(&buf)
	_, err = d.Decode()
	if err == nil {
		t.Fatal("expected error")
	}
	derr := err.(*Error)
	if derr.Processed != 3 {
		t.Fatal("expected error to be at byte 3 but was at byte", derr.Processed)
	}
}
