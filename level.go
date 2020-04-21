package level

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"github.com/beito123/nbt"
)

// Level is a simple level loader
type Level struct {
}

// Format is a simple interface for level formats
// This needs to be supported concurrency
type Format interface {

	// Name returns name of level
	Name() string

	// SetName sets the name of level
	SetName(name string)

	// GameType returns the default game mode of level
	GameType() GameType

	// SetGameType sets the game mode of level
	SetGameType(typ GameType)

	// Spawn returns the default spawn of level
	Spawn() (x, y, z int)

	// SetSpawn sets the default spawn of level
	SetSpawn(x, y, z int)

	// Property returns a property of level.dat
	Property(name string) (tag nbt.Tag, ok bool)

	// SetProperty sets a property
	SetProperty(tag nbt.Tag)

	// AllProperties returns all properties
	AllProperties() *nbt.Compound

	// SetAllProperties sets all properties
	SetAllProperties(com *nbt.Compound)

	// Close closes the level format
	// You must close after you use the format
	// It's should not run other functions after format is closed
	Close() error

	// Dimension return dimension of the level
	Dimension() Dimension

	// SetDimension set dimension of the level
	SetDimension(Dimension)

	// LoadChunk loads a chunk.
	// If create is enabled, generates a chunk if it doesn't exist
	LoadChunk(x, y int, create bool) error

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

	// Chunk returns a chunk.
	// If a chunk is not loaded, it will be loaded
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

	// Height returns the height of the highest block at chunk coordinate
	Height(x, y int) (height uint16)

	// Biome returns biome
	Biome(x, y int) byte

	// SetBiome set biome
	SetBiome(x, y int, biome byte)

	// Entities returns entities of nbt data
	Entities() []*nbt.Compound

	// SetEntities set entities of nbt data
	SetEntities(entities []*nbt.Compound)

	// BlockEntities returns block entities of nbt data
	BlockEntities() []*nbt.Compound

	// SetBlockEntities set block entities of nbt data
	SetBlockEntities(entities []*nbt.Compound)

	// GetBlock gets a BlockState at chunk coordinate
	GetBlock(x, y, z int) (BlockState, error)

	// SetBlock set a BlockState at chunk coordinate
	SetBlock(x, y, z int, state BlockState) error
}

// BlockState is a block information
type BlockState interface { // TODO: fix bad codes :P

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
