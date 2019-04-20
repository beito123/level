package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	"github.com/beito123/level"
)

// DefaultStorageIndex is the default index for StorageIndex
const DefaultStorageIndex = 0

// NewChunk returns new Chunk
func NewChunk(x, y int) *Chunk {
	return &Chunk{
		x:                   x,
		y:                   y,
		biomes:              make([]byte, 256),
		subChunks:           make([]*SubChunk, 16),
		DefaultStorageIndex: 0,
	}
}

// Chunk is a block area which splits a world by 16x16
// It has informations of block, biomes and etc...
type Chunk struct {
	x         int
	y         int
	biomes    []byte
	subChunks []*SubChunk

	DefaultStorageIndex int
}

// X returns x coordinate
func (chunk *Chunk) X() int {
	return chunk.x
}

// Y returns y coordinate
func (chunk *Chunk) Y() int {
	return chunk.y
}

// SetX set x coordinate
func (chunk *Chunk) SetX(x int) {
	chunk.x = x
}

// SetY set y coordinate
func (chunk *Chunk) SetY(y int) {
	chunk.y = y
}

// SubChunks returns sub chunks
func (chunk *Chunk) SubChunks() []*SubChunk {
	return chunk.subChunks
}

// GetSubChunk returns a sub chunk
func (chunk *Chunk) GetSubChunk(index int) (*SubChunk, bool) {
	if index >= len(chunk.subChunks) {
		return nil, false
	}

	return chunk.subChunks[index], true
}

// AtSubChunk returns a sub chunk
func (chunk *Chunk) AtSubChunk(y int) (*SubChunk, bool) {
	return chunk.GetSubChunk(y / 16)
}

// Vaild vailds a chunk coordinates
func (chunk *Chunk) Vaild(x, y, z int) bool {
	return x >= 0 && x <= 15 && y >= 0 && y <= 15 && z >= 0 && z <= 15
}

// GetBlock gets a BlockState at a chunk coordinate
// if both returned values are nil, maybe a block is air block
func (chunk *Chunk) GetBlock(x, y, z int) (level.BlockState, error) {
	return chunk.GetBlockAtStorage(x, y, z, chunk.DefaultStorageIndex)
}

// GetBlockAtStorage gets a BlockState at a chunk coordinate from storage of index
// if both returned values are nil, maybe a block is air block
func (chunk *Chunk) GetBlockAtStorage(x, y, z, index int) (*BlockState, error) {
	if chunk.Vaild(x, y, z) {
		return nil, fmt.Errorf("invaild chunk coordinate")
	}

	sub, ok := chunk.AtSubChunk(y)
	if !ok {
		return nil, nil // Air
	}

	return sub.AtBlock(x&15, y&15, z&15, index)
}

// SetBlock set a BlockState at chunk coordinate
func (chunk *Chunk) SetBlock(x, y, z int, state level.BlockState) error {
	return nil
}
