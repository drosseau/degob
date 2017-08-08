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

This should work on a lot of different and simple inputs but, since this isn't quite stable yet, some things probably won't work.

### Limitations

There are a few limitations that I can't really get around.

- gobs don't include information about the bit size of the type so all types are their largest possible (`int64`, `uint64`, `complex128`, `float64`) so as to be able to accept anything.
- `byte`s are received as `uint64`, but `[]byte` is correct.


### Sample output

Binary blob of a `map[interface{}]interface{}`

```sh
$ hexdump -C blob.bin 
00000000  0e ff 81 04 01 02 ff 82  00 01 10 01 10 00 00 62  |...............b|
00000010  ff 82 00 03 06 73 74 72  69 6e 67 0c 0e 00 0c 53  |.....string....S|
00000020  74 72 69 6e 67 54 6f 42  6f 6f 6c 04 62 6f 6f 6c  |tringToBool.bool|
00000030  02 02 00 00 06 73 74 72  69 6e 67 0c 0d 00 0b 53  |.....string....S|
00000040  74 72 69 6e 67 54 6f 49  6e 74 03 69 6e 74 04 02  |tringToInt.int..|
00000050  00 18 03 69 6e 74 04 04  00 fe 09 a4 06 73 74 72  |...int.......str|
00000060  69 6e 67 0c 0d 00 0b 49  6e 74 54 6f 53 74 72 69  |ing....IntToStri|
00000070  6e 67                                             |ng|
00000072
$ cat interfacemap.bin | degob 
// Types:
map[interface{}]interface{}

// Values:
map[interface{}]interface{}{"StringToBool": false,"StringToInt": 12,1234: "IntToString"}
```

Binary blob of nested user defined structs

```sh
$ hexdump -C blob.bin 
00000000  2b ff 83 03 01 01 04 54  65 73 74 01 ff 84 00 01  |+......Test.....|
00000010  04 01 01 57 01 ff 86 00  01 01 58 01 04 00 01 01  |...W......X.....|
00000020  59 01 06 00 01 01 5a 01  0c 00 00 00 25 ff 85 03  |Y.....Z.....%...|
00000030  01 01 05 49 6e 6e 65 72  01 ff 86 00 01 03 01 01  |...Inner........|
00000040  41 01 08 00 01 01 42 01  0e 00 01 01 43 01 0a 00  |A.....B.....C...|
00000050  00 00 28 ff 84 01 01 f8  1f 85 eb 51 b8 1e 09 40  |..(........Q...@|
00000060  01 fe 14 40 fe 08 40 01  05 01 02 03 04 05 00 01  |...@..@.........|
00000070  13 01 0a 01 05 48 65 6c  6c 6f 00                 |.....Hello.|
0000007b
$ cat blob.bin | degob 
// Types:
type Test struct {
	W Inner
	X int64
	Y uint64
	Z string
}

type Inner struct {
	A float64
	B complex128
	C []byte
}

// Values:
Test{W: Inner{A: 3.14, B: (5+3i), C: []byte{0x1, 0x2, 0x3, 0x4, 0x5}}, X: -10, Y: 10, Z: "Hello"}
```


## TODO

- Printing stylized output isn't complete yet.
- `display.go` and `type.go` are a mess
- Some tests are old and no longer passing
- Streaming no go still
- A lot more testing
- Documentation
- Errors maybe could be handled better and I haven't tested if the byte count on error is correct or not yet
- Figure out licensing; I heavily relied on the Go source for this tool