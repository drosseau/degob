package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/drosseau/degob"
)

var (
	outFile     = flag.String("ofile", "", "Output file (defaults to stdout)")
	inFile      = flag.String("ifile", "", "Input file (defaults to stdin)")
	truncateOut = flag.Bool("trunc", false, "Truncate output file")
	base64d     = flag.Bool("b64", false, "base64 input")
	base64urld  = flag.Bool("b64url", false, "base64url input")
)

func errorf(format string, v ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, v...)
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
			errorf("failed to open %s: %v\n", *outFile, err)
		}
		return f
	}
	return os.Stdout
}

func getReader() io.ReadCloser {
	if *inFile != "" {
		f, err := os.Open(*outFile)
		if err != nil {
			errorf("failed to open %s: %v\n", *inFile, err)
		}
		return f
	}
	return ioutil.NopCloser(os.Stdin)
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

	dec := degob.NewDecoder(in)
	gob, err := dec.Decode()
	if err != nil {
		errorf("failed to decode gob: %s\n", err)
	}
	errg := gob.WriteTypes(out)
	if errg != nil {
		errorf("failed to write types: %s\n", errg)
	}
	errg = gob.WriteValue(out, degob.SingleLine)
	if errg != nil {
		errorf("failed to write value: %s\n", errg)
	}
}
