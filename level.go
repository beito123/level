package level

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

// Level is a simple level loader
// I will implement later
type Level struct {
	// Manage level provider and more...
	// Load level.dat and manage level data (name, spaen and etc..)
}

// LevelFormat is a simple interface for level formats
type LevelFormat interface {

	// LoadChunk loads a chunk.
	// If create is true, generates a chunk.
	LoadChunk(x, y int, create bool) bool

	// UnloadChunk unloads a chunk.
	// If safe is false, unloads even when players are there.
	UnloadChunk(x, y int, safe bool) bool

	// IsLoadedChunk returns weather a chunk is loaded.
	IsLoadedChunk(x, y int) bool

	// SaveChunk saves a chunk.
	SaveChunk(x, y int) bool

	// SaveChunks saves all chunks.
	SaveChunks()

	// Chunk returns a chunk.
	// If it's not loaded, loads the chunks.
	// If create is true, generates a chunk.
	Chunk(x, y int, create bool) Chunk

	// Chunks retuns loaded chunks.
	Chunks() []Chunk
}

// Chunk is a simple interface for chunk
type Chunk interface {

	// X returns x coordinate
	X() int

	// Y returns y coordinate
	Y() int

	// SetX set x coordinate
	SetX(x int)

	// SetY set y coordinate
	SetY(y int)

	// GetBlockID gets a block id on xyz
	GetBlockID(x, y, z int) (id int, err error)

	// GetBlockData gets a block id on xyz
	GetBlockData(x, y, z int) (data int, err error)

	// SetBlockID set a block on xyz
	SetBlockID(x, y, z, id int) error

	// SetBlockData set a block on xyz
	SetBlockData(x, y, z, id int) error
}
