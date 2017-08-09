# degob

degob is a reversing tool for `gob`s. If you have the binary form of a gob, but you don't have to `struct` to decode it into, this library still allows you to get a more readable representation of the data.

## cmds/degob

The easiest way to use all of this is to just build the binary in `cmds/degob` and send gobs to it either through `stdin` or from files and then get the output to `stdout` or to a file. See its [README](cmds/degob/README.md) for more info.

## Usage

Create a new `Decoder` over your reader using `NewDecoder` and then decode that into a slice of `Gob`s with `Decode` or stream `Gob`s with `DecodeStream`. `DecodeStream` seems fairly stable, but it was difficult to test how it handles all error cases, so be wary of errors. Once you have `Gob`s you can either play with the types directly or just print them out to a writer using the `WriteTypes` and `WriteValues` methods.

The output from the Write methods on Gob should be close to valid Go source. One obvious instance that this isn't true is if the gob defines a type that isn't a struct (ie when sending a raw slice like `[]Foo` it first defines an unnamed type `[]Foo`).

The provided `degob` command provides a straightforward [sample usage](cmds/degob/main.go).

### Limitations

There are a few limitations that I can't really get around.

- gobs don't include information about the bit size of the type so all types are their largest possible (`int64`, `uint64`, `complex128`, `float64`) so as to be able to accept anything. This means that the representations you get aren't exactly the representations that the source was using with respect to bitsizes.
- `byte`s are received as `uint64`, but `[]byte` is correct. There is no type id for a single `byte` in the gob format.

## TODO

- Printing stylized output
- Some more testing (I'm around ~70%)
- Include more "bad gob" tests, but to be honest this tool shouldn't be being used on bad gobs so this is pretty low priority and part of the reason I'm OK
with some of the panics and lack of testing around bad gobs.