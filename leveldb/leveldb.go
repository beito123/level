package leveldb

/*
	level

	Copyright (c) 2019 beito

	This software is released under the MIT License.
	http://opensource.org/licenses/mit-license.php
*/

import (
	"fmt"
	"strings"
	"sync"

	lvldb "github.com/beito123/goleveldb/leveldb"
	"github.com/beito123/goleveldb/leveldb/filter"
	"github.com/beito123/goleveldb/leveldb/opt"
	"github.com/beito123/level"
)

// DefaultOptions is a default option for leveldb
// You can use at NewWithOptions() and LoadWithOptions()
var DefaultOptions = &opt.Options{
	Filter:      filter.NewBloomFilter(10),
	WriteBuffer: 4 * 1024 * 1024, // 4mb
}

// New returns new LevelDB
// The path is a directory for save
func New(path string) (*LevelDB, error) {
	return NewWithOptions(path, DefaultOptions)
}

// NewWithOptions returns new LevelDB with leveldb options
func NewWithOptions(path string, options *opt.Options) (*LevelDB, error) {
	db, err := lvldb.OpenFile(path, options)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		Database: db,
		Format:   &ChunkFormatV100{},
		chunks:   make(map[int]*Chunk),
		mutex:    new(sync.RWMutex),
	}, nil
}

// Load loads a leveldb level
func Load(path string) (*LevelDB, error) {
	return LoadWithOptions(path, DefaultOptions)
}

// LoadWithOptions loads a leveldb level with leveldb options
func LoadWithOptions(path string, options *opt.Options) (*LevelDB, error) {
	db, err := lvldb.OpenFile(path, options)
	if err != nil {
		return nil, err
	}

	return &LevelDB{
		Database: db,
		Format:   &ChunkFormatV100{},
		chunks:   make(map[int]*Chunk),
		mutex:    new(sync.RWMutex),
	}, nil
}

// NewRawBlockState returns new RawBlockState
func NewRawBlockState(name string, value int) *RawBlockState {
	return &RawBlockState{
		name:  strings.ToLower(name),
		value: value,
	}
}

// FromRawBlockState returns new RawBlockState
func FromRawBlockState(bs level.BlockState) (*RawBlockState, error) {
	name, meta, ok := bs.ToBlockNameMeta()
	if !ok {
		return nil, fmt.Errorf("level.leveldb: usable to convert from %s to RawBlockState", bs.Name())
	}

	return NewRawBlockState(name, meta), nil
}

// RawBlockState is a raw block information
type RawBlockState struct {
	name  string
	value int
}

// Name returns block name
func (block *RawBlockState) Name() string {
	return block.name
}

// Value returns block value
func (block *RawBlockState) Value() int {
	return block.value
}

// Equal returns whether block is equal b
func (block *RawBlockState) Equal(b *RawBlockState) bool {
	return block.name == b.name && block.value == b.value
}

// ToBlockNameProperties returns block name and properties
// If it's not supported, returns false for ok
func (block *RawBlockState) ToBlockNameProperties() (name string, properties map[string]string, ok bool) {
	return block.name, make(map[string]string), true
}

// ToBlockNameMeta returns block name and meta
// If it's not supported, returns false for ok
func (block *RawBlockState) ToBlockNameMeta() (name string, meta int, ok bool) {
	return block.name, block.value, true
}

// ToBlockIDMeta returns block id and meta
// If it's not supported, returns false for ok
func (block *RawBlockState) ToBlockIDMeta() (id int, meta int, ok bool) {
	return 0, 0, false
}

// LevelDB is a level format for mcbe
type LevelDB struct {
	Database *lvldb.DB

	Format    ChunkFormat
	Dimension level.Dimension

	chunks map[int]*Chunk

	mutex *sync.RWMutex
}

func (LevelDB) at(x, y int) int {
	return y<<16 | x
}

// Close closes database of leveldb
// You must close after you use the format
// It's should not run other functions after format is closed
func (lvl *LevelDB) Close() error {
	if lvl.Database != nil {
		lvl.mutex.Lock()
		err := lvl.Database.Close()
		lvl.mutex.Unlock()
		return err
	}

	return nil
}

// LoadChunk loads a chunk.
// If create is true, generates a chunk.
func (lvl *LevelDB) LoadChunk(x, y int) error {
	if lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: already loaded the chunk")
	}

	chunk, err := lvl.Format.Read(lvl.Database, x, y, lvl.Dimension)
	if err != nil {
		return err
	}

	lvl.mutex.Lock()
	lvl.chunks[lvl.at(x, y)] = chunk
	lvl.mutex.Unlock()

	return nil
}

// UnloadChunk unloads a chunk.
func (lvl *LevelDB) UnloadChunk(x, y int) error {
	if !lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: not loaded the chunk")
	}

	lvl.mutex.Lock()
	delete(lvl.chunks, lvl.at(x, y))
	lvl.mutex.Unlock()

	return nil
}

// GenerateChunk generates a chunk and loads
func (lvl *LevelDB) GenerateChunk(x, y int) error {
	if lvl.IsLoadedChunk(x, y) {
		return fmt.Errorf("level.leveldb: already loaded the chunk")
	}

	chunk := NewChunk(x, y)

	lvl.mutex.Lock()
	lvl.chunks[lvl.at(x, y)] = chunk
	lvl.mutex.Unlock()

	return nil
}

// HasGeneratedChunk returns whether the chunk is generaged
func (lvl *LevelDB) HasGeneratedChunk(x, y int) (bool, error) {
	return lvl.Format.Exist(lvl.Database, x, y, lvl.Dimension)
}

// IsLoadedChunk returns weather a chunk is loaded.
func (lvl *LevelDB) IsLoadedChunk(x, y int) bool {
	lvl.mutex.RLock()
	_, ok := lvl.chunks[lvl.at(x, y)]
	lvl.mutex.RUnlock()

	return ok
}

// SaveChunk saves a chunk.
func (lvl *LevelDB) SaveChunk(x, y int) error {
	lchunk, ok := lvl.Chunk(x, y)
	if !ok {
		return fmt.Errorf("level.leveldb: not loaded the chunk")
	}

	chunk := lchunk.(*Chunk)

	return lvl.Format.Write(lvl.Database, chunk, lvl.Dimension)
}

// SaveChunks saves all chunks.
func (lvl *LevelDB) SaveChunks() error {
	for key := range lvl.chunks {
		err := lvl.SaveChunk(key&0xffff, (key>>16)&0xffff)
		if err != nil {
			return err
		}
	}

	return nil
}

// Chunk returns a loaded chunk.
func (lvl *LevelDB) Chunk(x, y int) (level.Chunk, bool) {
	lvl.mutex.RLock()
	chunk, ok := lvl.chunks[lvl.at(x, y)]
	lvl.mutex.RUnlock()

	return chunk, ok
}

// LoadedChunks returns loaded chunks.
func (lvl *LevelDB) LoadedChunks() []level.Chunk {
	result := make([]level.Chunk, len(lvl.chunks))

	lvl.mutex.RLock()

	count := 0
	for _, chunk := range lvl.chunks {
		result[count] = chunk
		count++
	}

	lvl.mutex.RUnlock()

	return result
}
