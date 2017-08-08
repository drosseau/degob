# degob

degob is a reversing tool for `gob`s. If you have the binary form of a gob, but you don't have to `struct` to decode it into, this library still allows you to get a more readable representation of the data.

## cmds/degob

The easiest way to use all of this is to just build the binary in `cmds/degob` and send gobs to it either through `stdin` or from files and then get the output to `stdout` or to a file.

```sh
Usage of ./degob:
  -b64
    	base64 input
  -b64url
    	base64url input
  -ifile string
    	Input file (defaults to stdin)
  -ofile string
    	Output file (defaults to stdout)
  -trunc
    	Truncate output file
```

This should work on a ton of different inputs but since this isn't quite stable yet some things probably won't work.

## TODO

- Printing stylized output isn't complete yet.
- `display.go` and `type.go` are a mess
- Some tests are old and no longer passing
- Streaming no go still
- A lot more testing
- Documentation
- Errors maybe could be handled better and I haven't tested if the byte count on error is correct or not yet