package anvil

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"

	"github.com/beito123/level"
	"github.com/beito123/nbt"
)

// Anvil is a level format
// It often is used for minecraft java edition and server world
type Anvil struct {
	regions map[uint64]*Region
}

// getRegion returns a region with xy
// If the Region doesn't exist, ok is false
func (lvl *Anvil) getRegion(x, y int) (r *Region, ok bool) {
	r, ok = lvl.regions[lvl.toRegionKey(x, y)]

	return r, ok
}

// toCC returns id for container from region coordinate
func (Anvil) toRegionKey(x, y int) uint64 {
	return uint64(int64(x)<<32 | int64(y))
}

// chunkToRegion returns region coordinate from chunk coordinate
func (Anvil) chunkToRegion(x, y int) (cx, cy int) {
	return x >> 5, y >> 5
}

// LoadChunk loads a chunk.
// If create is true, generates a chunk.
func (lvl *Anvil) LoadChunk(x, y int, create bool) bool {
	return false
}

// UnloadChunk unloads a chunk.
// If safe is false, unloads even when players are there.
func (lvl *Anvil) UnloadChunk(x, y int, safe bool) bool {
	return false
}

// IsLoadedChunk returns weather a chunk is loaded.
func (lvl *Anvil) IsLoadedChunk(x, y int) bool {
	return false
}

// SaveChunk saves a chunk.
func (lvl *Anvil) SaveChunk(x, y int) bool {
	return false
}

// SaveChunks saves all chunks.
func (lvl *Anvil) SaveChunks() {

}

// Chunk returns a chunk.
// If it's not loaded, loads the chunks.
// If create is true, generates a chunk.
func (lvl *Anvil) Chunk(x, y int, create bool) level.Chunk {
	return nil
}

// Chunks retuns loaded chunks.
func (lvl *Anvil) Chunks() []level.Chunk {
	return nil
}

// NewChunk returns new Chunk
func NewChunk(x, y int) (*Chunk, error) {
	return &Chunk{
		x: x,
		y: y,
	}, nil
}

// ReadChunk returns new Chunk with CompoundTag
func ReadChunk(tag *nbt.Compound) (*Chunk, error) {
	dump, err := tag.ToString()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Dump: %s", dump)

	com, err := tag.GetCompound("Level")
	if err != nil {
		return nil, err
	}

	x, err := com.GetInt("xPos")
	if err != nil {
		return nil, err
	}

	y, err := com.GetInt("zPos")
	if err != nil {
		return nil, err
	}

	return NewChunk(int(x), int(y))
}

// Chunk is
type Chunk struct {
	x int
	y int
}

func (chunk *Chunk) X() int {
	return chunk.x
}

func (chunk *Chunk) Y() int {
	return chunk.y
}
