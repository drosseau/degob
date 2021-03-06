package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gitlab.com/drosseau/degob"
)

var (
	outFile     = flag.String("ofile", "", "Output file (defaults to stdout)")
	inFile      = flag.String("ifile", "", "Input file (defaults to stdin)")
	truncateOut = flag.Bool("trunc", false, "Truncate output file")
	base64d     = flag.Bool("b64", false, "base64 input")
	base64urld  = flag.Bool("b64url", false, "base64url input")
	noComments  = flag.Bool("nc", false, "don't print additional comments")
	noTypes     = flag.Bool("nt", false, "don't print type information")
	json        = flag.Bool("json", false, "show value as json")
	pkgName     = flag.String("pkg", "", "include a package definition in the output with the given name")
)

func errorf(s string, v ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, s, v...)
	os.Exit(1)
}

func getWriter() io.WriteCloser {
	if *outFile != "" {
		opts := os.O_WRONLY | os.O_CREATE
		if *truncateOut {
			opts |= os.O_TRUNC
		} else {
			opts |= os.O_APPEND
		}
		f, err := os.OpenFile(*outFile, opts, 0644)
		if err != nil {
			errorf("failed to open `%s` for writing: %v\n", *outFile, err)
		}
		return f
	}
	return os.Stdout
}

func getReader() io.ReadCloser {
	if *inFile != "" {
		f, err := os.Open(*outFile)
		if err != nil {
			errorf("failed to open `%s` for reading: %v\n", *inFile, err)
		}
		return f
	}
	return ioutil.NopCloser(os.Stdin)
}

type writer struct {
	w   io.Writer
	err error
}

func (w writer) Write(b []byte) (int, error) {
	if w.err != nil {
		errorf("error writing output: %v\n", w.err)
	}
	var n int
	n, w.err = w.w.Write(b)
	return n, w.err
}

func (w writer) writeComment(s string, v ...interface{}) {
	if *noComments {
		return
	}
	w.writeStr(s, v...)
}

func (w writer) writeStr(s string, v ...interface{}) {
	if w.err != nil {
		errorf("error writing output: %v\n")
	}
	_, w.err = fmt.Fprintf(w.w, s, v...)
}

func main() {
	flag.Parse()
	out := getWriter()
	defer out.Close()
	in := getReader()
	defer in.Close()

	if *base64d {
		in = ioutil.NopCloser(base64.NewDecoder(base64.StdEncoding, in))
	} else if *base64urld {
		in = ioutil.NopCloser(base64.NewDecoder(base64.URLEncoding, in))
	}

	w := writer{w: out}

	dec := degob.NewDecoder(in)
	gobs, err := dec.Decode()
	if err != nil {
		errorf("failed to decode gob: %s\n", err)
	}
	if *pkgName != "" {
		w.writeStr("package %s\n\n", *pkgName)
	}
	for i, g := range gobs {
		w.writeComment("// Decoded gob %d\n\n", i+1)
		if !*noTypes && g.Types != nil {
			w.writeComment("//Types\n")
			err = g.WriteTypes(w)
			if err != nil {
				errorf("error writing types: %v\n", err)
			}
		}
		w.writeComment("// Value: ")
		if *json {
			err = g.WriteValue(w, degob.JSON)
		} else {
			err = g.WriteValue(w, degob.SingleLine)
		}
		if err != nil {
			errorf("error writing values: %v\n", err)
		}
		w.writeComment("\n// End gob %d\n\n", i+1)
	}
}
