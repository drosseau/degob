package degob

import "io"

// gobBuf is a custom buffer type because I'd like to be able to print the
// whole gob in case of error. Using a `bytes.Buffer` I can't do that when
// also using the Read method because it effectively consumes the byte.
// I'm not sure how this effects performance but I don't think it is a huge
// deal
type gobBuf struct {
	data []byte // all of the data to be read
	pos  int    // the current position
}

func (g *gobBuf) Read(b []byte) (int, error) {
	n := copy(b, g.data[g.pos:])
	if n == 0 && len(b) != 0 {
		return 0, io.EOF
	}
	g.pos += n
	return n, nil
}

func (g *gobBuf) BytesRead() int {
	return g.pos
}

func (g *gobBuf) Write(b []byte) (int, error) {
	toWrite := len(b)
	g.Grow(toWrite)
	n := copy(g.data[g.pos:], b)
	if n != toWrite {
		return n, io.ErrShortWrite
	}
	return n, nil
}

func (g *gobBuf) Consumed(n int) {
	g.pos += n
}

func (g *gobBuf) Cap() int {
	return cap(g.data)
}

func (g *gobBuf) Len() int {
	return len(g.data) - g.pos
}

func (g *gobBuf) Bytes() []byte {
	return g.data[g.pos:]
}

// Data isn't a bytes.Buffer method, but it returns the entirety of the buffer
// which is pretty much the whole reason for this
func (g *gobBuf) Data() []byte {
	return g.data
}

func (g *gobBuf) Grow(n int) {
	if g.data == nil {
		g.data = make([]byte, n)
		return
	}
	if len(g.data)-g.pos >= n {
		return
	}
	for len(g.data)-g.pos < n {
		g.data = append(g.data, 0)
	}
}

func (g *gobBuf) ReadByte() (byte, error) {
	if g.pos >= len(g.data) {
		return 0, io.EOF
	}
	b := g.data[g.pos]
	g.pos++
	return b, nil
}

func (g *gobBuf) Reset() {
	g.data = g.data[:0]
	g.pos = 0
}
