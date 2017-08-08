# degob

degob is a reversing tool for `gob`s. If you have the binary form of a gob, but you don't have to `struct` to decode it into, this library still allows you to get a more readable representation of the data.

## cmds/degob

The easiest way to use all of this is to just build the binary in `cmds/degob` and send gobs to it either through `stdin` or from files and then get the output to `stdout` or to a file. See its [README](cmds/degob/README.md) for more info.

## Usage

Create a new `Decoder` over your reader using `NewDecoder` and then decode that into a slice of `Gob`s with `Decode` or stream `Gob`s with `DecodeStream`. `DecodeStream` isn't fully tested yet and will probably still fumble with errors. Once you have `Gob`s you can either play with the types directly or just print them out to a writer using the `WriteTypes` and `WriteValues` methods.

The provided `degob` command provides a straightforward [sample usage](cmds/degob/main.go).

## TODO

- Printing stylized output isn't complete yet.
- A lot more testing
- Documentation
- Errors maybe could be handled better and I haven't tested if the byte count on error is correct or not yet