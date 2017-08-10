# degob

Simple command line too for degobbing.

```
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

If you come up with a gob this doesn't work with I wouldn't be surprised but make an issue. (Currently an empty struct (`struct{}`) can cause issues).

### Sample output

`$ hexdump -C gob.bin`
```
00000000  0e ff 81 04 01 02 ff 82  00 01 10 01 10 00 00 62  |...............b|
00000010  ff 82 00 03 06 73 74 72  69 6e 67 0c 0e 00 0c 53  |.....string....S|
00000020  74 72 69 6e 67 54 6f 42  6f 6f 6c 04 62 6f 6f 6c  |tringToBool.bool|
00000030  02 02 00 00 06 73 74 72  69 6e 67 0c 0d 00 0b 53  |.....string....S|
00000040  74 72 69 6e 67 54 6f 49  6e 74 03 69 6e 74 04 02  |tringToInt.int..|
00000050  00 18 03 69 6e 74 04 04  00 fe 09 a4 06 73 74 72  |...int.......str|
00000060  69 6e 67 0c 0d 00 0b 49  6e 74 54 6f 53 74 72 69  |ing....IntToStri|
00000070  6e 67                                             |ng|
00000072
```
`$ cat gob.bin | degob`
```go
// Decoded gob #1

// Types:
map[interface{}]interface{}

// Values:
map[interface{}]interface{}{"StringToBool": false,"StringToInt": 12,1234: "IntToString"}

// End gob #1
```

`$ hexdump -C gob.bin`
```
00000000  2b ff 83 03 01 01 04 54  65 73 74 01 ff 84 00 01  |+......Test.....|
00000010  04 01 01 57 01 ff 86 00  01 01 58 01 04 00 01 01  |...W......X.....|
00000020  59 01 06 00 01 01 5a 01  0c 00 00 00 25 ff 85 03  |Y.....Z.....%...|
00000030  01 01 05 49 6e 6e 65 72  01 ff 86 00 01 03 01 01  |...Inner........|
00000040  41 01 08 00 01 01 42 01  0e 00 01 01 43 01 0a 00  |A.....B.....C...|
00000050  00 00 28 ff 84 01 01 f8  1f 85 eb 51 b8 1e 09 40  |..(........Q...@|
00000060  01 fe 14 40 fe 08 40 01  05 01 02 03 04 05 00 01  |...@..@.........|
00000070  13 01 0a 01 05 48 65 6c  6c 6f 00                 |.....Hello.|
0000007b
```
`$ cat gob.bin | degob `
```go
// Decoded gob #1

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


// End gob #1
```

Two gobs in the same file (note that the encoded gobs make the structs anonymous)

```
$ hexdump -C gob1.bin
00000000  0d ff 91 02 01 02 ff 92  00 01 ff 90 00 00 2a ff  |..............*.|
00000010  8f 03 01 01 0a 53 6c 69  63 65 49 6e 6e 65 72 01  |.....SliceInner.|
00000020  ff 90 00 01 02 01 04 55  69 6e 74 01 06 00 01 04  |.......Uint.....|
00000030  42 79 74 65 01 06 00 00  00 0c ff 92 00 02 02 30  |Byte...........0|
00000040  00 01 05 01 35 00                                 |....5.|

$ hexdump -C gob2.bin
00000000  0f ff 95 04 01 02 ff 96  00 01 0c 01 ff 94 00 00  |................|
00000010  22 ff 93 03 01 02 ff 94  00 01 02 01 07 43 6f 6d  |"............Com|
00000020  70 6c 65 78 01 0e 00 01  05 46 6c 6f 61 74 01 08  |plex.....Float..|
00000030  00 00 00 35 ff 96 00 02  07 6b 65 79 20 6f 6e 65  |...5.....key one|
00000040  01 ff c0 fe 08 40 01 f8  66 66 66 66 66 66 24 40  |.....@..ffffff$@|
00000050  00 07 6b 65 79 20 74 77  6f 01 40 fe 08 c0 01 f8  |..key two.@.....|
00000060  66 66 66 66 66 66 24 c0  00                       |ffffff$..|
```

`$ cat gob1.bin gob2.bin | degob`

```go
// Decoded gob #1

// Types:
[]SliceInner

type SliceInner struct {
  Uint uint64
  Byte uint64
}

// Values:
[]SliceInner{SliceInner{Uint: 0, Byte: 48}, SliceInner{Uint: 5, Byte: 53}}

// End gob #1

// Decoded gob #2

// Types:
map[string]Anon74_1da8c6d2

type Anon74_1da8c6d2 struct {
  Complex complex128
  Float float64
}

// Values:
map[string]Anon74_1da8c6d2{"key one": Anon74_1da8c6d2{Complex: (-2+3i), Float: 10.2},"key two": Anon74_1da8c6d2{Complex: (2-3i), Float: -10.2}}

// End gob #2
```

