package level

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

// Level is a simple level loader
// I will implement later // support high level access
type Level struct {
	// Manage level provider and more...
	// Load level.dat and manage level data (name, spaen and etc..)
}

// Format is a simple interface for level formats
type Format interface {

	// LoadChunk loads a chunk.
	LoadChunk(x, y int) error

	// UnloadChunk unloads a chunk.
	UnloadChunk(x, y int) error

	// GenerateChunk generates a chunk
	GenerateChunk(x, y int) error

	// HasGeneratedChunk returns whether the chunk is generaged
	HasGeneratedChunk(x, y int) bool

	// IsLoadedChunk returns weather a chunk is loaded.
	IsLoadedChunk(x, y int) bool

	// SaveChunk saves a chunk.
	SaveChunk(x, y int) error

	// SaveChunks saves all chunks.
	SaveChunks() error

	// Chunk returns a loaded chunk.
	Chunk(x, y int) (Chunk, error)

	// LoadedChunks returns loaded chunks.
	LoadedChunks() []Chunk
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

	// GetBlock gets a BlockState at chunk coordinate
	GetBlock(x, y, z int) (BlockState, error)

	// SetBlock set a BlockState at chunk coordinate
	SetBlock(x, y, z int, state BlockState) error
}

// BlockState is a block information
type BlockState interface {
	Name() string
}
