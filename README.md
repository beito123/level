# level

[WIP] I may throw away :P

## Installation

You can get the package with go get command.

    go get -u github.com/beito123/level

*-u is a option updating the package*

## License

This is licensed by MIT License. See LICENSE file.

## Examples

### Read

```go
func main() {
	// Load leveldb level for mcbe
	lvl, err := leveldb.Load("./db")
	if err != nil {
		panic(err)
	}

	// Chunk coordinates
	x := 0
	y := 0

	ok, err := lvl.HasGeneratedChunk(x, y)
	if err != nil {
		panic(err)
	}

	if !ok {
		panic("a chunk isn't generated")
	}

	// Get chunk
	chunk, err := lvl.Chunk(x, y)
	if err != nil {
		panic(err)
	}

	for y := 0; y < 256; y++ {
		for z := 0; z < 16; z++ {
			for x := 0; x < 16; x++ {
				b, err := chunk.GetBlock(x, y, z)
				if err != nil {
					panic(err)
				}

				if b.Name() == "minecraft:air" { // ignore air
					continue
				}

				fmt.Printf("%s (at %d, %d, %d)\n", b.Name(), x, y, z)
			}
		}
	}
}
```