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
// This needs to be supported goroutine (concurrency)
type Format interface {

	// Close closes the level format
	// You must close after you use the format
	// It's should not run other functions after format is closed
	Close() error

	// LoadChunk loads a chunk.
	LoadChunk(x, y int) error

	// UnloadChunk unloads a chunk.
	UnloadChunk(x, y int) error

	// GenerateChunk generates a chunk
	GenerateChunk(x, y int) error

	// HasGeneratedChunk returns whether the chunk is generaged
	HasGeneratedChunk(x, y int) (bool, error)

	// IsLoadedChunk returns weather a chunk is loaded.
	IsLoadedChunk(x, y int) bool

	// SaveChunk saves a chunk.
	SaveChunk(x, y int) error

	// SaveChunks saves all chunks.
	SaveChunks() error

	// Chunk returns a loaded chunk.
	Chunk(x, y int) (Chunk, bool)

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

	// Name returns block name
	Name() string

	// ToBlockNameProperties returns block name and properties
	// If it's not supported, returns false for ok
	ToBlockNameProperties() (name string, properties map[string]string, ok bool)

	// ToBlockNameMeta returns block name and meta
	// If it's not supported, returns false for ok
	ToBlockNameMeta() (name string, meta int, ok bool)

	// ToBlockIDMeta returns block id and meta
	// If it's not supported, returns false for ok
	ToBlockIDMeta() (id int, meta int, ok bool)
}
