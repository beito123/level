package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	"github.com/beito123/level/block"

	"github.com/beito123/nbt"
)

// NewChunk returns new Chunk
func NewChunk(x, y int, format ChunkFormat) *Chunk {
	return &Chunk{
		x:           x,
		y:           y,
		subChunks:   make([]*SubChunk, 16),
		ChunkFormat: format,
	}
}

// ReadChunk returns new Chunk with CompoundTag
func ReadChunk(x, y int, b []byte) (*Chunk, error) {
	stream := nbt.NewStreamBytes(nbt.BigEndian, b)

	tag, err := stream.ReadTag()
	if err != nil {
		return nil, err
	}

	com, ok := tag.(*nbt.Compound)
	if !ok {
		return nil, fmt.Errorf("level.anvil.region: expected to be CompoundTag, but it passed different tag(%sTag)", nbt.GetTagName(tag.ID()))
	}

	var format ChunkFormat = &ChunkFormatV112{}
	if com.Has("DataVersion") {
		ver, err := com.GetInt("DataVersion")
		if err != nil {
			return nil, err
		}

		switch {
		case ver >= 1519:
			format = &ChunkFormatV113{}
		}
	}
	return format.Read(com)
}

// Chunk is
type Chunk struct {
	x int
	y int

	lastUpdate    int64
	inhabitedTime int64
	biomes        []int
	subChunks     []*SubChunk

	ChunkFormat ChunkFormat
}

// X returns x coordinate
func (chunk *Chunk) X() int {
	return chunk.x
}

// Y returns y coordinate
func (chunk *Chunk) Y() int {
	return chunk.y
}

// SubChunks returns sub chunks
func (chunk *Chunk) SubChunks() []*SubChunk {
	return chunk.subChunks
}

// GetSubChunk returns a sub chunk at the y index
// you can set 0-15 at y
func (chunk *Chunk) GetSubChunk(y int) (*SubChunk, bool) {
	if len(chunk.subChunks) >= y {
		return nil, false
	}

	return chunk.subChunks[y], true
}

// AtSubChunk returns a sub chunk at the y (chunk coordinate)
func (chunk *Chunk) AtSubChunk(y int) (*SubChunk, bool) {
	index := y / 16
	if len(chunk.subChunks) >= index {
		return nil, false
	}

	return chunk.subChunks[index], true
}

// GetBlock gets a block at the xyz (chunk coordinate)
func (chunk *Chunk) GetBlock(x, y, z int) (*block.BlockData, error) {
	sub, ok := chunk.AtSubChunk(y)
	if !ok {
		return nil, nil
	}

	bl, err := sub.AtBlock(x, y, z)
	if err != nil {
		return nil, nil
	}

	return bl.ToBlockData(), nil
}

// Save saves the chunk, returns CompoundTag
func (chunk *Chunk) Save() (*nbt.Compound, error) {
	return nil, nil // TODO
}

// ChunkFormat is a chunk format for a version
type ChunkFormat interface {
	Read(tag *nbt.Compound) (*Chunk, error)
}

// ChunkFormatV112 is a chunk format for v1.12
type ChunkFormatV112 struct {
}

func (format *ChunkFormatV112) Read(tag *nbt.Compound) (*Chunk, error) {
	chunk := NewChunk(0, 0, format)

	com, err := tag.GetCompound("Level")
	if err != nil {
		return nil, err
	}

	xPos, err := com.GetInt("xPos")
	if err != nil {
		return nil, err
	}

	zPos, err := com.GetInt("zPos")
	if err != nil {
		return nil, err
	}

	chunk.x = int(xPos)
	chunk.y = int(zPos)

	// Biomes
	if com.Has("Biomes") {
		biomes, err := com.GetByteArray("Biomes")
		if err != nil {
			return nil, err
		}

		chunk.biomes = make([]int, len(biomes))
		for i, biome := range biomes {
			chunk.biomes[i] = int(biome)
		}
	}

	// Subchunks
	sections, err := com.GetList("Sections")
	if err != nil {
		return nil, err
	}

	subchunkFormat := &SubChunkFormatV112{}

	chunk.subChunks = make([]*SubChunk, 16)
	for _, entry := range sections {
		sec, ok := entry.(*nbt.Compound)
		if !ok {
			return nil, fmt.Errorf("couldn't convert to *Compound")
		}

		sub, err := subchunkFormat.Read(sec)
		if err != nil {
			return nil, err
		}

		chunk.subChunks[sub.Y] = sub
	}

	return chunk, nil
}

// ChunkFormatV113 is a chunk format for v1.13
type ChunkFormatV113 struct {
}

func (format *ChunkFormatV113) Read(tag *nbt.Compound) (*Chunk, error) {
	chunk := NewChunk(0, 0, format)

	com, err := tag.GetCompound("Level")
	if err != nil {
		return nil, err
	}

	xPos, err := com.GetInt("xPos")
	if err != nil {
		return nil, err
	}

	zPos, err := com.GetInt("zPos")
	if err != nil {
		return nil, err
	}

	chunk.x = int(xPos)
	chunk.y = int(zPos)

	// Biomes
	if com.Has("Biomes") {
		biomes, err := com.GetIntArray("Biomes")
		if err != nil {
			return nil, err
		}

		chunk.biomes = make([]int, len(biomes))
		for i, biome := range biomes {
			chunk.biomes[i] = int(biome)
		}
	}

	// Subchunks
	sections, err := com.GetList("Sections")
	if err != nil {
		return nil, err
	}

	subchunkFormat := &SubChunkFormatV113{}

	chunk.subChunks = make([]*SubChunk, 16)
	for _, entry := range sections {
		sec, ok := entry.(*nbt.Compound)
		if !ok {
			return nil, fmt.Errorf("couldn't convert to *Compound")
		}

		sub, err := subchunkFormat.Read(sec)
		if err != nil {
			return nil, err
		}

		chunk.subChunks[sub.Y] = sub
	}

	return nil, nil
}
